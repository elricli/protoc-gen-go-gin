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
var (
	DefaultResponseEncoder = func(ctx *gin.Context, reply any) {
		switch ctx.ContentType() {
		case "application/xml":
			ctx.XML(http.StatusOK, reply)
		case "application/yaml":
			ctx.YAML(http.StatusOK, reply)
		case "application/x-protobuf":
			ctx.ProtoBuf(http.StatusOK, reply)
		case "application/toml":
			ctx.TOML(http.StatusOK, reply)
		default:
			ctx.JSON(http.StatusOK, reply)
		}
	}
	DefaultErrorEncoder = func(ctx *gin.Context, err error) {
		resp := gin.H{"error": err.Error()}
		switch ctx.ContentType() {
		case "application/xml":
			ctx.XML(http.StatusInternalServerError, resp)
		case "application/yaml":
			ctx.YAML(http.StatusInternalServerError, resp)
		case "application/x-protobuf":
			ctx.ProtoBuf(http.StatusInternalServerError, resp)
		case "application/toml":
			ctx.TOML(http.StatusInternalServerError, resp)
		default:
			ctx.JSON(http.StatusInternalServerError, resp)
		}
	}
)

{{- range .Services }}

{{ .Comment }}type {{ .Name }} interface {
{{- range .Methods }}
	{{ .Comment }}	{{ .Name }}(context.Context, *{{ .Input }}) (*{{ .Output }}, error)
{{- end }}
}

func Register{{ .Name }}(eng *gin.Engine, svr {{ .Name }}) {
	init{{ .Name }}Router(eng, svr)
}

func init{{ .Name }}Router(eng *gin.Engine, svr {{ .Name }}) {
{{- range .Methods }}
	eng.{{ .HTTPRule.Method }}("{{ .HTTPRule.Path }}", func(ctx *gin.Context) {
		in := &{{ .Input }}{}
		{{- if .HTTPRule.HasBody }}
		if err := ctx.ShouldBind(in{{ .HTTPRule.Body }}); err != nil {
			DefaultErrorEncoder(ctx, err)
			return
		}
		{{- else }}
		if err := decoder.Decode(in, ctx.Request.URL.Query()); err != nil {
			DefaultErrorEncoder(ctx, err)
			return
		}
		{{- end }}
		out, err := svr.{{ .Name }}(ctx, in)
		if err != nil {
			DefaultErrorEncoder(ctx, err)
			return
		}
		DefaultResponseEncoder(ctx, out)
	})
{{- end }}
}

{{- end }}
`
