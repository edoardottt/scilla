/*

=======================
Scilla - Information Gathering Tool
=======================

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see http://www.gnu.org/licenses/.

	@Repository:  https://github.com/edoardottt/scilla

	@Author:      edoardottt, https://www.edoardoottavianelli.it

	@License: https://github.com/edoardottt/scilla/blob/main/LICENSE

*/

package input

import (
	"fmt"
	"os"
	"runtime"

	"gopkg.in/yaml.v3"
)

// Keys is a struct representing the format of the keys.yaml file.
type Keys struct {
	VirusTotal string `yaml:"VirusTotal,omitempty"`
	BuiltWith  string `yaml:"BuiltWith,omitempty"`
}

// ReadKeys gets as input a filename (keys.yaml) and returns a Keys object.
func ReadKeys(filename string) (Keys, error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return Keys{}, fmt.Errorf("error while reading file %s: %w", filename, err)
	}

	keys := Keys{}
	err = yaml.Unmarshal(buf, &keys)

	if err != nil {
		return Keys{}, fmt.Errorf("in file %q: %w", filename, err)
	}

	return keys, nil
}

// GetKey reads the key for the given api.
func GetKey(api string) string {
	filename := ""
	key := ""

	if runtime.GOOS == windows {
		filename = "keys.yaml"
	} else { // linux
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Cannot read api key.")
			os.Exit(1)
		}

		filename = home + "/.config/scilla/keys.yaml"
	}

	keys, err := ReadKeys(filename)

	if err != nil {
		fmt.Println("Cannot read api key.")
		os.Exit(1)
	}

	switch api {
	case "virustotal":
		if keys.VirusTotal == "" {
			fmt.Println("VirusTotal Api Key is empty.")
			os.Exit(1)
		}

		key = keys.VirusTotal
	case "builtwith":
		if keys.BuiltWith == "" {
			fmt.Println("BuiltWith Api Key is empty.")
			os.Exit(1)
		}

		key = keys.BuiltWith
	}

	return key
}
