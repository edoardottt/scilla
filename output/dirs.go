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
	"fmt"
	"os"
	"sync"

	"github.com/fatih/color"
)

//PrintDirs prints the results (only the resources not already printed).
//Also performs the checks based on the response status codes.
func PrintDirs(dirs map[string]Asset, ignore []string, outputFile string, mutex *sync.Mutex, plain bool) {
	mutex.Lock()
	for domain, asset := range dirs {
		if !asset.Printed {
			dir := Asset{
				Value:   asset.Value,
				Printed: true,
			}
			dirs[domain] = dir
			var resp = asset.Value
			if !plain {
				if string(resp[0]) == "2" || string(resp[0]) == "3" {
					fmt.Fprint(os.Stdout, "\r \r")
					fmt.Printf("[+]FOUND: %s ", domain)
					color.Green("%s\n", resp)
					if outputFile != "" {
						AppendWhere(domain, fmt.Sprint(resp), "DIR", "", outputFile)
					}
				} else if (resp[:3] != "404") || string(resp[0]) == "5" {
					fmt.Fprint(os.Stdout, "\r \r")
					fmt.Printf("[+]FOUND: %s ", domain)
					color.Red("%s\n", resp)

					if outputFile != "" {
						AppendWhere(domain, fmt.Sprint(resp), "DIR", "", outputFile)
					}
				}
			} else {
				if resp[:3] != "404" {
					fmt.Printf("%s\n", domain)
					if outputFile != "" {
						AppendWhere(domain, fmt.Sprint(resp), "DIR", "", outputFile)
					}
				}
			}
		}
	}
	mutex.Unlock()
}

//AddDirs adds the target found to the dirs map
func AddDirs(target string, value string, dirs map[string]Asset, mutex *sync.Mutex) {
	dir := Asset{
		Value:   value,
		Printed: false,
	}
	if !PresentDirs(target, dirs, mutex) {
		mutex.Lock()
		dirs[target] = dir
		mutex.Unlock()
	}
}

//PresentDirs checks if a directory is present inside the dirs map
func PresentDirs(input string, dirs map[string]Asset, mutex *sync.Mutex) bool {
	mutex.Lock()
	_, ok := dirs[input]
	mutex.Unlock()
	return ok
}
