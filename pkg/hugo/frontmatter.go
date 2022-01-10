package hugo

import (
	"io"
	"os"

	"github.com/gohugoio/hugo/parser/pageparser"
)

const (
	fmTitle      = "title"
	fmAuthor     = "author"
	fmCategories = "categories"
	fmTags       = "tags"

	fmDate        = "date"        // priority high
	fmLastmod     = "lastmod"     // priority middle
	fmPublishDate = "publishDate" // priority low
)

type FrontMatter struct {
	Title    string
	Author   string
	Category string
	Tags     []string
	Date     string
}

// ParseFrontMatter parses the frontmatter of the specified Hugo content.
func ParseFrontMatter(filename string) (*FrontMatter, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseFrontMatter(file)
}

func parseFrontMatter(r io.Reader) (*FrontMatter, error) {
	cfm, err := pageparser.ParseFrontMatterAndContent(r)
	if err != nil {
		return nil, err
	}

	fm := &FrontMatter{}
	if fm.Title, err = getString(&cfm, fmTitle); err != nil {
		return nil, err
	}
	if isArray := isArray(&cfm, fmAuthor); isArray {
		if fm.Author, err = getFirstStringItem(&cfm, fmAuthor); err != nil {
			return nil, err
		}
	} else {
		if fm.Author, err = getString(&cfm, fmAuthor); err != nil {
			return nil, err
		}
	}
	if fm.Category, err = getFirstStringItem(&cfm, fmCategories); err != nil {
		return nil, err
	}
	if fm.Tags, err = getAllStringItems(&cfm, fmTags); err != nil {
		return nil, err
	}
	if fm.Date, err = getString(&cfm, fmDate); err != nil {
		return nil, err
	}

	return fm, nil
}

func getString(cfm *pageparser.ContentFrontMatter, fmKey string) (string, error) {
	v, ok := cfm.FrontMatter[fmKey]
	if !ok {
		return "", NewFMNotExistError(fmKey)
	}

	switch s := v.(type) {
	case string:
		if s == "" {
			return "", NewFMNotExistError(fmKey)
		}
		return s, nil
	default:
		return "", NewFMInvalidTypeError(fmKey, "string", s)
	}
}

func getAllStringItems(cfm *pageparser.ContentFrontMatter, fmKey string) ([]string, error) {
	v, ok := cfm.FrontMatter[fmKey]
	if !ok {
		return nil, NewFMNotExistError(fmKey)
	}

	switch arr := v.(type) {
	case []interface{}:
		var strarr []string
		for _, item := range arr {
			switch s := item.(type) {
			case string:
				if s != "" {
					strarr = append(strarr, s)
				}
			default:
				return nil, NewFMInvalidTypeError(fmKey, "string", s)
			}
		}
		if len(strarr) < 1 {
			return nil, NewFMNotExistError(fmKey)
		}
		return strarr, nil

	default:
		return nil, NewFMInvalidTypeError(fmKey, "[]interface{}", arr)
	}
}

func getFirstStringItem(cfm *pageparser.ContentFrontMatter, fmKey string) (string, error) {
	arr, err := getAllStringItems(cfm, fmKey)
	if err != nil {
		return "", err
	}
	return arr[0], nil
}

func isArray(cfm *pageparser.ContentFrontMatter, fmKey string) bool {
	switch v := cfm.FrontMatter[fmKey]; v.(type) {
	case []interface{}:
		return true
	default:
		return false
	}
}
