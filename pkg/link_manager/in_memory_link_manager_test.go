package link_manager

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
)

var _ = Describe("In-memory link manager tests", func() {
	var err error
	var linkManager om.LinkManager
	var socialGraphManager *mockSocialGraphManager
	var eventSink *testEventsSink

	BeforeEach(func() {
		socialGraphManager = newMockSocialGraphManager([]string{"liat"})
		eventSink = newLinkManagerEventsSink()
		linkManager, err = NewLinkManager(newInMemoryLinkStore(),
			socialGraphManager,
			"",
			eventSink,
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

		// Verify link manager notified the event sink about a single added event for the follower "liat"
		Ω(eventSink.addLinkEvents).Should(HaveLen(1))
		Ω(eventSink.addLinkEvents["liat"]).Should(HaveLen(1))
		Ω(*eventSink.addLinkEvents["liat"][0]).Should(Equal(link))
		Ω(eventSink.updateLinkEvents).Should(HaveLen(0))
		Ω(eventSink.deletedLinkEvents).Should(HaveLen(0))
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

	It("should update link status when receiving OnLinkChecked() calls", func() {
		// Add a link
		r := om.AddLinkRequest{
			Username: "gigi",
			Url:      "https://golang.org/",
			Title:    "Golang",
			Tags:     map[string]bool{"programming": true},
		}
		err := linkManager.AddLink(r)
		Ω(err).Should(BeNil())

		// Should have 1 link in pending status
		r2 := om.GetLinksRequest{Username: "gigi"}
		res, err := linkManager.GetLinks(r2)
		Ω(err).Should(BeNil())
		Ω(res.Links).Should(HaveLen(1))
		Ω(res.Links[0].Status).Should(Equal(om.LinkStatusPending))

		// Call OnLinkChecked() with valid status on link manager (after type asserting to LinkCheckerEvents)
		linkCheckSink := linkManager.(om.LinkCheckerEvents)
		linkCheckSink.OnLinkChecked("gigi", r.Url, om.LinkStatusValid)

		// The link should have valid status
		res, err = linkManager.GetLinks(r2)
		Ω(err).Should(BeNil())
		Ω(res.Links).Should(HaveLen(1))
		Ω(res.Links[0].Status).Should(Equal(om.LinkStatusValid))

		// Call OnLinkChecked() with valid status again
		linkCheckSink.OnLinkChecked("gigi", r.Url, om.LinkStatusValid)

		// The link should still have valid status
		res, err = linkManager.GetLinks(r2)
		Ω(err).Should(BeNil())
		Ω(res.Links).Should(HaveLen(1))
		Ω(res.Links[0].Status).Should(Equal(om.LinkStatusValid))

		// Call OnLinkChecked() with invalid status
		linkCheckSink.OnLinkChecked("gigi", r.Url, om.LinkStatusInvalid)
		// The link should have invalid status now
		res, err = linkManager.GetLinks(r2)
		Ω(err).Should(BeNil())
		Ω(res.Links).Should(HaveLen(1))
		Ω(res.Links[0].Status).Should(Equal(om.LinkStatusInvalid))
	})

})
