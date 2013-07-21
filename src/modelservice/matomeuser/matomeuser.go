package matomeuser

import (
	"time"
	"appengine"
	"appengine/datastore"
)

type User struct {
	Id	string `json:"id" datastore:"-"`
	Email	string `json:"email" datastore:"email"`
	Nickname	string `json:"nickname" datastore:"nickname"`
	Admin	bool	`json:"admin" datastore:"admin"`
	Created	time.Time	`json:"created" datastore:"created"`
}

func (u *User) key(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "User", u.Email, 0, nil)
}

func (u *User) Save(c appengine.Context) (*User, error) {
	k, err := datastore.Put(c, u.key(c), u)
	if err != nil {
		return nil, err
	}
	u.Id = k.StringID()
	return u, nil
}
