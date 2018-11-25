package social_graph_manager

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
)

var _ = Describe("user manager tests", func() {
	var socialGraphManager om.SocialGraphManager
	var err error

	BeforeEach(func() {
		store := NewInMemorySocialGraphStore()
		socialGraphManager, err = NewSocialGraphManager(store)
		Ω(err).Should(BeNil())
		Ω(socialGraphManager).ShouldNot(BeNil())
	})

	It("should follow", func() {
		err := socialGraphManager.Follow("jack", "")
		Ω(err).ShouldNot(BeNil())

		err = socialGraphManager.Follow("", "jack")
		Ω(err).ShouldNot(BeNil())

		err = socialGraphManager.Follow("john", "jack")
		Ω(err).Should(BeNil())

		// Already following
		err = socialGraphManager.Follow("john", "jack")
		Ω(err).ShouldNot(BeNil())

	})

	It("should unfollow", func() {
		err = socialGraphManager.Unfollow("john", "jack")
		// Can't unfollow if not following
		Ω(err).ShouldNot(BeNil())

		// Follow
		err = socialGraphManager.Follow("john", "jack")
		Ω(err).Should(BeNil())

		// And then unfollow
		err = socialGraphManager.Unfollow("john", "jack")
		Ω(err).Should(BeNil())
	})

	It("should get followers and following", func() {
		// Follow
		err = socialGraphManager.Follow("john", "jack")
		Ω(err).Should(BeNil())

		followers := socialGraphManager.GetFollowers("john")
		Ω(followers).Should(HaveLen(1))

		following := socialGraphManager.GetFollowing("john")
		Ω(following).Should(HaveLen(0))

		followers = socialGraphManager.GetFollowers("jack")
		Ω(followers).Should(HaveLen(0))

		following = socialGraphManager.GetFollowing("jack")
		Ω(following).Should(HaveLen(1))
	})

})
