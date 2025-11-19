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

	@Author:      edoardottt, https://edoardottt.com

	@License: https://github.com/edoardottt/scilla/blob/main/LICENSE

*/

package output

import (
	"fmt"
	"os"
	"sync"

	urlUtils "github.com/edoardottt/scilla/internal/url"
	"github.com/fatih/color"
)

func PrintSubs(subs map[string]Asset, ignore []string, outputFileJSON, outputFileHTML, outputFileTXT string,
	mutex *sync.Mutex, plain bool) {
	mutex.Lock()

	for domain, asset := range subs {
		if !asset.Printed {
			sub := Asset{
				Value:   asset.Value,
				Printed: true,
			}
			subs[domain] = sub
			resp := asset.Value

			// Skip if 404
			if len(resp) >= 3 && resp[:3] == "404" {
				continue
			}

			subDomainFound := urlUtils.CleanProtocol(domain)

			// Save to files
			if outputFileJSON != "" {
				AppendWhere(domain, resp, "SUB", "", "json", outputFileJSON)
			}

			if outputFileHTML != "" {
				AppendWhere(domain, resp, "SUB", "", "html", outputFileHTML)
			}

			if outputFileTXT != "" {
				AppendWhere(domain, resp, "SUB", "", "txt", outputFileTXT)
			}

			// Terminal output
			PrintMutex.Lock()

			if plain {
				fmt.Printf("%s\n", subDomainFound)
			} else {
				fmt.Fprint(os.Stdout, "\r")
				fmt.Printf("%s ", subDomainFound)

				if len(resp) > 0 && resp[0] == '2' {
					color.Green("%s\n", resp)
				} else {
					color.Red("%s\n", resp)
				}
			}

			PrintMutex.Unlock()
		}
	}

	mutex.Unlock()
}

// AddSubs adds the target found to the subs map.
func AddSubs(target string, value string, subs map[string]Asset, mutex *sync.Mutex) {
	target = urlUtils.CleanProtocol(target)

	if !PresentSubs(target, subs, mutex) {
		mutex.Lock()
		defer mutex.Unlock()

		subs[target] = Asset{
			Value:   value,
			Printed: false,
		}
	}
}

// PresentSubs checks if a subdomain is present inside the subs map.
func PresentSubs(input string, subs map[string]Asset, mutex *sync.Mutex) bool {
	mutex.Lock()
	defer mutex.Unlock()

	_, ok := subs[input]

	return ok
}

// --- Helpers ---

func printColoredResponse(resp string) {
	if len(resp) > 0 && resp[0] == '2' {
		color.Green("%s\n", resp)
	} else {
		color.Red("%s\n", resp)
	}
}

func writeToAllOutputs(domain, resp, category, subType, jsonPath, htmlPath, txtPath string) {
	if jsonPath != "" {
		AppendWhere(domain, resp, category, subType, "json", jsonPath)
	}

	if htmlPath != "" {
		AppendWhere(domain, resp, category, subType, "html", htmlPath)
	}

	if txtPath != "" {
		AppendWhere(domain, resp, category, subType, "txt", txtPath)
	}
}
