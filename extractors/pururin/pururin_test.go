package pururin

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Single Gallery",
			url:  "https://pururin.to/gallery/53855/melty-yuel",
			want: 1,
		}, {
			name: "Tag",
			url:  "https://pururin.to/browse/tags/contents/1563/stockings.html",
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
			name: "Single Gallery",
			url:  "https://pururin.to/gallery/53855/melty-yuel",
			want: 1,
		}, {
			name: "Tag",
			url:  "https://pururin.to/browse/tags/contents/1563/stockings.html",
			want: 20,
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
