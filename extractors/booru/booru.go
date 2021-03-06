package booru

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gan-of-culture/go-hentai-scraper/config"
	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
)

const site = "https://booru.io/"
const postURL = "https://booru.io/p/"
const apiDataURL = "https://booru.io/api/data/"
const apiEntityURL = "https://booru.io/api/entity/"
const apiQueryURL = "https://booru.io/api/query/entity?query="

// Entity JSON type
type Entity struct {
	Key         string             `json:"key"`
	ContentType string             `json:"contentType"`
	Attributes  map[string]float32 `json:"attributes"`
	Tags        map[string]int     `json:"tags"`
	Transforms  map[string]string  `json:"transforms"`
}

// EntitySlice JSON type
type EntitySlice struct {
	Data   []Entity `json:"data"`
	Cursor string   `json:"cursor"`
}

type extractor struct{}

// New returns a booru.io extractor.
func New() static.Extractor {
	return &extractor{}
}

// Extract for booru pages
func (e *extractor) Extract(URL string) ([]*static.Data, error) {
	query, err := parseURL(URL)
	if err != nil {
		return nil, err
	}

	data, err := extractData(query)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// parseURL for danbooru pages
func parseURL(url string) (string, error) {
	if !strings.Contains(url, "%20") {
		re := regexp.MustCompile(`https://booru\.io/p/(.+)`)
		matchedID := re.FindStringSubmatch(url)
		if len(matchedID) > 2 {
			return "", static.ErrURLParseFailed
		}
		return fmt.Sprintf("%s%s", apiEntityURL, matchedID[1]), nil
	}

	tags := strings.Split(url, "https://booru.io/q/")

	return fmt.Sprintf("%s%s", apiQueryURL, tags[1]), nil
}

func extractData(queryURL string) ([]*static.Data, error) {
	jsonString, err := request.Get(queryURL)
	if err != nil {
		fmt.Println(queryURL)
		return nil, err
	}

	entitySlice := EntitySlice{}
	//single post
	if !strings.Contains(queryURL, "=") {
		entity := Entity{}
		err := json.Unmarshal([]byte(jsonString), &entity)
		if err != nil {
			fmt.Println(queryURL)
			return nil, err
		}
		entitySlice.Data = append(entitySlice.Data, entity)
	}

	if len(entitySlice.Data) == 0 {

		cursor := 0
		for {
			if config.Amount > 0 && config.Amount <= cursor {
				break
			}
			entitySliceTmp := EntitySlice{}
			err = json.Unmarshal([]byte(jsonString), &entitySliceTmp)
			if err != nil {
				fmt.Println(queryURL)
				fmt.Println("Cursor", cursor)
				fmt.Println(jsonString)
			}
			if len(entitySliceTmp.Data) == 0 && err == nil {
				break
			}
			entitySlice.Data = append(entitySlice.Data, entitySliceTmp.Data...)
			cursor += 50
			jsonString, err = request.Get(fmt.Sprintf("%s&cursor=%d", queryURL, cursor))
			fmt.Printf("%s&cursor=%d", queryURL, cursor)
			if err != nil {
				return nil, err
			}
			time.Sleep(100 * time.Millisecond)
		}
	}

	data := []*static.Data{}
	for _, e := range entitySlice.Data {
		tType, tVal := GetBestQualityImg(e.Transforms)
		ext := GetFileExt(tType)
		size, _ := request.Size(fmt.Sprintf("%s%s", apiDataURL, tVal), site)

		data = append(data, &static.Data{
			Site:  site,
			Title: e.Key,
			Type:  "image",
			Streams: map[string]*static.Stream{
				"0": {
					URLs: []*static.URL{
						{
							URL: fmt.Sprintf("%s%s", apiDataURL, tVal),
							Ext: ext,
						},
					},
					Quality: strings.Split(tType, ":")[0],
					Size:    size,
				},
			},
			Url: fmt.Sprintf("%s%s", postURL, e.Key),
		})
	}

	return data, nil
}

// GetBestQualityImg of transformation
func GetBestQualityImg(transformations map[string]string) (string, string) {
	re := regexp.MustCompile(`[0-9]+`)
	transformationType := ""
	transformationValue := ""
	currentBest := 0
	for key, val := range transformations {
		resString := re.FindString(key)
		resolution, _ := strconv.Atoi(resString)
		if resolution <= 0 {
			continue
		}
		if resolution < currentBest {
			continue
		}
		currentBest = resolution
		transformationType = key
		transformationValue = val
	}
	return transformationType, transformationValue
}

// GetFileExt of transformation
func GetFileExt(tranformation string) string {
	transSplit := strings.Split(tranformation, "/")
	if len(transSplit) > 1 {
		if transSplit[1] == "jpeg" {
			return "jpg"
		}
		return transSplit[1]
	}
	return ""
}
