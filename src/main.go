package main

import (
	"fmt"
	"net/http"
	"modelservice/rss"
	"time"
	"appengine"
	"appengine/user"
	"strings"
	"modelservice/matomeuser"
	"util"
	"strconv"
	"encoding/json"
)

func init() {
	http.HandleFunc("/createsamplerss/", sampleHandler)
	http.HandleFunc("/auth/", authHandler)
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request)  {
	if strings.Count(r.URL.Path[1:], "/") == 1 {
		path := strings.Split(r.URL.Path[1:], "/")
		if strings.HasSuffix(path[1], ".rss") {
			c := appengine.NewContext(r)
			id := strings.Split(path[1], ".rss")
			intId, _ := strconv.ParseInt(id[0], 10, 64)
			r1, _ := rss.Get(c, intId)
			ch := make(chan *util.Response, 100)
			for _, v := range r1.Urls{
				go util.FetchBodyFromUrl(c, v, ch)
			}
			var entries []util.Entry
			for _, _ = range r1.Urls {
				res := <-ch
				feed := util.ParseXmlToFeed(res.Body)
				//TODO Entryを混ぜて最新順に並び替えてRSSにしてfmt.Fprintfしてみる
				for _, v := range feed.Entries {
					entries = append(entries, v)
				}
			}
			rootFeed := util.Feed{"俺のフィード", r.URL.String(), entries}
			b, err := json.Marshal(rootFeed)
			if err != nil {
				fmt.Fprintln(w, "エラーだよ")
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			fmt.Fprintln(w, string(b))
		} else {
			fmt.Fprint(w, "そんなリクエストあたし知らない！")
		}
	}

	c := appengine.NewContext(r)
	u := user.Current(c)
	if u != nil {
		u1 := matomeuser.User{
			Email: u.Email,
			Admin: u.Admin,
			Created: time.Now(),
		}
		u1.Save(c)
	}
}

func sampleHandler(w http.ResponseWriter, r *http.Request) {
	r1 := rss.Rss{
		Created: time.Now(),
		Updated: time.Now(),
		Urls: []string{"http://blog.memoto.com/feed/", "http://blog.alexmaccaw.com/feed"},
	}
	r2, err := r1.Save(appengine.NewContext(r))
	if err != nil {
		return
	}
	fmt.Fprint(w, r2.Id)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, _ := user.LoginURL(c, "/")
		fmt.Fprint(w, url)
		return
	}
	url, _ := user.LogoutURL(c, "/")
	fmt.Fprint(w, url)
}
