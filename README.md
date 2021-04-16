# Core

## Getting Started

Initialize project with [hotshot](https://github.com/blackdreamers/hotshot)

```
hotshot new helloworld
```

## Example Service

env
```shell
ENV=dev
#ENV=prod
LOG_LEVEL=debug
ETCD_ADDRS=localhost:2379
ETCD_AUTH=false
ETCD_TLS=false
ETCD_USER=xxx
ETCD_PASSWORD=xxx
ETCD_CA_PATH=xxx
ETCD_CERT_PATH=xxx
ETCD_CERT_KEY_PATH=xxx
```

```go
package main

import (
	"github.com/blackdreamers/core/server"
	_ "github.com/blackdreamers/helloworld/db"
	_ "github.com/blackdreamers/helloworld/handler"
	_ "github.com/blackdreamers/helloworld/subscriber"
)

func main() {
	// Init server
	server.Init(
		server.Name("helloworld"),
		server.Type(server.SRV),
	)

	// Run server
	server.Run()
}

```
