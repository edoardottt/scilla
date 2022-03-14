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

//SpawnCrawler spawn a crawler that search for
//links with this characteristic:
//- only http, https or ftp protocols allowed
func SpawnCrawler(target string, scheme string, ignore []string, dirs map[string]output.Asset,
	subs map[string]output.Asset, outputFile string, mutex *sync.Mutex, what string, plain bool) {

	ignoreBool := len(ignore) != 0
	c := colly.NewCollector()
	if what == "dir" {
		c = colly.NewCollector(
			colly.URLFilters(
				regexp.MustCompile("(http://|https://|ftp://)" + "(www.)?" + target + "*"),
			),
		)
	} else {
		c = colly.NewCollector()
		targetTemp := "." + utils.GetRootHost(target)
		targetTemp = strings.ReplaceAll(targetTemp, ".", "\\.")
		targetRegex := "([-a-z0-9.]*)" + targetTemp + "([-a-z0-9.]*)"
		c.URLFilters = []*regexp.Regexp{regexp.MustCompile(targetRegex)}
	}
	c.AllowURLRevisit = false

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if link != "" {
			url := utils.CleanURL(e.Request.AbsoluteURL(link))
			if what == "dir" {
				if !output.PresentDirs(url, dirs, mutex) && url != target {

					e.Request.Visit(url)
				}
			} else {
				newDomain := utils.RetrieveHost(url)
				if !output.PresentSubs(newDomain, subs, mutex) && newDomain != target {

					e.Request.Visit(url)
				}
			}
		}
	})

	// On every script element which has src attribute call callback
	c.OnHTML("script[src]", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		if len(link) != 0 {
			url := utils.CleanURL(e.Request.AbsoluteURL(link))
			if what == "dir" {
				if !output.PresentDirs(url, dirs, mutex) && url != target {

					e.Request.Visit(url)
				}
			} else {
				newDomain := utils.RetrieveHost(url)
				if !output.PresentSubs(newDomain, subs, mutex) && newDomain != target {

					e.Request.Visit(url)
				}
			}
		}
	})

	// On every link element which has href attribute call callback
	c.OnHTML("link[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if len(link) != 0 {
			url := utils.CleanURL(e.Request.AbsoluteURL(link))
			if what == "dir" {
				if !output.PresentDirs(url, dirs, mutex) && url != target {

					e.Request.Visit(url)
				}
			} else {
				newDomain := utils.RetrieveHost(url)
				if !output.PresentSubs(newDomain, subs, mutex) && newDomain != target {

					e.Request.Visit(url)
				}
			}
		}
	})

	// On every iframe element which has src attribute call callback
	c.OnHTML("iframe[src]", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		if len(link) != 0 {
			url := utils.CleanURL(e.Request.AbsoluteURL(link))
			if what == "dir" {
				if !output.PresentDirs(url, dirs, mutex) && url != target {

					e.Request.Visit(url)
				}
			} else {
				newDomain := utils.RetrieveHost(url)
				if !output.PresentSubs(newDomain, subs, mutex) && newDomain != target {

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
					newDomain := utils.RetrieveHost(r.URL.String())
					output.AddSubs(newDomain, status, subs, mutex)
					output.PrintSubs(subs, ignore, outputFile, mutex, plain)
				}
			}
		} else {
			if what == "dir" {
				output.AddDirs(r.URL.String(), status, dirs, mutex)
				output.PrintDirs(dirs, ignore, outputFile, mutex, plain)
			} else {
				newDomain := utils.RetrieveHost(r.URL.String())
				output.AddSubs(newDomain, status, subs, mutex)
				output.PrintSubs(subs, ignore, outputFile, mutex, plain)
			}
		}
	})
	c.Visit(scheme + "://" + target)
}
