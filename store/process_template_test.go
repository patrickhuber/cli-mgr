package store_test

import (
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/store/memory"	
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/store/file"	

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ProcessTemplate", func() {
	It("can evaluate single store", func() {
		data := `
stores:
- name: one
  type: memory
processes:
- name: go
  stores:
  - one
  path: go
  args:
  - ((version))`
		cfg, err := config.DeserializeConfigString(data)
		Expect(err).To(BeNil())

		memoryStore := memory.NewMemoryStore("one")
		memoryStore.Set(store.NewItem("version", "value", "version"))
		provider := memory.NewMemoryStoreProvider(memoryStore)
		
		manager := store.NewManager()
		manager.Register(provider)

		template, err := store.NewProcessTemplate(cfg, manager)
		Expect(err).To(BeNil())

		processName := "go"
		evaluated, err := template.Evaluate(processName)
		Expect(err).To(BeNil())

		Expect(evaluated).ToNot(BeNil())
		Expect(evaluated.Args[0]).To(Equal("version"))
	})

	It("can resolve store parameters", func() {
		content := `
stores:
- name: one
  type: file
  stores:
  - two
  params:
    path: ((/file-name))
- name: two
  type: file
  params:
    path: /test2
processes:
- name: echo
  stores:
  - one
  args:
  - ((/key))
`
		configuration, err := config.DeserializeConfigString(content)
		Expect(err).To(BeNil())

		fileSystem := filesystem.NewMemory()
		fileSystem.Write("/test1", []byte("key: value"), 0644)
		fileSystem.Write("/test2", []byte("file-name: /test1"), 0644)

		manager := store.NewManager()
		manager.Register(file.NewFileStoreProvider(fileSystem, nil))
		
		template, err := store.NewProcessTemplate(configuration, manager)
		Expect(err).To(BeNil())
		environment, err := template.Evaluate("echo")
		Expect(err).To(BeNil())
		Expect(len(environment.Args)).To(Equal(1))
		Expect(environment.Args[0]).To(Equal("value"))
	})

	It("can resolve process args and vars", func() {
		content := `
stores:
- name: one
  type: file
  params:
    path: /test
processes:
- name: echo
  stores:
  - one
  args:
  - ((/key))
  env:
    prop: ((/prop))
`
		configuration, err := config.DeserializeConfigString(content)
		Expect(err).To(BeNil())

		fileSystem := filesystem.NewMemory()
		fileSystem.Write( "/test", []byte("key: 1\nprop: 2"), 0644)

		manager := store.NewManager()
		manager.Register(file.NewFileStoreProvider(fileSystem, nil))
		
		template, err := store.NewProcessTemplate(configuration, manager)
		Expect(err).To(BeNil())
		environment, err := template.Evaluate("echo")
		Expect(err).To(BeNil())
		Expect(len(environment.Args)).To(Equal(1))
		Expect(environment.Args[0]).To(Equal("1"))
		Expect(len(environment.Vars)).To(Equal(1))
		Expect(environment.Vars["prop"]).To(Equal("2"))
	})

	It("can cascade config stores", func() {
		content := `
stores:
- name: one
  type: file
  params:
    path: /test1
- name: two
  type: file
  params:
    path: /test2
- name: three
  type: file
  params:
    path: /test3
processes:
- name: echo
  stores:
  - one
  - two
  - three
  args:
  - ((/key1))
`
		configuration, err := config.DeserializeConfigString(content)
		Expect(err).To(BeNil())

		fileSystem := filesystem.NewMemory()
		fileSystem.Write("/test1", []byte("key1: ((/key2))"), 0644)
		fileSystem.Write("/test2", []byte("key2: ((/key3))"), 0644)
		fileSystem.Write("/test3", []byte("key3: value"), 0644)

		manager := store.NewManager()
		manager.Register(file.NewFileStoreProvider(fileSystem, nil))

		
		template, err := store.NewProcessTemplate(configuration, manager)
		Expect(err).To(BeNil())
		environment, err := template.Evaluate("echo")
		Expect(err).To(BeNil())
		Expect(len(environment.Args)).To(Equal(1))
		Expect(environment.Args[0]).To(Equal("value"))
	})

	It("can detect loops", func() {
		content := `
stores:
- name: one
  type: file
  stores:
  - two
  params:
    path: /test1
- name: two
  type: file
  stores:
  - three
  params:
    path: /test2
- name: three
  type: file
  stores:
  - one
  params:
    path: /test3
processes:
- name: echo
  stores:
  - one
  args:
  - ((/key1))
`
		configuration, err := config.DeserializeConfigString(content)
		Expect(err).To(BeNil())

		fileSystem := filesystem.NewMemory()
		fileSystem.Write( "/test1", []byte("key1: ((/key2))"), 0644)
		fileSystem.Write( "/test2", []byte("key2: ((/key3))"), 0644)
		fileSystem.Write( "/test3", []byte("key3: value"), 0644)

		manager := store.NewManager()
		manager.Register(file.NewFileStoreProvider(fileSystem, nil))

		
		_, err = store.NewProcessTemplate(configuration, manager)
		Expect(err).ToNot(BeNil())
	})

	It("can load variables from other store", func() {
		content := `
stores:
- name: one
  type: file
  params:
    path: /one
- name: two
  type: file
  stores:
  - one
  params:
    path: ((key))
processes:
- name: a
  stores:
  - two
  env:
    A: ((a))
    B: ((b))
    C: ((c))`

		configuration, err := config.DeserializeConfigString(content)
		Expect(err).To(BeNil())

		fileSystem := filesystem.NewMemory()
		fileSystem.Write("/one", []byte("key: /two"), 0666)
		fileSystem.Write("/two", []byte("a: a\nb: b\nc: c\n"), 0666)

		manager := store.NewManager()
		manager.Register(file.NewFileStoreProvider(fileSystem, nil))

		
		template, err := store.NewProcessTemplate(configuration, manager)
		Expect(err).To(BeNil())
		p, err := template.Evaluate("a")
		Expect(err).To(BeNil())

		Expect(p.Vars["A"]).To(Equal("a"))
		Expect(p.Vars["B"]).To(Equal("b"))
		Expect(p.Vars["C"]).To(Equal("c"))
	})
})
