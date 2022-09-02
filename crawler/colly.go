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

const (
	directory = "dir"
)

// SpawnCrawler spawn a crawler that search for
// links with this characteristic:
// - only http, https or ftp protocols allowed.
func SpawnCrawler(target string, scheme string, ignore []string, dirs map[string]output.Asset,
	subs map[string]output.Asset, outputFileJSON, outputFileHTML, outputFileTXT string,
	mutex *sync.Mutex, what string, plain bool) {
	ignoreBool := len(ignore) != 0
	//nolint:staticcheck // SA4006 ignore this!
	collector := colly.NewCollector()

	if what == directory {
		collector = colly.NewCollector(
			colly.URLFilters(
				regexp.MustCompile("(http://|https://|ftp://)" + "(www.)?" + target + "*"),
			),
		)
	} else {
		collector = colly.NewCollector()
		targetTemp := "." + utils.GetRootHost(target)
		targetTemp = strings.ReplaceAll(targetTemp, ".", "\\.")
		targetRegex := "([-a-z0-9.]*)" + targetTemp + "([-a-z0-9.]*)"
		collector.URLFilters = []*regexp.Regexp{regexp.MustCompile(targetRegex)}
	}

	collector.IgnoreRobotsTxt = true
	collector.AllowURLRevisit = false

	// Find and visit all links
	collector.OnHTML("a[href]", func(element *colly.HTMLElement) {
		link := element.Attr("href")
		if link != "" {
			url := utils.CleanURL(element.Request.AbsoluteURL(link))
			if what == directory {
				if !output.PresentDirs(url, dirs, mutex) && url != target {
					//nolint // errcheck ignore this!
					element.Request.Visit(url)
				}
			} else {
				newDomain := utils.RetrieveHost(url)
				if !output.PresentSubs(newDomain, subs, mutex) && newDomain != target {
					//nolint // errcheck ignore this!
					element.Request.Visit(url)
				}
			}
		}
	})

	// On every script element which has src attribute call callback
	collector.OnHTML("script[src]", func(element *colly.HTMLElement) {
		link := element.Attr("src")
		if len(link) != 0 {
			url := utils.CleanURL(element.Request.AbsoluteURL(link))
			if what == directory {
				if !output.PresentDirs(url, dirs, mutex) && url != target {
					//nolint // errcheck ignore this!
					element.Request.Visit(url)
				}
			} else {
				newDomain := utils.RetrieveHost(url)
				if !output.PresentSubs(newDomain, subs, mutex) && newDomain != target {
					//nolint // errcheck ignore this!
					element.Request.Visit(url)
				}
			}
		}
	})

	// On every link element which has href attribute call callback
	collector.OnHTML("link[href]", func(element *colly.HTMLElement) {
		link := element.Attr("href")
		if len(link) != 0 {
			url := utils.CleanURL(element.Request.AbsoluteURL(link))
			if what == directory {
				if !output.PresentDirs(url, dirs, mutex) && url != target {
					//nolint // errcheck ignore this!
					element.Request.Visit(url)
				}
			} else {
				newDomain := utils.RetrieveHost(url)
				if !output.PresentSubs(newDomain, subs, mutex) && newDomain != target {
					//nolint // errcheck ignore this!
					element.Request.Visit(url)
				}
			}
		}
	})

	// On every iframe element which has src attribute call callback
	collector.OnHTML("iframe[src]", func(element *colly.HTMLElement) {
		link := element.Attr("src")
		if len(link) != 0 {
			url := utils.CleanURL(element.Request.AbsoluteURL(link))
			if what == directory {
				if !output.PresentDirs(url, dirs, mutex) && url != target {
					//nolint // errcheck ignore this!
					element.Request.Visit(url)
				}
			} else {
				newDomain := utils.RetrieveHost(url)
				if !output.PresentSubs(newDomain, subs, mutex) && newDomain != target {
					//nolint // errcheck ignore this!
					element.Request.Visit(url)
				}
			}
		}
	})

	collector.OnRequest(func(r *colly.Request) {
		status, err := utils.HTTPGet(r.URL.String())
		if err == nil {
			if ignoreBool {
				statusArray := strings.Split(status, " ")
				statusInt, err := strconv.Atoi(statusArray[0])
				if err != nil {
					fmt.Fprintf(os.Stderr, "Could not get response status %s\n", status)
					os.Exit(1)
				}
				if !utils.IgnoreResponse(statusInt, ignore) {
					if what == directory {
						output.AddDirs(r.URL.String(), status, dirs, mutex)
						output.PrintDirs(dirs, ignore, outputFileJSON, outputFileHTML, outputFileTXT, mutex, plain)
					} else {
						newDomain := utils.RetrieveHost(r.URL.String())
						output.AddSubs(newDomain, status, subs, mutex)
						output.PrintSubs(subs, ignore, outputFileJSON, outputFileHTML, outputFileTXT, mutex, plain)
					}
				}
			} else {
				if what == directory {
					output.AddDirs(r.URL.String(), status, dirs, mutex)
					output.PrintDirs(dirs, ignore, outputFileJSON, outputFileHTML, outputFileTXT, mutex, plain)
				} else {
					newDomain := utils.RetrieveHost(r.URL.String())
					output.AddSubs(newDomain, status, subs, mutex)
					output.PrintSubs(subs, ignore, outputFileJSON, outputFileHTML, outputFileTXT, mutex, plain)
				}
			}
		}
	})
	//nolint // ignore err
	collector.Visit(scheme + "://" + target)
}
