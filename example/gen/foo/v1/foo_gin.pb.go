// Code generated by protoc-gen-go-gin. DO NOT EDIT.
// versions:
// - protoc-gen-go-gin v0.1.2
// - protoc            (unknown)
// source: foo/v1/foo.proto

package v1

import (
	context "context"
	gin "github.com/gin-gonic/gin"
	schema "github.com/gorilla/schema"
	http "net/http"
)

var _ = http.Request{}
var _ = context.TODO
var _ = gin.ContextKey
var decoder = schema.NewDecoder()

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

// FooService is a service that returns a FooResponse.
type FooService interface {
	// GetFoo returns a FooResponse.
	GetFoo(context.Context, *GetFooRequest) (*GetFooResponse, error)
	// SaveFoo returns a FooResponse.
	SaveFoo(context.Context, *SaveFooRequest) (*SaveFooResponse, error)
}

func RegisterFooService(eng *gin.Engine, svr FooService) {
	initFooServiceRouter(eng, svr)
}

func initFooServiceRouter(eng *gin.Engine, svr FooService) {
	eng.GET("/api/v1/foo", func(ctx *gin.Context) {
		in := &GetFooRequest{}
		if err := decoder.Decode(in, ctx.Request.URL.Query()); err != nil {
			DefaultErrorEncoder(ctx, err)
			return
		}
		out, err := svr.GetFoo(ctx, in)
		if err != nil {
			DefaultErrorEncoder(ctx, err)
			return
		}
		DefaultResponseEncoder(ctx, out)
	})
	eng.POST("/api/v1/foo", func(ctx *gin.Context) {
		in := &SaveFooRequest{}
		if err := ctx.ShouldBind(in); err != nil {
			DefaultErrorEncoder(ctx, err)
			return
		}
		out, err := svr.SaveFoo(ctx, in)
		if err != nil {
			DefaultErrorEncoder(ctx, err)
			return
		}
		DefaultResponseEncoder(ctx, out)
	})
}
