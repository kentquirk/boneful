package boneful

const (
	// PathParameterKind = indicator of Request parameter type "path"
	PathParameterKind = iota

	// QueryParameterKind = indicator of Request parameter type "query"
	QueryParameterKind

	// BodyParameterKind = indicator of Request parameter type "body"
	BodyParameterKind

	// HeaderParameterKind = indicator of Request parameter type "header"
	HeaderParameterKind

	// FormParameterKind = indicator of Request parameter type "form"
	FormParameterKind
)

// Parameter is for documententing the parameter used in a Http Request
type Parameter struct {
	D *ParameterData `json:"data"`
}

// ParameterData represents the state of a Parameter.
// It is made public to make it accessible to e.g. the Swagger package.
type ParameterData struct {
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	DataType        string            `json:"datatype"`
	DataFormat      string            `json:"dataformat"`
	Kind            int               `json:"kind"`
	Required        bool              `json:"required"`
	AllowableValues map[string]string `json:"allowablevalues"`
	AllowMultiple   bool              `json:"allowmultiple"`
	DefaultValue    string            `json:"defaultvalue"`
}

// Data returns the state of the Parameter
func (p *Parameter) Data() ParameterData {
	return *p.D
}

// ParameterKind returns the type of this parameter as a string.
func (p ParameterData) ParameterKind() string {
	switch p.Kind {
	case PathParameterKind:
		return "Path"
	case QueryParameterKind:
		return "Query"
	case BodyParameterKind:
		return "Body"
	case HeaderParameterKind:
		return "Header"
	case FormParameterKind:
		return "Form"
	default:
		return "Unknown"
	}
}

// PathParameter creates a new Parameter of kind Path for documentation purposes.
// It is initialized as required with string as its DataType.
func PathParameter(name, description string) *Parameter {
	p := &Parameter{&ParameterData{Name: name, Description: description, Required: true, DataType: "string"}}
	p.bePath()
	return p
}

// QueryParameter creates a new Parameter of kind Query for documentation purposes.
// It is initialized as not required with string as its DataType.
func QueryParameter(name, description string) *Parameter {
	p := &Parameter{&ParameterData{Name: name, Description: description, Required: false, DataType: "string"}}
	p.beQuery()
	return p
}

// BodyParameter creates a new Parameter of kind Body for documentation purposes.
// It is initialized as required without a DataType.
func BodyParameter(name, description string) *Parameter {
	p := &Parameter{&ParameterData{Name: name, Description: description, Required: true}}
	p.beBody()
	return p
}

// HeaderParameter creates a new Parameter of kind (Http) Header for documentation purposes.
// It is initialized as not required with string as its DataType.
func HeaderParameter(name, description string) *Parameter {
	p := &Parameter{&ParameterData{Name: name, Description: description, Required: false, DataType: "string"}}
	p.beHeader()
	return p
}

// FormParameter creates a new Parameter of kind Form (using application/x-www-form-urlencoded) for documentation purposes.
// It is initialized as required with string as its DataType.
func FormParameter(name, description string) *Parameter {
	p := &Parameter{&ParameterData{Name: name, Description: description, Required: false, DataType: "string"}}
	p.beForm()
	return p
}

// Kind returns the parameter type indicator (see const for valid values)
func (p *Parameter) Kind() int {
	return p.D.Kind
}

func (p *Parameter) bePath() *Parameter {
	p.D.Kind = PathParameterKind
	return p
}

func (p *Parameter) beQuery() *Parameter {
	p.D.Kind = QueryParameterKind
	return p
}

func (p *Parameter) beBody() *Parameter {
	p.D.Kind = BodyParameterKind
	return p
}

func (p *Parameter) beHeader() *Parameter {
	p.D.Kind = HeaderParameterKind
	return p
}

func (p *Parameter) beForm() *Parameter {
	p.D.Kind = FormParameterKind
	return p
}

// Required sets the required field and returns the receiver
func (p *Parameter) Required(required bool) *Parameter {
	p.D.Required = required
	return p
}

// AllowMultiple sets the allowMultiple field and returns the receiver
func (p *Parameter) AllowMultiple(multiple bool) *Parameter {
	p.D.AllowMultiple = multiple
	return p
}

// AllowableValues sets the allowableValues field and returns the receiver
func (p *Parameter) AllowableValues(values map[string]string) *Parameter {
	p.D.AllowableValues = values
	return p
}

// DataType sets the dataType field and returns the receiver
func (p *Parameter) DataType(typeName string) *Parameter {
	p.D.DataType = typeName
	return p
}

// DataFormat sets the dataFormat field for Swagger UI
func (p *Parameter) DataFormat(formatName string) *Parameter {
	p.D.DataFormat = formatName
	return p
}

// DefaultValue sets the default value field and returns the receiver
func (p *Parameter) DefaultValue(stringRepresentation string) *Parameter {
	p.D.DefaultValue = stringRepresentation
	return p
}

// Description sets the description value field and returns the receiver
func (p *Parameter) Description(doc string) *Parameter {
	p.D.Description = doc
	return p
}
