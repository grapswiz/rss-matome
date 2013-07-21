package util

import (
	"appengine"
	"appengine/urlfetch"
	"io/ioutil"
	"encoding/xml"
	"time"
)

type Auth struct {
	LoggedIn	bool	`json:"loggedIn",datastore:"loggedIn,noindex`
	LoginUrl	string	`json:"loginUrl",datastore:"loginUrl,noindex"`
	LogoutUrl	string	`json:"logoutUrl",datastore"logoutUrl,noindex`
}

type Response struct {
	Body	string
}

func FetchBodyFromUrl(c appengine.Context, url string, ch chan *Response) {
	client := urlfetch.Client(c)
	resp, err := client.Get(url)
	if err != nil {
		ch <- &Response{""}
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ch <- &Response{""}
		return
	}
	ch <- &Response{string(body)}
}

type Entry struct {
	Id	string	`json:"id" xml:"id"`
	Title	string	`json:"title" xml:"title"`
	Link	Link	`json:"link" xml:"link"`
	Published	time.Time	`json:"published" xml:"published"`
	Updated	time.Time	`json:"updated" xml:"updated"`
}

type Link struct {
	Rel	string	`json:"rel" xml:"rel,attr"`
	Type	string	`json:"type" xml"type,attr"`
	Href	string	`json:"link" xml:"href,attr"`
}

type Feed struct {
	Title	string	`json:"title"`
	Link	string	`json:"link"`
	Entries	[]Entry	`json:"entries" xml:"entry"`
}

func ParseXmlToFeed(StringXml string) Feed {
	var feed Feed
	byteXml := []byte(StringXml)
	xml.Unmarshal(byteXml, &feed)
	return feed
}

