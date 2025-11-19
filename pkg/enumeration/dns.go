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

package enumeration

import (
	"fmt"
	"net"

	"github.com/edoardottt/scilla/pkg/output"
	"github.com/fatih/color"
)

// LookupDNS prints the DNS informations for the inputted domain.
func LookupDNS(domain string, outputFileJSON, outputFileHTML, outputFileTXT string, plain bool) {
	if outputFileHTML != "" {
		output.HeaderHTML("DNS ENUMERATION", outputFileHTML)
	}
	// -- A RECORDS --
	ips, _ := net.LookupIP(domain)

	for _, ip := range ips {
		if !plain {
			fmt.Print("A: ")
			color.Green("%s\n", ip.String())
		} else {
			fmt.Printf("%s\n", ip.String())
		}

		appendDNSOutputHelper(outputFileJSON, outputFileHTML, outputFileTXT, ip.String(), "A")
	}
	// -- CNAME RECORD --
	cname, _ := net.LookupCNAME(domain)

	if !plain {
		fmt.Print("CNAME: ")
		color.Green("%s\n", cname)
	} else {
		fmt.Printf("%s\n", cname)
	}

	appendDNSOutputHelper(outputFileJSON, outputFileHTML, outputFileTXT, cname, "CNAME")

	// -- NS RECORDS --
	nameserver, _ := net.LookupNS(domain)

	for _, nsRecord := range nameserver {
		if !plain {
			fmt.Print("NS: ")
			color.Green("%s\n", nsRecord.Host)
		} else {
			fmt.Printf("%s\n", nsRecord.Host)
		}

		appendDNSOutputHelper(outputFileJSON, outputFileHTML, outputFileTXT, nsRecord.Host, "NS")
	}

	// -- MX RECORDS --
	mxrecords, _ := net.LookupMX(domain)

	for _, mxRecord := range mxrecords {
		if !plain {
			fmt.Print("MX: ")
			color.Green("%s %v\n", mxRecord.Host, mxRecord.Pref)
		} else {
			fmt.Printf("%s %v\n", mxRecord.Host, mxRecord.Pref)
		}

		appendDNSOutputHelper(outputFileJSON, outputFileHTML, outputFileTXT, mxRecord.Host, "MX")
	}

	// -- SRV SERVICE --
	services := []string{"ldap", "xmpp", "smpp-server", "xmpp-client"}
	srvResults := []*net.SRV{}

	for _, service := range services {
		_, srvs, _ := net.LookupSRV(service, "tcp", domain)
		srvResults = append(srvResults, srvs...)
	}

	for _, srv := range srvResults {
		if !plain {
			fmt.Print("SRV: ")
			color.Green("%v:%v:%d:%d\n", srv.Target, srv.Port, srv.Priority, srv.Weight)
		} else {
			fmt.Printf("%v:%v:%d:%d\n", srv.Target, srv.Port, srv.Priority, srv.Weight)
		}

		appendDNSOutputHelper(outputFileJSON, outputFileHTML, outputFileTXT, srv.Target, "SRV")
	}

	// -- TXT RECORDS --
	txtrecords, _ := net.LookupTXT(domain)
	for _, txt := range txtrecords {
		if !plain {
			fmt.Print("TXT: ")
			color.Green("%s\n", txt)
		} else {
			fmt.Printf("%s\n", txt)
		}

		appendDNSOutputHelper(outputFileJSON, outputFileHTML, outputFileTXT, txt, "TXT")
	}

	if outputFileHTML != "" {
		output.FooterHTML(outputFileHTML)
	}

	fmt.Println()
}

func appendDNSOutputHelper(outputFileJSON, outputFileHTML, outputFileTXT, data, rtype string) {
	if outputFileJSON != "" {
		output.AppendWhere(data, "", "DNS", rtype, "json", outputFileJSON)
	}

	if outputFileHTML != "" {
		output.AppendWhere(data, "", "DNS", rtype, "html", outputFileHTML)
	}

	if outputFileTXT != "" {
		output.AppendWhere(data, "", "DNS", rtype, "txt", outputFileTXT)
	}
}
