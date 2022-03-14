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

package enumeration

import (
	"fmt"
	"net"
	"os"

	"github.com/edoardottt/scilla/output"
	"github.com/fatih/color"
)

//LookupDNS prints the DNS informations for the inputted domain
func LookupDNS(domain string, outputFile string, plain bool) {
	if outputFile != "" {
		if outputFile[len(outputFile)-4:] == "html" {
			output.HeaderHTML("DNS SCANNING", outputFile)
		}
	}
	// -- A RECORDS --
	ips, err := net.LookupIP(domain)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
	}
	for _, ip := range ips {
		if !plain {
			fmt.Printf("[+]FOUND %s IN A: ", domain)
			color.Green("%s\n", ip.String())
		} else {
			fmt.Printf("%s\n", ip.String())
		}
		if outputFile != "" {
			output.AppendWhere(ip.String(), "", "DNS", "A", outputFile)
		}
	}
	// -- CNAME RECORD --
	cname, err := net.LookupCNAME(domain)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get CNAME: %v\n", err)
	}
	if !plain {
		fmt.Printf("[+]FOUND %s IN CNAME: ", domain)
		color.Green("%s\n", cname)
	} else {
		fmt.Printf("%s\n", cname)
	}
	if outputFile != "" {
		output.AppendWhere(cname, "", "DNS", "CNAME", outputFile)
	}
	// -- NS RECORDS --
	nameserver, err := net.LookupNS(domain)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get NSs: %v\n", err)
	}
	for _, ns := range nameserver {
		if !plain {
			fmt.Printf("[+]FOUND %s IN NS: ", domain)
			color.Green("%s\n", ns.Host)
		} else {
			fmt.Printf("%s\n", ns.Host)
		}
		if outputFile != "" {
			output.AppendWhere(ns.Host, "", "DNS", "NS", outputFile)
		}
	}
	// -- MX RECORDS --
	mxrecords, err := net.LookupMX(domain)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get MXs: %v\n", err)
	}
	for _, mx := range mxrecords {
		if !plain {
			fmt.Printf("[+]FOUND %s IN MX: ", domain)
			color.Green("%s %v\n", mx.Host, mx.Pref)
		} else {
			fmt.Printf("%s %v\n", mx.Host, mx.Pref)
		}
		if outputFile != "" {
			output.AppendWhere(mx.Host, "", "DNS", "MX", outputFile)
		}
	}
	// -- SRV SERVICE --
	_, srvs, err := net.LookupSRV("xmpp-server", "tcp", domain)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get SRVs: %v\n", err)
	}
	for _, srv := range srvs {
		if !plain {
			fmt.Printf("[+]FOUND %s IN SRV: ", domain)
			color.Green("%v:%v:%d:%d\n", srv.Target, srv.Port, srv.Priority, srv.Weight)
		} else {
			fmt.Printf("%v:%v:%d:%d\n", srv.Target, srv.Port, srv.Priority, srv.Weight)
		}
		if outputFile != "" {
			output.AppendWhere(srv.Target, "", "DNS", "SRV", outputFile)
		}
	}
	// -- TXT RECORDS --
	txtrecords, _ := net.LookupTXT(domain)
	for _, txt := range txtrecords {
		if !plain {
			fmt.Printf("[+]FOUND %s IN TXT: ", domain)
			color.Green("%s\n", txt)
		} else {
			fmt.Printf("%s\n", txt)
		}
		if outputFile != "" {
			output.AppendWhere(txt, "", "DNS", "TXT", outputFile)
		}
	}
	if outputFile != "" {
		if outputFile[len(outputFile)-4:] == "html" {
			output.FooterHTML(outputFile)
		}
	}
}
