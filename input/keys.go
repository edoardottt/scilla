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

*/

package input

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

//Keys is a struct representing the format of the keys.yaml file.
type Keys struct {
	Spyse      string `yaml:"Spyse,omitempty"`
	VirusTotal string `yaml:"VirusTotal,omitempty"`
}

//ReadKeys gets as input a filename (keys.yaml) and returns a Keys object.
func ReadKeys(filename string) (Keys, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return Keys{}, err
	}

	c := Keys{}
	err = yaml.Unmarshal(buf, &c)
	if err != nil {
		return Keys{}, fmt.Errorf("in file %q: %v", filename, err)
	}

	return c, nil
}
