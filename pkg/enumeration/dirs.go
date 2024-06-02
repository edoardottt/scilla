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

package enumeration

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	httpUtils "github.com/edoardottt/scilla/internal/http"
	ignoreUtils "github.com/edoardottt/scilla/internal/ignore"
	mathUtils "github.com/edoardottt/scilla/internal/math"
	"github.com/edoardottt/scilla/pkg/output"
)

// AsyncDir performs concurrent requests to the specified
// urls and prints the results.
func AsyncDir(urls []string, ignore []string, outputFileJSON, outputFileHTML, outputFileTXT string,
	dirs map[string]output.Asset, mutex *sync.Mutex, plain bool, redirect bool, ua string, rua bool) {
	ignoreBool := len(ignore) != 0
	total := len(urls)
	client := http.Client{}

	if !redirect {
		client = http.Client{
			Timeout: httpUtils.Seconds10,
		}
	} else {
		client = http.Client{
			Timeout: httpUtils.Seconds10,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	}

	channels := 30
	limiter := make(chan string, channels) // Limits simultaneous requests
	waitgroup := sync.WaitGroup{}          // Needed to not prematurely exit before all requests have been finished

	var count int

	for _, domain := range urls {
		limiter <- domain

		waitgroup.Add(1)

		if !plain {
			fmt.Fprint(os.Stdout, "\r")
		}

		output.PrintDirs(dirs, ignore, outputFileJSON, outputFileHTML, outputFileTXT, mutex, plain)

		go func(domain string) {
			defer waitgroup.Done()
			defer func() { <-limiter }()

			req, err := http.NewRequest("GET", domain, nil)
			if err != nil {
				return
			}

			if ua != "Go http/Client" {
				req.Header.Set("User-Agent", ua)
			}

			if rua {
				req.Header.Set("User-Agent", httpUtils.GenerateRandomUserAgent())
			}

			resp, err := client.Do(req)
			if err != nil {
				return
			}

			count++

			if ignoreBool {
				if ignoreUtils.IgnoreResponse(resp.StatusCode, ignore) {
					return
				}
			}

			output.AddDirs(domain, resp.Status, dirs, mutex)
			resp.Body.Close()
		}(domain)

		if !plain { // update counter
			fmt.Fprint(os.Stdout, "\r")
			fmt.Printf("%0.2f%% : %d / %d", mathUtils.Percentage(count, total), count, total)
		}
	}

	waitgroup.Wait()
	output.PrintDirs(dirs, ignore, outputFileJSON, outputFileHTML, outputFileTXT, mutex, plain)

	if !plain {
		fmt.Fprint(os.Stdout, "\r")
	}
}
