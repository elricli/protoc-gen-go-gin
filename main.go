package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

func main() {
	protogen.Options{}.Run(func(p *protogen.Plugin) error {
		// p.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_SUPPORTS_EDITIONS | pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		// p.SupportedEditionsMinimum = descriptorpb.Edition_EDITION_2024
		for _, f := range p.Files {
			if !f.Generate {
				continue
			}
			var hhr bool
			for _, s := range f.Services {
				if hasHTTPRule(s) {
					hhr = true
					break
				}
			}
			if !hhr {
				fmt.Fprintf(os.Stderr, "file %s has no http rules\n", f.Desc.Path())
				return nil
			}
			gf := p.NewGeneratedFile(f.GeneratedFilenamePrefix+"_gin.pb.go", f.GoImportPath)
			gf.P("// Code generated by protoc-gen-go-gin. DO NOT EDIT.")
			gf.P("// versions:")
			gf.P("// - protoc-gen-go-gin v1.0.0") //TODO: add version
			gf.P("// - protoc            ", protocVersion(p))
			gf.P("// source: ", f.Desc.Path())
			gf.P()
			gf.P("package ", f.GoPackageName)
			gf.P()
			gf.P("import (")
			gf.P(`	http "net/http"`)
			gf.P(`	context "context"`)
			gf.P()
			gf.P(`	gin "github.com/gin-gonic/gin"`)
			gf.P(`	schema "github.com/gorilla/schema"`)
			gf.P(")")
			gf.P()
			gf.P(`var (`)
			gf.P(`	decoder = schema.NewDecoder()`)
			gf.P(`)`)

			var tc TemplateContext
			for _, s := range f.Services {
				buildServics(&tc, s)
			}
			template.Must(template.New("protoc-gen-go-gin").Parse(tmpl)).Execute(gf, tc)
		}

		return nil
	})
}

func protocVersion(gen *protogen.Plugin) string {
	v := gen.Request.GetCompilerVersion()
	if v == nil {
		return "(unknown)"
	}
	var suffix string
	if s := v.GetSuffix(); s != "" {
		suffix = "-" + s
	}
	return fmt.Sprintf("v%d.%d.%d%s", v.GetMajor(), v.GetMinor(), v.GetPatch(), suffix)
}

func buildServics(tc *TemplateContext, s *protogen.Service) {
	if !hasHTTPRule(s) {
		fmt.Fprintf(os.Stderr, "service %s has no http rules\n", s.GoName)
		return
	}
	svc := Service{
		Comment: s.Comments.Leading.String(),
		Name:    s.GoName,
	}
	for _, m := range s.Methods {
		rule := proto.GetExtension(m.Desc.Options(), annotations.E_Http).(*annotations.HttpRule)
		meth := Method{
			Name:    m.GoName,
			Service: s.GoName,
			Input:   m.Input.GoIdent.GoName,
			Output:  m.Output.GoIdent.GoName,
			Comment: m.Comments.Leading.String(),
			HTTPRule: HTTPRule{
				ResponseBody: m.Output.GoIdent.GoName,
			},
		}
		switch pattern := rule.GetPattern().(type) {
		case *annotations.HttpRule_Get:
			meth.HTTPRule.Path = pattern.Get
			meth.HTTPRule.Method = "GET"
		case *annotations.HttpRule_Post:
			meth.HTTPRule.Path = pattern.Post
			meth.HTTPRule.Method = "POST"
			meth.HTTPRule.HasBody = true
			if rule.GetBody() == "" {
				fmt.Fprintf(os.Stderr, "service %s method %s has no body\n", s.GoName, m.GoName)
			}
		case *annotations.HttpRule_Put:
			meth.HTTPRule.Path = pattern.Put
			meth.HTTPRule.Method = "PUT"
			meth.HTTPRule.HasBody = true
			if rule.GetBody() == "" {
				fmt.Fprintf(os.Stderr, "service %s method %s has no body\n", s.GoName, m.GoName)
			}
		case *annotations.HttpRule_Delete:
			meth.HTTPRule.Path = pattern.Delete
			meth.HTTPRule.Method = "DELETE"
		case *annotations.HttpRule_Patch:
			meth.HTTPRule.Path = pattern.Patch
			meth.HTTPRule.Method = "PATCH"
			meth.HTTPRule.HasBody = true
			if rule.GetBody() == "" {
				fmt.Fprintf(os.Stderr, "service %s method %s has no body\n", s.GoName, m.GoName)
			}
		case *annotations.HttpRule_Custom:
			meth.HTTPRule.Path = pattern.Custom.Path
			meth.HTTPRule.Method = pattern.Custom.Kind
		}
		if rule.Body != "" && rule.Body != "*" {
			meth.HTTPRule.Body = "." + camelCaseVars(rule.Body)
		}
		svc.Methods = append(svc.Methods, meth)
	}
	tc.Services = append(tc.Services, svc)
}

func hasHTTPRule(s *protogen.Service) bool {
	for _, m := range s.Methods {
		if proto.HasExtension(m.Desc.Options(), annotations.E_Http) {
			return true
		}
	}
	return false
}

func camelCaseVars(s string) string {
	subs := strings.Split(s, ".")
	vars := make([]string, 0, len(subs))
	for _, sub := range subs {
		vars = append(vars, camelCase(sub))
	}
	return strings.Join(vars, ".")
}

// camelCase returns the CamelCased name.
// If there is an interior underscore followed by a lower case letter,
// drop the underscore and convert the letter to upper case.
// There is a remote possibility of this rewrite causing a name collision,
// but it's so remote we're prepared to pretend it's nonexistent - since the
// C++ generator lowercase names, it's extremely unlikely to have two fields
// with different capitalization.
// In short, _my_field_name_2 becomes XMyFieldName_2.
func camelCase(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'X')
		i++
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && isASCIILower(s[i+1]) {
			continue // Skip the underscore in s.
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if isASCIILower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		for i+1 < len(s) && isASCIILower(s[i+1]) {
			i++
			t = append(t, s[i])
		}
	}
	return string(t)
}

// Is c an ASCII lower-case letter?
func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// Is c an ASCII digit?
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}
