package user_manager

import (
	sq "github.com/Masterminds/squirrel"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
)

var _ = Describe("user manager tests with DB ", func() {
	var userStore *DbUserStore
	var userManager om.UserManager
	var err error
	var deleteAll = func() {
		sq.Delete("users").RunWith(userStore.db).Exec()
		sq.Delete("sessions").RunWith(userStore.db).Exec()
	}

	BeforeSuite(func() {
		userStore, err = NewDbUserStore("localhost", 5432, "postgres", "postgres")
		Ω(err).Should(BeNil())
		userManager, err = NewUserManager(userStore)
		Ω(err).Should(BeNil())
		Ω(userManager).ShouldNot(BeNil())
	})

	BeforeEach(func() {
		deleteAll()
	})

	AfterSuite(func() {
		deleteAll()
	})

	It("should register new user", func() {
		err := userManager.Register(om.User{"gg@gg.com", "user"})
		Ω(err).Should(BeNil())
	})

	It("should fail to register user with empty name", func() {
		err := userManager.Register(om.User{"gg@gg.com", ""})
		Ω(err).ShouldNot(BeNil())
	})

	It("should fail to register existing user", func() {
		err := userManager.Register(om.User{"gg@gg.com", "user"})
		Ω(err).Should(BeNil())

		err = userManager.Register(om.User{"gg@gg.com", "user"})
		Ω(err).ShouldNot(BeNil())
	})
})
