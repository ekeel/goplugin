# goplugin

*A simple RPC based framework for calling "plugins" (RPC servers).*

- [goplugin](#goplugin)
  * [Reference](#reference)
    + [Load](#load)
    + [LoadFromDirectory](#loadfromdirectory)
    + [NewPlugin](#newplugin)
    + [Plugin.ToString](#plugintostring)
    + [Plugin.StartServer](#pluginstartserver)
    + [Plugin.Dial](#plugindial)
    + [Plugin.Call](#plugincall)
    + [Plugin.Invoke](#plugininvoke)
  * [Load Complete Example](#load-complete-example)
  * [LoadFromDirectory Complete Example](#loadfromdirectory-complete-example)
  * [RPC Server Executable Plugin](#rpc-server-executable-plugin)

## Reference

### Load

*Load generates a Plugin object for a single executable.*

#### Parameters

> pluginFile : string    
> &nbsp;&nbsp;The path to the executable plugin file.

#### Returns

> plugin : *Plugin  
> &nbsp;&nbsp;A plugin reference.

> err : error  
> &nbsp;&nbsp;An error object.

#### Example

```go
plug, err := goplugin.Load("C:/tmp/testplugin.exe")
if err != nil {
  log.Fatal(err)
}
```

### LoadFromDirectory

*LoadFromDirectory finds all executables in a directory and generates Plugins for each one.*

#### Parameters

> directory : string  
> &nbsp;&nbsp;The path to the directory containing the plugins.

#### Returns

> plugins : map[string]*Plugin  
> &nbsp;&nbsp;A map of plugin references accessible by plugin name.

> err : error  
> &nbsp;&nbsp;An error object.

#### Example

```go
plugins, err := goplugin.LoadFromDirectory("C:/tmp/plugins")
if err != nil {
  log.Fatal(err)
}
```

### NewPlugin

*NewPlugin returns a reference to a new Plugin object.*

#### Parameters

> host : string  
> &nbsp;&nbsp;The host to execute the Plugin on (defaults to "localhost").

> name : string  
> &nbsp;&nbsp;The name of the plugin.

> pluginFile : string  
> &nbsp;&nbsp;The path to the plugins executable file.

> port : string  
> &nbsp;&nbsp;The port number the RPC server should listen on.

#### Returns

> *Plugin  
> &nbsp;&nbsp;A reference to the Plugin object.

#### Example

```go
plugin = NewPlugin(
  host,
  name,
  pluginFile,
  fmt.Sprintf("%v", port),
)
```

### Plugin.ToString

*ToString returns the string representation of the Plugin object.*

#### Returns

> error  
> &nbsp;&nbsp;An error object.

#### Example

```go
plug, _ := goplugin.Load("C:/tmp/testplugin.exe")

log.Print(plug.ToString())
```

### Plugin.StartServer

*StartServer starts the related RPC plugin server.*

#### Returns

> error  
> &nbsp;&nbsp;An error object.

#### Example

```go
plug, _ := goplugin.Load("C:/tmp/testplugin.exe")

_ = plug.StartServer()
```

### Plugin.Dial

*Dial connects the RPC client to the running server (running plugin executable).*

#### Returns

> error  
> &nbsp;&nbsp;An error object.

#### Example

```go
plug, _ := goplugin.Load("C:/tmp/testplugin.exe")

_ = plug.Dial()
```

### Plugin.Call

*Dial connects the RPC client to the running server (running plugin executable).*

#### Returns

> response : string  
> &nbsp;&nbsp;The string, preferably JSON, representation of the response.

> error  
> &nbsp;&nbsp;An error object.

#### Example

```go
response, err := plug.Call("Listener.ExecuteV2", "{\"woot\": \"woot\"}")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(response)
```

### Plugin.Invoke

*Invoke starts the plugin server, dials the plugin server, and calls a plugin server function.*

#### Returns

> response : string  
> &nbsp;&nbsp;The string, preferably JSON, representation of the response.

> error  
> &nbsp;&nbsp;An error object.

#### Example

```go
response, err := plug.Invoke("Listener.ExecuteV2", "{\"woot\": \"woot\"}")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(response)
```

## Load Complete Example

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

## LoadFromDirectory Complete Example

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

## RPC Server Executable Plugin

A plugin is a go RPC server application that registers execution functions against a "Listener" and registers those functions with the RPC server.

### RPC Server Executable Plugin Example

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