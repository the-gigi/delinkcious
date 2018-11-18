package user_manager

import (
	om "github.com/the-gigi/delinkcious/pkg/object_model"
	//om "../object_model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("user manager tests", func() {
	var userManager om.UserManager
	BeforeEach(func() {
		userManager = NewImMemoryUserManager()
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

