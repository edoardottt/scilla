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

package opendb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	httpUtils "github.com/edoardottt/scilla/internal/http"
)

// BufferOverrunSubdomains retrieves from the url below some known subdomains.
func BufferOverrunSubdomains(domain string, plain bool) []string {
	if !plain {
		fmt.Println("Pulling data from dns.bufferover.run")
	}

	client := http.Client{
		Timeout: httpUtils.Seconds30,
	}
	result := make([]string, 0)
	url := "https://dns.bufferover.run/dns?q=" + domain
	wrapper := struct {
		Records []string `json:"FDNS_A"`
	}{}
	resp, err := client.Get(url)

	if err != nil {
		return result
	}

	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)

	err = dec.Decode(&wrapper)
	if err != nil {
		return result
	}

	for _, r := range wrapper.Records {
		parts := strings.SplitN(r, ",", twoParts)
		if len(parts) != twoParts {
			continue
		}
		result = append(result, parts[1])
	}

	return result
}
