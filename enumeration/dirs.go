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

package enumeration

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/edoardottt/scilla/output"
	"github.com/edoardottt/scilla/utils"
)

//AsyncDir performs concurrent requests to the specified
//urls and prints the results
func AsyncDir(urls []string, ignore []string, outputFile string, dirs map[string]output.Asset, mutex *sync.Mutex, plain bool, redirect bool) {
	ignoreBool := len(ignore) != 0
	var count int = 0
	var total int = len(urls)
	client := http.Client{}
	if !redirect {
		client = http.Client{
			Timeout: 10 * time.Second,
		}
	} else {
		client = http.Client{
			Timeout: 10 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	}
	limiter := make(chan string, 30) // Limits simultaneous requests
	wg := sync.WaitGroup{}           // Needed to not prematurely exit before all requests have been finished
	for i, domain := range urls {
		limiter <- domain
		wg.Add(1)
		if count%50 == 0 { // update counter
			if !plain {
				fmt.Fprint(os.Stdout, "\r \r")
			}
			output.PrintDirs(dirs, ignore, outputFile, mutex, plain)
		}
		if !plain && count%100 == 0 { // update counter
			fmt.Fprint(os.Stdout, "\r \r")
			fmt.Printf("%0.2f%% : %d / %d", utils.Percentage(count, total), count, total)
		}
		go func(i int, domain string) {
			defer wg.Done()
			defer func() { <-limiter }()
			resp, err := client.Get(domain)
			count++
			if err != nil {
				return
			}
			if ignoreBool {
				if utils.IgnoreResponse(resp.StatusCode, ignore) {
					return
				}
			}
			output.AddDirs(domain, resp.Status, dirs, mutex)
			resp.Body.Close()
		}(i, domain)
	}
	output.PrintDirs(dirs, ignore, outputFile, mutex, plain)
	wg.Wait()
	output.PrintDirs(dirs, ignore, outputFile, mutex, plain)
}
