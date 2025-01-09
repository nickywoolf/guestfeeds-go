package internal

import (
	"encoding/json"
	"strings"

	"golang.org/x/net/html"
)

type Episode struct {
	titleFormatter
	rawTitle string
	Feed     string
	Title    string
}

type titleFormatter func(raw string) string

func defaultTitleFormatter(raw string) string {
	return raw
}

func (ep *Episode) FinalizeTitle() {
	ep.Title = ep.titleFormatter(ep.rawTitle)
}

func applePodcastsTitleFormatter(raw string) string {
	parts := strings.Split(raw, " - ")
	return parts[0]
}

func transistorFMTitleFormatter(raw string) string {
	parts := strings.Split(raw, " | ")
	if len(parts) > 0 {
		return parts[1]
	}
	return raw
}

type ApplePodcastServerData []struct {
	Data struct {
		Shelves []struct {
			Items []struct {
				ContextAction struct {
					EpisodeOffer struct {
						ShowOffer struct {
							FeedURL string `json:"feedUrl"`
						} `json:"showOffer"`
					} `json:"episodeOffer"`
				} `json:"contextAction"`
			} `json:"items"`
		} `json:"shelves"`
	} `json:"data"`
}

func ExtractFeed(source string) (*Episode, error) {
	doc, err := html.Parse(strings.NewReader(source))
	if err != nil {
		return nil, err
	}

	ep := &Episode{titleFormatter: defaultTitleFormatter}

	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "title" {
			if node.FirstChild != nil && node.FirstChild.Type == html.TextNode {
				ep.rawTitle = node.FirstChild.Data
				ep.titleFormatter = applePodcastsTitleFormatter
			}
		}

		if node.Type == html.ElementNode && node.Data == "link" {
			var linkRel, linkType, linkHref string
			for _, a := range node.Attr {
				switch a.Key {
				case "rel":
					linkRel = a.Val
				case "type":
					linkType = a.Val
				case "href":
					linkHref = a.Val
				}
			}
			if linkRel == "alternate" && linkType == "application/rss+xml" && strings.HasPrefix(linkHref, "https://feeds.transistor.fm") {
				ep.Feed = linkHref
				ep.titleFormatter = transistorFMTitleFormatter
			}
		}

		if node.Type == html.ElementNode && node.Data == "script" {
			for _, a := range node.Attr {
				if a.Key == "id" && a.Val == "serialized-server-data" {
					if node.FirstChild != nil && node.FirstChild.Type == html.TextNode {
						var serverData ApplePodcastServerData
						_ = json.Unmarshal([]byte(node.FirstChild.Data), &serverData)
						ep.Feed = serverData[0].Data.Shelves[0].Items[0].ContextAction.EpisodeOffer.ShowOffer.FeedURL
					}
				}
			}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			traverse(child)
		}
	}
	traverse(doc)

	ep.FinalizeTitle()

	return ep, nil
}
