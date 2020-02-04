# goplugin

*A simple RPC based framework for calling "plugins" (RPC servers).*

## Reference

### Load

*Load generates a Plugin object for a single executable.*

#### Parameters

> pluginFile : string  
>   The path to the executable plugin file.

#### Example

```go
// Load a plugin from the executable
plug, err := goplugin.Load("C:/tmp/testplugin.exe")
if err != nil {
  log.Fatal(err)
}
```

### LoadFromDirectory

*LoadFromDirectory finds all executables in a directory and generates Plugins for each one.*

#### Parameters

> directory : string
>   LoadFromDirectory finds all executables in a directory and generates Plugins for each one.

#### Example

```go
plugins, err := goplugin.LoadFromDirectory("C:/tmp/plugins")
if err != nil {
  log.Fatal(err)
}
```

## Example

The "client" Plugin type allows for easy communication and management of RPC server plugins.

```go
package main

import (
	"fmt"
	"log"

	"github.com/ekeel/goplugin"
)

func main() {
  // Load a plugin from the executable
	plug, err := goplugin.Load("C:/tmp/testplugin.exe")
	if err != nil {
		log.Fatal(err)
	}

  // Print the string representation of the loaded plugin.
	fmt.Println(plug.ToString())

  // Execute a function from the plugin and retrieve the response.
	response, err := plug.Invoke("Listener.ExecuteV2", "{\"payload\": \"json\"}")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(response)
}
```

## Example

The "client" Plugin type allows for easy communication and management of RPC server plugins.

```go
package main

import (
	"fmt"
	"log"

	"github.com/ekeel/goplugin"
)

func main() {
  // Load a plugins from the directory
	plugins, err := goplugin.LoadFromDirectory("C:/tmp/plugins")
  if err != nil {
    log.Fatal(err)
  }

  // Print the string representation of the loaded plugin.
	fmt.Println(plugins["testPlugin"].ToString())

  // Execute a function from the plugin and retrieve the response.
	response, err := plugins["testPlugin"].Invoke("Listener.ExecuteV2", "{\"payload\": \"json\"}")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(response)
}
```

## Plugin

A plugin is a go RPC server application that registers execution functions against a "Listener" and registers those functions with the RPC server.

### Example

```go
package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
)

const (
	// NAME contains the name of the plugin.
	NAME = "testPlugin"
)

// Listener contains the execute function for the plugin.
type Listener int

// Execute is a Plugin function that can be executed by the RPC client.
func (listener *Listener) Execute(payload string, response *string) error {
	*response = payload

	return nil
}

// ExecuteV2 is another Plugin function that can be executed by the RPC client.
func (listener *Listener) ExecuteV2(payload string, response *string) error {
	*response = payload

	return nil
}

func main() {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("0.0.0.0:%s", os.Args[1]))
	if err != nil {
		log.Fatal(err)
	}

	inbound, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	listener := new(Listener)
	rpc.Register(listener)
	rpc.Accept(inbound)
}
```