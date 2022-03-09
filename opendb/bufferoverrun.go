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

package opendb

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

//BufferOverrunSubdomains retrieves from the url below some known subdomains.
func BufferOverrunSubdomains(domain string) []string {
	client := http.Client{
		Timeout: 30 * time.Second,
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

	dec.Decode(&wrapper)
	if err != nil {
		return result
	}
	for _, r := range wrapper.Records {
		parts := strings.SplitN(r, ",", 2)
		if len(parts) != 2 {
			continue
		}
		result = append(result, parts[1])
	}
	return result
}
