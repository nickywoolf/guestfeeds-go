package internal

import (
	"testing"
)

func TestNewEpisode(t *testing.T) {
	tests := []struct {
		name  string
		input string
		feed  string
		title string
	}{
		{
			name: "Apple Podcasts Episode Page",
			input: `<!DOCTYPE html>
			<html>
				<head>
					<title>Episode Title - Show Title - Apple Podcasts</title>
				</head>
				<body>
					<script type="application/json" id="serialized-server-data">
						[{"data":{"shelves":[{"items":[{"contextAction":{"episodeOffer":{"showOffer":{"feedUrl":"https://feeds.podcasts.com/show-handle"}}}}]}]}}]
					</script>
				</body>
			</html>`,
			feed:  "https://feeds.podcasts.com/show-handle",
			title: "Episode Title",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ep, err := ExtractFeed(tc.input)

			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}

			if ep.Feed != tc.feed {
				t.Errorf("Expecting feed '%s' but got '%s'\n", tc.feed, ep.Feed)
			}

			if ep.Title != tc.title {
				t.Errorf("Expecting title '%s' but got '%s'\n", tc.title, ep.Title)
			}
		})
	}
}
