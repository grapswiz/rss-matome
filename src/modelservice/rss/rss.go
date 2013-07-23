package rss

import (
	"time"
	"appengine"
	"appengine/datastore"
	"appengine/user"
)

type Rss struct {
	Id	int64	`json:"id" datastore:"-"`
	Urls	[]string	`json:"urls" datastore:"urls,noindex"`
	Created	time.Time	`json:"created" datastore:"created"`
	Updated	time.Time	`json:"updated" datastore:"updated"`
}

func (r *Rss) key(c appengine.Context) *datastore.Key {
	u := user.Current(c)
	if u == nil {
		return nil
	}
	return datastore.NewIncompleteKey(c, "Rss", datastore.NewKey(c, "User", u.Email, 0, nil))
}

func (r *Rss) Save(c appengine.Context) (*Rss, error) {
	k, err := datastore.Put(c, r.key(c), r)
	if err != nil {
		return nil, err
	}
	r.Id = k.IntID()
	return r, nil
}

func Get(c appengine.Context, id int64, userId string) (Rss, error) {
	var r Rss
	key := datastore.NewKey(c, "Rss", "", id, datastore.NewKey(c, "User", userId, 0, nil))
	err := datastore.Get(c, key, &r)
	if err == nil {
		r.Id = key.IntID()
	}
	return r, err
}
