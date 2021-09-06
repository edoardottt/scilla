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
	"bufio"
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

//HackerTargetSubdomains retrieves from the url below some known subdomains.
func HackerTargetSubdomains(domain string) []string {
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	result := make([]string, 0)
	raw, err := client.Get("https://api.hackertarget.com/hostsearch/?q=" + domain)
	if err != nil {
		return result
	}
	res, err := ioutil.ReadAll(raw.Body)
	if err != nil {
		return result
	}
	raw.Body.Close()
	sc := bufio.NewScanner(bytes.NewReader(res))
	for sc.Scan() {
		parts := strings.SplitN(sc.Text(), ",", 2)
		if len(parts) != 2 {
			continue
		}
		result = append(result, parts[0])
	}
	return result
}
