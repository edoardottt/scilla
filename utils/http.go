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

package utils

import (
	"fmt"
	"net/http"
	"time"
)

const (
	Seconds10 = 10 * time.Second
	Seconds30 = 30 * time.Second
)

// HTTPGet performs a GET request (HTTP)
// and returns ERROR if it's not possible,
// the status string otherwise (e.g. "200 OK").
func HTTPGet(input string) (string, error) {
	resp, err := http.Get(input)
	if err != nil {
		return "", fmt.Errorf("error while getting %s: %w", input, err)
	}

	defer resp.Body.Close()

	return resp.Status, nil
}
