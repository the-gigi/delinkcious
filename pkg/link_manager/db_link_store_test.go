package link_manager

import (
	sq "github.com/Masterminds/squirrel"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
)

var _ = Describe("DB link store tests", func() {
	var linkStore *DbLinkStore
	var deleteAll = func() {
		sq.Delete("links").RunWith(linkStore.db).Exec()
		sq.Delete("tags").RunWith(linkStore.db).Exec()
	}
	BeforeSuite(func() {
		var err error
		linkStore, err = NewDbLinkStore("localhost", 5432, "postgres", "postgres")
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
	})

	It("should delete a link", func() {
	})

})
