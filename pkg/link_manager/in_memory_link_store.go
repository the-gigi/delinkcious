package link_manager

import (
	"errors"
	"fmt"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
	"regexp"
	"time"
)

// User links are a map of url:TaggedLink
type UserLinks map[string]*om.Link

// Link store is a map of username:UserLinks
type InMemoryLinkStore map[string]UserLinks

func (m *InMemoryLinkStore) GetLinks(request om.GetLinksRequest) (result om.GetLinksResult, err error) {
	userLinks := (*m)[request.Username]
	if userLinks == nil {
		return
	}

	// Prepare complied regexes
	var urlRegex *regexp.Regexp
	var titleRegex *regexp.Regexp
	var descriptionRegex *regexp.Regexp
	if request.UrlRegex != "" {
		urlRegex, err = regexp.Compile(request.UrlRegex)
		if err != nil {
			return
		}
	}

	if request.TitleRegex != "" {
		titleRegex, err = regexp.Compile(request.UrlRegex)
		if err != nil {
			return
		}
	}

	if request.DescriptionRegex != "" {
		descriptionRegex, err = regexp.Compile(request.UrlRegex)
		if err != nil {
			return
		}
	}

	for _, link := range userLinks {
		// Check wach link against the regular expressions
		if urlRegex != nil && !urlRegex.MatchString(link.Url) {
			continue
		}

		if titleRegex != nil && !titleRegex.MatchString(link.Title) {
			continue
		}

		if descriptionRegex != nil && !descriptionRegex.MatchString(link.Description) {
			continue
		}

		// If there no tag was requested add link immediately and continue
		if request.Tag == "" {
			result.Links = append(result.Links, *link)
			continue
		}

		// Add link only if it has the request tag
		if link.Tags[request.Tag] {
			result.Links = append(result.Links, *link)
		}
	}

	return
}

func (m *InMemoryLinkStore) AddLink(request om.AddLinkRequest) (link *om.Link, err error) {
	if request.Url == "" {
		err = errors.New("URL can't be empty")
		return
	}

	if request.Username == "" {
		err = errors.New("User name can't be empty")
		return
	}

	userLinks := (*m)[request.Username]
	if userLinks == nil {
		userLinks = UserLinks{}
	} else {
		if userLinks[request.Url] != nil {
			msg := fmt.Sprintf("User %s already has a link for %s", request.Username, request.Url)
			err = errors.New(msg)
			return
		}
	}

	link = &om.Link{
		Url:         request.Url,
		Title:       request.Title,
		Description: request.Description,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Tags:        request.Tags,
	}
	userLinks[request.Url] = link
	return
}

func (m *InMemoryLinkStore) UpdateLink(request om.UpdateLinkRequest) (link *om.Link, err error) {
	userLinks := (*m)[request.Username]
	if userLinks == nil || userLinks[request.Url] == nil {
		msg := fmt.Sprintf("User %s doesn't have a link for %s", request.Username, request.Url)
		err = errors.New(msg)
		return
	}

	link = userLinks[request.Url]
	if request.Title != "" {
		link.Title = request.Title
	}

	if request.Description != "" {
		link.Description = request.Description
	}

	newTags := request.AddTags
	for t, _ := range link.Tags {
		if request.RemoveTags[t] {
			continue
		}

		newTags[t] = true
	}

	return
}

func (m *InMemoryLinkStore) DeleteLink(username string, url string) error {
	if url == "" {
		return errors.New("URL can't be empty")
	}

	if username == "" {
		return errors.New("User name can't be empty")
	}

	userLinks := (*m)[username]
	if userLinks == nil || userLinks[url] == nil {
		msg := fmt.Sprintf("User %s doesn't have a link for %s", username, url)
		return errors.New(msg)
	}

	delete((*m)[username], url)
	return nil
}
