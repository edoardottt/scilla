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

package input

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/edoardottt/scilla/output"
)

const (
	PortScanTimeout = 3
	ArgsNumber      = 2
)

// Input is the struct containing the input parameters.
type Input struct {
	ReportTarget        string
	ReportWordDir       string
	ReportWordSub       string
	ReportOutputJSON    string
	ReportOutputHTML    string
	ReportOutputTXT     string
	ReportIgnoreDir     []string
	ReportIgnoreSub     []string
	ReportCrawlerDir    bool
	ReportCrawlerSub    bool
	ReportSubdomainDB   bool
	ReportCommon        bool
	ReportRedirect      bool
	ReportVirusTotal    bool
	ReportTimeoutPort   int
	DNSTarget           string
	DNSOutputJSON       string
	DNSOutputHTML       string
	DNSOutputTXT        string
	DNSPlain            bool
	SubdomainTarget     string
	SubdomainWord       string
	SubdomainOutputJSON string
	SubdomainOutputHTML string
	SubdomainOutputTXT  string
	SubdomainIgnore     []string
	SubdomainCrawler    bool
	SubdomainDB         bool
	SubdomainPlain      bool
	SubdomainNoCheck    bool
	SubdomainVirusTotal bool
	DirTarget           string
	DirWord             string
	DirOutputJSON       string
	DirOutputHTML       string
	DirOutputTXT        string
	DirIgnore           []string
	DirCrawler          bool
	DirPlain            bool
	DirRedirect         bool
	PortTarget          string
	PortOutputJSON      string
	PortOutputHTML      string
	PortOutputTXT       string
	StartPort           int
	EndPort             int
	PortArrayBool       bool
	PortsArray          []int
	PortCommon          bool
	PortPlain           bool
	PortTimeout         int
}

// ReadArgs reads arguments/options from stdin.
// Subcommands:
// 		report		==> Full report
// 		dns			==> Dns records enumeration
// 		subdomains	==> Subdomains enumeration
// 		port		==> ports enumeration
//		dir			==> directiories enumeration
// 		help		==> doc
//		examples	==> examples
func ReadArgs() Input {
	reportCommand := flag.NewFlagSet("report", flag.ExitOnError)
	dnsCommand := flag.NewFlagSet("dns", flag.ExitOnError)
	subdomainCommand := flag.NewFlagSet("subdomain", flag.ExitOnError)
	portCommand := flag.NewFlagSet("port", flag.ExitOnError)
	dirCommand := flag.NewFlagSet("dir", flag.ExitOnError)
	helpCommand := flag.NewFlagSet("help", flag.ExitOnError)
	examplesCommand := flag.NewFlagSet("examples", flag.ExitOnError)

	// report subcommand flag pointers
	reportTargetPtr := reportCommand.String("target", "", "Target {URL/IP} (Required)")

	// report subcommand flag pointers
	reportPortsPtr := reportCommand.String("p", "", "ports range <start-end>")

	// report subcommand flag pointers
	reportWordlistDirPtr := reportCommand.String("wd", "", "wordlist to use for directories (default enabled)")

	// report subcommand flag pointers
	reportWordlistSubdomainPtr := reportCommand.String("ws", "", "wordlist to use for subdomains (default enabled)")

	// report subcommand flag pointers
	reportOutputJSONPtr := reportCommand.String("oj", "", "JSON output path where save the results to")

	// report subcommand flag pointers
	reportOutputHTMLPtr := reportCommand.String("oh", "", "HTML output path where save the results to")

	// report subcommand flag pointers
	reportOutputTXTPtr := reportCommand.String("ot", "", "TXT output path where save the results to")

	// report subcommand flag pointers
	reportIgnoreDirPtr := reportCommand.String("id", "", "Ignore response code(s) for directories scanning")
	reportIgnoreDir := []string{}

	// report subcommand flag pointers
	reportIgnoreSubPtr := reportCommand.String("is", "", "Ignore response code(s) for subdomains scanning")
	reportIgnoreSub := []string{}

	// report subcommand flag pointers
	reportCrawlerDirPtr := reportCommand.Bool("cd", false, "Use also a web crawler for directories enumeration")

	// report subcommand flag pointers
	reportCrawlerSubdomainPtr := reportCommand.Bool("cs", false, "Use also a web crawler for subdomains enumeration")

	// report subcommand flag pointers
	reportSubdomainDBPtr := reportCommand.Bool("db", false, "Use also a public database for subdomains enumeration")

	// report subcommand flag pointers
	reportCommonPtr := reportCommand.Bool("common", false, "Scan common ports")

	// report subcommand flag pointers
	reportRedirectPtr := reportCommand.Bool("nr", false, "No follow redirects")

	// report subcommand flag pointers
	reportVirusTotalPtr := reportCommand.Bool("vt", false, "Use VirusTotal as a subdomain source")

	// report subcommand flag pointers
	reportTimeoutPortPtr := reportCommand.Int("tp", PortScanTimeout, "Port Scan timeout")

	// dns subcommand flag pointers
	dnsTargetPtr := dnsCommand.String("target", "", "Target {URL/IP} (Required)")

	// dns subcommand flag pointers
	dnsOutputJSONPtr := dnsCommand.String("oj", "", "JSON output path where save the results to")

	// dns subcommand flag pointers
	dnsOutputHTMLPtr := dnsCommand.String("oh", "", "HTML output path where save the results to")

	// dns subcommand flag pointers
	dnsOutputTXTPtr := dnsCommand.String("ot", "", "TXT output path where save the results to")

	// dns subcommand flag pointers
	dnsPlainPtr := dnsCommand.Bool("plain", false, "Print only results")

	// subdomains subcommand flag pointers
	subdomainTargetPtr := subdomainCommand.String("target", "", "Target {URL} (Required)")

	// subdomains subcommand wordlist
	subdomainWordlistPtr := subdomainCommand.String("w", "", "wordlist to use (default enabled)")

	// subdomains subcommand flag pointers
	subdomainOutputJSONPtr := subdomainCommand.String("oj", "", "JSON output path where save the results to")

	// subdomains subcommand flag pointers
	subdomainOutputHTMLPtr := subdomainCommand.String("oh", "", "HTML output path where save the results to")

	// subdomains subcommand flag pointers
	subdomainOutputTXTPtr := subdomainCommand.String("ot", "", "TXT output path where save the results to")

	// subdomains subcommand flag pointers
	subdomainIgnorePtr := subdomainCommand.String("i", "", "Ignore response code(s)")
	subdomainIgnore := []string{}

	// subdomains subcommand flag pointers
	subdomainCrawlerPtr := subdomainCommand.Bool("c", false, "Use also a web crawler")

	// subdomains subcommand flag pointers
	subdomainDBPtr := subdomainCommand.Bool("db", false, "Use also public databases")

	// subdomains subcommand flag pointers
	subdomainPlainPtr := subdomainCommand.Bool("plain", false, "Print only results")

	// subdomains subcommand flag pointers
	subdomainNoCheckPtr := subdomainCommand.Bool("no-check", false, "Don't check status codes for subdomains.")

	// subdomains subcommand flag pointers
	subdomainVirusTotalPtr := subdomainCommand.Bool("vt", false, "Use VirusTotal as a subdomain source")

	// dir subcommand flag pointers
	dirTargetPtr := dirCommand.String("target", "", "Target {URL/IP} (Required)")

	// dir subcommand wordlist
	dirWordlistPtr := dirCommand.String("w", "", "wordlist to use (default enabled)")

	// dir subcommand flag pointers
	dirOutputJSONPtr := dirCommand.String("oj", "", "JSON output path where save the results to")

	// dir subcommand flag pointers
	dirOutputHTMLPtr := dirCommand.String("oh", "", "HTML output path where save the results to")

	// dir subcommand flag pointers
	dirOutputTXTPtr := dirCommand.String("ot", "", "TXT output path where save the results to")

	// dir subcommand flag pointers
	dirIgnorePtr := dirCommand.String("i", "", "Ignore response code(s)")
	dirIgnore := []string{}

	// dir subcommand flag pointers
	dirCrawlerPtr := dirCommand.Bool("c", false, "Use also a web crawler")

	// dir subcommand flag pointers
	dirPlainPtr := dirCommand.Bool("plain", false, "Print only results")

	// dir subcommand flag pointers
	dirRedirectPtr := dirCommand.Bool("nr", false, "No follow redirects")

	// port subcommand flag pointers
	portTargetPtr := portCommand.String("target", "", "Target {URL/IP} (Required)")

	// port subcommand flag pointers
	portOutputJSONPtr := portCommand.String("oj", "", "JSON output path where save the results to")

	// port subcommand flag pointers
	portOutputHTMLPtr := portCommand.String("oh", "", "HTML output path where save the results to")

	// port subcommand flag pointers
	portOutputTXTPtr := portCommand.String("ot", "", "TXT output path where save the results to")

	// port subcommand flag pointers
	portsPtr := portCommand.String("p", "", "ports range <start-end>")

	// port subcommand flag pointers
	portCommonPtr := portCommand.Bool("common", false, "Scan common ports")

	// port subcommand flag pointers
	portPlainPtr := portCommand.Bool("plain", false, "Print only results")

	// port subcommand flag pointers
	portTimeoutPtr := portCommand.Int("t", PortScanTimeout, "Port scan timeout")

	// Default ports
	StartPort := 1
	EndPort := 65535
	portsArray := []int{}
	portArrayBool := false

	// Verify that a subcommand has been provided
	// os.Args[0] is the main command
	// os.Args[1] will be the subcommand
	if len(os.Args) < ArgsNumber {
		output.Intro()
		fmt.Println("[ERROR] subcommand is required.")
		fmt.Println("	Type: scilla help      - Full overview of the commands.")
		fmt.Println("	Type: scilla examples  - Some explanatory examples.")
		os.Exit(1)
	}

	// Switch on the subcommand
	// Parse the flags for appropriate FlagSet
	switch os.Args[1] {
	case "report":
		err := reportCommand.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}
	case "dns":
		err := dnsCommand.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}
	case "subdomain":
		err := subdomainCommand.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}
	case "port":
		err := portCommand.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}
	case "dir":
		err := dirCommand.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}
	case "help":
		output.Intro()

		err := helpCommand.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}
	case "examples":
		output.Intro()

		err := examplesCommand.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}
	default:
		output.Intro()
		output.Help()
		os.Exit(1)
	}

	// REPORT subcommand
	if reportCommand.Parsed() {
		if *reportOutputJSONPtr != "" {
			*reportOutputJSONPtr = output.AppendExtension(*reportOutputJSONPtr, "json")
		}

		if *reportOutputHTMLPtr != "" {
			*reportOutputHTMLPtr = output.AppendExtension(*reportOutputHTMLPtr, "html")
		}

		if *reportOutputTXTPtr != "" {
			*reportOutputTXTPtr = output.AppendExtension(*reportOutputTXTPtr, "txt")
		}

		StartPort, EndPort, portsArray, portArrayBool, reportIgnoreDir, reportIgnoreSub =
			ReportSubcommandCheckFlags(*reportCommand,
				reportTargetPtr, reportPortsPtr, reportCommonPtr,
				reportVirusTotalPtr, reportSubdomainDBPtr, StartPort,
				EndPort, reportIgnoreDirPtr, reportIgnoreSubPtr, reportTimeoutPortPtr,
				reportOutputJSONPtr, reportOutputHTMLPtr, reportOutputTXTPtr)
	}

	// DNS subcommand
	if dnsCommand.Parsed() {
		if *dnsOutputJSONPtr != "" {
			*dnsOutputJSONPtr = output.AppendExtension(*dnsOutputJSONPtr, "json")
		}

		if *dnsOutputHTMLPtr != "" {
			*dnsOutputHTMLPtr = output.AppendExtension(*dnsOutputHTMLPtr, "html")
		}

		if *dnsOutputTXTPtr != "" {
			*dnsOutputTXTPtr = output.AppendExtension(*dnsOutputTXTPtr, "txt")
		}

		DNSSubcommandCheckFlags(*dnsCommand, dnsTargetPtr, dnsOutputJSONPtr, dnsOutputHTMLPtr, dnsOutputTXTPtr)
	}

	// SUBDOMAIN subcommand
	if subdomainCommand.Parsed() {
		if *subdomainOutputJSONPtr != "" {
			*subdomainOutputJSONPtr = output.AppendExtension(*subdomainOutputJSONPtr, "json")
		}

		if *subdomainOutputHTMLPtr != "" {
			*subdomainOutputHTMLPtr = output.AppendExtension(*subdomainOutputHTMLPtr, "html")
		}

		if *subdomainOutputTXTPtr != "" {
			*subdomainOutputTXTPtr = output.AppendExtension(*subdomainOutputTXTPtr, "txt")
		}

		subdomainIgnore = SubdomainSubcommandCheckFlags(*subdomainCommand, subdomainTargetPtr,
			subdomainNoCheckPtr, subdomainDBPtr, subdomainWordlistPtr, subdomainIgnorePtr,
			subdomainCrawlerPtr, subdomainVirusTotalPtr,
			subdomainOutputJSONPtr, subdomainOutputHTMLPtr, subdomainOutputTXTPtr)
	}

	// PORT subcommand
	if portCommand.Parsed() {
		if *portOutputJSONPtr != "" {
			*portOutputJSONPtr = output.AppendExtension(*portOutputJSONPtr, "json")
		}

		if *portOutputHTMLPtr != "" {
			*portOutputHTMLPtr = output.AppendExtension(*portOutputHTMLPtr, "html")
		}

		if *portOutputTXTPtr != "" {
			*portOutputTXTPtr = output.AppendExtension(*portOutputTXTPtr, "txt")
		}

		StartPort, EndPort, portsArray, portArrayBool = PortSubcommandCheckFlags(*portCommand, portTargetPtr, portsPtr,
			portCommonPtr, StartPort, EndPort, portTimeoutPtr, portOutputJSONPtr, portOutputHTMLPtr, portOutputTXTPtr)
	}

	// DIR subcommand
	if dirCommand.Parsed() {
		if *dirOutputJSONPtr != "" {
			*dirOutputJSONPtr = output.AppendExtension(*dirOutputJSONPtr, "json")
		}

		if *dirOutputHTMLPtr != "" {
			*dirOutputHTMLPtr = output.AppendExtension(*dirOutputHTMLPtr, "html")
		}

		if *dirOutputTXTPtr != "" {
			*dirOutputTXTPtr = output.AppendExtension(*dirOutputTXTPtr, "txt")
		}

		dirIgnore = DirSubcommandCheckFlags(*dirCommand, dirTargetPtr, dirIgnorePtr,
			dirOutputJSONPtr, dirOutputHTMLPtr, dirOutputTXTPtr)
	}

	// HELP subcommand
	if helpCommand.Parsed() {
		// Print help
		output.Help()
		os.Exit(0)
	}

	// EXAMPLES subcommand
	if examplesCommand.Parsed() {
		// Print examples
		output.Examples()
		os.Exit(0)
	}

	result := Input{
		*reportTargetPtr,
		*reportWordlistDirPtr,
		*reportWordlistSubdomainPtr,
		*reportOutputJSONPtr,
		*reportOutputHTMLPtr,
		*reportOutputTXTPtr,
		reportIgnoreDir,
		reportIgnoreSub,
		*reportCrawlerDirPtr,
		*reportCrawlerSubdomainPtr,
		*reportSubdomainDBPtr,
		*reportCommonPtr,
		*reportRedirectPtr,
		*reportVirusTotalPtr,
		*reportTimeoutPortPtr,
		*dnsTargetPtr,
		*dnsOutputJSONPtr,
		*dnsOutputHTMLPtr,
		*dnsOutputTXTPtr,
		*dnsPlainPtr,
		*subdomainTargetPtr,
		*subdomainWordlistPtr,
		*subdomainOutputJSONPtr,
		*subdomainOutputHTMLPtr,
		*subdomainOutputTXTPtr,
		subdomainIgnore,
		*subdomainCrawlerPtr,
		*subdomainDBPtr,
		*subdomainPlainPtr,
		*subdomainNoCheckPtr,
		*subdomainVirusTotalPtr,
		*dirTargetPtr,
		*dirWordlistPtr,
		*dirOutputJSONPtr,
		*dirOutputHTMLPtr,
		*dirOutputTXTPtr,
		dirIgnore,
		*dirCrawlerPtr,
		*dirPlainPtr,
		*dirRedirectPtr,
		*portTargetPtr,
		*portOutputJSONPtr,
		*portOutputHTMLPtr,
		*portOutputTXTPtr,
		StartPort,
		EndPort,
		portArrayBool,
		portsArray,
		*portCommonPtr,
		*portPlainPtr,
		*portTimeoutPtr,
	}

	return result
}
