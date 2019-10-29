package packages

import (
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/patrickhuber/wrangle/collections"
	"github.com/patrickhuber/wrangle/feed"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/tasks"
	"github.com/patrickhuber/wrangle/templates"
)

type manager struct {
	fileSystem      filesystem.FileSystem
	taskProviders   tasks.ProviderRegistry
	contextProvider ContextProvider
	feedService     feed.FeedService
}

// Manager defines a manager interface
type Manager interface {
	Install(p Package) error
	Load(packageName, packageVersion string) (Package, error)
}

// NewManager creates a new package manager
func NewManager(fileSystem filesystem.FileSystem, feedService feed.FeedService, contextProvider ContextProvider, taskProviders tasks.ProviderRegistry) Manager {
	return &manager{
		fileSystem:      fileSystem,
		feedService:     feedService,
		taskProviders:   taskProviders,
		contextProvider: contextProvider}
}

func (manager *manager) Install(p Package) error {
	for _, task := range p.Tasks() {
		provider, err := manager.taskProviders.Get(task.Type())
		if err != nil {
			return err
		}
		err = provider.Execute(task, p.Context())
		if err != nil {
			return err
		}
	}
	return nil
}

func (manager *manager) Load(packageName, packageVersion string) (Package, error) {
	resp, err := manager.feedService.Get(&feed.FeedGetRequest{
		Name:           packageName,
		Version:        packageVersion,
		IncludeContent: true,
	})

	if err != nil {
		return nil, err
	}

	if resp == nil || resp.Package == nil || resp.Package.Versions == nil || len(resp.Package.Versions) == 0 {
		return nil, fmt.Errorf("unable to find package %s version %s", packageName, packageVersion)
	}

	version := resp.Package.Versions[0]
	if resp.Package.Latest != "" {
		for _, v := range resp.Package.Versions {
			if v.Version == resp.Package.Latest {
				version = v
			}
		}
	}

	if version == nil {
		return nil, fmt.Errorf("package is missing latest version")
	}
	if version.Manifest == nil {
		return nil, fmt.Errorf("package is missing latest version manifest")
	}
	if version.Manifest.Content == "" {
		return nil, fmt.Errorf("package is missing latest version manifest content")
	}

	content := version.Manifest.Content
	manifest, err := NewYamlInterfaceReader(strings.NewReader(content)).Read()
	if err != nil {
		return nil, err
	}

	// validate?

	// interpolate package
	manifest, err = manager.interpolatePackageManifest(manifest, map[string]string{
		"/version": packageVersion,
	})
	if err != nil {
		return nil, err
	}

	packageContext, err := manager.contextProvider.Get(packageName, packageVersion)
	if err != nil {
		return nil, err
	}

	// turn package manifest into packages.Package
	// return package
	return manager.convertManifestToPackage(manifest, packageContext)
}

func (manager *manager) interpolatePackageManifest(pkg interface{}, values map[string]string) (interface{}, error) {

	template := templates.NewTemplate(pkg)
	dictionary := collections.NewDictionaryFromMap(values)
	resolver := templates.NewDictionaryResolver(dictionary)

	return template.Evaluate(resolver)
}

func (manager *manager) convertManifestToPackage(manifest interface{}, packageContext PackageContext) (Package, error) {
	pkg := &Manifest{}

	// convert to config structure
	err := mapstructure.Decode(manifest, pkg)
	if err != nil {
		return nil, err
	}

	// convert task list
	taskList := []tasks.Task{}
	for _, target := range pkg.Targets {
		for _, task := range target.Tasks {
			tsk, err := manager.taskProviders.Decode(task)
			if err != nil {
				return nil, err
			}
			taskList = append(taskList, tsk)
		}
	}

	// convert package metadata
	return New(pkg.Name, pkg.Version, packageContext, taskList...), nil
}
