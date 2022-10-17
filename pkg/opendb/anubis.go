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
	"io/ioutil"
	"net/http"

	httpUtils "github.com/edoardottt/scilla/internal/http"
)

// AnubisSubdomains retrieves from the url below some known subdomains.
func AnubisSubdomains(target string, plain bool) []string {
	if !plain {
		fmt.Println("Pulling data from AnubisDB")
	}

	client := http.Client{
		Timeout: httpUtils.Seconds30,
	}

	var arr []string

	resp, err := client.Get("https://jonlu.ca/anubis/subdomains/" + target)
	if err != nil {
		return arr
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return arr
		}

		bodyString := string(bodyBytes)
		err = json.Unmarshal([]byte(bodyString), &arr)

		if err != nil {
			return arr
		}
	}

	for index, elem := range arr {
		arr[index] = "http://" + elem
	}

	return arr
}
