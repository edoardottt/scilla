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
			fmt.Printf("[+]FOUND %s IN A: ", domain)
			color.Green("%s\n", ip.String())
		} else {
			fmt.Printf("%s\n", ip.String())
		}

		if outputFileJSON != "" {
			output.AppendWhere(ip.String(), "", "DNS", "A", "json", outputFileJSON)
		}

		if outputFileHTML != "" {
			output.AppendWhere(ip.String(), "", "DNS", "A", "html", outputFileHTML)
		}

		if outputFileTXT != "" {
			output.AppendWhere(ip.String(), "", "DNS", "A", "txt", outputFileTXT)
		}
	}
	// -- CNAME RECORD --
	cname, _ := net.LookupCNAME(domain)

	if !plain {
		fmt.Printf("[+]FOUND %s IN CNAME: ", domain)
		color.Green("%s\n", cname)
	} else {
		fmt.Printf("%s\n", cname)
	}

	if outputFileJSON != "" {
		output.AppendWhere(cname, "", "DNS", "CNAME", "json", outputFileJSON)
	}

	if outputFileHTML != "" {
		output.AppendWhere(cname, "", "DNS", "CNAME", "html", outputFileHTML)
	}

	if outputFileTXT != "" {
		output.AppendWhere(cname, "", "DNS", "CNAME", "txt", outputFileTXT)
	}

	// -- NS RECORDS --
	nameserver, _ := net.LookupNS(domain)

	for _, nsRecord := range nameserver {
		if !plain {
			fmt.Printf("[+]FOUND %s IN NS: ", domain)
			color.Green("%s\n", nsRecord.Host)
		} else {
			fmt.Printf("%s\n", nsRecord.Host)
		}

		if outputFileJSON != "" {
			output.AppendWhere(nsRecord.Host, "", "DNS", "NS", "json", outputFileJSON)
		}

		if outputFileHTML != "" {
			output.AppendWhere(nsRecord.Host, "", "DNS", "NS", "html", outputFileHTML)
		}

		if outputFileTXT != "" {
			output.AppendWhere(nsRecord.Host, "", "DNS", "NS", "txt", outputFileTXT)
		}
	}

	// -- MX RECORDS --
	mxrecords, _ := net.LookupMX(domain)

	for _, mxRecord := range mxrecords {
		if !plain {
			fmt.Printf("[+]FOUND %s IN MX: ", domain)
			color.Green("%s %v\n", mxRecord.Host, mxRecord.Pref)
		} else {
			fmt.Printf("%s %v\n", mxRecord.Host, mxRecord.Pref)
		}

		if outputFileJSON != "" {
			output.AppendWhere(mxRecord.Host, "", "DNS", "MX", "json", outputFileJSON)
		}

		if outputFileHTML != "" {
			output.AppendWhere(mxRecord.Host, "", "DNS", "MX", "html", outputFileHTML)
		}

		if outputFileTXT != "" {
			output.AppendWhere(mxRecord.Host, "", "DNS", "MX", "txt", outputFileTXT)
		}
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
			fmt.Printf("[+]FOUND %s IN SRV: ", domain)
			color.Green("%v:%v:%d:%d\n", srv.Target, srv.Port, srv.Priority, srv.Weight)
		} else {
			fmt.Printf("%v:%v:%d:%d\n", srv.Target, srv.Port, srv.Priority, srv.Weight)
		}

		if outputFileJSON != "" {
			output.AppendWhere(srv.Target, "", "DNS", "SRV", "json", outputFileJSON)
		}

		if outputFileHTML != "" {
			output.AppendWhere(srv.Target, "", "DNS", "SRV", "html", outputFileHTML)
		}

		if outputFileTXT != "" {
			output.AppendWhere(srv.Target, "", "DNS", "SRV", "txt", outputFileTXT)
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

		if outputFileJSON != "" {
			output.AppendWhere(txt, "", "DNS", "TXT", "json", outputFileJSON)
		}

		if outputFileHTML != "" {
			output.AppendWhere(txt, "", "DNS", "TXT", "html", outputFileHTML)
		}

		if outputFileTXT != "" {
			output.AppendWhere(txt, "", "DNS", "TXT", "txt", outputFileTXT)
		}
	}

	if outputFileHTML != "" {
		output.FooterHTML(outputFileHTML)
	}

	fmt.Println()
}
