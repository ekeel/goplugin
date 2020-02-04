package goplugin

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/phayes/freeport"
)

// Load generates a Plugin object for a single executable.
func Load(pluginFile string) (plugin *Plugin, err error) {
	host := "localhost"
	name := (strings.Split(filepath.Base(pluginFile), "."))[0]
	port, err := freeport.GetFreePort()
	if err != nil {
		return plugin, err
	}

	plugin = NewPlugin(
		host,
		name,
		pluginFile,
		fmt.Sprintf("%v", port),
	)

	return plugin, nil
}

// LoadFromDirectory finds all executables in a directory and generates Plugins for each one.
func LoadFromDirectory(directory string) (plugins map[string]*Plugin, err error) {
	plugins = make(map[string]*Plugin)

	fileInfo, err := ioutil.ReadDir(directory)
	if err != nil {
		return plugins, err
	}

	for _, fi := range fileInfo {
		if len(fi.Name()) > 0 {
			host := "localhost"
			name := (strings.Split(fi.Name(), "."))[0]
			port, err := freeport.GetFreePort()
			if err != nil {
				return plugins, err
			}

			plugin := NewPlugin(
				host,
				name,
				path.Join(directory, fi.Name()),
				fmt.Sprintf("%v", port),
			)

			plugins[name] = plugin
		}
	}

	return plugins, nil
}
