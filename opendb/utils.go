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

package opendb

import (
	"math/rand"
	"strings"
	"time"

	"github.com/edoardottt/scilla/utils"
)

//AppendDBSubdomains appends to the subdomains in the list
//the subdomains found with the open DBs.
func AppendDBSubdomains(dbsubs []string, urls []string) []string {
	if len(dbsubs) == 0 {
		return urls
	}
	var result []string
	dbsubs = utils.RemoveDuplicateValues(dbsubs)
	result = append(dbsubs, urls...)
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(result), func(i, j int) { result[i], result[j] = result[j], result[i] })
	return result
}

//CleanSubdomainsOk checks if the subdomain is valid
func CleanSubdomainsOk(target string, input []string) []string {
	var result []string
	for _, elem := range input {
		if strings.Contains(elem, "."+target) && elem[len(elem)-len(target):] == target {
			if strings.Contains(elem, "\n") {
				splits := strings.Split(elem, "\n")
				elem = splits[1]
			}
			result = append(result, elem)
		}
	}
	return result
}
