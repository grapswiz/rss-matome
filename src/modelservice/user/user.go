package user

import (
	"time"
	"appengine"
	"appengine/datastore"
)

type User struct {
	Id	string `json:"id" datastore:"-"`
	Admin	bool	`json:"admin"`
	Created	time.Time	`json:"created`
}

func (u *User) key(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "User", u.Id, 0, nil)
}
