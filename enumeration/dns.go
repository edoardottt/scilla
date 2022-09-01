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

// LookupDNS prints the DNS informations for the inputted domain
func LookupDNS(domain string, outputFileJson, outputFileHtml, outputFileTxt string, plain bool) {
	if outputFileHtml != "" {
		output.HeaderHTML("DNS ENUMERATION", outputFileHtml)
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

		if outputFileJson != "" {
			output.AppendWhere(ip.String(), "", "DNS", "A", "json", outputFileJson)
		}

		if outputFileHtml != "" {
			output.AppendWhere(ip.String(), "", "DNS", "A", "html", outputFileHtml)
		}

		if outputFileTxt != "" {
			output.AppendWhere(ip.String(), "", "DNS", "A", "txt", outputFileTxt)
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

	if outputFileJson != "" {
		output.AppendWhere(cname, "", "DNS", "CNAME", "json", outputFileJson)
	}

	if outputFileHtml != "" {
		output.AppendWhere(cname, "", "DNS", "CNAME", "html", outputFileHtml)
	}

	if outputFileTxt != "" {
		output.AppendWhere(cname, "", "DNS", "CNAME", "txt", outputFileTxt)
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

		if outputFileJson != "" {
			output.AppendWhere(ns.Host, "", "DNS", "NS", "json", outputFileJson)
		}

		if outputFileHtml != "" {
			output.AppendWhere(ns.Host, "", "DNS", "NS", "html", outputFileHtml)
		}

		if outputFileTxt != "" {
			output.AppendWhere(ns.Host, "", "DNS", "NS", "txt", outputFileTxt)
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

		if outputFileJson != "" {
			output.AppendWhere(mx.Host, "", "DNS", "MX", "json", outputFileJson)
		}

		if outputFileHtml != "" {
			output.AppendWhere(mx.Host, "", "DNS", "MX", "html", outputFileHtml)
		}

		if outputFileTxt != "" {
			output.AppendWhere(mx.Host, "", "DNS", "MX", "txt", outputFileTxt)
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

		if outputFileJson != "" {
			output.AppendWhere(srv.Target, "", "DNS", "SRV", "json", outputFileJson)
		}

		if outputFileHtml != "" {
			output.AppendWhere(srv.Target, "", "DNS", "SRV", "html", outputFileHtml)
		}

		if outputFileTxt != "" {
			output.AppendWhere(srv.Target, "", "DNS", "SRV", "txt", outputFileTxt)
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

		if outputFileJson != "" {
			output.AppendWhere(txt, "", "DNS", "TXT", "json", outputFileJson)
		}

		if outputFileHtml != "" {
			output.AppendWhere(txt, "", "DNS", "TXT", "html", outputFileHtml)
		}

		if outputFileTxt != "" {
			output.AppendWhere(txt, "", "DNS", "TXT", "txt", outputFileTxt)
		}
	}

	if outputFileHtml != "" {
		output.FooterHTML(outputFileHtml)
	}

	fmt.Println()
}
