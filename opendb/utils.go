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

package opendb

import (
	"math/rand"
	"strings"
	"time"

	"github.com/edoardottt/scilla/utils"
)

// AppendDBSubdomains appends to the subdomains in the list
// the subdomains found with the open DBs.
func AppendDBSubdomains(dbsubs []string, urls []string) []string {
	if len(dbsubs) == 0 {
		return urls
	}
	dbsubs = utils.RemoveDuplicateValues(dbsubs)
	dbsubs = append(dbsubs, urls...)
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(dbsubs), func(i, j int) { dbsubs[i], dbsubs[j] = dbsubs[j], dbsubs[i] })
	return dbsubs
}

// CleanSubdomainsOk checks if the subdomains found are well formatted:
// - contain ".domain.tld"
// - ".domain.tld" at the end
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
