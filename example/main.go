package main

import (
	"context"
	"encoding/json"
	v1 "example/gen/foo/v1"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	v1.RegisterFooService(router, &FooService{}, func(rw http.ResponseWriter, reply any, err error) {
		var resp struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
			Data any    `json:"data"`
		}
		if err != nil {
			resp.Code = -1
			resp.Msg = err.Error()
		} else {
			resp.Data = reply
		}
		json.NewEncoder(rw).Encode(resp)
	})
	router.Run(":8080")
}

type FooService struct{}

// GetFoo implements v1.FooService.
func (f *FooService) GetFoo(_ context.Context, req *v1.GetFooRequest) (*v1.GetFooResponse, error) {
	resp := &v1.GetFooResponse{
		Foo: req.Foo,
	}
	return resp, nil
}

// SaveFoo implements v1.FooService.
func (f *FooService) SaveFoo(_ context.Context, req *v1.SaveFooRequest) (*v1.SaveFooResponse, error) {
	resp := &v1.SaveFooResponse{
		Foo: req.Foo,
	}
	return resp, nil
}
