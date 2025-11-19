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

package opendb

import (
	"fmt"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"

	httpUtils "github.com/edoardottt/scilla/internal/http"
)

func scrape(body io.ReadCloser) []string {
	var result = make([]string, 0)

	tableIndex := 1  // Looks for the second table from DNSRepoNoc html page
	columnIndex := 0 // Looks for the first index that contains list of subdomains

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		fmt.Println(err)
	}

	var row string

	doc.Find(".table-responsive").Each(func(index int, tablehtml *goquery.Selection) {
		if index == tableIndex {
			tablehtml.Find("tr").Each(func(indextr int, rowhtml *goquery.Selection) {
				rowhtml.Find("td").Each(func(indexth int, tablecell *goquery.Selection) {
					if indexth == columnIndex {
						row = tablecell.Text()
					}
				})

				result = append(result, row)
			})
		}
	})

	return result
}

// DNSRepoNoc retrieves from the url below some known subdomains - without API Key.
func DNSRepoNocSubdomains(domain string, plain bool) []string {
	if !plain {
		fmt.Println("Pulling data from Dns Repo Noc")
	}

	client := http.Client{
		Timeout: httpUtils.Seconds30,
	}

	url := "https://dnsrepo.noc.org/?domain=" + domain

	resp, err := client.Get(url)
	if err != nil {
		return []string{}
	}

	defer resp.Body.Close()

	output := scrape(resp.Body)

	return output
}
