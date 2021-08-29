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

package crawler

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/edoardottt/scilla/output"
	"github.com/edoardottt/scilla/utils"
	"github.com/gocolly/colly"
)

//SpawnCrawler spawn a crawler that search for
//links with this characteristic:
//- only http, https or ftp protocols allowed
func SpawnCrawler(target string, scheme string, ignore []string, dirs map[string]output.Asset, subs map[string]output.Asset, outputFile string, mutex *sync.Mutex, what string, plain bool) {
	ignoreBool := len(ignore) != 0
	c := colly.NewCollector()
	if what == "dir" {
		c = colly.NewCollector(
			colly.URLFilters(
				regexp.MustCompile("(http://|https://|ftp://)" + "(www.)?" + target + "*"),
			),
		)
	} else {
		c = colly.NewCollector(
			colly.URLFilters(
				regexp.MustCompile("(http://|https://|ftp://)" + "+." + target),
			),
		)
	}
	c.AllowURLRevisit = false
	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if e.Attr("href") != "" {
			url := utils.CleanURL(e.Attr("href"))
			if what == "dir" {
				if !output.PresentDirs(url, dirs) && url != target {

					e.Request.Visit(url)
				}
			} else {
				if !output.PresentSubs(url, subs) && url != target {

					e.Request.Visit(url)
				}
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		var status = utils.HttpGet(r.URL.String())
		if ignoreBool {
			statusArray := strings.Split(status, " ")
			statusInt, err := strconv.Atoi(statusArray[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not get response status %s\n", status)
				os.Exit(1)
			}
			if !utils.IgnoreResponse(statusInt, ignore) {
				if what == "dir" {
					output.AddDirs(r.URL.String(), status, dirs, mutex)
					output.PrintDirs(dirs, ignore, outputFile, mutex, plain)
				} else {
					output.AddSubs(r.URL.String(), status, subs, mutex)
					output.PrintSubs(subs, ignore, outputFile, mutex, plain)
				}
			}
		} else {
			if what == "dir" {
				output.AddDirs(r.URL.String(), status, dirs, mutex)
				output.PrintDirs(dirs, ignore, outputFile, mutex, plain)
			} else {
				output.AddSubs(r.URL.String(), status, subs, mutex)
				output.PrintSubs(subs, ignore, outputFile, mutex, plain)
			}
		}
	})
	c.Visit(scheme + "://" + target)
}
