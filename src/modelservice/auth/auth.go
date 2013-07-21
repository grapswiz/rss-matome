package auth

type Auth struct {
	LoggedIn	bool	`json:"loggedIn",datastore:"loggedIn,noindex`
	LoginUrl	string	`json:"loginUrl",datastore:"loginUrl,noindex"`
	LogoutUrl	string	`json:"logoutUrl",datastore"logoutUrl,noindex`
}
