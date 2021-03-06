package credentials_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/credentials"
	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/store/memory"
)

var _ = Describe("CredentialService", func() {
	const (
		sourceKey          = "source"
		sourceOnlyKey      = "sourceOnly"
		destinationKey     = "destination"
		destinationOnlyKey = "destinationOnly"
		bothKey            = "both"
	)
	var (
		credentialService credentials.Service
		sourceStore       store.Store
		destinationStore  store.Store
	)
	BeforeEach(func() {

		cfg := &config.Config{
			Stores: []config.Store{
				config.Store{
					Name:      sourceKey,
					StoreType: "memory",
				},
				config.Store{
					Name:      destinationKey,
					StoreType: "memory",
				},
			},
		}
		graph, err := config.NewConfigurationGraph(cfg)
		Expect(err).To(BeNil())

		manager := store.NewManager()

		sourceStore = memory.NewMemoryStore(sourceKey)
		sourceStore.Set(store.NewValueItem(sourceOnlyKey, sourceOnlyKey))
		sourceStore.Set(store.NewValueItem(bothKey, sourceOnlyKey))

		destinationStore = memory.NewMemoryStore(destinationKey)
		destinationStore.Set(store.NewValueItem(destinationOnlyKey, destinationOnlyKey))
		destinationStore.Set(store.NewValueItem(bothKey, destinationOnlyKey))

		manager.Register(memory.NewMemoryStoreProvider(sourceStore, destinationStore))

		credentialService, err = credentials.NewService(cfg, graph, manager)

		Expect(err).To(BeNil())
		Expect(credentialService).ToNot(BeNil())
	})
	Describe("Move", func() {
		Context("when source exists and target exists", func() {
			Context("when source credential exists and target credential doesn't exist", func() {
				It("moves the credential", func() {
					err := credentialService.Move(sourceKey, sourceOnlyKey, destinationKey, destinationOnlyKey)
					Expect(err).To(BeNil())

					item, err := destinationStore.Get(destinationOnlyKey)
					Expect(err).To(BeNil())
					Expect(item).ToNot(BeNil())
					Expect(item.Value()).To(Equal(sourceOnlyKey))

					_, err = sourceStore.Get(sourceOnlyKey)
					Expect(err).ToNot(BeNil())
				})
			})
			Context("when source credential exists and target credential exists", func() {
				It("overwrites the target", func() {
					err := credentialService.Move(sourceKey, bothKey, destinationKey, bothKey)
					Expect(err).To(BeNil())

					item, err := destinationStore.Get(bothKey)
					Expect(err).To(BeNil())
					Expect(item).ToNot(BeNil())
					Expect(item.Value()).To(Equal(sourceOnlyKey))

					_, err = sourceStore.Get(bothKey)
					Expect(err).ToNot(BeNil())
				})
			})
		})
		Context("when source exists and target doesn't exist", func() {
			It("fails", func() {

			})
		})
		Context("When source is empty", func() {
			It("fails", func() {
				err := credentialService.Move("", "", "test2", "")
				Expect(err).ToNot(BeNil())
			})
		})
		Context("when source doesn't exist and target exists", func() {
			It("fails", func() {
				err := credentialService.Move("s", "some/path", "destination", "some/path")
				Expect(err).ToNot(BeNil())
			})
		})
	})
	Describe("Copy", func() {
		Context("when source exists and target exists", func() {
			It("copies the item", func() {

			})
		})
		Context("when source exists and target doesn't exist", func() {
			It("fails", func() {

			})
		})
		Context("When source is empty", func() {
			It("fails", func() {
				err := credentialService.Copy("", "", "test2", "")
				Expect(err).ToNot(BeNil())
			})
		})
		Context("when source doesn't exist and target exists", func() {
			It("fails", func() {
				err := credentialService.Copy("test1", "", "test2", "")
				Expect(err).ToNot(BeNil())
			})
		})
	})
	Describe("Set", func() {
		It("sets the item", func() {
			err := credentialService.Set(sourceKey, store.NewValueItem("test", "something"))
			Expect(err).To(BeNil())
			item, err := sourceStore.Get("test")
			Expect(err).To(BeNil())
			s, ok := item.Value().(string)
			Expect(ok).To(BeTrue())
			Expect(s).To(Equal("something"))
		})
	})
	Describe("Get", func() {
		It("gets the item", func() {
			item, err := credentialService.Get(sourceKey, sourceOnlyKey)
			Expect(err).To(BeNil())
			s, ok := item.Value().(string)
			Expect(ok).To(BeTrue())
			Expect(s).To(Equal(sourceOnlyKey))
		})
	})
	Describe("List", func() {
		It("lists the items", func() {
			items, err := credentialService.List(sourceKey, "/")
			Expect(err).To(BeNil())
			Expect(len(items)).To(Equal(2))
		})
	})
})
