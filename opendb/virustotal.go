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
	"time"
)

// VirusTotalSubdomains retrieves from the url below some known subdomains.
func VirusTotalSubdomains(target string, apikey string) []string {

	var result []string
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	fetchURL := fmt.Sprintf(
		"https://www.virustotal.com/vtapi/v2/domain/report?domain=%s&apikey=%s",
		target, apikey,
	)

	wrapper := struct {
		Subdomains []string `json:"subdomains"`
	}{}

	resp, err := client.Get(fetchURL)
	if err != nil {
		return result
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)

	err = dec.Decode(&wrapper)
	if err != nil {
		return result
	}
	result = append(result, wrapper.Subdomains...)
	return result
}
