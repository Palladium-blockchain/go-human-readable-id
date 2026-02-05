# Go Human Readable ID

A lightweight and flexible Go library for generating human-readable identifiers based on customizable templates.

## Overview

`go-human-readable-id` (hid) allows you to generate strings by replacing tokens in a template with values produced by generators. It comes with built-in generators for adjectives, nouns, verbs, and digits, and can be easily extended with your own custom logic.

## Quick Start

```go
package main

import (
	"fmt"
	"github.com/Palladium-blockchain/go-human-readable-id/pkg/hid"
)

func main() {
	// Generate a simple ID using default generators
	id, _ := hid.Generate("{adj}-{noun}-{digit}", hid.WithDefaultGenerators())
	fmt.Println(id) // Output example: "fluffy-cat-7"
}
```

## Installation

```bash
go get github.com/Palladium-blockchain/go-human-readable-id
```

## Usage

The library uses a template string where tokens are enclosed in curly braces `{}`.

### Using Default Generators

The library provides several built-in generators:
- `{adj}`: A random adjective
- `{noun}`: A random noun
- `{verb}`: A random verb
- `{digit}`: A random digit (0-9)
- `{2-digit}`: A random 2-digit number (10-99)
- `{3-digit}`: A random 3-digit number (100-999)

```go
id, err := hid.Generate("user-{adj}-{noun}", hid.WithDefaultGenerators())
```

### Context Support

You can also use `GenerateContext` if you need to pass a context to your generators:

```go
id, err := hid.GenerateContext(ctx, "{adj}-{noun}", hid.WithDefaultGenerators())
```

## Configuration with Options

The `Generate` and `GenerateContext` functions accept various options to customize behavior.

### Custom Generators

You can define and use your own generators:

```go
myGen := func(ctx context.Context, cfg *hid.Config) (string, error) {
    return "custom-value", nil
}

id, err := hid.Generate("{custom}-id", hid.WithGenerator("custom", myGen))
```

### Strict Mode

By default, strict mode is enabled (`true`). If a token in the template has no matching generator, or if a token is not closed, an error is returned.

You can disable it to keep unknown tokens as-is:

```go
// Returns "hello-{unknown}" instead of an error
id, err := hid.Generate("hello-{unknown}", hid.WithStrict(false))
```

### Custom Seed

You can provide a specific seed for the random number generator:

```go
id, err := hid.Generate("{adj}", hid.WithDefaultGenerators(), hid.WithSeed(12345))
```

### Combining Options

Options can be combined. Note that `WithDefaultGenerators()` will not override any custom generators you've already registered with the same name.

```go
id, err := hid.Generate("{adj}-{custom}", 
    hid.WithGenerator("custom", myGen),
    hid.WithDefaultGenerators(),
    hid.WithStrict(false),
)
```
