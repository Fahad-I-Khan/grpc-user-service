# Learn Protobuf gRPC and Kubernetes

Install 

1️⃣ Protocol Buffers compiler

``` Bash
brew install protobuf
```

Check:
``` Bash
protoc --version
```
---

2️⃣ Install Go protobuf plugins


``` Bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```
``` Bash
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
---

### Generate Go Code

``` Bash
protoc --go_out=. --go-grpc_out=. proto/user.proto
```
```
.
├── github.com
│   └── Fahad-I-Khan
│       └── grpc-user-service
│           ├── user_grpc.pb.go
│           └── user.pb.go
├── go.mod
├── go.sum
├── proto
│   └── user.proto
├── Readme.md
└── server
    └── main.go
```

this cmd created file's at wrong place.

**Correct Path** 

``` Bash
protoc \
  --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  proto/user.proto
```
```
.
├── go.mod
├── go.sum
├── proto
│   ├── user_grpc.pb.go
│   ├── user.pb.go
│   └── user.proto
├── Readme.md
└── server
    └── main.go

3 directories, 7 files
```
this one create correctly. and when we do changes in user.proto we have do run this cmd again.

After generating code it was showing import errors I just installed 

``` Bash
go get google.golang.org/grpc
```

And all import errors gone. 

---

**Run CMD**

``` Bash
go mod tidy
```

This will install dependencies.

**Run Locally**

``` Bash
go run server/main.go
```

---

### Test with grpcurl (like curl for gRPC)

Install

``` Bash
brew install grpcurl
```

Test

``` bash
grpcurl -plaintext -d '{"id":"123"}' \
localhost:50051 user.UserService/GetUser
```

---

**Reflection**

```go
reflection.Register(grpcServer) 
```
because grpcurl needs server reflection enabled to discover services automatically.

#### 🧠 Why Reflection Matters

Reflection lets tools like:

* grpcurl
* Postman (gRPC mode)
* BloomRPC

discover:

* Services
* Methods
* Request/Response schemas

Without reflection, you’d need to manually pass the `.proto` file to grpcurl.

---

## Tag & Run Docker Image

After adding distroless image

``` Bash
docker build -t grpc-user-service:1.0 .
```

``` Bash
docker images | grep grpc-user-service
```

``` Bash
docker run -p 50051:50051 grpc-user-service:1.0
```

``` Bash
grpcurl -plaintext localhost:50051 list
```

**Important Debug Note**

Distroless has no shell.

So this WON’T work:

``` Bash
docker exec -it container sh
```

Production containers should not have shells.

---

