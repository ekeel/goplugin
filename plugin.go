package plugin

import (
	"bytes"
	"errors"
	"fmt"
	"net/rpc"
	"os/exec"
)

// Plugin is a struct containing the necessary data and methods to execute a plugin.
type Plugin struct {
	Client     *rpc.Client
	Host       string
	Name       string
	PluginFile string
	Port       string
	Meta       map[string]interface{}
}

// NewPlugin returns a reference to a new Plugin object.
func NewPlugin(host, name, pluginFile, port string) *Plugin {
	plugin := Plugin{
		Host:       host,
		Name:       name,
		Port:       port,
		PluginFile: pluginFile,
	}

	return &plugin
}

// ToString returns the string representation of the Plugin object.
func (plugin *Plugin) ToString() string {
	return fmt.Sprintf(
		"{\"Host\": \"%s\", \"Name\": \"%s\", \"PluginFile\": \"%s\", \"Port\": \"%s\"}",
		plugin.Host,
		plugin.Name,
		plugin.PluginFile,
		plugin.Port,
	)
}

// StartServer starts the related RPC plugin server.
func (plugin *Plugin) StartServer() error {
	var err error

	go func(plug *Plugin) {
		cmd := exec.Command(plug.PluginFile, plug.Port)

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		e := cmd.Run()
		if e != nil {
			err = e
		}

		if stderr.Len() > 0 {
			err = errors.New(stderr.String())
		}
	}(plugin)

	if err != nil {
		return err
	}

	return nil
}

// Dial connects the RPC client to the running server (running plugin executable).
func (plugin *Plugin) Dial() error {
	rpcClient, err := rpc.Dial("tcp", fmt.Sprintf("%s:%s", plugin.Host, plugin.Port))
	if err != nil {
		return err
	}

	plugin.Client = rpcClient

	return nil
}

// Call invokes a function in the Plugin over RPC.
func (plugin *Plugin) Call(function, payload string) (response string, err error) {
	err = plugin.Client.Call(function, payload, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

// Invoke starts the plugin server, dials the plugin server, and calls a plugin server function.
func (plugin *Plugin) Invoke(function, payload string) (response string, err error) {
	err = plugin.StartServer()
	if err != nil {
		return response, err
	}

	err = plugin.Dial()
	if err != nil {
		return response, err
	}

	err = plugin.Client.Call(function, payload, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}
