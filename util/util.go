/*Package util provides commonly used function and types. This could probably be
 * split up into a few separate packages. This package sets up common command
 * arguments for all URL Shortener services when it's imported.*/
package util

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"strings"
)

var (
	configPath string
)

// init set up common command line arguments for URL shortener services.
func init() {
	const (
		defaultConfig = "config.json"
		configUsage   = "the location of the configuration file"
	)
	flag.StringVar(&configPath, "config", defaultConfig, configUsage)
	flag.StringVar(&configPath, "c", defaultConfig, configUsage+" (shorthand)")
}

// BaseConfig shared config directives by every service.
type BaseConfig struct {
	ListenAddr string `json:"listen_addr"`
}

// Redirect represents a shortened URL. Key is the shortener version of the url
// and Target is the orgional URL
type Redirect struct {
	Key, Target string
}

// ErrorResult josn API error result.
type ErrorResult struct {
	ID      int
	Message string
}

// HandleError by panicing.
func HandleError(result interface{}, err error) (r interface{}) {
	if err != nil {
		panic(err)
	}
	return result
}

// EnvHandler http handler for printing debugging and environment information.
func EnvHandler(rw http.ResponseWriter, req *http.Request) {
	environment := make(map[string]string)
	for _, item := range os.Environ() {
		splits := strings.Split(item, "=")
		key := splits[0]
		val := strings.Join(splits[1:], "=")
		environment[key] = val
	}

	envJSON := HandleError(json.MarshalIndent(environment, "", "  ")).([]byte)
	rw.Write(envJSON)
}

// WriteErrorResp writer a JSON api error response.
func WriteErrorResp(rw http.ResponseWriter, status int, msg string) {
	rw.WriteHeader(status)

	encoder := json.NewEncoder(rw)
	encoder.Encode(ErrorResult{status, msg})
}

// WriterJSONResponse writer a JSON object as a response with HTTP
// status OK
func WriterJSONResponse(rw http.ResponseWriter, obj interface{}) {
	respJSON, err := json.Marshal(obj)
	if err != nil {
		WriteErrorResp(rw, http.StatusInternalServerError, err.Error())
		return
	}
	rw.Write(respJSON)
}

// LoadConfig load a configuration file using the config command line argument.
func LoadConfig(config interface{}) error {
	fd, err := os.Open(configPath)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(fd)
	err = decoder.Decode(config)
	if err != nil {
		return err
	}
	return nil
}
