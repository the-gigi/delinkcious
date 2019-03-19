package news_manager

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
)

var _ = Describe("In-memory link manager tests", func() {
	var newsManager *NewsManager

	BeforeEach(func() {
		nm, err := NewNewsManager(NewInMemoryNewsStore(), "", "")
		Ω(err).Should(BeNil())
		newsManager = nm.(*NewsManager)
		Ω(newsManager).ShouldNot(BeNil())
	})

	It("should get news", func() {
		// No news initially
		r := om.GetNewsRequest{
			Username: "gigi",
		}
		res, err := newsManager.GetNews(r)
		Ω(err).Should(BeNil())
		Ω(res.Events).Should(HaveLen(0))

		// Add a link
		link := &om.Link{
			Url:   "http://123.com",
			Title: "123",
		}
		newsManager.OnLinkAdded("gigi", link)
		res, err = newsManager.GetNews(r)
		Ω(err).Should(BeNil())
		Ω(res.Events).Should(HaveLen(1))
		event := res.Events[0]
		Ω(event.EventType).Should(Equal(om.LinkAdded))
		Ω(event.Url).Should(Equal("http://123.com"))

		// Update a link
		link.Title = "New Title"
		newsManager.OnLinkUpdated("gigi", link)
		res, err = newsManager.GetNews(r)
		Ω(err).Should(BeNil())
		Ω(res.Events).Should(HaveLen(2))
		event = res.Events[0]
		Ω(event.EventType).Should(Equal(om.LinkAdded))
		Ω(event.Url).Should(Equal("http://123.com"))

		event = res.Events[1]
		Ω(event.EventType).Should(Equal(om.LinkUpdated))
		Ω(event.Url).Should(Equal("http://123.com"))

		// Delete a link
		newsManager.OnLinkDeleted("gigi", link.Url)
		res, err = newsManager.GetNews(r)
		Ω(err).Should(BeNil())
		Ω(res.Events).Should(HaveLen(3))
		event = res.Events[0]
		Ω(event.EventType).Should(Equal(om.LinkAdded))
		Ω(event.Url).Should(Equal("http://123.com"))

		event = res.Events[1]
		Ω(event.EventType).Should(Equal(om.LinkUpdated))
		Ω(event.Url).Should(Equal("http://123.com"))

		event = res.Events[2]
		Ω(event.EventType).Should(Equal(om.LinkDeleted))
		Ω(event.Url).Should(Equal("http://123.com"))
	})
})
