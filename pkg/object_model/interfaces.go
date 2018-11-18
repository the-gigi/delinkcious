package object_model

type LinkManager interface {
	GetLinks(request GetLinksRequest) (GetLinksResult, error)
	AddLink(request AddLinkRequest) error
	UpdateLink(request UpdateLinkRequest) error
	DeleteLink(username string, url string) error
}

type UserManager interface {
	Register(user User) error
	Login(username string, authToken string) (session string, err error)
	Logout(username string, session string) error
}

type SocialGraphManager interface {
	Follow(followed string, follower string) error
	Unfollow(followed string, follower string) error

	AcceptFollowRequest(followed string, follower string) error
	RejectFollowRequest(followed string, follower string) error

	KickFollower(followed string, follower string) error

	GetFollowing(username string) map[string]bool
	GetFollowers(username string) map[string]bool
}

type LinkEvents interface {
	OnLinkAdded(username string, links *Link)
	OnLinkUpdated(username string, links *Link)
}
