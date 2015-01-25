# goji-static

Middleware for Goji (https://github.com/zenazn/goji) to serve static content
from a directory.

It’s based on the static middleware from Martini (https://github.com/go-martini/martini).

## Usage

```
package main

import (
    “github.com/hypebeast/gojistatic”
    “github.com/zenazn/goji”
)

func main() {
    // Serve static files from “public”
    goji.Use(gojistatic.Static(“public”, gojistatic.StaticOptions{SkipLogging: true}))
    
    // Run Goji
    goji.Serve()
}
```

## Options

```
goji.Use(“public”, gojistatic.StaticOptions{
	// Prefix is the optional prefix used to serve the static directory content
	Prefix: “static”,
	// Indexfile is the name of the index file that should be served if it exists
	IndexFile: “index.html”,
	// SkipLogging will disable [Static] log messages when a static file is served
	SkipLogging: true,
	// Expires defines which user-defined function to use for producing a HTTP Expires Header
	// https://developers.google.com/speed/docs/insights/LeverageBrowserCaching
	Expires: ExpiresFunc,
})
```

## Default Options

```
var options []gojistatic.StaticOptions

goji.Use(gojistatic.Static(“public”, options))

// Is the same as the default configuration options:

goji.Use(gojistatic.Static(“public”, gojistatic.StaticOptions{
	Prefix: “”,
	IndexFile: “index.html”,
	SkipLogging: false,
	Expires: nil,
}))

```
