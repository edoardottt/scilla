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
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/edoardottt/scilla/utils"
)

// CrtShResult is the struct containing crt.sh results
type CrtShResult struct {
	Name string `json:"name_value"`
}

// CrtshSubdomains retrieves from the url below some known subdomains.
func CrtshSubdomains(domain string) []string {
	client := http.Client{
		Timeout: utils.Seconds30,
	}

	var results []CrtShResult

	url := "https://crt.sh/?q=%25." + domain + "&output=json"
	resp, err := client.Get(url)

	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()

	output := make([]string, 0)
	body, _ := ioutil.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &results); err != nil {
		return []string{}
	}

	for _, res := range results {
		out := strings.ReplaceAll(res.Name, "{", "")
		out = strings.ReplaceAll(out, "}", "")
		output = append(output, out)
	}

	return output
}
