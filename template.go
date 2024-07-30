package main

type TemplateContext struct {
	Services []Service
}

type Service struct {
	Comment string
	Name    string
	Methods []Method
}

type Method struct {
	Name     string
	Service  string
	Input    string
	Output   string
	Comment  string
	HTTPRule HTTPRule
}

type HTTPRule struct {
	Method       string
	Path         string
	HasBody      bool
	Body         string
	ResponseBody string
	Additional   []HTTPRule
}

const tmpl = `
type ServiceResponseHandler func(resp http.ResponseWriter, reply any, err error)

{{- range .Services }}

type {{ .Name }}ResponseHandler ServiceResponseHandler

{{ .Comment }}type {{ .Name }} interface {
{{- range .Methods }}
	{{ .Comment }}	{{ .Name }}(context.Context, *{{ .Input }}) (*{{ .Output }}, error)
{{- end }}
}

func Register{{ .Name }}(eng *gin.Engine, svr {{ .Name }}, rh {{ .Name }}ResponseHandler) {
	if rh == nil {
		panic("response handler is nil")
	}
	init{{ .Name }}Router(eng, svr, rh)
}

func init{{ .Name }}Router(eng *gin.Engine, svr {{ .Name }}, rh {{ .Name }}ResponseHandler) {
{{- range .Methods }}
	eng.{{ .HTTPRule.Method }}("{{ .HTTPRule.Path }}", func(ctx *gin.Context) {
		in := &{{ .Input }}{}
		{{- if .HTTPRule.HasBody }}
		if err := ctx.ShouldBind(in{{ .HTTPRule.Body }}); err != nil {
			rh(ctx.Writer, nil, err)
			return
		}
		{{- else }}
		if err := decoder.Decode(in, ctx.Request.URL.Query()); err != nil {
			rh(ctx.Writer, nil, err)
			return
		}
		{{- end }}
		out, err := svr.{{ .Name }}(ctx, in)
		if err != nil {
			rh(ctx.Writer, nil, err)
			return
		}
		rh(ctx.Writer, out, nil)
	})
{{- end }}
}

{{- end }}
`
