# example

## Run example

```shell
go run main.go
```

## cURL

### cURL with JSON

```shell
curl http://localhost:8080/api/v1/foo -X POST -H "Content-Type: application/json" -d '{"foo": "this is a message from cURL"}'
```

### cURL with Protobuf Payload

```shell
protoc --encode=foo.v1.SaveFooRequest -I {google/api/annotation.proto base path} -I . ./proto/foo/v1/foo.proto  < message.txt | curl -X POST -H "Content-Type: application/x-protobuf" --data-binary @- http://localhost:8080/api/v1/foo
```

message.txt content:
```text
foo: "this is a message from message.txt"
```