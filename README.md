# Core

## Getting Started

Initialize project with [hotshot](https://github.com/blackdreamers/hotshot)

```
hotshot new helloworld
```

## Example Service

```go
package main

import (
	coredb "github.com/blackdreamers/core/db"
	"github.com/blackdreamers/core/server"
)

func main() {
	// DB repositories
	coredb.Repositories(&db.Example{})

	// Init server
	server.Init(
		server.Name("example"),
		server.Type(server.SRV),
	)

	// Register handles
	server.Handles(
		new(handler.Example),
	)

	// Register subscribers
	server.Subscribers(
		new(subscriber.Example),
	)

	// Run server
	server.Run()
}

```
