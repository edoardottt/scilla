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
	"context"
	"net"
	"time"
)

const (
	DNSTimeout = 10000
)

// SimpleDNSLookup performs a DNS lookup using the default system DNS.
func SimpleDNSLookup(domain string) bool {
	ips, err := net.LookupIP(domain)
	return err == nil && len(ips) != 0 && !ips[0].IsLoopback()
}

// CustomDNSLookup performs a DNS lookup using the provided custom DNS.
func CustomDNSLookup(r *net.Resolver, domain string) bool {
	ips, err := r.LookupHost(context.Background(), domain)
	return err == nil && len(ips) != 0 && !net.ParseIP(ips[0]).IsLoopback()
}

// NewCustomResolver returns a DNS resolver using the provided DNS IP.
func NewCustomResolver(customDNS string) *net.Resolver {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(DNSTimeout),
			}
			return d.DialContext(ctx, network, customDNS+":53")
		},
	}

	return r
}
