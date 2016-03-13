//frontend handles all incoming user requests. Including creating 
// new shortened urls, redirection, and statistics.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/mssola/user_agent"
	"github.com/rlguarino/shortener/stats/types"
	"github.com/rlguarino/shortener/util"
)

var (
	config *Config
)

// ShortenForm html form binding.
type ShortenForm struct {
	Url string
}

// Config the Frontend service configuration directives.
type Config struct {
	util.BaseConfig
	URLServiceAddr  string `json:"url_service_addr"`
	StatServiceAddr string `json:"stat_service_addr"`
}

// RedirectHandler handler redirection requests by directiong the http
// client to the target address. Record the click after the client has
// been sent to the target.
func RedirectHandler(rw http.ResponseWriter, req *http.Request) {
	key := mux.Vars(req)["key"]
	resp, err := http.Get(config.URLServiceAddr + "/v1/route/" + key)
	if err != nil {
		http.Error(rw, "An error occured please try again later", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		http.NotFound(rw, req)
		return
	} else if resp.StatusCode != http.StatusOK {
		http.Error(rw, "An error occured please try again later", http.StatusInternalServerError)
		return
	}

	d := json.NewDecoder(resp.Body)
	var redirect util.Redirect
	err = d.Decode(&redirect)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(rw, req, redirect.Target, http.StatusFound)

	recordClick(req)
}

// recordClick extract click information from the request and send it
// to the statistics service.
func recordClick(req *http.Request) {
	c := types.Click{}
	c.UserAgent = types.UserAgent{}
	c.Key = mux.Vars(req)["key"]
	c.ClientIP = strings.Split(req.RemoteAddr, ":")[0]
	c.Time = time.Now()
	c.Referer = req.Referer()
	ua := user_agent.New(req.UserAgent())
	c.UserAgent.OS = ua.OS()
	c.UserAgent.Str = req.UserAgent()
	c.UserAgent.Bot = ua.Bot()
	c.UserAgent.Mobile = ua.Mobile()
	c.UserAgent.Engine, c.UserAgent.EngineVersion = ua.Engine()
	c.UserAgent.Browser, c.UserAgent.BrowserVersion = ua.Browser()

	jsonStr, err := json.Marshal(c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Click Not Recorded: %s\n", err.Error())
	}
	resp, err := http.Post(config.StatServiceAddr+"/v1/click", "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

// Redirect the client to a failed page.
func failRedirect(rw http.ResponseWriter, req *http.Request) {
	http.Redirect(rw, req, "/failure.html", http.StatusFound)
}

// CreateRouteHandler handles request to create a shortened URL.
func CreateRouteHandler(rw http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		failRedirect(rw, req)
		return
	}
	decoder := schema.NewDecoder()
	var sf ShortenForm
	err = decoder.Decode(&sf, req.PostForm)
	if err != nil {
		failRedirect(rw, req)
		return
	}
	// TODO(rlg): Validate target url here.
	redirect := util.Redirect{Target: sf.Url}
	jsonStr, err := json.Marshal(redirect)
	if err != nil {
		failRedirect(rw, req)
		return
	}
	resp, err := http.Post(config.URLServiceAddr+"/v1/route/", "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		failRedirect(rw, req)
		return
	}
	if resp.StatusCode != http.StatusOK {
		failRedirect(rw, req)
		return
	}
	defer resp.Body.Close()
	respDecoder := json.NewDecoder(resp.Body)
	err = respDecoder.Decode(&redirect)
	if err != nil {
		failRedirect(rw, req)
		return
	}

	t, _ := template.ParseFiles("template/success.html")

	t.Execute(rw, struct {
		LongURL, ShortURL, StatsURL string
	}{
		redirect.Target,
		fmt.Sprintf("/r/%s", redirect.Key),
		fmt.Sprintf("/s/%s", redirect.Key),
	})
}

// LinkStatsHandler fulfils request for statistical information about
// a shortened URL.
func LinkStatsHandler(rw http.ResponseWriter, req *http.Request) {
	key := mux.Vars(req)["key"]
	resp, err := http.Get(config.StatServiceAddr + "/v1/stats/" + key)
	if err != nil {
		fmt.Println("Failed Get Stats")
		failRedirect(rw, req)
		return
	}
	defer resp.Body.Close()
	var stats types.RouteStats
	respDecoder := json.NewDecoder(resp.Body)
	err = respDecoder.Decode(&stats)
	if err != nil {
		failRedirect(rw, req)
		return
	}

	resp, err = http.Get(config.URLServiceAddr + "/v1/route/" + key)
	if err != nil {
		http.Error(rw, "An error occured please try again later", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		http.NotFound(rw, req)
		return
	} else if resp.StatusCode != http.StatusOK {
		http.Error(rw, "An error occured please try again later", http.StatusInternalServerError)
		return
	}

	d := json.NewDecoder(resp.Body)
	var redirect util.Redirect
	err = d.Decode(&redirect)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	t, _ := template.ParseFiles("template/stats.html")

	t.Execute(rw, struct {
		LongURL, ShortURL string
		Clicks            int
	}{
		redirect.Target,
		fmt.Sprintf("/r/%s", stats.Key),
		stats.Clicks,
	})
}

func main() {
	flag.Parse()
	config = new(Config)
	if err := util.LoadConfig(config); err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.Path("/r/{key}").Methods("GET").HandlerFunc(RedirectHandler)
	r.Path("/s/{key}").Methods("GET").HandlerFunc(LinkStatsHandler)
	r.Path("/new").Methods("POST").HandlerFunc(CreateRouteHandler)
	r.Path("/env").Methods("GET").HandlerFunc(util.EnvHandler)

	n := negroni.Classic()
	n.UseHandler(r)
	n.Run(config.ListenAddr)
}
