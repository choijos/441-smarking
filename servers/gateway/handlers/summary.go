package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

//PreviewImage represents a preview image for a page
type PreviewImage struct {
	URL       string `json:"url,omitempty"`
	SecureURL string `json:"secureURL,omitempty"`
	Type      string `json:"type,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Alt       string `json:"alt,omitempty"`
}

//PageSummary represents summary properties for a web page
type PageSummary struct {
	Type        string          `json:"type,omitempty"`
	URL         string          `json:"url,omitempty"`
	Title       string          `json:"title,omitempty"`
	SiteName    string          `json:"siteName,omitempty"`
	Description string          `json:"description,omitempty"`
	Author      string          `json:"author,omitempty"`
	Keywords    []string        `json:"keywords,omitempty"`
	Icon        *PreviewImage   `json:"icon,omitempty"`
	Images      []*PreviewImage `json:"images,omitempty"`
}

//SummaryHandler handles requests for the page summary API.
//This API expects one query string parameter named `url`,
//which should contain a URL to a web page. It responds with
//a JSON-encoded PageSummary struct containing the page summary
//meta-data.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	/*TODO: add code and additional functions to do the following:
	- Add an HTTP header to the response with the name
	 `Access-Control-Allow-Origin` and a value of `*`. This will
	  allow cross-origin AJAX requests to your server.
	- Get the `url` query string parameter value from the request.
	  If not supplied, respond with an http.StatusBadRequest error.
	- Call fetchHTML() to fetch the requested URL. See comments in that
	  function for more details.
	- Call extractSummary() to extract the page summary meta-data,
	  as directed in the assignment. See comments in that function
	  for more details
	- Close the response HTML stream so that you don't leak resources.
	- Finally, respond with a JSON-encoded version of the PageSummary
	  struct. That way the client can easily parse the JSON back into
	  an object. Remember to tell the client that the response content
	  type is JSON.

	Helpful Links:
	https://golang.org/pkg/net/http/#Request.FormValue
	https://golang.org/pkg/net/http/#Error
	https://golang.org/pkg/encoding/json/#NewEncoder
	*/

	// w.Header().Add("Access-Control-Allow-Origin", "*")

	paramValue := r.FormValue("url")
	if len(paramValue) == 0 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return

	}

	bodyStream, err := fetchHTML(paramValue)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	pgSummary, err := extractSummary(paramValue, bodyStream)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	bodyStream.Close()

	// w.Header().Add("Content-Type", "application/json")

	enc := json.NewEncoder(w)
	if err := enc.Encode(pgSummary); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

}

//fetchHTML fetches `pageURL` and returns the body stream or an error.
//Errors are returned if the response status code is an error (>=400),
//or if the content type indicates the URL is not an HTML page.
func fetchHTML(pageURL string) (io.ReadCloser, error) {
	/*TODO: Do an HTTP GET for the page URL. If the response status
	code is >= 400, return a nil stream and an error. If the response
	content type does not indicate that the content is a web page, return
	a nil stream and an error. Otherwise return the response body and
	no (nil) error.

	To test your implementation of this function, run the TestFetchHTML
	test in summary_test.go. You can do that directly in Visual Studio Code,
	or at the command line by running:
		go test -run TestFetchHTML

	Helpful Links:
	https://golang.org/pkg/net/http/#Get
	*/
	resp, err := http.Get(pageURL)
	ctype := resp.Header.Get("Content-Type")

	if err != nil {
		return nil, fmt.Errorf("error: %v", err)

	} else if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("bad status request")

	} else if !strings.HasPrefix(ctype, "text/html") {
		return nil, fmt.Errorf("not a webpage")

	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)

	}

	return resp.Body, nil

}

func fillPrevImg(imgMap map[string]string, pageURL string, imgType string) (*PreviewImage, error) {
	imgStruct := &PreviewImage{}

	map_url := ""
	map_type := ""
	width := ""
	height := ""

	if imgType == "image" {
		map_url = "og:image"
		map_type = "og:image:type"
		width = imgMap["og:image:width"]
		height = imgMap["og:image:height"]

		if imgMap["og:image:secure_url"] != "" {
			imgStruct.SecureURL = imgMap["og:image:secure_url"]

		}

		if imgMap["og:image:alt"] != "" {
			imgStruct.Alt = imgMap["og:image:alt"]

		}

	} else {
		map_url = "href"
		map_type = "type"

		if imgMap["sizes"] == "any" {
			height = "any"
			width = "any"

		} else if imgMap["sizes"] != "" {
			heightWidth := strings.Split(imgMap["sizes"], "x")
			for i := 0; i < 2; i++ {
				if i == 0 {
					height = heightWidth[i]

				} else {
					width = heightWidth[i]

				}

			}

		}

	}

	if height != "any" && height != "" && width != "any" && width != "" {
		num, convErr := strconv.Atoi(height)
		if convErr != nil {
			return nil, fmt.Errorf("error converting the image/icon height: %v", convErr)

		}

		imgStruct.Height = num

		num, convErr = strconv.Atoi(width)
		if convErr != nil {
			return nil, fmt.Errorf("error converting the image/icon width: %v", convErr)

		}

		imgStruct.Width = num

	}

	u, err := url.Parse(imgMap[map_url])
	if err != nil {
		return nil, fmt.Errorf("error parsing provided url %v", err)

	}
	base, err := url.Parse(pageURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing base/page url: %v", err)

	}

	imgStruct.URL = base.ResolveReference(u).String()

	if imgMap[map_type] != "" {
		imgStruct.Type = imgMap[map_type]

	}

	return imgStruct, nil

}

//extractSummary tokenizes the `htmlStream` and populates a PageSummary
//struct with the page's summary meta-data.
func extractSummary(pageURL string, htmlStream io.ReadCloser) (*PageSummary, error) {
	/*TODO: tokenize the `htmlStream` and extract the page summary meta-data
	according to the assignment description.

	To test your implementation of this function, run the TestExtractSummary
	test in summary_test.go. You can do that directly in Visual Studio Code,
	or at the command line by running:
		go test -run TestExtractSummary

	Helpful Links:
	https://drstearns.github.io/tutorials/tokenizing/
	http://ogp.me/
	https://developers.facebook.com/docs/reference/opengraph/
	https://golang.org/pkg/net/url/#URL.ResolveReference  */

	tokenizer := html.NewTokenizer(htmlStream)
	propertyContents := map[string]string{}
	var keywords []string
	var images []*PreviewImage
	imageMap := map[string]string{}
	lastIcon := &PreviewImage{}
	pgSum := &PageSummary{}

	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			if tokenizer.Err() == io.EOF {
				break

			}
			return nil, tokenizer.Err()

		}

		if tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken {
			token := tokenizer.Token()
			if "meta" == token.Data {
				key := ""
				val := ""

				for _, attr := range token.Attr {
					if attr.Key == "property" || (attr.Key == "name" && (attr.Val == "keywords" || key == "")) {
						if attr.Val == "og:image" {
							if len(imageMap) != 0 {
								prevImgStruct, err := fillPrevImg(imageMap, pageURL, "image")

								if err != nil {
									return nil, fmt.Errorf("error: %v", err)

								}

								images = append(images, prevImgStruct)
								imageMap = make(map[string]string)

							}

						}

						key = attr.Val

					} else if attr.Key == "content" {
						val = attr.Val

					}

				}

				if key == "keywords" {
					splitKeywords := strings.Split(val, ",")
					for _, word := range splitKeywords {
						keywords = append(keywords, strings.TrimSpace(word))

					}

				} else if strings.Contains(key, "og:image") {
					imageMap[key] = val

				} else {
					if key == "og:type" {
						pgSum.Type = val

					} else if key == "og:url" {
						pgSum.URL = val

					} else if key == "og:site_name" {
						pgSum.SiteName = val

					} else if key == "author" {
						pgSum.Author = val

					} else {
						propertyContents[key] = val

					}

				}

			} else if "title" == token.Data {
				tokenType = tokenizer.Next()
				if tokenType == html.TextToken {
					propertyContents["title"] = tokenizer.Token().Data

				}

			} else if "link" == token.Data {
				isIcon := false
				for _, attr := range token.Attr {
					if attr.Key == "rel" {
						if attr.Val == "icon" {
							isIcon = !isIcon

						}

						break

					}
				}

				if isIcon {
					iconCon := map[string]string{}
					for _, attr := range token.Attr {
						if attr.Key != "rel" {
							iconCon[attr.Key] = attr.Val

						}
					}

					prevImgStruct, err := fillPrevImg(iconCon, pageURL, "icon")

					if err != nil {
						return nil, fmt.Errorf("error: %v", err)
					}

					lastIcon = prevImgStruct

				}

			}

		} else if tokenType == html.EndTagToken && tokenizer.Token().Data == "head" {
			break

		}

	}

	if len(imageMap) != 0 {
		prevImgStruct, err := fillPrevImg(imageMap, pageURL, "image")

		if err != nil {
			return nil, fmt.Errorf("error: %v", err)

		}

		images = append(images, prevImgStruct)

	}

	if propertyContents["og:title"] == "" {
		pgSum.Title = propertyContents["title"]

	} else {
		pgSum.Title = propertyContents["og:title"]

	}

	if propertyContents["og:description"] == "" {
		pgSum.Description = propertyContents["description"]

	} else {
		pgSum.Description = propertyContents["og:description"]

	}

	if images != nil {
		pgSum.Images = images

	}

	if keywords != nil {
		pgSum.Keywords = keywords

	}

	if lastIcon.URL != "" {
		pgSum.Icon = lastIcon

	}

	return pgSum, nil

}
