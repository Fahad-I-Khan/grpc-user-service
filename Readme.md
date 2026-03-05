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

## Deploy to DockerHub

1. Create a new repository in DockerHub. 
2. Create personal access token on DockerHub for jenkins to login.

Both called **grpc-user-service**

3. Add DockerHub Credentials
```text
In Jenkins:

Manage Jenkins → Credentials → Add Credentials

Type: Username/Password

ID: dockerhub-creds-k8

Username: your DockerHub username

Password: DockerHub access token (not your password)
```

Username will be fahadkhan2105 as this is DockerHub username. 

4. Write Jenkins pipeline and run it.

---

## Deploy image to kind cluster

* Deployment
* Resource limits
* Liveness probe
* Readiness probe
* Service (ClusterIP)
* Then test gRPC inside cluster

### 🧱 Step 1 — Create k8s Folder

In your repo:

```Code
grpc-user-service/
 ├── k8s/
 │    ├── deployment.yaml
 │    └── service.yaml
```

``` Bash
mkdir k8s
```
cd into k8s directory.
``` Bash
touch deployment.yaml
touch service.yaml
```

After writting yaml.

### We need to start kind cluster

CMD's

``` Bash
kubectl create cluster --name <name>
```

``` Bash
kubectl create cluster --name dev-cluster
```

Then check

``` Bash
kubectl cluster-info
```

``` Bash
kubectl get nodes
```

This should show like:

```
NAME                        STATUS   ROLES           AGE   VERSION
dev-cluster-control-plane   Ready    control-plane   38s   v1.35.0
```

``` Bash
docker ps -a
```
This should also show docker container running with name: **dev-cluster-control-plane**

### Apply k8s

First check Nodes, Clusters are present if not create. 

``` Bash
 kubectl get nodes
```

``` Bash
kind get clusters
```

``` Bash
kubectl apply -f k8s/
```

Then Check pods, svc

``` Bash
kubectl get pods
```

``` Bash
kubectl get svc
```

Port-Forward

``` Bash
kubectl port-forward svc/grpc-user 50051:50051
```

Now add Ingress to Kind

``` Bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
```
Wait until pods are ready:
``` Bash
kubectl get pods -n ingress-nginx
```

You should see something like:
```
ingress-nginx-controller   Running
```

---

Add Ingress yaml

Then apply:
``` Bash
kubectl apply -f k8s/ingress.yaml
```

### Add Local DNS

Edit hosts file:

Mac:

``` bash
sudo nano /etc/hosts
```

Add:

```
127.0.0.1 grpc.local
```

### Test gRPC via Ingress

Now instead of:

``` Bash
localhost:50051
```

You use:

``` Bash
grpc.local
```

Now we will port-forward

``` Bash
kubectl port-forward -n ingress-nginx service/ingress-nginx-controller 8080:80
```

Now test:

``` Bash
grpcurl -plaintext grpc.local:8080 list
```

## Shell-based debug container

Since our pod don't have sh we can't sh into them. 

This won't work:

``` Bash
 kubectl exec -it grpc-user-59d55cc89f-484qw -- sh
```

```
error: Internal error occurred: Internal error occurred: error executing command in container: failed to exec in container: failed to start exec "2867aa18fff8993df8d2891918a91fcabf4ac5460e96f69244bf950d163eece3": OCI runtime exec failed: exec failed: unable to start container process: exec: "sh": executable file not found in $PATH
```

So we will create one. 

``` Bash
kubectl run net-debug \
  --rm -it \
  --restart=Never \
  --image=busybox \
  -- sh
```

``` Bash
nslookup grpc-user
```

``` Bash
wget grpc-user:50051
```
----

## How to check Ingress

``` Bash
kubectl get ingress
```

``` Bash
kubectl describe ingress grpc-user-ingress
```

``` Bash
kubectl get ingress grpc-user-ingress -o yaml
```

``` Bash
kubectl get pods -A | grep ingress
```

``` Bash
kubectl get svc -A | grep ingress
```

``` Bash
kubectl get svc -A
```

``` Bash
kubectl get pods -A
```

``` Bash
grpcurl -plaintext grpc.local:8080 list
```

```
                grpc.local
                    │
                    ▼
        NGINX Ingress Controller
                    │
                    ▼
            grpc-user Service
                    │
        ┌───────────┼───────────┐
        ▼           ▼           ▼
     Pod 1       Pod 2       Pod 3
```