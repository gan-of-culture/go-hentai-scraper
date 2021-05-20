package miohentai

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
)

const site = "https://miohentai.com/"

func ParseURL(URL string) []string {
	if strings.HasPrefix(URL, site+"video/") || strings.HasPrefix(URL, site+"image-library/") {
		return []string{URL}
	}

	if !strings.HasPrefix(URL, site+"tag/") {
		return []string{}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return []string{}
	}

	re := regexp.MustCompile(`post.*"([^"]*/video/[^"]*)`)
	episodes := re.FindAllStringSubmatch(htmlString, -1)
	URLs := []string{}
	for _, v := range episodes {
		URLs = append(URLs, v[1])
	}
	return URLs
}

func Extract(URL string) ([]static.Data, error) {
	URLs := ParseURL(URL)
	if len(URLs) == 0 {
		return nil, fmt.Errorf("[Miohentai] No scrapable URL found for %s", URL)
	}

	data := []static.Data{}
	for _, u := range URLs {
		d, err := extractData(u)
		if err != nil {
			return nil, err
		}
		data = append(data, d)
	}

	return data, nil
}

func extractData(URL string) (static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return static.Data{}, err
	}

	videoRe := regexp.MustCompile(`[^"]*cdn\.miohentai[^"]*`)
	srcURL := videoRe.FindString(htmlString)
	if srcURL == "" {
		imageRe := regexp.MustCompile(`async data-src="([^"]*)`)
		srcURL = imageRe.FindStringSubmatch(htmlString)[1]
	}

	headers, err := request.Headers(srcURL, URL)
	if err != nil {
		return static.Data{}, err
	}
	size, err := request.GetSizeFromHeaders(&headers)
	if err != nil {
		return static.Data{}, err
	}

	ext := utils.GetLastItemString(strings.Split(srcURL, "."))
	if strings.Contains(srcURL, "index.php") {
		ext = strings.Split(headers.Get("content-type"), "/")[1]
	}

	return static.Data{
		Site:  site,
		Title: strings.Split(utils.GetMeta(&htmlString, "og:title"), " | ")[0],
		Type:  headers.Get("content-type"),
		Streams: map[string]static.Stream{
			"0": {
				URLs: []static.URL{
					{
						URL: srcURL,
						Ext: ext,
					},
				},
				Quality: "best",
				Size:    size,
			},
		},
		Url: URL,
	}, nil
}
