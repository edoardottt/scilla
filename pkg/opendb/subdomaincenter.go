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

package opendb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	httpUtils "github.com/edoardottt/scilla/internal/http"
)

// SubdomainCenter retrieves from the url below some known subdomains.
func SubdomainCenterSubdomains(domain string, plain bool) []string {
	if !plain {
		fmt.Println("Pulling data from Subdomain Center")
	}

	client := http.Client{
		Timeout: httpUtils.Seconds30,
	}

	result := make([]string, 0)
	url := "http://api.subdomain.center/?domain=" + domain

	resp, err := client.Get(url)
	if err != nil {
		return result
	}
	defer resp.Body.Close()

	// read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result
	}

	// Decode the response body as list of string
	var response []string

	err = json.Unmarshal(body, &response)
	if err != nil {
		return result
	}

	result = append(result, response...)

	return result
}
