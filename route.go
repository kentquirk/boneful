package boneful

import (
	"encoding/json"
	"net/http"

	"github.com/go-zoo/bone"
)

// Route binds a Method and Path to a handler.
// It also holds the documentation for the route.
type Route struct {
	Method  string           `json:"method"`
	Path    string           `json:"path"` // webservice root path + described path
	Handler http.HandlerFunc `json:"-"`
	muxfunc func(string, http.HandlerFunc) *bone.Route

	// documentation
	Doc            string                `json:"doc"`
	Notes          string                `json:"notes"`
	Operation      string                `json:"operation"`
	Consumes       []string              `json:"consumes"`
	Produces       []string              `json:"produces"`
	ParameterDocs  []*Parameter          `json:"parms"`
	ResponseErrors map[int]ResponseError `json:"-"`
	ReadSample     interface{}           `json:"-"` // models an example request payload
	WriteSample    interface{}           `json:"-"` // models an example response payload
}

func (r *Route) postBuild() {
}

func (r Route) String() string {
	return r.Method + " " + r.Path
}

// CodeFormat generates the marker for file contents for a code
// block in Markdown
func (r Route) CodeFormat() string {
	for _, c := range r.Consumes {
		switch c {
		case "text/plain":
			return "text"
		case "text/markdown":
			return "markdown"
		case "text/html":
			return "html"
		case "application/json":
			return "json"
		}
	}
	return ""
}

// Reads returns formatted example content for a Reads value
func (r Route) Reads() string {
	for _, c := range r.Consumes {
		switch c {
		case "text/plain", "text/markdown", "text/html":
			return r.ReadSample.(string)
		case "application/json":
			b, err := json.MarshalIndent(r.ReadSample, "        ", "  ")
			if err != nil {
				continue
			}
			return string(b)
		}
	}
	return ""
}

// Writes returns formatted example content for a Writes value
func (r Route) Writes() string {
	for _, p := range r.Produces {
		switch p {
		case "text/plain":
			if r.WriteSample != nil {
				return r.WriteSample.(string)
			}
			return ""
		case "application/json":
			b, err := json.MarshalIndent(r.WriteSample, "        ", "  ")
			if err != nil {
				continue
			}
			return string(b)
		}
	}
	return ""
}
