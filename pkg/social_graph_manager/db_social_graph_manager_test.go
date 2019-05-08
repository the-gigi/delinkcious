package social_graph_manager

import (
	sq "github.com/Masterminds/squirrel"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/the-gigi/delinkcious/pkg/db_util"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
	"log"
)

var _ = Describe("social graph manager tests with DB", func() {
	var socialGraphStore *DbSocialGraphStore
	var socialGraphManager om.SocialGraphManager
	var err error

	var deleteAll = func() {
		sq.Delete("social_graph").RunWith(socialGraphStore.db).Exec()
	}

	BeforeSuite(func() {
		var err error
		dbHost, dbPort, err := db_util.GetDbEndpoint("scoial_graph_manager")
		Ω(err).Should(BeNil())

		socialGraphStore, err = NewDbSocialGraphStore(dbHost, dbPort, "postgres", "postgres")
		if err != nil {
			_, err = db_util.RunLocalDB("postgres")
			Ω(err).Should(BeNil())
			if err != nil {
				log.Fatal(err)
			}

			socialGraphStore, err = NewDbSocialGraphStore(dbHost, dbPort, "postgres", "postgres")
			Ω(err).Should(BeNil())
			if err != nil {
				log.Fatal(err)
			}
		}

		Ω(err).Should(BeNil())
		Ω(socialGraphStore).ShouldNot(BeNil())

		socialGraphManager, err = NewSocialGraphManager(socialGraphStore)
		Ω(err).Should(BeNil())
		Ω(socialGraphManager).ShouldNot(BeNil())
	})

	BeforeEach(func() {
		deleteAll()
	})

	AfterSuite(func() {
		if socialGraphStore != nil && socialGraphStore.db != nil {
			deleteAll()
		}
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

		followers, err := socialGraphManager.GetFollowers("john")
		Ω(err).Should(BeNil())
		Ω(followers).Should(HaveLen(1))

		following, err := socialGraphManager.GetFollowing("john")
		Ω(err).Should(BeNil())
		Ω(following).Should(HaveLen(0))

		followers, err = socialGraphManager.GetFollowers("jack")
		Ω(err).Should(BeNil())
		Ω(followers).Should(HaveLen(0))

		following, err = socialGraphManager.GetFollowing("jack")
		Ω(err).Should(BeNil())
		Ω(following).Should(HaveLen(1))
	})
})
