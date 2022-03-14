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
	"log"
	"net"
	"strings"
)

//IsIP checks if the input is a
//proper ip address (formatted in a good manner)
func IsIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

//IpToHostname tries to translate an ip
//address to a hostname.
func IpToHostname(ip string) string {
	addr, err := net.LookupAddr(ip)
	if err != nil || len(addr) == 0 {
		log.Fatalf("Failed to resolve ip address %s", ip)
	}
	return strings.TrimSuffix(addr[0], ".")
}
