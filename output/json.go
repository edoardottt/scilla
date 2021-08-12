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

package output

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

//OutputFile struct helping json output
type OutputFile struct {
	Port      []string            `json:"port,omitempty"`
	Dns       map[string][]string `json:"dns,omitempty"`
	Subdomain []string            `json:"subdomain,omitempty"`
	Dir       []string            `json:"dir,omitempty"`
}

//AppendOutputToJSON >
func AppendOutputToJSON(output string, key string, record string, filename string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
	}
	data := OutputFile{}

	_ = json.Unmarshal([]byte(file), &data)
	if key == "PORT" {
		data.Port = append(data.Port, output)
	} else if key == "SUB" {
		data.Subdomain = append(data.Subdomain, output)
	} else if key == "DIR" {
		data.Dir = append(data.Dir, output)
	} else if key == "DNS" {
		if data.Dns == nil {
			data.Dns = make(map[string][]string)
		}
		if _, ok := data.Dns[record]; !ok {
			data.Dns[record] = make([]string, 0)
		}
		data.Dns[record] = append(data.Dns[record], output)
	}

	file, err = json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println(err)
	} else {
		ioutil.WriteFile(filename, file, 0644)
	}

}
