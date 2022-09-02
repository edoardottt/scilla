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
	"strings"

	"github.com/edoardottt/scilla/utils"
)

// ReportSubcommandCheckFlags performs all the necessary checks on the flags
// for the report subcommand
func ReportSubcommandCheckFlags(reportCommand flag.FlagSet, reportTargetPtr *string,
	reportPortsPtr *string, reportCommonPtr *bool, reportVirusTotalPtr *bool, reportSubdomainDBPtr *bool,
	startPort int, endPort int, reportIgnoreDirPtr *string,
	reportIgnoreSubPtr *string, reportTimeoutPort *int,
	reportOutputJson *string, reportOutputHtml *string, reportOutputTxt *string) (int, int,
	[]int, bool, []string, []string) {
	// Required Flags
	if *reportTargetPtr == "" {
		reportCommand.PrintDefaults()
		os.Exit(1)
	}

	// Verify good inputs
	if !utils.IsURL(*reportTargetPtr) {
		fmt.Println("The inputted target is not valid.")
		os.Exit(1)
	}

	// output files all different
	if *reportOutputJson != "" {
		if *reportOutputJson == *reportOutputTxt || *reportOutputJson == *reportOutputHtml {
			fmt.Println("The output paths must be all different.")
			os.Exit(1)
		}
	}

	if *reportOutputHtml != "" {
		if *reportOutputHtml == *reportOutputTxt || *reportOutputJson == *reportOutputHtml {
			fmt.Println("The output paths must be all different.")
			os.Exit(1)
		}
	}

	if *reportOutputTxt != "" {
		if *reportOutputJson == *reportOutputTxt || *reportOutputTxt == *reportOutputHtml {
			fmt.Println("The output paths must be all different.")
			os.Exit(1)
		}
	}

	// common and p not together
	if *reportPortsPtr != "" && *reportCommonPtr {
		fmt.Println("You can't specify a port range and common option together.")
		os.Exit(1)
	}

	if *reportVirusTotalPtr && !*reportSubdomainDBPtr {
		fmt.Println("You can't specify VirusTotal and not the Open Database option.")
		fmt.Println("If you want to use VirusTotal Api, set also -db option.")
		os.Exit(1)
	}

	var portsArray []int

	var portArrayBool bool

	if *reportPortsPtr != "" {
		if strings.Contains(*reportPortsPtr, "-") && strings.Contains(*reportPortsPtr, ",") {
			fmt.Println("You can specify a ports range or an array, not both.")
			os.Exit(1)
		}

		switch {
		case strings.Contains(*reportPortsPtr, "-"):
			{
				portsRange := *reportPortsPtr
				startPort, endPort = utils.CheckPortsRange(portsRange, startPort, endPort)
				portArrayBool = false
				break
			}
		case strings.Contains(*reportPortsPtr, ","):
			{
				portsArray = utils.CheckPortsArray(*reportPortsPtr)
				portArrayBool = true
				break
			}
		default:
			{
				portsRange := *reportPortsPtr
				startPort, endPort = utils.CheckPortsRange(portsRange, startPort, endPort)
				portArrayBool = false
			}
		}
	}

	var reportIgnoreDir, reportIgnoreSub []string

	if *reportIgnoreDirPtr != "" {
		toBeIgnored := *reportIgnoreDirPtr
		reportIgnoreDir = utils.CheckIgnore(toBeIgnored)
	}

	if *reportIgnoreSubPtr != "" {
		toBeIgnored := *reportIgnoreSubPtr
		reportIgnoreSub = utils.CheckIgnore(toBeIgnored)
	}

	if *reportTimeoutPort < 1 || *reportTimeoutPort > 100 {
		fmt.Println("Port Scan timeout must be an integer between 1 and 100.")
		os.Exit(1)
	}

	return startPort, endPort, portsArray, portArrayBool, reportIgnoreDir, reportIgnoreSub
}

// DNSSubcommandCheckFlags performs all the necessary checks on the flags
// for the dns subcommand
func DNSSubcommandCheckFlags(dnsCommand flag.FlagSet, dnsTargetPtr, dnsOutputJson,
	dnsOutputHtml, dnsOutputTxt *string) {
	// Required Flags
	if *dnsTargetPtr == "" {
		dnsCommand.PrintDefaults()
		os.Exit(1)
	}

	// output files all different
	if *dnsOutputJson != "" {
		if *dnsOutputJson == *dnsOutputTxt || *dnsOutputJson == *dnsOutputHtml {
			fmt.Println("The output paths must be all different.")
			os.Exit(1)
		}
	}

	if *dnsOutputHtml != "" {
		if *dnsOutputHtml == *dnsOutputTxt || *dnsOutputJson == *dnsOutputHtml {
			fmt.Println("The output paths must be all different.")
			os.Exit(1)
		}
	}

	if *dnsOutputTxt != "" {
		if *dnsOutputJson == *dnsOutputTxt || *dnsOutputTxt == *dnsOutputHtml {
			fmt.Println("The output paths must be all different.")
			os.Exit(1)
		}
	}

	// Verify good inputs
	if !utils.IsURL(*dnsTargetPtr) {
		fmt.Println("The inputted target is not valid.")
		os.Exit(1)
	}
}

// SubdomainSubcommandCheckFlags performs all the necessary checks on the flags
// for the subdomain subcommand
func SubdomainSubcommandCheckFlags(subdomainCommand flag.FlagSet, subdomainTargetPtr *string,
	subdomainNoCheckPtr *bool, subdomainDBPtr *bool, subdomainWordlistPtr *string,
	subdomainIgnorePtr *string, subdomainCrawlerPtr *bool, subdomainVirusTotalPtr *bool,
	subdomainOutputJSON, subdomainOutputHTML, subdomainOutputTXT *string) []string {

	// Required Flags
	if *subdomainTargetPtr == "" {
		subdomainCommand.PrintDefaults()
		os.Exit(1)
	}

	// Verify good inputs
	if !utils.IsURL(*subdomainTargetPtr) {
		fmt.Println("The inputted target is not valid.")
		os.Exit(1)
	}

	// no-check checks
	if *subdomainNoCheckPtr && !*subdomainDBPtr {
		fmt.Println("You can use no-check only with db option.")
		os.Exit(1)
	}

	if *subdomainNoCheckPtr && *subdomainWordlistPtr != "" {
		fmt.Println("You can't use no-check with wordlist option.")
		os.Exit(1)
	}

	if *subdomainNoCheckPtr && *subdomainIgnorePtr != "" {
		fmt.Println("You can't use no-check with ignore option.")
		os.Exit(1)
	}

	if *subdomainNoCheckPtr && *subdomainCrawlerPtr {
		fmt.Println("You can't use no-check with crawler option.")
		os.Exit(1)
	}

	// output files all different
	if *subdomainOutputJSON != "" {
		if *subdomainOutputJSON == *subdomainOutputTXT || *subdomainOutputJSON == *subdomainOutputHTML {
			fmt.Println("The output paths must be all different.")
			os.Exit(1)
		}
	}

	if *subdomainOutputHTML != "" {
		if *subdomainOutputHTML == *subdomainOutputTXT || *subdomainOutputJSON == *subdomainOutputHTML {
			fmt.Println("The output paths must be all different.")
			os.Exit(1)
		}
	}

	if *subdomainOutputTXT != "" {
		if *subdomainOutputJSON == *subdomainOutputTXT || *subdomainOutputTXT == *subdomainOutputHTML {
			fmt.Println("The output paths must be all different.")
			os.Exit(1)
		}
	}

	if *subdomainVirusTotalPtr && !*subdomainDBPtr {
		fmt.Println("You can't specify VirusTotal and not the Open Database option.")
		fmt.Println("If you want to use VirusTotal Api, set also -db option.")
		os.Exit(1)
	}

	var subdomainIgnore []string

	if *subdomainIgnorePtr != "" {
		toBeIgnored := *subdomainIgnorePtr
		subdomainIgnore = utils.CheckIgnore(toBeIgnored)
	}

	return subdomainIgnore
}

// PortSubcommandCheckFlags performs all the necessary checks on the flags
// for the port subcommand
func PortSubcommandCheckFlags(portCommand flag.FlagSet, portTargetPtr *string, portsPtr *string,
	portCommonPtr *bool, startPort int, endPort int, portTimeout *int,
	portOutputJSON, portOutputHTML, portOutputTXT *string) (int, int, []int, bool) {
	// Required Flags
	if *portTargetPtr == "" {
		portCommand.PrintDefaults()
		os.Exit(1)
	}

	// common and p not together
	if *portsPtr != "" && *portCommonPtr {
		fmt.Println("You can't specify a port range and common option together.")
		os.Exit(1)
	}

	var portArrayBool bool

	var portsArray []int

	if *portsPtr != "" {
		if strings.Contains(*portsPtr, "-") && strings.Contains(*portsPtr, ",") {
			fmt.Println("You can specify a ports range or an array, not both.")
			os.Exit(1)
		}

		switch {
		case strings.Contains(*portsPtr, "-"):
			{
				portsRange := *portsPtr
				startPort, endPort = utils.CheckPortsRange(portsRange, startPort, endPort)
				portArrayBool = false
				break
			}
		case strings.Contains(*portsPtr, ","):
			{
				portsArray = utils.CheckPortsArray(*portsPtr)
				portArrayBool = true
				break
			}
		default:
			{
				portsRange := *portsPtr
				startPort, endPort = utils.CheckPortsRange(portsRange, startPort, endPort)
				portArrayBool = false
			}
		}
	}

	// output files all different
	if *portOutputJSON != "" {
		if *portOutputJSON == *portOutputTXT || *portOutputJSON == *portOutputHTML {
			fmt.Println("The output paths must be all different.")
			os.Exit(1)
		}
	}

	if *portOutputHTML != "" {
		if *portOutputHTML == *portOutputTXT || *portOutputJSON == *portOutputHTML {
			fmt.Println("The output paths must be all different.")
			os.Exit(1)
		}
	}

	if *portOutputTXT != "" {
		if *portOutputJSON == *portOutputTXT || *portOutputTXT == *portOutputHTML {
			fmt.Println("The output paths must be all different.")
			os.Exit(1)
		}
	}

	// Verify good inputs
	if !utils.IsURL(*portTargetPtr) {
		fmt.Println("The inputted target is not valid.")
		os.Exit(1)
	}

	if *portTimeout < 1 || *portTimeout > 100 {
		fmt.Println("Port Scan timeout must be an integer between 1 and 100.")
		os.Exit(1)
	}

	return startPort, endPort, portsArray, portArrayBool
}

// DirSubcommandCheckFlags performs all the necessary checks on the flags
// for the dir subcommand
func DirSubcommandCheckFlags(dirCommand flag.FlagSet, dirTargetPtr *string,
	dirIgnorePtr *string, dirOutputJSON, dirOutputHTML, dirOutputTXT *string) []string {
	// Required Flags
	if *dirTargetPtr == "" {
		dirCommand.PrintDefaults()
		os.Exit(1)
	}

	// Verify good inputs
	if !utils.IsURL(*dirTargetPtr) {
		fmt.Println("The inputted target is not valid.")
		os.Exit(1)
	}

	// output files all different
	if *dirOutputJSON != "" {
		if *dirOutputJSON == *dirOutputTXT || *dirOutputJSON == *dirOutputHTML {
			fmt.Println("The output paths must be all different.")
			os.Exit(1)
		}
	}

	if *dirOutputHTML != "" {
		if *dirOutputHTML == *dirOutputTXT || *dirOutputJSON == *dirOutputHTML {
			fmt.Println("The output paths must be all different.")
			os.Exit(1)
		}
	}

	if *dirOutputTXT != "" {
		if *dirOutputJSON == *dirOutputTXT || *dirOutputTXT == *dirOutputHTML {
			fmt.Println("The output paths must be all different.")
			os.Exit(1)
		}
	}

	var dirIgnore []string

	if *dirIgnorePtr != "" {
		toBeIgnored := *dirIgnorePtr
		dirIgnore = utils.CheckIgnore(toBeIgnored)
	}

	return dirIgnore
}
