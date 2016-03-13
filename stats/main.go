//StatsService logs clicks from the Frontend and produces statistics.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/rlguarino/shortener/stats/types"
	"github.com/rlguarino/shortener/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	session    *mgo.Session
	configPath string
	config     *Config
)

// Config the Configuration of the StatService
type Config struct {
	util.BaseConfig
	RootURL    string      `json:"root_url"`
	Mongo      MongoConfig `json:"mongo"`
	URLService string      `json:"url_service_addr"`
}

// MongoConfig the connection information for the mongo cluster.
type MongoConfig struct {
	Addrs  string `json:"addrs"`
	DBName string `json:"db_name"`
}

// RecordClickHandler saves clicks to the mongo db.
func RecordClickHandler(rw http.ResponseWriter, req *http.Request) {
	var c types.Click
	reqDecoder := json.NewDecoder(req.Body)
	err := reqDecoder.Decode(&c)
	if err != nil {
		util.WriteErrorResp(rw, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}
	if c.Key == "" {
		util.WriteErrorResp(rw, http.StatusBadRequest, "Must Specifiy Key")
		return
	}
	s := session.Copy()
	err = s.DB(config.Mongo.DBName).C("clicks").Insert(c)
	if err != nil {
		util.WriteErrorResp(rw, http.StatusInternalServerError, err.Error())
		return
	}

	util.WriterJSONResponse(rw, c)
}

// StatsHandler returns statics for a specific shortened url.
func StatsHandler(rw http.ResponseWriter, req *http.Request) {
	key := mux.Vars(req)["key"]

	resp, err := http.Get(config.URLService + "/v1/route/" + key)
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

	s := session.Copy()
	defer s.Close()
	c := s.DB(config.Mongo.DBName).C("clicks")
	count, err := c.Find(bson.M{"key": key}).Count()
	fmt.Println(count)
	if err != nil {
		util.WriteErrorResp(rw, http.StatusInternalServerError, err.Error())
		return
	}
	util.WriterJSONResponse(rw, types.RouteStats{Key: key, Clicks: count})
}

func main() {
	var err error
	flag.Parse()
	config = new(Config)
	if err = util.LoadConfig(config); err != nil {
		panic(err)
	}
	session, err = mgo.Dial(config.Mongo.Addrs)
	if err != nil {
		panic(err)
	}
	r := mux.NewRouter()
	r.Path("/v1/click").Methods("POST").HandlerFunc(RecordClickHandler)
	r.Path("/v1/stats/{key}").Methods("GET").HandlerFunc(StatsHandler)
	r.Path("/env").Methods("GET").HandlerFunc(util.EnvHandler)

	n := negroni.Classic()
	n.UseHandler(r)
	n.Run(config.ListenAddr)
}
