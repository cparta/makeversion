![build](https://github.com/cparta/makeversion/actions/workflows/go.yml/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/cparta/makeversion/badge.svg?branch=main)](https://coveralls.io/github/cparta/makeversion?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/cparta/makeversion)](https://goreportcard.com/report/github.com/cparta/makeversion)

# makeversion
Create a project version string from Git tags and build counters.

```go
//go:generate go run github.com/cparta/makeversion/v2/cmd/mkver@latest -name packagename -out version.gen.go
```
