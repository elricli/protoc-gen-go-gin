package main

import (
	"context"
	v1 "example/gen/foo/v1"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	v1.RegisterFooService(router, &FooService{})
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
