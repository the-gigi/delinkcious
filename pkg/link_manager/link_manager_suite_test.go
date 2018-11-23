package link_manager_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLinkManager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LinkManager Suite")
}
