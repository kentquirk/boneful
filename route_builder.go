package boneful

// RouteBuilder is a helper to construct Routes.
import (
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

type RouteBuilder struct {
	rootPath    string
	currentPath string
	produces    []string
	consumes    []string
	httpMethod  string           // required
	handler     http.HandlerFunc // required
	// documentation
	doc         string
	notes       string
	operation   string
	readSample  interface{}
	writeSample interface{}
	parameters  []*Parameter
	errorMap    map[int]ResponseError
}

type ResponseError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Model   interface{} `json:"model"`
}

func NewRouteBuilder() *RouteBuilder {
	return &RouteBuilder{
		parameters: make([]*Parameter, 0),
		errorMap:   make(map[int]ResponseError),
	}
}

func (b *RouteBuilder) servicePath(path string) *RouteBuilder {
	b.rootPath = path
	return b
}

// To bind the route to a function.
// If this route is matched with the incoming Http Request then call this function with the *Request,*Response pair. Required.
func (b *RouteBuilder) To(function http.HandlerFunc) *RouteBuilder {
	b.handler = function
	return b
}

// Method specifies what HTTP method to match. Required.
func (b *RouteBuilder) Method(method string) *RouteBuilder {
	b.httpMethod = method
	return b
}

// Produces specifies what MIME types can be produced ; the matched one will appear in the Content-Type Http header.
func (b *RouteBuilder) Produces(mimeTypes ...string) *RouteBuilder {
	b.produces = mimeTypes
	return b
}

// Consumes specifies what MIME types can be consumes ; the Accept Http header must matched any of these
func (b *RouteBuilder) Consumes(mimeTypes ...string) *RouteBuilder {
	b.consumes = mimeTypes
	return b
}

// Path specifies the relative (w.r.t WebService root path) URL path to match. Default is "/".
func (b *RouteBuilder) Path(subPath string) *RouteBuilder {
	b.currentPath = subPath
	return b
}

// Doc tells what this route is all about. Optional.
// Both Doc and Notes will remove leading whitespace after a newline
// so that you can indent the doc text for prettier source code.
func (b *RouteBuilder) Doc(documentation string) *RouteBuilder {
	re := regexp.MustCompile("\n[ \t]+")
	b.doc = re.ReplaceAllString(documentation, "\n")
	return b
}

// A verbose explanation of the operation behavior. Optional.
func (b *RouteBuilder) Notes(notes string) *RouteBuilder {
	re := regexp.MustCompile("\n[ \t]+")
	b.notes = re.ReplaceAllString(notes, "\n")
	return b
}

// Reads tells what resource type will be read from the request payload. Optional.
// A parameter of type "body" is added ,required is set to true and the dataType is set to the qualified name of the sample's type.
func (b *RouteBuilder) Reads(sample interface{}) *RouteBuilder {
	b.readSample = sample
	typeAsName := reflect.TypeOf(sample).String()
	bodyParameter := &Parameter{&ParameterData{Name: "body"}}
	bodyParameter.beBody()
	bodyParameter.Required(true)
	bodyParameter.DataType(typeAsName)
	b.Param(bodyParameter)
	return b
}

// ParameterNamed returns a Parameter already known to the RouteBuilder. Returns nil if not.
// Use this to modify or extend information for the Parameter (through its Data()).
func (b RouteBuilder) ParameterNamed(name string) (p *Parameter) {
	for _, each := range b.parameters {
		if each.Data().Name == name {
			return each
		}
	}
	return p
}

// Writes tells what resource type will be written as the response payload. Optional.
func (b *RouteBuilder) Writes(sample interface{}) *RouteBuilder {
	b.writeSample = sample
	return b
}

// Param allows you to document the parameters of the Route. It adds a new Parameter (does not check for duplicates).
func (b *RouteBuilder) Param(parameter *Parameter) *RouteBuilder {
	b.parameters = append(b.parameters, parameter)
	return b
}

// Operation allows you to document what the acutal method/function call is of the Route.
// Unless called, the operation name is derived from the http.Handler set using To(..).
func (b *RouteBuilder) Operation(name string) *RouteBuilder {
	b.operation = name
	return b
}

// Returns allows you to document what responses (errors or regular) can be expected.
// The model parameter is optional ; either pass a struct instance or use nil if not applicable.
func (b *RouteBuilder) Returns(code int, message string, model interface{}) *RouteBuilder {
	err := ResponseError{
		Code:    code,
		Message: message,
		Model:   model,
	}
	b.errorMap[code] = err
	return b
}

// Build creates a new Route using the specification details collected by the RouteBuilder
func (b *RouteBuilder) Build() Route {
	if b.handler == nil {
		panic("[boneful] No function specified for route:" + b.currentPath)
	}
	route := Route{
		Method:         b.httpMethod,
		Path:           concatPath(b.rootPath, b.currentPath),
		Produces:       b.produces,
		Consumes:       b.consumes,
		Handler:        b.handler,
		Doc:            b.doc,
		Notes:          b.notes,
		Operation:      b.operation,
		ParameterDocs:  b.parameters,
		ResponseErrors: b.errorMap,
		ReadSample:     b.readSample,
		WriteSample:    b.writeSample,
	}
	route.postBuild()
	return route
}

func concatPath(path1, path2 string) string {
	return strings.TrimRight(path1, "/") + "/" + strings.TrimLeft(path2, "/")
}
