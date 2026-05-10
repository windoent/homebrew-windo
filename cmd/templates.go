package cmd

const (
	goModTemplate = `module {{modulePath}}

go 1.19

require (
	git.code.tencent.com/windo-/bd v0.0.0-20260428033354-69a804b5451b
	github.com/go-kratos/kratos/contrib/registry/consul/v2 v2.0.0-20230728114959-072b2ee7693c
	github.com/go-kratos/kratos/v2 v2.6.3
	github.com/google/wire v0.6.0
	github.com/hashicorp/consul/api v1.23.0
	github.com/spf13/viper v1.16.0
	go.uber.org/automaxprocs v1.5.1
	gorm.io/gorm v1.25.2
)

require github.com/google/uuid v1.3.0
`

	makefileTemplate = `GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)

ifeq ($(GOHOSTOS), windows)
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find internal -name *.proto")
	API_PROTO_FILES=$(shell $(Git_Bash) -c "find api -name *.proto")
else
	INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
	API_PROTO_FILES=$(shell find api -name *.proto)
endif

.PHONY: init
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest

.PHONY: config
config:
	protoc --proto_path=./internal \
	       --proto_path=./third_party \
	       --go_out=paths=source_relative:./internal \
	       $(INTERNAL_PROTO_FILES)

.PHONY: api
api:
	protoc --proto_path=./api \
	       --proto_path=./third_party \
	       --go_out=paths=source_relative:./api \
	       --go-http_out=paths=source_relative:./api \
	       --go-grpc_out=paths=source_relative:./api \
	       --openapi_out=fq_schema_naming=true,default_response=false:. \
	       $(API_PROTO_FILES)

.PHONY: build
build:
	 mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: generate
generate:
	go mod tidy
	go get github.com/google/wire/cmd/wire@latest
	go generate ./...

.PHONY: all
all:
	make api;
	make config;
	make generate;

help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
`

	dockerfileTemplate = `FROM golang:1.19 AS builder
COPY . /src
WORKDIR /src

RUN umask 0007 \
	    && mkdir -p ~/.ssh \
	    && ssh-keyscan git.code.tencent.com >> ~/.ssh/known_hosts \
	    && git config --global url."ssh://git@git.code.tencent.com/".insteadOf  "https://git.code.tencent.com/"

ADD .ssh/id_rsa /root/.ssh/id_rsa
ADD .ssh/id_rsa.pub /root/.ssh/id_rsa.pub
ADD .ssh/config /root/.ssh/config
RUN chmod 600 /root/.ssh/id_rsa

RUN go env -w GO111MODULE=on
RUN go env -w GOPRIVATE=*.code.tencent.com
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod tidy

RUN GOPROXY=https://goproxy.cn make build

FROM ubuntu

RUN apt-get update
RUN apt-get install -y --no-install-recommends ca-certificates curl tzdata
RUN rm /etc/localtime && ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

COPY --from=builder /src/bin /app
COPY --from=builder /src/configs/config.yaml /data/conf/config.yaml

WORKDIR /app

EXPOSE 9000
VOLUME /data/conf

CMD ["./{{projectName}}", "-conf", "/data/conf"]
`

	readmeTemplate = "# {{projectName}} grpc 服务\n\n\n## migrate 迁移数据库\n\n\n```\nmake build\n./bin/{{projectName}} -conf=<config path> -migrate=true\neg:\n./bin/{{projectName}} -conf=/path/to/configs -migrate=true\n```\n\n## Docker\n```bash\n# build\ndocker build -t <your-docker-image-name> .\n\n# run\ndocker run --rm  -p 9000:9000 -v </path/to/your/configs>:/data/conf <your-docker-image-name>\n\n```\n"

	gitignoreTemplate = "# Reference https://github.com/github/gitignore/blob/master/Go.gitignore\n# Binaries for programs and plugins\n*.exe\n*.exe~\n*.dll\n*.dylib\n\n# Test binary, built with `go test -c`\n*.test\n\n# Output of the go coverage tool, specifically when used with LiteIDE\n*.out\n\n# Dependency directories (remove the comment below to include it)\nvendor/\n\n# Go workspace file\ngo.work\n\n# Compiled Object files, Static and Dynamic libs (Shared Objects)\n*.o\n*.a\n*.so\n\n# OS General\nThumbs.db\n.DS_Store\n\n# project\n*.cert\n*.key\n*.log\nbin/\n\n# Develop tools\n.vscode/\n.idea/\n*.swp\n\ngo.sum\n\nconfigs/config.yaml\n"

	gitlabCiTemplate = `stages:
  - build-master
  - build-release
  - build-dev

build-master:
  stage: build-master
  script:
    - mkdir -p configs
    - cp profile/prod.yaml configs/config.yaml
    - /usr/local/bin/deploy_yd_docker_image bd-services-{{projectName}}
  only:
    - master
  tags:
    - windoent

build-release:
  stage: build-release
  script:
    - mkdir -p configs
    - cp profile/release.yaml configs/config.yaml
    - /usr/local/bin/deploy_yd_docker_image release-bd-services-{{projectName}}
  only:
    - release
  tags:
    - windoent

build-dev:
  stage: build-dev
  script:
    - mkdir -p configs
    - cp profile/dev.yaml configs/config.yaml
    - /usr/local/bin/deploy_yd_docker_image dev-bd-services-{{projectName}}
  only:
    - develop
  tags:
    - windoent
`

	cmdMainTemplate = `package main

import (
	"flag"
	"git.code.tencent.com/windo-/bd/common"
	config2 "git.code.tencent.com/windo-/bd/config"
	"git.code.tencent.com/windo-/bd/config/conf"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	_ "go.uber.org/automaxprocs"
	"os"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	Name       string
	Version    string
	flagconf   string
	Flagmigrate bool

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
	flag.BoolVar(&Flagmigrate, "migrate", false, "gorm migrate, eg: -migrate true")
}

func newApp(logger log.Logger, gs *grpc.Server, rr registry.Registrar) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
		),
		kratos.Registrar(rr),
	)
}

func main() {
	flag.Parse()

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	var err interface{}

	if err = c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err = c.Scan(&bc); err != nil {
		panic(err)
	}
	config2.Bootstrap = &bc
	Name = bc.Server.Name

	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	app, cleanup, err := wireApp(bc.Server, bc.Consul, bc.Data, logger)
	if err != nil {
		panic(err)
	}

	if Flagmigrate {
		return
	}

	err = common.TracerProvider(bc.Trace.Endpoint, bc.Trace.Token, Name, id, Version)
	if err != nil {
		panic(err)
	}

	defer cleanup()

	if err = app.Run(); err != nil {
		panic(err)
	}
}
`

	wireTemplate = `//go:build wireinject
// +build wireinject

package main

import (
	"git.code.tencent.com/windo-/bd/config/conf"
	"{{modulePath}}/internal/biz"
	"{{modulePath}}/internal/data"
	"{{modulePath}}/internal/server"
	"{{modulePath}}/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

func wireApp(*conf.Server, *conf.Consul, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
`

	bizTemplate = `package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewGreeterUsecase)
`

	bizGreeterTemplate = `package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

// GreeterUsecase is a Greeter usecase.
type GreeterUsecase struct {
	repo Repo
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewGreeterUsecase(repo Repo, logger log.Logger) *GreeterUsecase {
	return &GreeterUsecase{repo: repo, log: log.NewHelper(logger)}
}

// Repo is a interface for data operations
type Repo interface {
	// Add your methods here
}
`

	dataTemplate = `package data

import (
	"git.code.tencent.com/windo-/bd/config/conf"
	"git.code.tencent.com/windo-/bd/db"
	"git.code.tencent.com/windo-/bd/db/redis"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)

// Data .
type Data struct {
	// TODO wrapped database client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	_, err := db.Resolver()
	if err != nil {
		return nil, nil, err
	}
	_, err = db.OwnerResolver()
	if err != nil {
		return nil, nil, err
	}
	_, err = db.BResolver()
	if err != nil {
		return nil, nil, err
	}
	_, err = db.PointsResolver()
	if err != nil {
		return nil, nil, err
	}
	_, err = db.BackUpResolver()
	if err != nil {
		return nil, nil, err
	}

	err = redis.Init()
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
		sqlDB, err := db.Source.DB()
		if err != nil {
			sqlDB.Close()
		}
		sqlDB, err = db.OwnerSource.DB()
		if err != nil {
			sqlDB.Close()
		}
		sqlDB, err = db.BSource.DB()
		if err != nil {
			sqlDB.Close()
		}
		sqlDB, err = db.Points.DB()
		if err != nil {
			sqlDB.Close()
		}
		sqlDB, err = db.BackUp.DB()
		if err != nil {
			sqlDB.Close()
		}
	}
	return &Data{}, cleanup, nil
}
`

	dataRepoTemplate = `package data

import (
	"context"
	"{{modulePath}}/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
)

// GreeterRepo implements the biz.Repo interface
type GreeterRepo struct {
	log *log.Helper
}

// NewGreeterRepo new a Greeter repo.
func NewGreeterRepo(logger log.Logger) *GreeterRepo {
	return &GreeterRepo{log: log.NewHelper(logger)}
}

// Ensure GreeterRepo implements biz.Repo
var _ biz.Repo = (*GreeterRepo)(nil)
`

	serverTemplate = `package server

import (
	"git.code.tencent.com/windo-/bd/config/conf"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/google/wire"
	consulAPI "github.com/hashicorp/consul/api"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewGRPCServer, NewHTTPServer, NewRegistrar)

// NewRegistrar introduces consul
func NewRegistrar(conf *conf.Consul) registry.Registrar {
	var errs interface{}
	c := consulAPI.DefaultConfig()
	c.Address = conf.Address
	c.Scheme = conf.Scheme
	cli, err := consulAPI.NewClient(c)
	if err != nil {
		errs = err
		panic(errs)
	}
	r := consul.New(cli, consul.WithHealthCheck(false))
	return r
}
`

	grpcTemplate = `package server

import (
	v1 "git.code.tencent.com/windo-/bd/api/{{projectName}}/v1"
	"git.code.tencent.com/windo-/bd/config/conf"
	"{{modulePath}}/internal/service"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, greeter *service.GreeterService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
			tracing.Server(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.Register{{projectName}}Server(srv, greeter)
	return srv
}
`

	httpTemplate = `package server

import (
	"git.code.tencent.com/windo-/bd/config/conf"
	"{{modulePath}}/internal/service"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	return srv
}
`

	serviceTemplate = `package service

import "github.com/google/wire"

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewGreeterService)
`

	serviceGreeterTemplate = `package service

import (
	"context"

	v1 "git.code.tencent.com/windo-/bd/api/{{projectName}}/v1"
	"{{modulePath}}/internal/biz"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.Unimplemented{{projectName}}Server
	uc *biz.GreeterUsecase
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase) *GreeterService {
	return &GreeterService{uc: uc}
}

// TODO: Implement your service methods here
`

	configTemplate = `server:
  env: loc
  name: bd.{{projectName}}.provider
  http:
    addr: 0.0.0.0:8000
    timeout: 1s
  grpc:
    addr: 0.0.0.0:9021
    timeout: 1s
data:
  database:
    master: root:123456@tcp(127.0.0.1:3306)/bd_base
    slave: root:123456@tcp(127.0.0.1:3306)/bd_base
    bdb: root:ey9zE4cyQNQl@tcp(114.117.4.43:3306)/shop_base
    pdb: root:123456@tcp(127.0.0.1:3306)/point_base
    slavePdb: root:123456@tcp(127.0.0.1:3306)/point_base
    backup: root:123456@tcp(127.0.0.1:3306)/bd_base_backup
    slaveBackup: root:123456@tcp(127.0.0.1:3306)/bd_base_backup
  redis:
    addr: 127.0.0.1:6379
    password: "123123"
    db: 7
    queue: 8
    read_timeout: 0.2s
    write_timeout: 0.2s
trace:
  endpoint: http://127.0.0.1:14268/api/traces
  token: KdydzIcSfIpSxUOmiyjn
consul:
  address: 127.0.0.1:8500
  scheme: http
`

	profileDevTemplate = `server:
  env: dev
  name: bd.{{projectName}}.provider
  http:
    addr: 0.0.0.0:8000
    timeout: 1s
  grpc:
    addr: 0.0.0.0:9021
    timeout: 1s
data:
  database:
    master: root:123456@tcp(127.0.0.1:3306)/bd_base
    slave: root:123456@tcp(127.0.0.1:3306)/bd_base
    bdb: root:ey9zE4cyQNQl@tcp(114.117.4.43:3306)/shop_base
    pdb: root:123456@tcp(127.0.0.1:3306)/point_base
    slavePdb: root:123456@tcp(127.0.0.1:3306)/point_base
    backup: root:123456@tcp(127.0.0.1:3306)/bd_base_backup
    slaveBackup: root:123456@tcp(127.0.0.1:3306)/bd_base_backup
  redis:
    addr: 127.0.0.1:6379
    password: "123123"
    db: 7
    queue: 8
    read_timeout: 0.2s
    write_timeout: 0.2s
trace:
  endpoint: http://127.0.0.1:14268/api/traces
  token: KdydzIcSfIpSxUOmiyjn
consul:
  address: 127.0.0.1:8500
  scheme: http
`

	profileProdTemplate = `server:
  env: prod
  name: bd.{{projectName}}.provider
  http:
    addr: 0.0.0.0:8000
    timeout: 1s
  grpc:
    addr: 0.0.0.0:9021
    timeout: 1s
data:
  database:
    master: ${DB_MASTER}
    slave: ${DB_SLAVE}
    bdb: ${DB_BDB}
    pdb: ${DB_PDB}
    slavePdb: ${DB_SLAVE_PDB}
    backup: ${DB_BACKUP}
    slaveBackup: ${DB_SLAVE_BACKUP}
  redis:
    addr: ${REDIS_ADDR}
    password: "${REDIS_PASSWORD}"
    db: 7
    queue: 8
    read_timeout: 0.2s
    write_timeout: 0.2s
trace:
  endpoint: ${TRACE_ENDPOINT}
  token: ${TRACE_TOKEN}
consul:
  address: ${CONSUL_ADDRESS}
  scheme: http
`
)
