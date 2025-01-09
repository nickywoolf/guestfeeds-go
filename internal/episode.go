package internal

import (
	"encoding/json"
	"strings"

	"golang.org/x/net/html"
)

type Episode struct {
	rawTitle string
	Feed     string
	Title    string
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

	ep := &Episode{}

	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "title" {
			if node.FirstChild != nil && node.FirstChild.Type == html.TextNode {
				ep.rawTitle = node.FirstChild.Data
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

	titleParts := strings.Split(ep.rawTitle, " - ")
	ep.Title = titleParts[0]

	return ep, nil
}
