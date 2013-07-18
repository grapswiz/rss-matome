package auth

type Auth struct {
	LoggedIn	bool	`json:"loggedIn",datastore:",noindex`
	LoginUrl	string	`json:"loginUrl",datastore:",noindex"`
	LogoutUrl	string	`json:"logoutUrl",datastore",noindex`
}
