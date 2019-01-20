package boneful

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"
	"text/template"

	"github.com/go-zoo/bone"
)

// Service is the base type for what users of the API will manage
type Service struct {
	rootPath      string
	routes        []Route
	documentation string
}

// GenerateDocumentation is used to spit out markdown format of docs.
func (s *Service) GenerateDocumentation(w io.Writer) {
	funcMap := template.FuncMap{
		// The name "lower" is what the function will be called in the template text.
		"lower": strings.ToLower,
	}
	tmpl := template.Must(template.New("md").Funcs(funcMap).Parse(mdTemplate))
	tmpl.Execute(w, s)
}

// GenerateJSONDoc emits JSON-formatted documentation info
func (s *Service) GenerateJSONDoc(w io.Writer) {
	// fmt.Println(len(s.routes))
	// r := s.routes[0]
	// fmt.Printf("%#v\n", r)
	// j, err := json.Marshal(r)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(j)
	json.NewEncoder(w).Encode(s.routes)
}

// Mux returns a multiplexer that can be used as a master handler to
// route requests to the appropriate handler.
func (s *Service) Mux() *bone.Mux {
	mux := bone.New()
	for _, r := range s.routes {
		switch r.Method {
		case "HEAD":
			mux.HeadFunc(r.Path, r.Handler)
		case "GET":
			mux.GetFunc(r.Path, r.Handler)
		case "POST":
			mux.PostFunc(r.Path, r.Handler)
		case "PUT":
			mux.PutFunc(r.Path, r.Handler)
		case "PATCH":
			mux.PatchFunc(r.Path, r.Handler)
		case "DELETE":
			mux.DeleteFunc(r.Path, r.Handler)
		}
	}

	// helper function to make sure we don't overwrite
	// routes provided by our caller
	hasRoute := func(rt string) bool {
		for _, routes := range mux.Routes {
			for _, r := range routes {
				if r.Path == rt {
					return true
				}
			}
		}
		return false
	}

	mdpath := concatPath(s.RootPath(), "/md")
	if !hasRoute(mdpath) {
		mux.GetFunc(mdpath, s.GetDocMD)
	}
	jsondoc := concatPath(s.RootPath(), "/jsondoc")
	if !hasRoute(jsondoc) {
		mux.GetFunc(jsondoc, s.GetJSONDoc)
	}
	health := concatPath(s.RootPath(), "/health")
	if !hasRoute(health) {
		mux.GetFunc(health, s.HealthCheck)
	}
	// for verb, routes := range mux.Routes {
	// 	for _, r := range routes {
	// 		fmt.Printf("%s %#v\n", verb, *r)
	// 	}
	// }

	return mux
}

// GetDocMD is a handler for markdown documentation.
func (s *Service) GetDocMD(rw http.ResponseWriter, req *http.Request) {
	s.GenerateDocumentation(rw)
}

// GetJSONDoc is a handler to return JSON documentation
func (s *Service) GetJSONDoc(rw http.ResponseWriter, req *http.Request) {
	s.GenerateJSONDoc(rw)
}

// HealthCheck is a rudimentary endpoint that simply returns "OK".
// If you want something more sophisticated, simply write your own
// and put it under the /health endpoint.
func (s *Service) HealthCheck(rw http.ResponseWriter, req *http.Request) {
	json.NewEncoder(rw).Encode("OK")
}

// Path specifies the root URL template path of the WebService.
// All Routes will be relative to this path.
func (s *Service) Path(root string) *Service {
	s.rootPath = root
	return s
}

// Route creates a new Route using the RouteBuilder and add to the ordered list of Routes.
func (s *Service) Route(builder *RouteBuilder) *Service {
	s.routes = append(s.routes, builder.Build())
	return s
}

// Method creates a new RouteBuilder and initializes its http method
func (s *Service) Method(httpMethod string) *RouteBuilder {
	return NewRouteBuilder().servicePath(s.rootPath).Method(httpMethod)
}

// RootPath returns the RootPath associated with this WebService. Default "/"
func (s *Service) RootPath() string {
	return s.rootPath
}

// Doc is used to set the documentation text of this service.
// It strips leading whitespace from every line so that you can
// have attractive indenting in your source files. Markdown will
// mess with it anyway.
func (s *Service) Doc(plainText string) *Service {
	re := regexp.MustCompile("\n[ \t]+")
	s.documentation = re.ReplaceAllString(plainText, "\n")
	return s
}

// Documentation returns it.
func (s *Service) Documentation() string {
	return s.documentation
}

// Routes returns the array of routes defined for this service.
func (s *Service) Routes() []Route {
	return s.routes
}

/*
   Convenience methods
*/

// HEAD is a shortcut for .Method("HEAD").Path(subPath)
func (s *Service) HEAD(subPath string) *RouteBuilder {
	return NewRouteBuilder().servicePath(s.rootPath).Method("HEAD").Path(subPath)
}

// GET is a shortcut for .Method("GET").Path(subPath)
func (s *Service) GET(subPath string) *RouteBuilder {
	return NewRouteBuilder().servicePath(s.rootPath).Method("GET").Path(subPath)
}

// POST is a shortcut for .Method("POST").Path(subPath)
func (s *Service) POST(subPath string) *RouteBuilder {
	return NewRouteBuilder().servicePath(s.rootPath).Method("POST").Path(subPath)
}

// PUT is a shortcut for .Method("PUT").Path(subPath)
func (s *Service) PUT(subPath string) *RouteBuilder {
	return NewRouteBuilder().servicePath(s.rootPath).Method("PUT").Path(subPath)
}

// PATCH is a shortcut for .Method("PATCH").Path(subPath)
func (s *Service) PATCH(subPath string) *RouteBuilder {
	return NewRouteBuilder().servicePath(s.rootPath).Method("PATCH").Path(subPath)
}

// OPTIONS is a shortcut for .Method("OPTIONS").Path(subPath)
func (s *Service) OPTIONS(subPath string) *RouteBuilder {
	return NewRouteBuilder().servicePath(s.rootPath).Method("OPTIONS").Path(subPath)
}

// DELETE is a shortcut for .Method("DELETE").Path(subPath)
func (s *Service) DELETE(subPath string) *RouteBuilder {
	return NewRouteBuilder().servicePath(s.rootPath).Method("DELETE").Path(subPath)
}
