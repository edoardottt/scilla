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

package utils

import (
	"fmt"
	"math/rand"
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

// genOsString generates a random OS string for a User Agent.
func genOsString() string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Operating system.
	var OsStrings = []string{
		"Macintosh; Intel Mac OS X 10_10",
		"Windows NT 10.0",
		"Windows NT 5.1",
		"Windows NT 6.1; WOW64",
		"Windows NT 6.1; Win64; x64",
		"X11; Linux x86_64",
	}

	return OsStrings[rng.Intn(len(OsStrings))]
}

// genFirefoxUA generates a random Firefox User Agent.
func genFirefoxUA() string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Firefox versions.
	var FirefoxVersions = []float32{
		58.0,
		57.0,
		56.0,
		52.0,
		48.0,
		40.0,
		35.0,
	}

	version := FirefoxVersions[rng.Intn(len(FirefoxVersions))]

	return fmt.Sprintf("Mozilla/5.0 (%s; rv:%.1f) Gecko/20100101 Firefox/%.1f", genOsString(), version, version)
}

// genChromeUA generates a random Chrome User Agent.
func genChromeUA() string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Chrome versions.
	var ChromeVersions = []string{
		"65.0.3325.146",
		"64.0.3282.0",
		"41.0.2228.0",
		"40.0.2214.93",
		"37.0.2062.124",
	}

	version := ChromeVersions[rng.Intn(len(ChromeVersions))]

	return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36",
		genOsString(), version)
}

// GenerateRandomUserAgent generates a random user agent
// (can be Chrome or Firefox).
func GenerateRandomUserAgent() string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	decision := rng.Intn(100)

	var ua string
	if decision%2 == 0 {
		ua = genChromeUA()
	} else {
		ua = genFirefoxUA()
	}

	return ua
}
