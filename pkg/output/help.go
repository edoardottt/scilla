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

package output

import "fmt"

// Help prints in stdout scilla usage.
func Help() {
	fmt.Println("Information Gathering tool - DNS / Subdomain / Ports / Directories enumeration")
	fmt.Println("")
	fmt.Println("usage: scilla subcommand { options }")
	fmt.Println("")
	fmt.Println("   Available subcommands:")
	fmt.Println("       - dns [-oj JSON output file]")
	fmt.Println("             [-oh HTML output file]")
	fmt.Println("             [-ot TXT output file]")
	fmt.Println("             [-plain Print only results]")
	fmt.Println("             -target <target (URL/IP)> REQUIRED")
	fmt.Println("       - port [-p <start-end> or ports divided by comma]")
	fmt.Println("              [-oj JSON output file]")
	fmt.Println("              [-oh HTML output file]")
	fmt.Println("              [-ot TXT output file]")
	fmt.Println("              [-common scan common ports]")
	fmt.Println("              [-plain Print only results]")
	fmt.Println("              [-t Port scan timeout]")
	fmt.Println("              -target <target (URL/IP)> REQUIRED")
	fmt.Println("       - subdomain [-w wordlist]")
	fmt.Println("                   [-oj JSON output file]")
	fmt.Println("                   [-oh HTML output file]")
	fmt.Println("                   [-ot TXT output file]")
	fmt.Println("                   [-i ignore status codes]")
	fmt.Println("                   [-c use also a web crawler]")
	fmt.Println("                   [-db use also public databases]")
	fmt.Println("                   [-plain Print only results]")
	fmt.Println("                   [-db -no-check Don't check status codes for subdomains]")
	fmt.Println("                   [-db -vt Use VirusTotal as subdomains source]")
	fmt.Println("                   [-db -bw Use BuiltWith as subdomains source]")
	fmt.Println("                   [-ua Set the User Agent]")
	fmt.Println("                   [-rua Generate a random user agent for each request]")
	fmt.Println("                   [-dns Set DNS IP to resolve the subdomains]")
	fmt.Println("                   [-alive Check also if the subdomains are alive]")
	fmt.Println("                   -target <target (URL)> REQUIRED")
	fmt.Println("       - dir [-w wordlist]")
	fmt.Println("             [-oj JSON output file]")
	fmt.Println("             [-oh HTML output file]")
	fmt.Println("             [-ot TXT output file]")
	fmt.Println("             [-i ignore status codes]")
	fmt.Println("             [-c use also a web crawler]")
	fmt.Println("             [-plain Print only results]")
	fmt.Println("             [-nr No follow redirects]")
	fmt.Println("             [-ua Set the User Agent]")
	fmt.Println("             [-rua Generate a random user agent for each request]")
	fmt.Println("             -target <target (URL)> REQUIRED")
	fmt.Println("       - report [-p <start-end> or ports divided by comma]")
	fmt.Println("                [-ws subdomains wordlist]")
	fmt.Println("                [-wd directories wordlist]")
	fmt.Println("                [-oj JSON output file]")
	fmt.Println("                [-oh HTML output file]")
	fmt.Println("                [-ot TXT output file]")
	fmt.Println("                [-id ignore status codes in directories scanning]")
	fmt.Println("                [-is ignore status codes in subdomains scanning]")
	fmt.Println("                [-cd use also a web crawler for directories scanning]")
	fmt.Println("                [-cs use also a web crawler for subdomains scanning]")
	fmt.Println("                [-db use also public databases for subdomains scanning]")
	fmt.Println("                [-common scan common ports]")
	fmt.Println("                [-nr No follow redirects]")
	fmt.Println("                [-db -vt Use VirusTotal as subdomains source]")
	fmt.Println("                [-tp Port scan timeout]")
	fmt.Println("                [-ua Set the User Agent]")
	fmt.Println("                [-rua Generate a random user agent for each request]")
	fmt.Println("                [-dns Set DNS IP to resolve the subdomains]")
	fmt.Println("                [-alive Check also if the subdomains are alive]")
	fmt.Println("                -target <target (URL/IP)> REQUIRED")
	fmt.Println("       - help")
	fmt.Println("       - examples")
	fmt.Println()
}
