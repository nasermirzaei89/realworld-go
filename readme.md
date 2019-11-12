# ![RealWorld Example App](logo.png)

> ### pure go codebase containing real world examples (CRUD, auth, advanced patterns, etc) that adheres to the [RealWorld](https://github.com/gothinkster/realworld) spec and API.

[![Build Status](https://travis-ci.org/nasermirzaei89/realworld-go.svg?branch=master)](https://travis-ci.org/nasermirzaei89/realworld-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/nasermirzaei89/realworld-go)](https://goreportcard.com/report/github.com/nasermirzaei89/realworld-go)
[![GolangCI](https://golangci.com/badges/github.com/nasermirzaei89/realworld-go.svg)](https://golangci.com/r/github.com/nasermirzaei89/realworld-go)
[![Codecov](https://codecov.io/gh/nasermirzaei89/realworld-go/branch/master/graph/badge.svg)](https://codecov.io/gh/nasermirzaei89/realworld-go)
[![GoDoc](https://godoc.org/github.com/nasermirzaei89/realworld-go?status.svg)](https://godoc.org/github.com/nasermirzaei89/realworld-go)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://raw.githubusercontent.com/nasermirzaei89/realworld-go/master/LICENSE)

### [Demo](https://github.com/gothinkster/realworld)&nbsp;&nbsp;&nbsp;&nbsp;[RealWorld](https://github.com/gothinkster/realworld)

This codebase was created to demonstrate a fully fledged fullstack application built with **pure go** including CRUD operations, authentication, routing, pagination, and more.

We've gone to great lengths to adhere to the **pure go** community styleguides & best practices.

For more information on how to this works with other frontends/backends, head over to the [RealWorld](https://github.com/gothinkster/realworld) repo.


# How it works

> Describe the general architecture of your app here

# Getting started

## Run:

```bash
make run
```

## Test:

You should be in [this](https://github.com/gothinkster/realworld/tree/master/api) path

```bash
export APIURL=http://localhost:8080
./run-api-tests.sh
```

## Environments

1. `JWT_SECRET` with default value `secret` for sign jwt token with `HS256` algorithm
1. `API_ADDRESS` with default value `0.0.0.0:8080` for host and port of the API
