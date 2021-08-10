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

	"github.com/edoardottt/scilla/utils"
	"github.com/fatih/color"
)

//PrintSubs prints the results (only the resources not already printed).
//Also performs the checks based on the response status codes.
func PrintSubs(subs map[string]Asset, ignore []string, outputFile string, mutex *sync.Mutex, plain bool) {
	mutex.Lock()
	for domain, asset := range subs {
		if !asset.Printed {
			sub := Asset{
				Value:   asset.Value,
				Printed: true,
			}
			subs[domain] = sub
			var resp = asset.Value
			if !plain {
				fmt.Fprint(os.Stdout, "\r \r")
				if resp[:3] != "404" {
					subDomainFound := utils.CleanProtocol(domain)
					fmt.Printf("[+]FOUND: %s ", subDomainFound)
					if string(resp[0]) == "2" {
						if outputFile != "" {
							AppendWhere(domain, fmt.Sprint(resp), "SUB", "", outputFile)
						}
						color.Green("%s\n", resp)
					} else {
						if outputFile != "" {
							AppendWhere(domain, fmt.Sprint(resp), "SUB", "", outputFile)
						}
						color.Red("%s\n", resp)
					}
				}
			} else {
				if resp[:3] != "404" {
					subDomainFound := utils.CleanProtocol(domain)
					fmt.Printf("%s\n", subDomainFound)
					if string(resp[0]) == "2" {
						if outputFile != "" {
							AppendWhere(domain, fmt.Sprint(resp), "SUB", "", outputFile)
						}
					} else {
						if outputFile != "" {
							AppendWhere(domain, fmt.Sprint(resp), "SUB", "", outputFile)
						}
					}
				}
			}
		}
	}
	mutex.Unlock()
}

//AddSubs adds the target found to the subs map
func AddSubs(target string, value string, subs map[string]Asset, mutex *sync.Mutex) {
	sub := Asset{
		Value:   value,
		Printed: false,
	}
	target = utils.CleanProtocol(target)
	mutex.Lock()
	if !PresentSubs(target, subs) {
		subs[target] = sub
	}
	mutex.Unlock()
}

//PresentSubs
func PresentSubs(input string, subs map[string]Asset) bool {
	_, ok := subs[input]
	return ok
}
