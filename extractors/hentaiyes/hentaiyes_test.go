package hentaiyes

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Single Episode",
			url:  "https://hentaiyes.com/watch/hime-sama-love-life-episode-03/",
			want: 1,
		}, {
			name: "Series",
			url:  "https://hentaiyes.com/series/hime-sama-love-life/",
			want: 3,
		}, {
			name: "Tag",
			url:  "https://hentaiyes.com/tag/1080p/",
			want: 20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urls := parseURL(tt.url)
			if len(urls) > tt.want || len(urls) == 0 {
				t.Errorf("Got: %v - want: %v", len(urls), tt.want)
			}
		})
	}
}

func TestExtract(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Single Episode",
			url:  "https://hentaiyes.com/watch/shikijou-kyoudan-episode-01/",
			want: 1,
		}, {
			name: "Series",
			url:  "https://hentaiyes.com/series/akiko/",
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := New().Extract(tt.url)
			if err != nil {
				t.Error(err)
			}
			if len(data) > tt.want || len(data) == 0 {
				t.Errorf("Got: %v - want: %v", len(data), tt.want)
			}
		})
	}
}
