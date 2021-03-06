package boneful

var mdTemplate = `
---
# ` + "`" + `{{.RootPath}}` + "`" + `

{{.Documentation}}


{{range .Routes}}
* [{{.Operation}}](#{{lower .Operation}})
{{end}}


{{range .Routes}}
---
## {{.Operation}}

### ` + "`" + `{{.Method}} {{.Path}}` + "`" + `

_{{.Doc}}_

{{.Notes}}

{{if .ParameterDocs}}
_**Parameters:**_

Name | Kind | Description | DataType
---- | ---- | ----------- | --------
{{range .ParameterDocs}} {{.Data.Name}} | {{.Data.ParameterKind}} | {{.Data.Description}} | {{.Data.DataType}}
{{end}}
{{end}}

{{if .Consumes}}
_**Consumes:**_ ` + "`" + `{{.Consumes}}` + "`" + `
{{end}}
{{if .Reads}}
_**Reads:**_
` + "```{{.CodeFormat}}" + `
        {{.Reads}}
` + "```" + `
{{end}}
{{if .Produces}}
_**Produces:**_ ` + "`" + `{{.Produces}}` + "`" + `
{{end}}
{{if .Writes}}
_**Writes:**_
` + "```{{.CodeFormat}}" + `
        {{.Writes}}
` + "```" + `
{{end}}
{{if .ResponseErrors}}
_**Error returns:**_

Code | Meaning
---- | --------
{{range .ResponseErrors}} {{.Code}} | {{.Message}}
{{end}}
{{end}}
{{end}}
`
