package user_manager

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestUserManager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UserManager Suite")
}
