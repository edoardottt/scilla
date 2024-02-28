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
	"io"
	"net/http"

	httpUtils "github.com/edoardottt/scilla/internal/http"
)

// ThreatMinerResult is the struct containing ThreatMiner results.
type ThreatMinerResult struct {
	StatusCode    string   `json:"status_code"`
	StatusMessage string   `json:"status_message"`
	Results       []string `json:"results"`
}

// ThreatMinerSubdomains retrieves from the url below some known subdomains.
func ThreatMinerSubdomains(domain string, plain bool) []string {
	if !plain {
		fmt.Println("Pulling data from ThreatMiner.org")
	}

	client := http.Client{
		Timeout: httpUtils.Seconds30,
	}

	var result ThreatMinerResult

	url := "https://api.threatminer.org/v2/domain.php?q=" + domain + "&rt=5"

	resp, err := client.Get(url)

	if err != nil {
		return []string{}
	}

	defer resp.Body.Close()

	output := make([]string, 0)
	body, _ := io.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &result); err != nil {
		return []string{}
	}

	output = append(output, result.Results...)

	return output
}
