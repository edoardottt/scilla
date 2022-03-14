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
	"time"

	"github.com/edoardottt/scilla/output"
	"github.com/edoardottt/scilla/utils"
)

//AsyncGet performs concurrent requests to the specified
//urls and prints the results
func AsyncGet(protocol string, urls []string, ignore []string, outputFile string, subs map[string]output.Asset, mutex *sync.Mutex, plain bool) {
	ignoreBool := len(ignore) != 0
	var count int = 0
	var total int = len(urls)
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	limiter := make(chan string, 10) // Limits simultaneous requests
	wg := sync.WaitGroup{}           // Needed to not prematurely exit before all requests have been finished
	for i, domain := range urls {
		limiter <- domain
		wg.Add(1)
		if count%50 == 0 { // update counter
			if !plain {
				fmt.Fprint(os.Stdout, "\r \r")
			}
			output.PrintSubs(subs, ignore, outputFile, mutex, plain)
		}
		if !plain && count%100 == 0 { // update counter
			fmt.Fprint(os.Stdout, "\r \r")
			fmt.Printf("%0.2f%% : %d / %d", utils.Percentage(count, total), count, total)
		}
		go func(i int, domain string) {
			defer wg.Done()
			defer func() { <-limiter }()
			resp, err := client.Get(protocol + "://" + domain)
			count++
			if err != nil {
				return
			}
			if ignoreBool {
				if utils.IgnoreResponse(resp.StatusCode, ignore) {
					return
				}
			}
			output.AddSubs(domain, resp.Status, subs, mutex)
			resp.Body.Close()
		}(i, domain)
	}
	output.PrintSubs(subs, ignore, outputFile, mutex, plain)
	wg.Wait()
	output.PrintSubs(subs, ignore, outputFile, mutex, plain)
}
