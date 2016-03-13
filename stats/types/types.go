/*Package types contains type declerations necessary to communicate with the statistics
 * package API */
package types

import (
	"time"
)

// Click represents a redirect of a user to a target
type Click struct {
	Key       string    `json:"key,omitempty"`
	Time      time.Time `json:"time,omitempty"`
	ClientIP  string    `json:"ip,omitempty"`
	Referer   string    `json:"referer,omitempty"`
	UserAgent UserAgent `json:"user_agent,omitempty"`
}

// UserAgent represents information about the user.
type UserAgent struct {
	Str            string `json:"str,omitempty"`
	Platform       string `json:"platform,omitempty"`
	OS             string `json:"os,omitempty"`
	Engine         string `json:"engine,omitempty"`
	EngineVersion  string `json:"engine_version,omitempty"`
	Browser        string `json:"browser,omitempty"`
	BrowserVersion string `json:"browser_version,omitempty"`
	Bot            bool   `json:"bot,omitempty"`
	Mobile         bool   `json:"mobile,omitempty"`
}

// RouteStats represents the statistics associated with a shortened URL.
type RouteStats struct {
	Key    string
	Clicks int
}
