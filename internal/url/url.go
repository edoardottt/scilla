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

package utils

import (
	"net/url"
	"strings"

	"github.com/bobesa/go-domain-util/domainutil"
)

// ProtocolExists checks if the protocol is present in the URL
// taken as input.
func ProtocolExists(target string) bool {
	u, err := url.Parse(target)
	if err != nil {
		return false
	}

	return u.Scheme != ""
}

// CleanProtocol removes the protocol from the url.
func CleanProtocol(target string) string {
	if ProtocolExists(target) {
		res := strings.Index(target, "://")
		if res >= 0 {
			return target[res+3:]
		}
	}

	return target
}

// CleanURL takes as input a string and it tries to
// remove the fragment and the query.
// Example: https://example.com/directory1/?id=abcdef&path=ok#fragment1
// Output: https://example.com/directory1/
func CleanURL(input string) string {
	u, err := url.Parse(input)
	if err != nil {
		return input
	}

	if u.Scheme == "" {
		u.Scheme = "http"
	}

	return u.Scheme + "://" + u.Host + u.Path
}

// IsURL checks if the inputted url is valid.
func IsURL(str string) bool {
	if !ProtocolExists(str) {
		str = "http://" + str
	}

	u, err := url.Parse(str)

	return err == nil && u.Host != ""
}

// BuildURL returns full URL with the subdomain.
func BuildURL(scheme string, subdomain string, domain string) string {
	return scheme + "://" + subdomain + "." + domain
}

// AppendDir returns full URL with the directory.
func AppendDir(scheme string, domain string, dir string) (string, string) {
	return scheme + "://" + domain + "/" + dir + "/", scheme + "://" + domain + "/" + dir
}

// CleanSubdomainsOk takes as input a slice of subdomains and remove
// from the input slice all the 'wrong' subdomains.
func CleanSubdomainsOk(target string, input []string) []string {
	result := []string{}

	for _, elem := range input {
		if strings.Contains(elem, "."+target) && elem != "."+target &&
			elem[len(elem)-len(target):] == target {
			if strings.Contains(elem, "\n") {
				splits := strings.Split(elem, "\n")
				elem = splits[1]
			}
			result = append(result, CleanProtocol(elem))
		}
	}

	return result
}

// RetrieveProtocol remove from the url the protocol scheme.
func RetrieveProtocol(target string) string {
	u, err := url.Parse(target)
	if err != nil {
		return ""
	}

	return u.Scheme
}

// AbsoluteURL takes as input a path and returns the full
// absolute URL with protocol + host + path.
func AbsoluteURL(scheme string, target string, path string) string {
	// if the path variable starts with a scheme, it means that the
	// path is itself an absolute path.
	if ProtocolExists(path) {
		return path
	}

	if path[0] == '/' {
		return scheme + "://" + target + path
	}

	return scheme + "://" + target + "/" + path
}

// RetrieveHost takes as input a URL and returns
// as output the domain (host).
func RetrieveHost(input string) string {
	u, err := url.Parse(input)

	if err != nil {
		return input
	}

	return u.Host
}

// GetRootHost returns the root host (domain.tld).
func GetRootHost(input string) string {
	_, err := url.Parse(input)

	if err != nil {
		return ""
	}

	return domainutil.Domain(input)
}
