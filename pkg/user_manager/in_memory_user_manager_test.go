package user_manager

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
)

var _ = Describe("user manager tests", func() {
	var userManager om.UserManager
	var err error
	BeforeEach(func() {
		store := NewInMemoryUserStore()
		userManager, err = NewUserManager(store)
		Ω(err).Should(BeNil())
		Ω(userManager).ShouldNot(BeNil())
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
