package object_model

import "time"

type Link struct {
	Url         string
	Title       string
	Description string
	Tags        map[string]bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type GetLinksRequest struct {
	UrlRegex         string
	TitleRegex       string
	DescriptionRegex string
	Username         string
	Tag              string
	StartToken       string
}

type GetLinksResult struct {
	Links         []Link
	NextPageToken string
}

type AddLinkRequest struct {
	Url         string
	Title       string
	Description string
	Username    string
	Tags        map[string]bool
}

type UpdateLinkRequest struct {
	Url         string
	Title       string
	Description string
	Username    string
	AddTags     map[string]bool
	RemoveTags  map[string]bool
}

type User struct {
	Email string
	Name  string
}
