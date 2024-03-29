package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/npmania/bong/internal/bong"
	"github.com/npmania/bong/internal/config"
	tg "github.com/npmania/bong/internal/server/tmplgen"
)

type SearchHandler struct {
	Config  config.Config
	BongMap bong.BongMap
}

func (h SearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var query string

	if err := r.ParseForm(); err != nil {
		//TODO: when does this error get triggered?
		fmt.Println("wrong post data!")
	}

	query = r.FormValue("q")
	if ok := h.bongRedirect(w, r, query); ok {
		return
	}

	data := tg.SearchParams{
		Title: "bong",
		Query: query,
	}

	if err := tg.Search(w, data); err != nil {
		fmt.Println(err)
	}
}

func (sh SearchHandler) bongRedirect(w http.ResponseWriter, r *http.Request, query string) bool {
	var (
		realQuery string
		target    string
		bongus    string
	)

	splited := strings.Split(query, " ")
	if strings.HasPrefix(string(query), sh.Config.DefaultPrefix) {
		bongus = splited[0][len(sh.Config.DefaultPrefix):]
	}

	b, ok := sh.BongMap[bongus]

	if !ok {
		if sh.Config.Fallback == "" {
			return false
		} else if sh.Config.Fallback != "" {
			b, ok = sh.BongMap[sh.Config.Fallback]
			if !ok {
				return false
			}
			realQuery = strings.Join(splited, " ")
		}
	} else {
		realQuery = strings.Join(splited[1:], " ")
	}

	realQuery = url.QueryEscape(realQuery)
	realQuery = strings.ReplaceAll(realQuery, "+", "%20")

	if realQuery != "" && strings.Contains(b.BongUrl, "%[1]s") {
		target = fmt.Sprintf(b.BongUrl, realQuery)
	} else {
		target = b.MainUrl
	}

	http.Redirect(w, r, target, http.StatusMovedPermanently)
	return true
}
