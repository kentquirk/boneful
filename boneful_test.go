package boneful

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func SampleHandler(rw http.ResponseWriter, req *http.Request) {
}

type readstr struct{}

func TestRegistryBasics(t *testing.T) {
	s := new(Service).Path("/test")

	s.Route(s.POST("/test").To(SampleHandler).
		Operation("Operation").
		Doc(`Documentation`).
		Notes(`Notes`).
		Param(PathParameter("hash", "The hash returned by the registration").DataType("string")).
		Consumes("application/json").
		Reads(readstr{}).
		Produces("application/json").
		Writes("docid").
		Returns(http.StatusBadRequest, "document is invalid", nil))

	assert.NotNil(t, s)
}

func TestDocGen(t *testing.T) {
	s := new(Service).Path("/test").
		Doc(`This is a test service`)

	s.Route(s.POST("/test").To(SampleHandler).
		Operation("Operation").
		Doc(`Documentation`).
		Notes(`Notes`).
		Param(PathParameter("hash", "The hash returned by the registration").DataType("string")).
		Consumes("application/json").
		Reads(readstr{}).
		Produces("application/json").
		Writes("docid").
		Returns(http.StatusBadRequest, "document is invalid", nil))

	assert.NotNil(t, s)
	buf := &bytes.Buffer{}
	s.GenerateDocumentation(buf)
	assert.True(t, buf.Len() > 500)
}

func TestSaveDocs(t *testing.T) {
	s := new(Service).Path("/").
		Doc(`This is a test service. It is designed to show you how to generate
			documentation automatically by running tests.`)

	s.Route(s.POST("/foo").To(SampleHandler).
		Operation("Operation").
		Doc(`Documentation`).
		Notes(`Notes`).
		Param(PathParameter("hash", "The hash returned by the registration").DataType("string")).
		Consumes("application/json").
		Reads(readstr{}).
		Produces("application/json").
		Writes("docid").
		Returns(http.StatusBadRequest, "document is invalid", nil))

	mux := s.Mux()

	req, _ := http.NewRequest("GET", "/md", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	f, err := os.Create("SAMPLE_DOCS.md")
	if err != nil {
		panic("Couldn't create SAMPLE_DOCS.md")
	}
	io.Copy(f, w.Body)
	f.Close()
}
