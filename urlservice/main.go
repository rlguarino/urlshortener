/*UrlService is a MicroSercvice that stores and generates keys which corespond to
 * urls. The service is stateless and relies on a Redis Cluster for storage of the
 * redirection data.
 */
package main

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/rlguarino/shortener/util"
	"gopkg.in/redis.v3"
)

var (
	client *redis.Client
	master string
	config *Config
)

// Config the Configuration of the UrlService
type Config struct {
	util.BaseConfig
	Redis RedisConfig `json:"redis"`
}

// RedisConfig the Connection Information for Redis
type RedisConfig struct {
	SentinelAddrs string `json:"sentinel_addrs"`
	Master        string `json:"master_name"`
}

// GetRouteHandler handles requests for retreiving routes given a key.
func GetRouteHandler(rw http.ResponseWriter, req *http.Request) {
	key := mux.Vars(req)["key"]
	// Retreive key from Redis.
	resp := client.Get(key)
	url, err := resp.Result()
	if err != nil {
		util.WriteErrorResp(rw, http.StatusNotFound, "Not Found!")
		return
	}
	redirect := util.Redirect{Key: key, Target: url}
	util.WriterJSONResponse(rw, redirect)
}

// CreateRouteHandler create a new redirection. Given a target it generates
// a unique key and returns the key in the response body.
func CreateRouteHandler(rw http.ResponseWriter, req *http.Request) {
	var r util.Redirect
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&r)
	if err != nil {
		util.WriteErrorResp(rw, http.StatusBadRequest, "Invalid Request")
		return
	}

	key, err := getKey()
	if err != nil {
		util.WriteErrorResp(rw, http.StatusInternalServerError, "Could not generate id")
		return
	}

	err = client.Set(key, r.Target, 0).Err()
	if err != nil {
		util.WriteErrorResp(rw, http.StatusInternalServerError, "failed to write to db")
		return
	}

	r.Key = key
	util.WriterJSONResponse(rw, r)
}

// getKey generates a new key using cryto/random and checks to make sure it does
// not already exist in the database.
func getKey() (string, error) {
	retries := 100
	i := 0
	b := make([]byte, 10)
	rand.Read(b)
	for i < retries {
		key := fmt.Sprintf("%x", md5.Sum(b))[:10]
		if exists := client.Exists(key).Val(); !exists {
			return key, nil
		}
	}
	return "", errors.New("max retry limit reached")
}

// InfoHandler display information regarding the current Redis Master.
func InfoHandler(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("Info GetMaster")
	info, err := client.Info().Bytes()
	if err != nil {
		util.WriteErrorResp(rw, http.StatusInternalServerError, err.Error())
		return
	}
	rw.Write(info)
}

func main() {
	flag.Parse()
	config = new(Config)
	if err := util.LoadConfig(config); err != nil {
		panic(err)
	}
	client = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    config.Redis.Master,
		SentinelAddrs: []string{config.Redis.SentinelAddrs},
	})

	r := mux.NewRouter()
	r.Path("/v1/route/{key}").Methods("GET").HandlerFunc(GetRouteHandler)
	r.Path("/v1/route/").Methods("POST").HandlerFunc(CreateRouteHandler)
	r.Path("/info").Methods("GET").HandlerFunc(InfoHandler)
	r.Path("/env").Methods("GET").HandlerFunc(util.EnvHandler)

	n := negroni.Classic()
	n.UseHandler(r)
	n.Run(config.ListenAddr)
}
