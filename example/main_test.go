package main

import (
	"bytes"
	v1 "example/gen/foo/v1"
	"io"
	"net/http"
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestPBRequest(t *testing.T) {
	bs, err := proto.Marshal(&v1.SaveFooRequest{
		Foo: "bar-from-http-client",
	})
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}
	resp, err := http.Post("http://localhost:8080/api/v1/foo", "application/x-protobuf", bytes.NewReader(bs))
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	body := new(bytes.Buffer)
	io.Copy(body, resp.Body)
	var out v1.SaveFooResponse
	if err := proto.Unmarshal(body.Bytes(), &out); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if out.Foo != "bar-from-http-client" {
		t.Fatalf("unexpected response: %s", out.String())
	}
}
