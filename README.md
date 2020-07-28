# 🐐 

![👀 Lint](https://github.com/banditml/goat/workflows/Goat%20Inspector/badge.svg)
![🛶 Dev](https://github.com/banditml/goat/workflows/%F0%9F%9B%B6%20Ship%20It/badge.svg?branch=dev)
![🛶 Prod](https://github.com/banditml/goat/workflows/%F0%9F%9B%B6%20Ship%20It/badge.svg?branch=prod)

## 🐳 Docker things

We're using [multi-stage-builds](https://docs.docker.com/develop/develop-images/multistage-build/) here to trim the crap out of the final artifact.

1. *dev* - this stage is used for local development.

* `docker build --target dev --tag goat:dev && docker run --env PORT=9000 goat:dev`
* `docker-compose run goat`
* `docker-compose up`
* `make up`

All do (appx) the same thing:
* create a dev image with a code mounted volume
* runs container with cmd `air`.

2. *prod* - this stage is used for deployment.

It takes the built binary from `dev` and puts it in a minimal `alpine` image.
That's it.

## 📁 Interesting files

*.air.conf*
https://github.com/cosmtrek/air

Live reloading configuration for Go apps.

*.golangci.yml*
Master linter for Golang.  Runs a lot of checks right now which I thought might
be helpful.  If it becomes to opinionated we can tone down the settings in this file.

*Makefile*
* `up`
* `shell`

## Go stuff

[`main.go`](./main.go) is always the main entrypoint in Go, it exists by
itself, has the line `package main` in it, and has a single `func main()`
defined.  This is where all the dependencies are collected and the http server
is started.

### Dependencies

This service uses [FX](https://github.com/uber-go/fx) as the dependency
injection framework.  Convention is packages that are wrapped in FX-compatible
packages are suffixed with `fx`.  For example, I wrapped the `zap` package in
an FX compatible package and called it [`zapfx`](./zapfx) (same with `envfx`).

FX compatible packages have a single file called `module.go`, in which, a
single exported variable called `Module` is defined.  This allows very clean
dependency management in [`main.go`](./main.go):

```go
app := fx.New(
    zapfx.Module,
    envfx.Module,
    ...
)
```

Dependencies added to the app lifecycle are long-lived, which means they are
shared by all requests concurrently.  This is very good use of resources, but
means you should be careful about dependency-owned variables, and keep anything
dangerous allocated in thread scope (just make new stuff in the function call
and everything will be ok).

### Logging

The logging framework we use here is [zap](https://github.com/uber-go/zap).  It
is purpose built for structured backend logging which allows it to be not
magical, thus fast.
