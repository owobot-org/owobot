package builtins

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"runtime/debug"
	"strings"

	"golang.org/x/net/publicsuffix"
)

// Options contains options for the JavaScript fetch function
type Options struct {
	Method        string
	Body          string
	Headers       map[string]any
	HandleCookies *bool
}

// Response contains the response object for the JavaScript fetch function
type Response struct {
	Status     string
	StatusCode int
	Headers    http.Header
	body       []byte
}

func (r Response) JSON() (v any, err error) {
	err = json.Unmarshal(r.body, &v)
	return v, err
}

func (r Response) String() string {
	return string(r.body)
}

// FetchFunc is the fetch function signature
type FetchFunc = func(string, *Options) (*Response, error)

func fetch(pluginName, pluginVersion string) FetchFunc {
	// cookiejar.New always returns a nil error
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})

	return func(url string, opts *Options) (*Response, error) {
		if opts == nil {
			t := true
			opts = &Options{HandleCookies: &t}
		}

		if opts.HandleCookies == nil {
			t := true
			opts.HandleCookies = &t
		}

		if opts.Method == "" {
			opts.Method = http.MethodGet
		}

		req, err := http.NewRequest(opts.Method, url, strings.NewReader(opts.Body))
		if err != nil {
			return nil, err
		}

		for key, value := range opts.Headers {
			req.Header.Add(key, value.(string))
		}

		if req.Header.Get("User-Agent") == "" {
			req.Header.Set("User-Agent", getUserAgent(pluginName, pluginVersion))
		}

		client := &http.Client{}
		if *opts.HandleCookies {
			client.Jar = jar
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return &Response{
			Status:     resp.Status,
			StatusCode: resp.StatusCode,
			Headers:    resp.Header,
			body:       responseBody,
		}, nil
	}
}

// getUserAgent uses the built in vcs information to generate a user agent string
func getUserAgent(pluginName, pluginVersion string) string {
	commit := "unknown"
	modified := "unmodified"
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				commit = setting.Value[:8]
			case "vcs.modified":
				if setting.Value == "true" {
					modified = "modified"
				}
			}
		}
	}

	return fmt.Sprintf("owobot/%s (%s; %s/%s)", commit, modified, pluginName, pluginVersion)
}
