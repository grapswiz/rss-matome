package main

import (
	"fmt"
	"net/http"
	"modelservice/rss"
	"time"
	"appengine"
	"appengine/user"
)

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/createsamplerss", sampleHandler)
	http.HandleFunc("/auth", authHandler)
}

func handler(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprint(w, "Hello, World")
}

func sampleHandler(w http.ResponseWriter, r *http.Request) {
	var urls []string
	urls[0] = "http://blog.memoto.com"
	urls[1] = "http://blog.alexmaccaw.com"
	r1 := rss.Rss{
		Created: time.Now(),
		Updated: time.Now(),
		Urls: urls,
	}
	r2, _ := r1.Save(appengine.NewContext(r))
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
