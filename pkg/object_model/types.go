package object_model

import "time"

type EventTypeEnum int

const (
	LinkAdded EventTypeEnum = iota
	LinkUpdated
	LinkDeleted
)

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

type Event struct {
	EventType EventTypeEnum
	Username  string
	Url       string
	Timestamp time.Time
}

type GetNewsRequest struct {
	Username   string
	StartToken string
}

type GetNewsResult struct {
	Events    []*Event
	NextToken string
}
