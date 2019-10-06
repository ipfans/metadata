# metadata

[![GoDoc](https://godoc.org/github.com/ipfans/metadata?status.svg)](https://godoc.org/github.com/ipfans/metadata)
[![Release](https://img.shields.io/github/release/ipfans/metadata.svg?style=flat-square)](https://github.com/ipfans/metadata/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/ipfans/metadata)](https://goreportcard.com/report/github.com/ipfans/metadata)

`metadata` provides interfaces to contracts with metadata informations with `context.Context`. 

## why?

`context.Context.Value()` stores data in linked list. To proformace reasons, we provide `MetaData`
struct to store more informations to context.

## HOWTO

```go
// Store to context
ctx = metadata.NewContext(ctx, md)

// Extract from context
md, ok := metadata.FromContext(ctx)
```

