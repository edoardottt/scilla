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

type OutputFile struct {
	port      []string `json: "dns,omitempty"`
	dns       []string `json: "dns,omitempty"`
	subdomain []string `json: "subdomain,omitempty"`
	dir       []string `json: "dir,omitempty"`
}

//AppendOutputToJSON >
func AppendOutputToJSON(output string, key string, filename string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
	}
	data := OutputFile{}

	_ = json.Unmarshal([]byte(file), &data)
	if key == "PORT" {
		data.port = append(data.port, output)
	}

	file, err = json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println(err)
	}
	ioutil.WriteFile("filename.json", file, 0644)
	log.Println("Written")

}
