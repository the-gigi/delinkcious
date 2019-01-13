package link_manager

import (
	sq "github.com/Masterminds/squirrel"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/the-gigi/delinkcious/pkg/db_util"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
	"log"
)

var _ = Describe("DB link store tests", func() {
	var linkStore *DbLinkStore
	var deleteAll = func() {
		sq.Delete("links").RunWith(linkStore.db).Exec()
		sq.Delete("tags").RunWith(linkStore.db).Exec()
	}
	BeforeSuite(func() {
		var err error
		dbHost, dbPort, err := db_util.GetDbEndpoint("social_graph")
		Ω(err).Should(BeNil())

		linkStore, err = NewDbLinkStore(dbHost, dbPort, "postgres", "postgres")
		if err != nil {
			log.Fatal(err)
		}

		Ω(err).Should(BeNil())
		Ω(linkStore).ShouldNot(BeNil())
	})

	BeforeEach(deleteAll)
	AfterSuite(deleteAll)

	It("should add and get links", func() {
		// No links initially
		r := om.GetLinksRequest{
			Username: "gigi",
		}
		res, err := linkStore.GetLinks(r)
		Ω(err).Should(BeNil())
		Ω(res.Links).Should(HaveLen(0))

		// Add a link
		r2 := om.AddLinkRequest{
			Username: "gigi",
			Url:      "https://golang.org/",
			Title:    "Golang",
			Tags:     map[string]bool{"programming": true},
		}
		_, err = linkStore.AddLink(r2)
		Ω(err).Should(BeNil())

		res, err = linkStore.GetLinks(r)
		Ω(err).Should(BeNil())
		Ω(res.Links).Should(HaveLen(1))
		link := res.Links[0]
		Ω(link.Url).Should(Equal(r2.Url))
		Ω(link.Title).Should(Equal(r2.Title))

	})

	It("should update a link", func() {
		// Add a link
		r := om.AddLinkRequest{
			Username: "gigi",
			Url:      "https://golang.org/",
			Title:    "Golang",
			Tags:     map[string]bool{"programming": true},
		}
		_, err := linkStore.AddLink(r)
		Ω(err).Should(BeNil())

		r2 := om.UpdateLinkRequest{
			Username:    r.Username,
			Url:         r.Url,
			Description: "The main web site for the Go programming language",
			RemoveTags:  map[string]bool{"programming": true},
		}
		_, err = linkStore.UpdateLink(r2)
		Ω(err).Should(BeNil())

		r3 := om.GetLinksRequest{Username: "gigi"}
		res, err := linkStore.GetLinks(r3)
		Ω(err).Should(BeNil())
		Ω(res.Links).Should(HaveLen(1))
		link := res.Links[0]
		Ω(link.Url).Should(Equal(r.Url))
		Ω(link.Description).Should(Equal(r2.Description))
	})

	It("should delete a link", func() {
		// Add a link
		r := om.AddLinkRequest{
			Username: "gigi",
			Url:      "https://golang.org/",
			Title:    "Golang",
			Tags:     map[string]bool{"programming": true},
		}
		_, err := linkStore.AddLink(r)
		Ω(err).Should(BeNil())

		// Should have 1 link
		r2 := om.GetLinksRequest{Username: "gigi"}
		res, err := linkStore.GetLinks(r2)
		Ω(err).Should(BeNil())
		Ω(res.Links).Should(HaveLen(1))

		// Delete the link
		err = linkStore.DeleteLink("gigi", r.Url)
		Ω(err).Should(BeNil())

		// There should be no more links
		res, err = linkStore.GetLinks(r2)
		Ω(err).Should(BeNil())
		Ω(res.Links).Should(HaveLen(0))
	})
})
