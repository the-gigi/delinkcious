package social_graph_manager

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSocialGraphManager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SocialGraphManager Suite")
}
