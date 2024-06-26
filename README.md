# protoc-gen-go-gin

# Usage

## Buf(recommended)

```yaml
# buf.gen.yaml
version: v2
plugins:
  - remote: buf.build/protocolbuffers/go
    out: gen
    opt: paths=source_relative
  # using binary. Install: go install github.com/elricli/protoc-gen-go-gin@latest
  - local: protoc-gen-go-gin 
    out: gen
    opt: paths=source_relative
  # using remote.
  - local: ["go", "run", "github.com/elricli/protoc-gen-go-gin@latest"]
    out: gen
    opt: paths=source_relative
```