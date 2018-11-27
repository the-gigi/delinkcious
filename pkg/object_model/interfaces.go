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

	GetFollowing(username string) (map[string]bool, error)
	GetFollowers(username string) (map[string]bool, error)

	//AcceptFollowRequest(followed string, follower string) error
	//RejectFollowRequest(followed string, follower string) error
	//KickFollower(followed string, follower string) error
}

type LinkManagerEvents interface {
	OnLinkAdded(username string, link *Link)
	OnLinkUpdated(username string, link *Link)
	OnLinkDeleted(username string, url string)
}
