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

type ServiceResponseHandler func(resp http.ResponseWriter, reply any, err error)

type FooServiceResponseHandler ServiceResponseHandler

// FooService is a service that returns a FooResponse.
type FooService interface {
	// GetFoo returns a FooResponse.
	GetFoo(context.Context, *GetFooRequest) (*GetFooResponse, error)
	// SaveFoo returns a FooResponse.
	SaveFoo(context.Context, *SaveFooRequest) (*SaveFooResponse, error)
}

func RegisterFooService(eng *gin.Engine, svr FooService, rh FooServiceResponseHandler) {
	if rh == nil {
		panic("response handler is nil")
	}
	initFooServiceRouter(eng, svr, rh)
}

func initFooServiceRouter(eng *gin.Engine, svr FooService, rh FooServiceResponseHandler) {
	eng.GET("/api/v1/foo", func(ctx *gin.Context) {
		in := &GetFooRequest{}
		if err := decoder.Decode(in, ctx.Request.URL.Query()); err != nil {
			rh(ctx.Writer, nil, err)
			return
		}
		out, err := svr.GetFoo(ctx, in)
		if err != nil {
			rh(ctx.Writer, nil, err)
			return
		}
		rh(ctx.Writer, out, nil)
	})
	eng.POST("/api/v1/foo", func(ctx *gin.Context) {
		in := &SaveFooRequest{}
		if err := ctx.ShouldBind(in); err != nil {
			rh(ctx.Writer, nil, err)
			return
		}
		out, err := svr.SaveFoo(ctx, in)
		if err != nil {
			rh(ctx.Writer, nil, err)
			return
		}
		rh(ctx.Writer, out, nil)
	})
}
