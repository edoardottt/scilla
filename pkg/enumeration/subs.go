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
	"net"
	"net/http"
	"os"
	"sync"

	dnsUtils "github.com/edoardottt/scilla/internal/dns"
	httpUtils "github.com/edoardottt/scilla/internal/http"
	ignoreUtils "github.com/edoardottt/scilla/internal/ignore"
	mathUtils "github.com/edoardottt/scilla/internal/math"
	"github.com/edoardottt/scilla/pkg/output"
)

// AsyncGet performs concurrent requests to the specified
// urls and prints the results.
func AsyncGet(protocol string, urls []string, ignore []string, outputFileJSON, outputFileHTML, outputFileTXT string,
	subs map[string]output.Asset, mutex *sync.Mutex, plain bool, ua string, rua bool, alive bool, custom string) {
	ignoreBool := len(ignore) != 0

	var count int

	var total = len(urls)

	client := http.Client{
		Timeout: httpUtils.Seconds10,
	}

	var r *net.Resolver
	if custom != "" {
		r = dnsUtils.NewCustomResolver(custom)
	}

	channels := 10
	limiter := make(chan string, channels) // Limits simultaneous requests
	waitgroup := sync.WaitGroup{}          // Needed to not prematurely exit before all requests have been finished

	for _, domain := range urls {
		limiter <- domain

		waitgroup.Add(1)

		if count%50 == 0 { // update counter
			if !plain {
				fmt.Fprint(os.Stdout, "\r \r")
			}

			output.PrintSubs(subs, ignore, outputFileJSON, outputFileHTML, outputFileTXT, mutex, plain)
		}

		if !plain && count%10 == 0 { // update counter
			fmt.Fprint(os.Stdout, "\r \r")
			fmt.Printf("%0.2f%% : %d / %d", mathUtils.Percentage(count, total), count, total)
		}

		go func(domain string) {
			defer waitgroup.Done()
			defer func() { <-limiter }()

			count++

			if !alive {
				found := false
				if custom != "" {
					found = dnsUtils.CustomDNSLookup(r, domain)
				} else {
					found = dnsUtils.SimpleDNSLookup(domain)
				}

				if found {
					output.AddSubs(domain, "", subs, mutex)
				}
			} else {
				req, err := http.NewRequest("GET", protocol+"://"+domain, nil)
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

				if ignoreBool {
					if ignoreUtils.IgnoreResponse(resp.StatusCode, ignore) {
						return
					}
				}

				output.AddSubs(domain, resp.Status, subs, mutex)
				resp.Body.Close()
			}
		}(domain)
	}

	waitgroup.Wait()
	output.PrintSubs(subs, ignore, outputFileJSON, outputFileHTML, outputFileTXT, mutex, plain)

	if !plain {
		fmt.Fprint(os.Stdout, "\r \r")
		fmt.Println()
	}
}
