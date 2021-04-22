package json_response

import (
	"fmt"
	"math"
	"net/url"
	"regexp"
	"strings"
)

type Response struct {
	Data interface{} `json:"data"`
}

type ListResponse struct {
	Data  interface{} `json:"data"`
	Links Links       `json:"links"`
}

type Links struct {
	Self  *string `json:"self"`
	Next  *string `json:"next"`
	Prev  *string `json:"prev"`
	First *string `json:"first"`
	Last  *string `json:"last"`
}

func NewResponse(data interface{}) *Response {
	return &Response{data}
}

func NewListResponseAndPageLinks(items interface{}, path string, total uint64,
	pageNumber, pageSize uint64,
) (*ListResponse, error) {
	unescapedPath, err := url.PathUnescape(path)
	if err != nil {
		return nil, err
	}

	links := buildPageLinks(unescapedPath, pageNumber, pageSize, total)

	return &ListResponse{
		Data:  items,
		Links: links,
	}, nil
}

func buildPageLinks(path string, number, size, total uint64) Links {
	links := Links{
		First: buildFirstPath(path, size),
		Last:  buildLastPath(path, size, total),
		Next:  buildNextPath(path, number, size, total),
		Prev:  buildPrevPath(path, number, size, total),
		Self:  buildSelfPath(path, number, size),
	}

	return links
}

func buildFirstPath(path string, size uint64) *string {
	if size == 0 {
		return &path
	}
	return changePageNumber(path, 1)
}

func buildLastPath(path string, size, total uint64) *string {
	if size < 1 {
		return nil
	}
	return changePageNumber(path, getLastPageNumber(size, total))
}

func buildNextPath(path string, number, size, total uint64) *string {
	if size == 0 {
		return &path
	}

	lastPageNumber := getLastPageNumber(size, total)
	if number >= lastPageNumber {
		return nil
	}
	return changePageNumber(path, number+1)
}

func buildSelfPath(path string, number, size uint64) *string {
	if size == 0 {
		return &path
	}
	return changePageNumber(path, number)
}

func buildPrevPath(path string, number, size, total uint64) *string {
	if number <= 1 {
		return nil
	}

	lastPageNumber := getLastPageNumber(size, total)
	if number > lastPageNumber {
		return changePageNumber(path, lastPageNumber)
	}

	return changePageNumber(path, number-1)
}

func changePageNumber(input string, to uint64) *string {
	re := regexp.MustCompile(`(page\[number\])(=\d+)`)
	strToParse := input
	if !re.MatchString(input) {
		strToParse = addParameterToPath(input, "page[number]=0")

	}
	result := re.ReplaceAllString(strToParse, fmt.Sprintf(`$1=%v`, to))
	return &result
}

func addParameterToPath(path, parameter string) string {
	var delimeter string
	if strings.Contains(path, "?") {
		delimeter = "&"
	} else {
		delimeter = "?"
	}
	return path + delimeter + parameter
}

func getLastPageNumber(size, total uint64) uint64 {
	count := float64(total) / float64(size)
	return uint64(math.Ceil(count))
}
