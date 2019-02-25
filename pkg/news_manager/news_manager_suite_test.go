package news_manager

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLinkManager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "NewsManager Suite")
}
