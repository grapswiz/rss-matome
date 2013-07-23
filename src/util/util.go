package util

import (
	"appengine"
	"appengine/urlfetch"
	"io/ioutil"
	"encoding/xml"
	"time"
	"math/rand"
	"strconv"
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
	PubDate	string	`json:"pubDate" xml:"updated"`
}

func ParseXmlToFeed(StringXml string) Feed {
	var feed Feed
	byteXml := []byte(StringXml)
	xml.Unmarshal(byteXml, &feed)
	return feed
}

type Guid struct {
	XMLName	xml.Name	`xml:"guid"`
	Guid	string	`xml:",chardata"`
	IsPermaLink	bool	`xml:"isPermaLink,attr"`
}

type Item struct {
	XMLName	xml.Name	`xml:"item"`
	Title	string	`xml:"title"`
	Description	string	`xml:"description"`
	Link	string	`xml:"link"`
	GUID	Guid	`xml:"guid"`
	PubDate	string	`xml:"pubDate"`
}

type Channel struct {
	XMLName	xml.Name	`xml:"channel"`
	Title	string	`xml:"title"`
	Description	string	`xml:"description"`
	Link	string	`xml:"link"`
	LastBuildDate	string	`xml:"lastBuildDate"`
	PubDate	string	`xml:"pubDate"`
	TTL	int	`xml:"ttl"`
	Items	[]Item
}

type Rss struct {
	XMLName	xml.Name	`xml:"rss"`
	Version	string	`xml:"version,attr"`
	Channel	Channel
}

func FeedToRss(f Feed) Rss {
	pubDate, err := time.Parse(time.RFC3339, f.PubDate)
	var pubDateString string
	if err != nil {
		pubDateString = ""
	} else {
		pubDateString = pubDate.Format(time.RFC1123Z)
	}
	var rss = Rss{
		Version: "2.0",
		Channel: Channel{
			Title: f.Title,
			Description: "RSS's Description",
			Link: f.Link,
			LastBuildDate: time.Now().Format(time.RFC1123Z),
			TTL: 1800,
			Items: nil,
		},
	}
	rss.Channel.PubDate = pubDateString
	var items = make([]Item, len(f.Entries))
	for i, v := range f.Entries {
		items[i].Title = v.Title
		items[i].Description = "Feed's Description"
		items[i].Link = v.Link.Href
		items[i].GUID = Guid{
			Guid: strconv.FormatInt(rand.Int63(), 10),
			IsPermaLink: false,
		}
		items[i].PubDate = v.Published.Format(time.RFC1123Z)
	}
	rss.Channel.Items = items
	return rss
}
