package duck

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/npmania/bong/internal/bong"
)

type duckBang struct {
	Category     string `json:"c"`
	MainUrl      string `json:"d"`
	SomeWeirdNum int    `json:"r"`
	Title        string `json:"s"`
	SubCategory  string `json:"sc"`
	Bang         string `json:"t"`
	BangUrl      string `json:"u"`
}

func parseBang(raw []byte) (bangs []duckBang, err error) {
	err = json.Unmarshal(raw, &bangs)
	return
}

func fixBangs(bangs []duckBang) (fixed []duckBang, err error) {
	for _, b := range bangs {
		if string(b.Title[0]) == " " {
			b.Title = b.Title[1:]
		}
		if string(b.Title[len(b.Title)-1]) == " " {
			b.Title = b.Title[:len(b.Title)-1]
		}

		// add duckduckgo address to self redirected bangs
		if string(b.BangUrl[0]) == "/" {
			b.BangUrl = duck + b.BangUrl
		}

		if strings.HasPrefix(b.BangUrl, duck) {
			if b.BangUrl == duck+"/newbang" {
				b.MainUrl = b.BangUrl
			} else {
				b.MainUrl = strings.Replace(b.BangUrl, "{{{s}}}", "", 1)
			}
		}

		b.MainUrl, err = url.QueryUnescape(b.MainUrl)
		if err != nil {
			return nil, err
		}

		b.BangUrl, err = url.QueryUnescape(b.BangUrl)
		if err != nil {
			return nil, err
		}

		// in case percent mark is in path
		// TODO: fix hacky way of unescaping percent mark
		b.BangUrl = strings.ReplaceAll(b.BangUrl, "%", "%%25")

		b.MainUrl = strings.ReplaceAll(b.MainUrl, "{{{s}}}", "%s")
		b.BangUrl = strings.ReplaceAll(b.BangUrl, "{{{s}}}", "%s")

		if !strings.HasPrefix(b.MainUrl, "http://") && !strings.HasPrefix(b.MainUrl, "https://") {
			b.MainUrl = "http://" + b.MainUrl
		}

		fixed = append(fixed, b)
	}

	return fixed, nil
}

func toBongMap(bangs []duckBang) (bong.BongMap, error) {
	bangs, err := fixBangs(bangs)
	if err != nil {
		return nil, err
	}

	bm := make(bong.BongMap)

	for _, bang := range bangs {
		bm[bang.Bang] = bong.Bong{
			Title:   bang.Title,
			MainUrl: bang.MainUrl,
			BongUrl: bang.BangUrl,
			Bongus:  bang.Bang,
		}
	}

	return bm, nil
}
