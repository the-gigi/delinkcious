package link_checker

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLinkChecker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LinkChecker Suite")
}
