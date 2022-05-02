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
	"os"

	"github.com/edoardottt/scilla/output"
)

//Input is the struct containing the input parameters
type Input struct {
	ReportTarget        string
	ReportWordDir       string
	ReportWordSub       string
	ReportOutputJson    string
	ReportOutputHtml    string
	ReportOutputTxt     string
	ReportIgnoreDir     []string
	ReportIgnoreSub     []string
	ReportCrawlerDir    bool
	ReportCrawlerSub    bool
	ReportSubdomainDB   bool
	ReportCommon        bool
	ReportRedirect      bool
	ReportSpyse         bool
	ReportVirusTotal    bool
	ReportTimeoutPort   int
	DNSTarget           string
	DNSOutputJson       string
	DNSOutputHtml       string
	DNSOutputTxt        string
	DNSPlain            bool
	SubdomainTarget     string
	SubdomainWord       string
	SubdomainOutputJson string
	SubdomainOutputHtml string
	SubdomainOutputTxt  string
	SubdomainIgnore     []string
	SubdomainCrawler    bool
	SubdomainDB         bool
	SubdomainPlain      bool
	SubdomainNoCheck    bool
	SubdomainSpyse      bool
	SubdomainVirusTotal bool
	DirTarget           string
	DirWord             string
	DirOutputJson       string
	DirOutputHtml       string
	DirOutputTxt        string
	DirIgnore           []string
	DirCrawler          bool
	DirPlain            bool
	DirRedirect         bool
	PortTarget          string
	PortOutputJson      string
	PortOutputHtml      string
	PortOutputTxt       string
	StartPort           int
	EndPort             int
	PortArrayBool       bool
	PortsArray          []int
	PortCommon          bool
	PortPlain           bool
	PortTimeout         int
}

//ReadArgs reads arguments/options from stdin
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
	reportOutputJsonPtr := reportCommand.String("oj", "", "JSON output path where save the results to")

	// report subcommand flag pointers
	reportOutputHtmlPtr := reportCommand.String("oh", "", "HTML output path where save the results to")

	// report subcommand flag pointers
	reportOutputTxtPtr := reportCommand.String("ot", "", "TXT output path where save the results to")

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
	reportSpysePtr := reportCommand.Bool("spyse", false, "Use Spyse as a subdomain source")

	// report subcommand flag pointers
	reportVirusTotalPtr := reportCommand.Bool("vt", false, "Use VirusTotal as a subdomain source")

	// report subcommand flag pointers
	reportTimeoutPortPtr := reportCommand.Int("tp", 3, "Scan Port timeout")

	// dns subcommand flag pointers
	dnsTargetPtr := dnsCommand.String("target", "", "Target {URL/IP} (Required)")

	// dns subcommand flag pointers
	dnsOutputJsonPtr := dnsCommand.String("oj", "", "JSON output path where save the results to")

	// dns subcommand flag pointers
	dnsOutputHtmlPtr := dnsCommand.String("oh", "", "HTML output path where save the results to")

	// dns subcommand flag pointers
	dnsOutputTxtPtr := dnsCommand.String("ot", "", "TXT output path where save the results to")

	// dns subcommand flag pointers
	dnsPlainPtr := dnsCommand.Bool("plain", false, "Print only results")

	// subdomains subcommand flag pointers
	subdomainTargetPtr := subdomainCommand.String("target", "", "Target {URL} (Required)")

	// subdomains subcommand wordlist
	subdomainWordlistPtr := subdomainCommand.String("w", "", "wordlist to use (default enabled)")

	// subdomains subcommand flag pointers
	subdomainOutputJsonPtr := subdomainCommand.String("oj", "", "JSON output path where save the results to")

	// subdomains subcommand flag pointers
	subdomainOutputHtmlPtr := subdomainCommand.String("oh", "", "HTML output path where save the results to")

	// subdomains subcommand flag pointers
	subdomainOutputTxtPtr := subdomainCommand.String("ot", "", "TXT output path where save the results to")

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
	subdomainSpysePtr := subdomainCommand.Bool("spyse", false, "Use Spyse as a subdomain source")

	// subdomains subcommand flag pointers
	subdomainVirusTotalPtr := subdomainCommand.Bool("vt", false, "Use VirusTotal as a subdomain source")

	// dir subcommand flag pointers
	dirTargetPtr := dirCommand.String("target", "", "Target {URL/IP} (Required)")

	// dir subcommand wordlist
	dirWordlistPtr := dirCommand.String("w", "", "wordlist to use (default enabled)")

	// dir subcommand flag pointers
	dirOutputJsonPtr := dirCommand.String("oj", "", "JSON output path where save the results to")

	// dir subcommand flag pointers
	dirOutputHtmlPtr := dirCommand.String("oh", "", "HTML output path where save the results to")

	// dir subcommand flag pointers
	dirOutputTxtPtr := dirCommand.String("ot", "", "TXT output path where save the results to")

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
	portOutputJsonPtr := portCommand.String("oj", "", "JSON output path where save the results to")

	// port subcommand flag pointers
	portOutputHtmlPtr := portCommand.String("oh", "", "HTML output path where save the results to")

	// port subcommand flag pointers
	portOutputTxtPtr := portCommand.String("ot", "", "TXT output path where save the results to")

	// port subcommand flag pointers
	portsPtr := portCommand.String("p", "", "ports range <start-end>")

	// port subcommand flag pointers
	portCommonPtr := portCommand.Bool("common", false, "Scan common ports")

	// port subcommand flag pointers
	portPlainPtr := portCommand.Bool("plain", false, "Print only results")

	// port subcommand flag pointers
	portTimeoutPtr := portCommand.Int("t", 3, "Port scan timeout")

	// Default ports
	StartPort := 1
	EndPort := 65535
	portsArray := []int{}
	portArrayBool := false

	// Verify that a subcommand has been provided
	// os.Arg[0] is the main command
	// os.Arg[1] will be the subcommand
	if len(os.Args) < 2 {
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
		reportCommand.Parse(os.Args[2:])
	case "dns":
		dnsCommand.Parse(os.Args[2:])
	case "subdomain":
		subdomainCommand.Parse(os.Args[2:])
	case "port":
		portCommand.Parse(os.Args[2:])
	case "dir":
		dirCommand.Parse(os.Args[2:])
	case "help":
		output.Intro()
		helpCommand.Parse(os.Args[2:])
	case "examples":
		output.Intro()
		examplesCommand.Parse(os.Args[2:])
	default:
		output.Intro()
		output.Help()
		os.Exit(1)
	}

	// REPORT subcommand
	if reportCommand.Parsed() {
		StartPort, EndPort, portsArray, portArrayBool, reportIgnoreDir, reportIgnoreSub = ReportSubcommandCheckFlags(*reportCommand,
			reportTargetPtr, reportPortsPtr, reportCommonPtr,
			reportSpysePtr, reportVirusTotalPtr, reportSubdomainDBPtr, StartPort,
			EndPort, reportIgnoreDirPtr, reportIgnoreSubPtr, reportTimeoutPortPtr)
	}

	// DNS subcommand
	if dnsCommand.Parsed() {
		DNSSubcommandCheckFlags(*dnsCommand, dnsTargetPtr)
	}

	// SUBDOMAIN subcommand
	if subdomainCommand.Parsed() {
		subdomainIgnore = SubdomainSubcommandCheckFlags(*subdomainCommand, subdomainTargetPtr,
			subdomainNoCheckPtr, subdomainDBPtr, subdomainWordlistPtr, subdomainIgnorePtr,
			subdomainCrawlerPtr, subdomainSpysePtr, subdomainVirusTotalPtr)
	}

	// PORT subcommand
	if portCommand.Parsed() {
		StartPort, EndPort, portsArray, portArrayBool = PortSubcommandCheckFlags(*portCommand, portTargetPtr, portsPtr,
			portCommonPtr, StartPort, EndPort, portTimeoutPtr)
	}

	// DIR subcommand
	if dirCommand.Parsed() {
		dirIgnore = DirSubcommandCheckFlags(*dirCommand, dirTargetPtr, dirIgnorePtr)
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
		*reportOutputJsonPtr,
		*reportOutputHtmlPtr,
		*reportOutputTxtPtr,
		reportIgnoreDir,
		reportIgnoreSub,
		*reportCrawlerDirPtr,
		*reportCrawlerSubdomainPtr,
		*reportSubdomainDBPtr,
		*reportCommonPtr,
		*reportRedirectPtr,
		*reportSpysePtr,
		*reportVirusTotalPtr,
		*reportTimeoutPortPtr,
		*dnsTargetPtr,
		*dnsOutputJsonPtr,
		*dnsOutputHtmlPtr,
		*dnsOutputTxtPtr,
		*dnsPlainPtr,
		*subdomainTargetPtr,
		*subdomainWordlistPtr,
		*subdomainOutputJsonPtr,
		*subdomainOutputHtmlPtr,
		*subdomainOutputTxtPtr,
		subdomainIgnore,
		*subdomainCrawlerPtr,
		*subdomainDBPtr,
		*subdomainPlainPtr,
		*subdomainNoCheckPtr,
		*subdomainSpysePtr,
		*subdomainVirusTotalPtr,
		*dirTargetPtr,
		*dirWordlistPtr,
		*dirOutputJsonPtr,
		*dirOutputHtmlPtr,
		*dirOutputTxtPtr,
		dirIgnore,
		*dirCrawlerPtr,
		*dirPlainPtr,
		*dirRedirectPtr,
		*portTargetPtr,
		*portOutputJsonPtr,
		*portOutputHtmlPtr,
		*portOutputTxtPtr,
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
