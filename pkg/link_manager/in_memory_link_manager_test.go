package link_manager

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
)

var _ = Describe("In-memory link manager tests", func() {
	var err error
	var linkManager om.LinkManager

	BeforeEach(func() {
		linkManager, err = NewLinkManager(NewInMemoryLinkStore(),
			nil,
			nil,
			10)
		Ω(err).Should(BeNil())
	})

	It("should add and get links", func() {
		// No links initially
		r := om.GetLinksRequest{
			Username: "gigi",
		}
		res, err := linkManager.GetLinks(r)
		Ω(err).Should(BeNil())
		Ω(res.Links).Should(HaveLen(0))

		// Add a link
		r2 := om.AddLinkRequest{
			Username: "gigi",
			Url:      "https://golang.org/",
			Title:    "Golang",
			Tags:     map[string]bool{"programming": true},
		}
		err = linkManager.AddLink(r2)
		Ω(err).Should(BeNil())

		res, err = linkManager.GetLinks(r)
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
		err := linkManager.AddLink(r)
		Ω(err).Should(BeNil())

		r2 := om.UpdateLinkRequest{
			Username:    r.Username,
			Url:         r.Url,
			Description: "The main web site for the Go programming language",
			RemoveTags:  map[string]bool{"programming": true},
		}
		err = linkManager.UpdateLink(r2)
		Ω(err).Should(BeNil())

		r3 := om.GetLinksRequest{Username: "gigi"}
		res, err := linkManager.GetLinks(r3)
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
		err := linkManager.AddLink(r)
		Ω(err).Should(BeNil())

		// Should have 1 link
		r2 := om.GetLinksRequest{Username: "gigi"}
		res, err := linkManager.GetLinks(r2)
		Ω(err).Should(BeNil())
		Ω(res.Links).Should(HaveLen(1))

		// Delete the link
		err = linkManager.DeleteLink("gigi", r.Url)
		Ω(err).Should(BeNil())

		// There should be no more links
		res, err = linkManager.GetLinks(r2)
		Ω(err).Should(BeNil())
		Ω(res.Links).Should(HaveLen(0))
	})
})
