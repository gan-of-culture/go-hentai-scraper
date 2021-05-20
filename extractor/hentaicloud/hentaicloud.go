package hentaicloud

import (
	"fmt"
	"regexp"

	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
)

const site = "https://www.hentaicloud.com/"

func ParseURL(URL string) []string {
	if ok, _ := regexp.MatchString(`https://www.hentaicloud.com/video/\d*/[^/]*/episode\d*/`, URL); ok {
		return []string{URL}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return []string{}
	}
	re := regexp.MustCompile(`video/\d*/[^/]*/episode\d*/[^"]*`)
	URLs := []string{}
	for i, v := range re.FindAllString(htmlString, -1) {
		if i%2 == 0 {
			URLs = append(URLs, site+v)
		}
	}

	return URLs
}

func Extract(URL string) ([]static.Data, error) {
	URLs := ParseURL(URL)
	if len(URLs) == 0 {
		return nil, fmt.Errorf("[Hentaicloud] No scrapable URL found for %s", URL)
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
	title := utils.GetMeta(&htmlString, "og:title")

	re := regexp.MustCompile(`(https://www.hentaicloud.com/media/videos/hd/\d*\.([^"]*)).+res="([^"]*)`)
	srcTag := re.FindStringSubmatch(htmlString) //1=URL 2=ext 3=resolution
	if len(srcTag) != 4 {
		return static.Data{}, fmt.Errorf("[Hentaicloud] No scrapable source tag found for %s", URL)
	}

	size, _ := request.Size(srcTag[1], URL)

	return static.Data{
		Site:  site,
		Title: title,
		Type:  utils.GetMediaType(srcTag[2]),
		Streams: map[string]static.Stream{
			"0": {
				URLs: []static.URL{
					{
						URL: srcTag[1],
						Ext: srcTag[2],
					},
				},
				Quality: srcTag[3],
				Size:    size,
			},
		},
		Url: URL,
	}, nil
}
