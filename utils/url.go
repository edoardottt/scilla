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

package utils

import (
	"net/url"
	"strings"
)

//ProtocolExists checks if the protocol is present in the URL
func ProtocolExists(target string) bool {
	res := strings.Index(target, "://")
	return res >= 0
}

//CleanProtocol remove from the url the protocol scheme
func CleanProtocol(target string) string {
	res := strings.Index(target, "://")
	if res >= 0 {
		return target[res+3:]
	} else {
		return target
	}
}

//CleanURL takes as input a string and it tries to
//remove the fragment and the query
//Example: https://example.com/directory1/?id=abcdef&path=ok#fragment1
//Output: https://example.com/directory1/
func CleanURL(input string) string {
	u, err := url.Parse(input)
	if err != nil {
		return input
	}
	return u.Scheme + "://" + u.Host + u.Path
}

//IsUrl checks if the inputted Url is valid
func IsURL(str string) bool {
	if !ProtocolExists(str) {
		str = "http://" + str
	}
	u, err := url.Parse(str)
	return err == nil && u.Host != ""
}

//BuildUrl returns full URL with the subdomain
func BuildURL(scheme string, subdomain string, domain string) string {
	return scheme + "://" + subdomain + "." + domain
}

//AppendDir returns full URL with the directory
func AppendDir(domain string, dir string) (string, string) {
	return domain + "/" + dir + "/", domain + "/" + dir
}

//CleanSubdomainsOk >
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

//RetrieveProtocol remove from the url the protocol scheme
func RetrieveProtocol(target string) string {
	res := strings.Index(target, "://")
	if res >= 0 {
		return target[:res]
	} else {
		return target
	}
}
