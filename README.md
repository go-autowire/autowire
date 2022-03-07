# 🔌 autowire [![GoDoc][doc-img]][doc] [![Github release][release-img]][release] [![Go Report Card][report-card-img]][report-card]

[doc-img]: http://img.shields.io/badge/GoDoc-Reference-blue.svg
[doc]: https://godoc.org/github.com/go-autowire/autowire

[release-img]: https://img.shields.io/github/release/go-autowire/autowire.svg
[release]: https://github.com/go-autowire/autowire

[report-card-img]: https://goreportcard.com/badge/github.com/go-autowire/autowire
[report-card]: https://goreportcard.com/report/github.com/go-autowire/autowire

Autowire is reflection based dependency-injection library for Golang.

This README is in working in progress state.

## Installation

The whole project build with go modules.
To get the latest version, use go1.16+ and fetch it using the go get command. For example:

```bash
go get github.com/go-autowire/autowire
```

To get the specific version, use go1.16+ and fetch it using the go get command. For example:

```bash
go get github.com/go-autowire/autowire@v1.0.6
```

## Documentation

### Overview

Autowire is a simple golang module that automatically connects components using dependency injection. Autowire works using Go's reflection package and struct field tags. Dependencies inside the components should be annotated with "autowire" tag.

### Quick Start
