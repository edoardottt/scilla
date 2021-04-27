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
*/

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/gocolly/colly"
)

//intro prints the banner when the program is started
func intro() {
	banner1 := "                  _ _ _\n"
	banner2 := "         ___  ___(_) | | __ _\n"
	banner3 := "        / __|/ __| | | |/ _` |\n"
	banner4 := "        \\__ \\ (__| | | | (_| |\n"
	banner5 := "        |___/\\___|_|_|_|\\__,_| v1.0\n"
	banner6 := " > github.com/edoardottt/scilla\n"
	banner7 := " > edoardoottavianelli.it"
	bannerPart1 := banner1 + banner2 + banner3 + banner4 + banner5
	bannerPart2 := banner6 + banner7
	color.Cyan("%s\n", bannerPart1)
	fmt.Println(bannerPart2)
	fmt.Println("========================================")
}

//help prints in stdout scilla usage
func help() {
	fmt.Println("Information Gathering tool - DNS / Subdomain / Ports / Directories enumeration")
	fmt.Println("")
	fmt.Println("usage: scilla subcommand { options }")
	fmt.Println("")
	fmt.Println("   Available subcommands:")
	fmt.Println("       - dns [-o output-format]")
	fmt.Println("             [-plain Print only results]")
	fmt.Println("             -target <target (URL/IP)> REQUIRED")
	fmt.Println("       - port [-p <start-end> or ports divided by comma]")
	fmt.Println("              [-o output-format]")
	fmt.Println("              [-common scan common ports]")
	fmt.Println("              [-plain Print only results]")
	fmt.Println("              -target <target (URL/IP)> REQUIRED")
	fmt.Println("       - subdomain [-w wordlist]")
	fmt.Println("                   [-o output-format]")
	fmt.Println("                   [-i ignore status codes]")
	fmt.Println("                   [-c use also a web crawler]")
	fmt.Println("                   [-db use also a public database]")
	fmt.Println("                   [-plain Print only results]")
	fmt.Println("                   -target <target (URL)> REQUIRED")
	fmt.Println("       - dir [-w wordlist]")
	fmt.Println("             [-o output-format]")
	fmt.Println("             [-i ignore status codes]")
	fmt.Println("             [-c use also a web crawler]")
	fmt.Println("             [-plain Print only results]")
	fmt.Println("             -target <target (URL)> REQUIRED")
	fmt.Println("       - report [-p <start-end> or ports divided by comma]")
	fmt.Println("                [-ws subdomains wordlist]")
	fmt.Println("                [-wd directories wordlist]")
	fmt.Println("                [-o output-format]")
	fmt.Println("                [-id ignore status codes in directories scanning]")
	fmt.Println("                [-is ignore status codes in subdomains scanning]")
	fmt.Println("                [-cd use also a web crawler for directories scanning]")
	fmt.Println("                [-cs use also a web crawler for subdomains scanning]")
	fmt.Println("                [-db use also a public database for subdomains scanning]")
	fmt.Println("                [-common scan common ports]")
	fmt.Println("                -target <target (URL/IP)> REQUIRED")
	fmt.Println("       - help")
	fmt.Println("       - examples")
	fmt.Println()
}

//examples prints some examples
func examples() {
	fmt.Println("	Examples:")
	fmt.Println("		- scilla dns -target target.domain")
	fmt.Println("		- scilla dns -target -o txt target.domain")
	fmt.Println("		- scilla dns -target -o html target.domain")
	fmt.Println("		- scilla dns -target -plain target.domain")
	fmt.Println()
	fmt.Println("		- scilla subdomain -target target.domain")
	fmt.Println("		- scilla subdomain -w wordlist.txt -target target.domain")
	fmt.Println("		- scilla subdomain -o txt -target target.domain")
	fmt.Println("		- scilla subdomain -o html -target target.domain")
	fmt.Println("		- scilla subdomain -i 400 -target target.domain")
	fmt.Println("		- scilla subdomain -i 4** -target target.domain")
	fmt.Println("		- scilla subdomain -c -target target.domain")
	fmt.Println("		- scilla subdomain -db -target target.domain")
	fmt.Println("		- scilla subdomain -plain -target target.domain")
	fmt.Println()
	fmt.Println("		- scilla port -p -450 -target target.domain")
	fmt.Println("		- scilla port -p 90- -target target.domain")
	fmt.Println("		- scilla port -p 10-1000 -target target.domain")
	fmt.Println("		- scilla port -o txt -target target.domain")
	fmt.Println("		- scilla port -o html -target target.domain")
	fmt.Println("		- scilla port -p 21,25,80 -target target.domain")
	fmt.Println("		- scilla port -common -target target.domain")
	fmt.Println("		- scilla port -plain -target target.domain")
	fmt.Println()
	fmt.Println("		- scilla dir -target target.domain")
	fmt.Println("		- scilla dir -o txt -target target.domain")
	fmt.Println("		- scilla dir -o html -target target.domain")
	fmt.Println("		- scilla dir -w wordlist.txt -target target.domain")
	fmt.Println("		- scilla dir -i 500,401 -target target.domain")
	fmt.Println("		- scilla dir -i 5**,401 -target target.domain")
	fmt.Println("		- scilla dir -c -target target.domain")
	fmt.Println("		- scilla dir -plain -target target.domain")
	fmt.Println()
	fmt.Println("		- scilla report -p 80 -target target.domain")
	fmt.Println("		- scilla report -o txt -target target.domain")
	fmt.Println("		- scilla report -o html -target target.domain")
	fmt.Println("		- scilla report -p 50-200 -target target.domain")
	fmt.Println("		- scilla report -wd dirs.txt -target target.domain")
	fmt.Println("		- scilla report -ws subdomains.txt -target target.domain")
	fmt.Println("		- scilla report -id 500,501,502 -target target.domain")
	fmt.Println("		- scilla report -is 500,501,502 -target target.domain")
	fmt.Println("		- scilla report -id 5**,4** -target target.domain")
	fmt.Println("		- scilla report -is 5**,4** -target target.domain")
	fmt.Println("		- scilla report -cd -target target.domain")
	fmt.Println("		- scilla report -cs -target target.domain")
	fmt.Println("		- scilla report -db -target target.domain")
	fmt.Println("		- scilla report -p 21,25,80 -target target.domain")
	fmt.Println("		- scilla report -common -target target.domain")
	fmt.Println("")
}

//main function
func main() {
	input := readArgs()
	// common assets found (only subdomain and dir)
	subs := make(map[string]Asset)
	dirs := make(map[string]Asset)
	//common ports
	common := []int{13, 20, 21, 22, 23, 25, 42, 50, 51, 53, 67, 68,
		69, 70, 79, 80, 88, 102, 107, 109, 110, 111, 113, 115, 118,
		119, 123, 135, 136, 137, 138, 139, 143, 156, 161, 162, 179,
		194, 220, 311, 389, 443, 445, 464, 500, 512, 513, 514, 515,
		530, 543, 546, 547, 556, 587, 631, 636, 660, 749, 802, 853,
		873, 902, 989, 990, 992, 993, 994, 995, 1000, 1025, 1080,
		1194, 1241, 1293, 1337, 1417, 1433, 1434, 1527, 1755, 1812,
		1813, 1880, 1883, 2000, 2049, 2095, 2096, 2222, 2483, 2484,
		2638, 3000, 3268, 3283, 3333, 3306, 3389, 4000, 4444, 5000,
		5432, 5555, 5938, 6000, 6666, 7000, 7071, 7777, 8000, 8001,
		8002, 8003, 8004, 8005, 8080, 8200, 8888, 9050, 10000}
	execute(input, subs, dirs, common)
}

//Asset gives information about the asset found
type Asset struct {
	Value   string
	Printed bool
}

//execute reads inputs and starts the correct procedure
func execute(input Input, subs map[string]Asset, dirs map[string]Asset, common []int) {

	var mutex = &sync.Mutex{}
	if input.ReportTarget != "" {
		intro()
		target := cleanProtocol(input.ReportTarget)
		var targetIP string
		fmt.Printf("target: %s\n", target)
		fmt.Println("=============== FULL REPORT ===============")
		outputFile := ""
		if input.ReportOutput != "" {
			outputFile = createOutputFile(input.ReportTarget, "report", input.ReportOutput)
			if outputFile[len(outputFile)-4:] == "html" {
				bannerHTML(input.ReportTarget, outputFile)
			}
		}

		fmt.Println("=============== SUBDOMAINS SCANNING ===============")
		var strings1 []string
		// change from ip to Hostname
		if isIP(target) {
			targetIP = target
			target = ipToHostname(target)
		}
		if outputFile != "" {
			if outputFile[len(outputFile)-4:] == "html" {
				headerHTML("SUBDOMAIN SCANNING", outputFile)
			}
		}
		if input.ReportCrawlerSub {
			go spawnCrawler(target, input.ReportIgnoreSub, dirs, subs, outputFile, mutex, "sub", false)
		}
		strings1 = createSubdomains(input.ReportWordSub, target)
		if input.ReportSubdomainDB {
			sonar := sonarSubdomains(target)
			strings1 = appendDBSubdomains(sonar, strings1)
			hackerTarget := hackerTargetSubdomains(target)
			strings1 = appendDBSubdomains(hackerTarget, strings1)
			bufferOverrun := bufferOverrunSubdomains(target)
			strings1 = appendDBSubdomains(bufferOverrun, strings1)
		}
		// be sure to not scan duplicate values
		strings1 = removeDuplicateValues(strings1)
		asyncGet(strings1, input.ReportIgnoreSub, outputFile, subs, mutex, false)
		if outputFile != "" {
			if outputFile[len(outputFile)-4:] == "html" {
				footerHTML(outputFile)
			}
		}

		if targetIP != "" {
			target = targetIP
		}
		fmt.Println("=============== PORT SCANNING ===============")

		asyncPort(input.portsArray, input.portArrayBool, input.StartPort, input.EndPort, target, outputFile, input.ReportCommon, common, false)

		fmt.Println("=============== DNS SCANNING ===============")
		lookupDNS(target, outputFile, false)

		fmt.Println("=============== DIRECTORIES SCANNING ===============")
		var strings2 []string
		strings2 = createUrls(input.ReportWordDir, target)
		if outputFile != "" {
			if outputFile[len(outputFile)-4:] == "html" {
				headerHTML("DIRECTORY SCANNING", outputFile)
			}
		}
		if input.ReportCrawlerDir {
			go spawnCrawler(target, input.ReportIgnoreDir, dirs, subs, outputFile, mutex, "dir", false)
		}
		asyncDir(strings2, input.ReportIgnoreDir, outputFile, dirs, mutex, false)
		if outputFile != "" {
			if outputFile[len(outputFile)-4:] == "html" {
				footerHTML(outputFile)
			}
		}
		if input.ReportOutput != "" {
			if outputFile[len(outputFile)-4:] == "html" {
				bannerFooterHTML(outputFile)
			}
		}
	}

	if input.DNSTarget != "" {
		if !input.DNSPlain {
			intro()
		}
		target := cleanProtocol(input.DNSTarget)
		// change from ip to Hostname
		if isIP(target) {
			target = ipToHostname(target)
		}
		if !input.DNSPlain {
			fmt.Printf("target: %s\n", target)
			fmt.Println("=============== DNS SCANNING ===============")
		}
		outputFile := ""
		if input.DNSOutput != "" {
			outputFile = createOutputFile(input.DNSTarget, "dns", input.DNSOutput)

			if outputFile[len(outputFile)-4:] == "html" {
				bannerHTML(input.DNSTarget, outputFile)
			}
		}
		lookupDNS(target, outputFile, input.DNSPlain)
		if input.DNSOutput != "" {
			if outputFile[len(outputFile)-4:] == "html" {
				bannerFooterHTML(outputFile)
			}
		}
	}

	if input.SubdomainTarget != "" {

		if !input.SubdomainPlain {
			intro()
		}

		target := cleanProtocol(input.SubdomainTarget)
		// change from ip to Hostname
		if isIP(target) {
			target = ipToHostname(target)
		}
		if !input.SubdomainPlain {
			fmt.Printf("target: %s\n", target)
			fmt.Println("=============== SUBDOMAINS SCANNING ===============")
		}
		outputFile := ""
		if input.SubdomainOutput != "" {
			outputFile = createOutputFile(input.SubdomainTarget, "subdomain", input.SubdomainOutput)
			if outputFile[len(outputFile)-4:] == "html" {
				bannerHTML(input.SubdomainTarget, outputFile)
			}
		}
		var strings1 []string
		strings1 = createSubdomains(input.SubdomainWord, target)
		if input.SubdomainDB {
			sonar := sonarSubdomains(target)
			strings1 = appendDBSubdomains(sonar, strings1)
			hackerTarget := hackerTargetSubdomains(target)
			strings1 = appendDBSubdomains(hackerTarget, strings1)
			bufferOverrun := bufferOverrunSubdomains(target)
			strings1 = appendDBSubdomains(bufferOverrun, strings1)
		}
		if outputFile != "" {
			if outputFile[len(outputFile)-4:] == "html" {
				headerHTML("SUBDOMAIN SCANNING", outputFile)
			}
		}
		if input.SubdomainCrawler {
			go spawnCrawler(target, input.SubdomainIgnore, dirs, subs, outputFile, mutex, "sub", input.SubdomainPlain)
		}
		// be sure to not scan duplicate values
		strings1 = removeDuplicateValues(strings1)
		asyncGet(strings1, input.SubdomainIgnore, outputFile, subs, mutex, input.SubdomainPlain)
		if outputFile != "" {
			if outputFile[len(outputFile)-4:] == "html" {
				footerHTML(outputFile)
			}
		}
		if input.SubdomainOutput != "" {
			if outputFile[len(outputFile)-4:] == "html" {
				bannerFooterHTML(outputFile)
			}
		}
	}

	if input.DirTarget != "" {

		if !input.DirPlain {
			intro()
		}

		target := cleanProtocol(input.DirTarget)
		if !input.DirPlain {
			fmt.Printf("target: %s\n", target)
			fmt.Println("=============== DIRECTORIES SCANNING ===============")
		}
		outputFile := ""
		if input.DirOutput != "" {
			outputFile = createOutputFile(input.DirTarget, "dir", input.DirOutput)

			if outputFile[len(outputFile)-4:] == "html" {
				bannerHTML(input.DirTarget, outputFile)
			}
		}
		var strings2 []string
		strings2 = createUrls(input.DirWord, target)
		if outputFile != "" {
			if outputFile[len(outputFile)-4:] == "html" {
				headerHTML("DIRECTORY SCANNING", outputFile)
			}
		}
		if input.DirCrawler {
			go spawnCrawler(target, input.DirIgnore, dirs, subs, outputFile, mutex, "dir", input.DirPlain)
		}
		asyncDir(strings2, input.DirIgnore, outputFile, dirs, mutex, input.DirPlain)
		if outputFile != "" {
			if outputFile[len(outputFile)-4:] == "html" {
				footerHTML(outputFile)
			}
		}
		if input.DirOutput != "" {
			if outputFile[len(outputFile)-4:] == "html" {
				bannerFooterHTML(outputFile)
			}
		}
	}

	if input.PortTarget != "" {
		if !input.PortPlain {
			intro()
		}
		target := input.PortTarget
		if isURL(target) {
			target = cleanProtocol(input.PortTarget)
		}
		outputFile := ""
		if input.PortOutput != "" {
			outputFile = createOutputFile(input.PortTarget, "port", input.PortOutput)
			if outputFile[len(outputFile)-4:] == "html" {
				bannerHTML(input.PortTarget, outputFile)
			}
		}
		if !input.PortPlain {
			fmt.Printf("target: %s\n", target)
			fmt.Println("=============== PORT SCANNING ===============")
		}
		asyncPort(input.portsArray, input.portArrayBool, input.StartPort, input.EndPort, target, outputFile, input.PortCommon, common, input.PortPlain)

		if input.PortOutput != "" {
			if outputFile[len(outputFile)-4:] == "html" {
				bannerFooterHTML(outputFile)
			}
		}
	}
}

//cleanProtocol remove from the url the protocol scheme
// (http - https - tls)
func cleanProtocol(target string) string {
	if len(target) > 6 {
		// clean protocols and go ahead
		if target[:6] == "tls://" {
			target = target[6:]
		}
	}
	if len(target) > 7 {
		if target[:7] == "http://" {
			target = target[7:]
		}
	}
	if len(target) > 8 {
		if target[:8] == "https://" {
			target = target[8:]
		}
	}
	if target[len(target)-1:] == "/" {
		return target[:len(target)-1]
	}
	return target
}

//outputFormatIsOk (txt or html)
func outputFormatIsOk(input string) bool {
	if input == "" {
		return true
	}
	acceptedOutput := [2]string{"txt", "html"}
	input = strings.ToLower(input)
	for _, output := range acceptedOutput {
		if output == input {
			return true
		}
	}
	return false
}

//Input struct contains the input parameters
type Input struct {
	ReportTarget      string
	ReportWordDir     string
	ReportWordSub     string
	ReportOutput      string
	ReportIgnoreDir   []string
	ReportIgnoreSub   []string
	ReportCrawlerDir  bool
	ReportCrawlerSub  bool
	ReportSubdomainDB bool
	ReportCommon      bool
	DNSTarget         string
	DNSOutput         string
	DNSPlain          bool
	SubdomainTarget   string
	SubdomainWord     string
	SubdomainOutput   string
	SubdomainIgnore   []string
	SubdomainCrawler  bool
	SubdomainDB       bool
	SubdomainPlain    bool
	DirTarget         string
	DirWord           string
	DirOutput         string
	DirIgnore         []string
	DirCrawler        bool
	DirPlain          bool
	PortTarget        string
	PortOutput        string
	StartPort         int
	EndPort           int
	portArrayBool     bool
	portsArray        []int
	PortCommon        bool
	PortPlain         bool
}

//readArgs reads arguments/options from stdin
// Subcommands:
// 		report		==> Full report
// 		dns			==> Dns records enumeration
// 		subdomains	==> Subdomains enumeration
// 		port		==> ports enumeration
//		dir			==> directiories enumeration
// 		help		==> doc
//		examples	==> examples
func readArgs() Input {
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
	reportOutputPtr := reportCommand.String("o", "", "output format (txt/html)")

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
	reportSubdomainDBPtr := reportCommand.Bool("cdb", false, "Use also a public database for subdomains enumeration")

	// report subcommand flag pointers
	reportCommonPtr := reportCommand.Bool("common", false, "Scan common ports")

	// dns subcommand flag pointers
	dnsTargetPtr := dnsCommand.String("target", "", "Target {URL/IP} (Required)")

	// dns subcommand flag pointers
	dnsOutputPtr := dnsCommand.String("o", "", "output format (txt/html)")

	// dns subcommand flag pointers
	dnsPlainPtr := dnsCommand.Bool("plain", false, "Print only results")

	// subdomains subcommand flag pointers
	subdomainTargetPtr := subdomainCommand.String("target", "", "Target {URL} (Required)")

	// subdomains subcommand wordlist
	subdomainWordlistPtr := subdomainCommand.String("w", "", "wordlist to use (default enabled)")

	// subdomains subcommand flag pointers
	subdomainOutputPtr := subdomainCommand.String("o", "", "output format (txt/html)")

	// subdomains subcommand flag pointers
	subdomainIgnorePtr := subdomainCommand.String("i", "", "Ignore response code(s)")
	subdomainIgnore := []string{}

	// subdomains subcommand flag pointers
	subdomainCrawlerPtr := subdomainCommand.Bool("c", false, "Use also a web crawler")

	// subdomains subcommand flag pointers
	subdomainDBPtr := subdomainCommand.Bool("db", false, "Use also a public database")

	// subdomains subcommand flag pointers
	subdomainPlainPtr := subdomainCommand.Bool("plain", false, "Print only results")

	// dir subcommand flag pointers
	dirTargetPtr := dirCommand.String("target", "", "Target {URL/IP} (Required)")

	// dir subcommand wordlist
	dirWordlistPtr := dirCommand.String("w", "", "wordlist to use (default enabled)")

	// dir subcommand flag pointers
	dirOutputPtr := dirCommand.String("o", "", "output format (txt/html)")

	// dir subcommand flag pointers
	dirIgnorePtr := dirCommand.String("i", "", "Ignore response code(s)")
	dirIgnore := []string{}

	// dir subcommand flag pointers
	dirCrawlerPtr := dirCommand.Bool("c", false, "Use also a web crawler")

	// dir subcommand flag pointers
	dirPlainPtr := dirCommand.Bool("plain", false, "Print only results")

	// port subcommand flag pointers
	portTargetPtr := portCommand.String("target", "", "Target {URL/IP} (Required)")

	// port subcommand flag pointers
	portOutputPtr := portCommand.String("o", "", "output format (txt/html)")

	// port subcommand flag pointers
	portsPtr := portCommand.String("p", "", "ports range <start-end>")

	// port subcommand flag pointers
	portCommonPtr := portCommand.Bool("common", false, "Scan common ports")

	// port subcommand flag pointers
	portPlainPtr := portCommand.Bool("plain", false, "Print only results")

	// Default ports
	StartPort := 1
	EndPort := 65535
	portsArray := []int{}
	portArrayBool := false

	// Verify that a subcommand has been provided
	// os.Arg[0] is the main command
	// os.Arg[1] will be the subcommand
	if len(os.Args) < 2 {
		intro()
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
		intro()
		helpCommand.Parse(os.Args[2:])
	case "examples":
		intro()
		examplesCommand.Parse(os.Args[2:])
	default:
		intro()
		flag.PrintDefaults()
		os.Exit(1)
	}

	// REPORT subcommand
	if reportCommand.Parsed() {

		// Required Flags
		if *reportTargetPtr == "" {
			reportCommand.PrintDefaults()
			os.Exit(1)
		}
		//Verify good inputs
		if !isURL(*reportTargetPtr) {
			fmt.Println("The inputted target is not valid.")
			os.Exit(1)
		}
		if !outputFormatIsOk(*reportOutputPtr) {
			fmt.Println("The output format is not valid.")
			os.Exit(1)
		}
		//common and p not together
		if *reportPortsPtr != "" && *reportCommonPtr {
			fmt.Println("You can't specify a port range and common option together.")
			os.Exit(1)
		}

		if *reportPortsPtr != "" {
			if strings.Contains(*reportPortsPtr, "-") && strings.Contains(*reportPortsPtr, ",") {
				fmt.Println("You can specify a ports range or an array, not both.")
				os.Exit(1)
			}
			if strings.Contains(*reportPortsPtr, "-") {
				portsRange := string(*reportPortsPtr)
				StartPort, EndPort = checkPortsRange(portsRange, StartPort, EndPort)
				portArrayBool = false
			} else if strings.Contains(*reportPortsPtr, ",") {
				portsArray = checkPortsArray(*reportPortsPtr)
				portArrayBool = true
			} else {
				portsRange := string(*reportPortsPtr)
				StartPort, EndPort = checkPortsRange(portsRange, StartPort, EndPort)
				portArrayBool = false
			}
		}
		if *reportIgnoreDirPtr != "" {
			toBeIgnored := string(*reportIgnoreDirPtr)
			reportIgnoreDir = checkIgnore(toBeIgnored)
		}
		if *reportIgnoreSubPtr != "" {
			toBeIgnored := string(*reportIgnoreSubPtr)
			reportIgnoreSub = checkIgnore(toBeIgnored)
		}
	}

	// DNS subcommand
	if dnsCommand.Parsed() {

		// Required Flags
		if *dnsTargetPtr == "" {
			dnsCommand.PrintDefaults()
			os.Exit(1)
		}
		//Verify good inputs
		if !isURL(*dnsTargetPtr) {
			fmt.Println("The inputted target is not valid.")
			os.Exit(1)
		}
		if !outputFormatIsOk(*dnsOutputPtr) {
			fmt.Println("The output format is not valid.")
			os.Exit(1)
		}
	}

	// SUBDOMAIN subcommand
	if subdomainCommand.Parsed() {

		// Required Flags
		if *subdomainTargetPtr == "" {
			subdomainCommand.PrintDefaults()
			os.Exit(1)
		}
		//Verify good inputs
		if !isURL(*subdomainTargetPtr) {
			fmt.Println("The inputted target is not valid.")
			os.Exit(1)
		}
		if !outputFormatIsOk(*subdomainOutputPtr) {
			fmt.Println("The output format is not valid.")
			os.Exit(1)
		}
		if *subdomainIgnorePtr != "" {
			toBeIgnored := string(*subdomainIgnorePtr)
			subdomainIgnore = checkIgnore(toBeIgnored)
		}
	}

	// PORT subcommand
	if portCommand.Parsed() {

		// Required Flags
		if *portTargetPtr == "" {
			portCommand.PrintDefaults()
			os.Exit(1)
		}
		//common and p not together
		if *portsPtr != "" && *portCommonPtr {
			fmt.Println("You can't specify a port range and common option together.")
			os.Exit(1)
		}
		if *portsPtr != "" {
			if strings.Contains(*portsPtr, "-") && strings.Contains(*portsPtr, ",") {
				fmt.Println("You can specify a ports range or an array, not both.")
				os.Exit(1)
			}
			if strings.Contains(*portsPtr, "-") {
				portsRange := string(*portsPtr)
				StartPort, EndPort = checkPortsRange(portsRange, StartPort, EndPort)
				portArrayBool = false
			} else if strings.Contains(*portsPtr, ",") {
				portsArray = checkPortsArray(*portsPtr)
				portArrayBool = true
			} else {
				portsRange := string(*portsPtr)
				StartPort, EndPort = checkPortsRange(portsRange, StartPort, EndPort)
				portArrayBool = false
			}
		}
		//Verify good inputs
		if !isURL(*portTargetPtr) {
			fmt.Println("The inputted target is not valid.")
			os.Exit(1)
		}
		if !outputFormatIsOk(*portOutputPtr) {
			fmt.Println("The output format is not valid.")
			os.Exit(1)
		}
	}

	// DIR subcommand
	if dirCommand.Parsed() {

		// Required Flags
		if *dirTargetPtr == "" {
			dirCommand.PrintDefaults()
			os.Exit(1)
		}
		//Verify good inputs
		if !isURL(*dirTargetPtr) {
			fmt.Println("The inputted target is not valid.")
			os.Exit(1)
		}
		if !outputFormatIsOk(*dirOutputPtr) {
			fmt.Println("The output format is not valid.")
			os.Exit(1)
		}
		if *dirIgnorePtr != "" {
			toBeIgnored := string(*dirIgnorePtr)
			dirIgnore = checkIgnore(toBeIgnored)
		}
	}

	// HELP subcommand
	if helpCommand.Parsed() {
		// Print help
		help()
		os.Exit(0)
	}

	// EXAMPLES subcommand
	if examplesCommand.Parsed() {
		// Print examples
		examples()
		os.Exit(0)
	}

	result := Input{
		*reportTargetPtr,
		*reportWordlistDirPtr,
		*reportWordlistSubdomainPtr,
		*reportOutputPtr,
		reportIgnoreDir,
		reportIgnoreSub,
		*reportCrawlerDirPtr,
		*reportCrawlerSubdomainPtr,
		*reportSubdomainDBPtr,
		*reportCommonPtr,
		*dnsTargetPtr,
		*dnsOutputPtr,
		*dnsPlainPtr,
		*subdomainTargetPtr,
		*subdomainWordlistPtr,
		*subdomainOutputPtr,
		subdomainIgnore,
		*subdomainCrawlerPtr,
		*subdomainDBPtr,
		*subdomainPlainPtr,
		*dirTargetPtr,
		*dirWordlistPtr,
		*dirOutputPtr,
		dirIgnore,
		*dirCrawlerPtr,
		*dirPlainPtr,
		*portTargetPtr,
		*portOutputPtr,
		StartPort,
		EndPort,
		portArrayBool,
		portsArray,
		*portCommonPtr,
		*portPlainPtr,
	}
	return result
}

//checkIgnore checks the inputted status code to be ignored
func checkIgnore(input string) []string {
	result := []string{}
	temp := strings.Split(input, ",")
	temp = removeDuplicateValues(temp)
	for _, elem := range temp {
		elem := strings.TrimSpace(elem)
		if len(elem) != 3 {
			fmt.Println("The status code you entered is invalid (It should consist of three characters).")
			os.Exit(1)
		}
		if ignoreInt, err := strconv.Atoi(elem); err == nil {
			// if it is a valid status code without * (e.g. 404)
			if 100 <= ignoreInt && ignoreInt <= 599 {
				result = append(result, elem)
			} else {
				fmt.Println("The status code you entered is invalid (100 <= code <= 599).")
				os.Exit(1)
			}
		} else if strings.Contains(elem, "*") {
			// if it is a valid status code without * (e.g. 4**)
			if ignoreClassOk(elem) {
				result = append(result, elem)
			} else {
				fmt.Println("The status code you entered is invalid. You can enter * only as 1**,2**,3**,4**,5**.")
				os.Exit(1)
			}
		}
	}
	result = removeDuplicateValues(result)
	result = deleteUnusefulIgnoreresponses(result)
	return result
}

//deleteUnusefulIgnoreresponses removes from to-be-ignored arrays
//the responses included yet with * as classes
func deleteUnusefulIgnoreresponses(input []string) []string {
	var result []string
	toberemoved := []string{}
	classes := []string{}
	for _, elem := range input {
		if strings.Contains(elem, "*") {
			classes = append(classes, elem)
		}
	}
	for _, class := range classes {
		for _, elem := range input {
			if class[0] == elem[0] && elem[1] != '*' {
				toberemoved = append(toberemoved, elem)
			}
		}
	}
	result = Difference(input, toberemoved)
	return result
}

//Difference A - B
func Difference(a, b []string) (diff []string) {
	m := make(map[string]bool)
	for _, item := range b {
		m[item] = true
	}
	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}

//ignoreClass states if the class of ignored status codes
//is correct or not (4**,2**...)
func ignoreClassOk(input string) bool {
	if strings.Contains(input, "*") {
		if _, err := strconv.Atoi(string(input[0])); err == nil {
			i, err := strconv.Atoi(string(input[0]))
			if err != nil {
				// handle error
				fmt.Println(err)
				os.Exit(2)
			}
			if i >= 1 && i <= 5 {
				if input[1] == byte('*') && input[2] == byte('*') {
					return true
				}
			}
		}
	}
	return false
}

//checkPortsArray checks the basic rules to
//be valid and then returns the ports array to scan.
func checkPortsArray(input string) []int {
	delimiter := byte(',')
	sliceOfPorts := strings.Split(input, string(delimiter))
	sliceOfPorts = removeDuplicateValues(sliceOfPorts)
	result := []int{}
	for _, elem := range sliceOfPorts {
		try, err := strconv.Atoi(elem)
		if err != nil {
			fmt.Println("The inputted ports array is not valid.")
			os.Exit(1)
		}
		if try > 0 && try < 65536 {
			result = append(result, try)
		}
	}
	return result
}

//checkPortsRange checks the basic rules to
//be valid and then returns the starting port and the ending port.
func checkPortsRange(portsRange string, StartPort int, EndPort int) (int, int) {
	// If there's ports range, define it as inputs for the struct
	delimiter := byte('-')
	//If there is only one number

	// If starting port isn't specified
	if portsRange[0] == delimiter {
		maybeEnd, err := strconv.Atoi(portsRange[1:])
		if err != nil {
			fmt.Println("The inputted port range is not valid.")
			os.Exit(1)
		}
		if maybeEnd >= 1 && maybeEnd <= EndPort {
			EndPort = maybeEnd
		}
	} else if portsRange[len(portsRange)-1] == delimiter {
		// If ending port isn't specified
		maybeStart, err := strconv.Atoi(portsRange[:len(portsRange)-1])
		if err != nil {
			fmt.Println("The inputted port range is not valid.")
			os.Exit(1)
		}
		if maybeStart > 0 && maybeStart < EndPort {
			StartPort = maybeStart
		}
	} else if !strings.Contains(portsRange, string(delimiter)) {
		// If a single port is specified
		maybePort, err := strconv.Atoi(portsRange)
		if err != nil {
			fmt.Println("The inputted port range is not valid.")
			os.Exit(1)
		}
		if maybePort > 0 && maybePort < EndPort {
			StartPort = maybePort
			EndPort = maybePort
		}
	} else {
		// If a range is specified
		sliceOfPorts := strings.Split(portsRange, string(delimiter))
		if len(sliceOfPorts) != 2 {
			fmt.Println("The inputted port range is not valid.")
			os.Exit(1)
		}
		maybeStart, err := strconv.Atoi(sliceOfPorts[0])
		if err != nil {
			fmt.Println("The inputted port range is not valid.")
			os.Exit(1)
		}
		maybeEnd, err := strconv.Atoi(sliceOfPorts[1])
		if err != nil {
			fmt.Println("The inputted port range is not valid.")
			os.Exit(1)
		}
		if maybeStart > maybeEnd || maybeStart < 1 || maybeEnd > EndPort {
			fmt.Println("The inputted port range is not valid.")
			os.Exit(1)
		}
		StartPort = maybeStart
		EndPort = maybeEnd
	}
	return StartPort, EndPort
}

//sonarSubdomains retrieves from the below url some known subdomains.
func sonarSubdomains(target string) []string {
	var arr []string
	resp, err := http.Get("https://sonar.omnisint.io/subdomains/" + target)
	if err != nil {
		return arr
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return arr
		}
		bodyString := string(bodyBytes)
		_ = json.Unmarshal([]byte(bodyString), &arr)
	}
	for index, elem := range arr {
		arr[index] = "http://" + elem
	}
	return arr
}

//appendDBSubdomains appends to the subdomains in the list
//the subdomains found with the open DBs.
func appendDBSubdomains(dbsubs []string, urls []string) []string {
	if len(dbsubs) == 0 {
		return urls
	}
	var result = []string{}
	dbsubs = removeDuplicateValues(dbsubs)
	result = append(dbsubs, urls...)
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(result), func(i, j int) { result[i], result[j] = result[j], result[i] })
	return result
}

//hackerTargetSubdomain retrieves from the below url some known subdomains.
func hackerTargetSubdomains(domain string) []string {
	result := make([]string, 0)
	raw, err := http.Get("https://api.hackertarget.com/hostsearch/?q=" + domain)
	if err != nil {
		return result
	}
	res, err := ioutil.ReadAll(raw.Body)
	if err != nil {
		return result
	}
	raw.Body.Close()
	sc := bufio.NewScanner(bytes.NewReader(res))
	for sc.Scan() {
		parts := strings.SplitN(sc.Text(), ",", 2)
		if len(parts) != 2 {
			continue
		}
		result = append(result, parts[0])
	}
	return result
}

//bufferOverrunSubdomains retrieves from the below url some known subdomains.
func bufferOverrunSubdomains(domain string) []string {
	result := make([]string, 0)
	url := "https://dns.bufferover.run/dns?q=" + domain
	wrapper := struct {
		Records []string `json:"FDNS_A"`
	}{}
	resp, err := http.Get(url)
	if err != nil {
		return result
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)

	dec.Decode(&wrapper)
	if err != nil {
		return result
	}
	for _, r := range wrapper.Records {
		parts := strings.SplitN(r, ",", 2)
		if len(parts) != 2 {
			continue
		}
		result = append(result, parts[1])
	}
	return result
}

//replaceBadCharacterOutput
func replaceBadCharacterOutput(input string) string {
	result := strings.ReplaceAll(input, "/", "-")
	return result
}

//createOutputFolder
func createOutputFolder() {
	//Create a folder/directory at a full qualified path
	err := os.Mkdir("output-scilla", 0755)
	if err != nil {
		fmt.Println("Can't create output folder.")
		os.Exit(1)
	}
}

//createOutputFile
func createOutputFile(target string, subcommand string, format string) string {
	target = replaceBadCharacterOutput(target)
	filename := "output-scilla" + "/" + target + "." + subcommand + "." + format
	_, err := os.Stat(filename)

	if os.IsNotExist(err) {
		if _, err := os.Stat("output-scilla/"); os.IsNotExist(err) {
			createOutputFolder()
		}
		// If the file doesn't exist, create it.
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Can't create output file.")
			os.Exit(1)
		}
		f.Close()
	} else {
		// The file already exists, check what the user want.
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("The output file already esists, do you want to overwrite? (Y/n): ")
		text, _ := reader.ReadString('\n')
		answer := strings.ToLower(text)
		answer = strings.TrimSpace(answer)

		if answer == "y" || answer == "yes" || answer == "" {
			f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println("Can't create output file.")
				os.Exit(1)
			}
			err = f.Truncate(0)
			if err != nil {
				fmt.Println("Can't create output file.")
				os.Exit(1)
			}
			f.Close()
		} else {
			os.Exit(1)
		}
	}
	return filename
}

//isUrl checks if the inputted Url is valid
func isURL(str string) bool {
	target := cleanProtocol(str)
	str = "http://" + target
	u, err := url.Parse(str)
	return err == nil && u.Host != ""
}

//isIP
func isIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

//ipToHostname
func ipToHostname(ip string) string {
	addr, err := net.LookupAddr(ip)
	if err != nil || len(addr) == 0 {
		log.Fatalf("Failed to resolve ip address %s", ip)
	}
	return strings.TrimSuffix(addr[0], ".")
}

//get performs an HTTP GET request to the target
func get(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return true
}

//buildUrl returns full URL with the subdomain
func buildURL(subdomain string, domain string) string {
	return "http://" + subdomain + "." + domain
}

//appendDir returns full URL with the directory
func appendDir(domain string, dir string) (string, string) {
	return "http://" + domain + "/" + dir + "/", "http://" + domain + "/" + dir
}

//readDict scan all the possible subdomains from file
func readDictSubs(inputFile string) []string {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("failed to open %s ", inputFile)
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	file.Close()
	return text
}

//readDict scan all the possible dirs from file
func readDictDirs(inputFile string) []string {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("failed to open %s ", inputFile)
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string
	var dir = ""
	for scanner.Scan() {
		dir = scanner.Text()
		if string(dir[len(dir)-1:]) == "/" {
			dir = dir[:len(dir)-1]
		}
		text = append(text, dir)
	}
	file.Close()
	text = removeDuplicateValues(text)
	return text
}

//removeDuplicateValues
func removeDuplicateValues(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

//createSubdomains returns a list of subdomains
//from the default file lists/subdomains.txt.
func createSubdomains(filename string, url string) []string {
	var subs []string
	if filename == "" {
		if runtime.GOOS == "windows" {
			subs = readDictSubs("lists/subdomains.txt")
		} else { // linux
			subs = readDictSubs("/usr/bin/lists/subdomains.txt")
		}
	} else {
		subs = readDictSubs(filename)
	}
	result := []string{}
	for _, sub := range subs {
		path := buildURL(sub, url)
		result = append(result, path)
	}
	return result
}

//createUrls returns a list of directories
//from the default file lists/dirs.txt.
func createUrls(filename string, url string) []string {
	var dirs []string
	if filename == "" {
		if runtime.GOOS == "windows" {
			dirs = readDictDirs("lists/dirs.txt")
		} else { // linux
			dirs = readDictDirs("/usr/bin/lists/dirs.txt")
		}
	} else {
		dirs = readDictDirs(filename)
	}
	result := []string{}
	for _, dir := range dirs {
		path, path2 := appendDir(url, dir)
		result = append(result, path)
		result = append(result, path2)
	}
	return result
}

//appendOutputToTxt
func appendOutputToTxt(output string, filename string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	if _, err := file.WriteString(cleanProtocol(output) + "\n"); err != nil {
		log.Fatal(err)
	}
	file.Close()
}

//bannerHTML
func bannerHTML(target string, filename string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	file.WriteString("<html><body><div style='" + "background-color:#4adeff;color:white" + "'><h1>Scilla - Information Gathering Tool</h1>")
	file.WriteString("<ul>")
	file.WriteString("<li><a href='" + "https://github.com/edoardottt/scilla'" + ">github.com/edoardottt/scilla</a></li>")
	file.WriteString("<li>edoardottt, <a href='" + "https://www.edoardoottavianelli.it'" + ">edoardoottavianelli.it</a></li>")
	file.WriteString("<li>Released under <a href='" + "http://www.gnu.org/licenses/gpl-3.0.html'" + ">GPLv3 License</a></li></ul></div>")
	file.WriteString("<h4>target: " + target + "</h4>")
	file.Close()
}

//appendOutputToHtml
func appendOutputToHTML(output string, status string, filename string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	var statusColor string
	if status != "" {
		if string(status[0]) == "2" || string(status[0]) == "3" {
			statusColor = "<p style='color:green;display:inline'>" + status + "</p>"
		} else {
			statusColor = "<p style='color:red;display:inline'>" + status + "</p>"
		}
	} else {
		statusColor = status
	}
	if _, err := file.WriteString("<li><a target='_blank' href='" + output + "'>" + cleanProtocol(output) + "</a> " + statusColor + "</li>"); err != nil {
		log.Fatal(err)
	}
	file.Close()
}

//headerHtml
func headerHTML(header string, filename string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	if _, err := file.WriteString("<h3>" + header + "</h3><ul>"); err != nil {
		log.Fatal(err)
	}
	file.Close()
}

//footerHTML
func footerHTML(filename string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	if _, err := file.WriteString("</ul>"); err != nil {
		log.Fatal(err)
	}
	file.Close()
}

//bannerFooterHTML
func bannerFooterHTML(filename string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	file.WriteString("<div style='" + "background-color:#4adeff;color:white" + "'>")
	file.WriteString("<ul><li><a href='" + "https://github.com/edoardottt/scilla'" + ">Contribute to scilla</a></li>")
	file.WriteString("<li>Released under <a href='" + "http://www.gnu.org/licenses/gpl-3.0.html'" + ">GPLv3 License</a></li></ul></div>")
	file.Close()
}

//percentage
func percentage(done, total int) float64 {
	result := (float64(done) / float64(total)) * 100
	return result
}

//ignoreResponse returns a boolean if the response
//should be ignored or not.
func ignoreResponse(response int, ignore []string) bool {
	responseString := strconv.Itoa(response)
	// if I don't have to ignore responses, just return true
	if len(ignore) == 0 {
		return false
	}
	for _, ignorePort := range ignore {
		if strings.Contains(ignorePort, "*") {
			if responseString[0] == ignorePort[0] {
				return true
			}
		}
		if responseString == ignorePort {
			return true
		}
	}
	return false
}

//asyncGet performs concurrent requests to the specified
//urls and prints the results
func asyncGet(urls []string, ignore []string, outputFile string, subs map[string]Asset, mutex *sync.Mutex, plain bool) {
	ignoreBool := len(ignore) != 0
	var count int = 0
	var total int = len(urls)
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	limiter := make(chan string, 10) // Limits simultaneous requests
	wg := sync.WaitGroup{}           // Needed to not prematurely exit before all requests have been finished

	for i, domain := range urls {
		limiter <- domain
		wg.Add(1)
		if count%50 == 0 { // update counter
			if !plain {
				fmt.Fprint(os.Stdout, "\r \r")
			}
			printSubs(subs, ignore, outputFile, mutex, plain)
		}
		if !plain && count%100 == 0 { // update counter
			fmt.Fprint(os.Stdout, "\r \r")
			fmt.Printf("%0.2f%% : %d / %d", percentage(count, total), count, total)
		}
		go func(i int, domain string) {
			defer wg.Done()
			defer func() { <-limiter }()
			resp, err := client.Get(domain)
			count++
			if err != nil {
				return
			}
			if ignoreBool {
				if ignoreResponse(resp.StatusCode, ignore) {
					return
				}
			}
			addSubs(domain, resp.Status, subs, mutex)
			resp.Body.Close()
		}(i, domain)
	}
	printSubs(subs, ignore, outputFile, mutex, plain)
	wg.Wait()
	printSubs(subs, ignore, outputFile, mutex, plain)
}

//appendWhere (html or txt file)
func appendWhere(what string, status string, outputFile string) {
	if outputFile[len(outputFile)-4:] == "html" {
		appendOutputToHTML(what, status, outputFile)
	} else {
		appendOutputToTxt(what, outputFile)
	}
}

//isOpenPort scans if a port is open
func isOpenPort(host string, port string) bool {
	timeout := 3 * time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		return false
	}
	if conn != nil {
		defer conn.Close()
		return true
	}
	return false
}

//asyncPort performs concurrent requests to the specified
//ports range and, if someone is open it prints the results
func asyncPort(portsArray []int, portsArrayBool bool, StartingPort int, EndingPort int, host string, outputFile string, common bool, commonPorts []int, plain bool) {
	var count int = 0
	var total int = (EndingPort - StartingPort) + 1
	if portsArrayBool {
		total = len(portsArray)
	}
	limiter := make(chan string, 200) // Limits simultaneous requests
	wg := sync.WaitGroup{}            // Needed to not prematurely exit before all requests have been finished
	if outputFile != "" {
		if outputFile[len(outputFile)-4:] == "html" {
			headerHTML("PORT SCANNING", outputFile)
		}
	}
	ports := []int{}
	if !common {
		if portsArrayBool {
			ports = portsArray
		} else {
			for port := StartingPort; port <= EndingPort; port++ {
				ports = append(ports, port)
			}
		}
	} else {
		ports = commonPorts
	}
	for _, port := range ports {
		wg.Add(1)
		portStr := fmt.Sprint(port)
		limiter <- portStr
		if !plain && count%100 == 0 { // update counter
			fmt.Fprint(os.Stdout, "\r \r")
			fmt.Printf("%0.2f%% : %d / %d", percentage(count, total), count, total)
		}
		go func(portStr string, host string) {
			defer func() { <-limiter }()
			defer wg.Done()
			resp := isOpenPort(host, portStr)
			count++
			if resp {
				if !plain {
					fmt.Fprint(os.Stdout, "\r \r")
					fmt.Printf("[+]FOUND: %s ", host)
					color.Green("%s\n", portStr)
				} else {
					fmt.Printf("%s\n", portStr)
				}
				if outputFile != "" {
					appendWhere("http://"+host+":"+portStr, "", outputFile)
				}
			}
		}(portStr, host)
	}
	wg.Wait()
	fmt.Fprint(os.Stdout, "\r \r")
	if outputFile != "" {
		if outputFile[len(outputFile)-4:] == "html" {
			footerHTML(outputFile)
		}
	}
}

//lookupDNS prints the DNS informations for the inputted domain
func lookupDNS(domain string, outputFile string, plain bool) {
	if outputFile != "" {
		if outputFile[len(outputFile)-4:] == "html" {
			headerHTML("DNS SCANNING", outputFile)
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
			appendWhere(ip.String(), "", outputFile)
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
		appendWhere(cname, "", outputFile)
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
			appendWhere(ns.Host, "", outputFile)
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
			appendWhere(mx.Host, "", outputFile)
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
			appendWhere(srv.Target, "", outputFile)
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
			appendWhere(txt, "", outputFile)
		}
	}
	if outputFile != "" {
		if outputFile[len(outputFile)-4:] == "html" {
			footerHTML(outputFile)
		}
	}
}

//asyncDir performs concurrent requests to the specified
//urls and prints the results
func asyncDir(urls []string, ignore []string, outputFile string, dirs map[string]Asset, mutex *sync.Mutex, plain bool) {
	ignoreBool := len(ignore) != 0
	var count int = 0
	var total int = len(urls)
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	limiter := make(chan string, 30) // Limits simultaneous requests
	wg := sync.WaitGroup{}           // Needed to not prematurely exit before all requests have been finished
	for i, domain := range urls {
		limiter <- domain
		wg.Add(1)
		if count%50 == 0 { // update counter
			if !plain {
				fmt.Fprint(os.Stdout, "\r \r")
			}
			printDirs(dirs, ignore, outputFile, mutex, plain)
		}
		if !plain && count%100 == 0 { // update counter
			fmt.Fprint(os.Stdout, "\r \r")
			fmt.Printf("%0.2f%% : %d / %d", percentage(count, total), count, total)
		}
		go func(i int, domain string) {
			defer wg.Done()
			defer func() { <-limiter }()
			resp, err := client.Get(domain)
			count++
			if err != nil {
				return
			}
			if ignoreBool {
				if ignoreResponse(resp.StatusCode, ignore) {
					return
				}
			}
			addDirs(domain, resp.Status, dirs, mutex)
			resp.Body.Close()
		}(i, domain)
	}
	printDirs(dirs, ignore, outputFile, mutex, plain)
	wg.Wait()
	printDirs(dirs, ignore, outputFile, mutex, plain)
}

//printSubs prints the results (only the resources not already printed).
//Also performs the checks based on the response status codes.
func printSubs(subs map[string]Asset, ignore []string, outputFile string, mutex *sync.Mutex, plain bool) {
	mutex.Lock()
	for domain, asset := range subs {
		if !asset.Printed {
			sub := Asset{
				Value:   asset.Value,
				Printed: true,
			}
			subs[domain] = sub
			var resp = asset.Value
			if !plain {
				fmt.Fprint(os.Stdout, "\r \r")
				if resp[:3] != "404" {
					subDomainFound := cleanProtocol(domain)
					fmt.Printf("[+]FOUND: %s ", subDomainFound)
					if string(resp[0]) == "2" {
						if outputFile != "" {
							appendWhere(domain, fmt.Sprint(resp), outputFile)
						}
						color.Green("%s\n", resp)
					} else {
						if outputFile != "" {
							appendWhere(domain, fmt.Sprint(resp), outputFile)
						}
						color.Red("%s\n", resp)
					}
				}
			} else {
				if resp[:3] != "404" {
					subDomainFound := cleanProtocol(domain)
					fmt.Printf("%s\n", subDomainFound)
					if string(resp[0]) == "2" {
						if outputFile != "" {
							appendWhere(domain, fmt.Sprint(resp), outputFile)
						}
					} else {
						if outputFile != "" {
							appendWhere(domain, fmt.Sprint(resp), outputFile)
						}
					}
				}
			}
		}
	}
	mutex.Unlock()
}

//printDirs prints the results (only the resources not already printed).
//Also performs the checks based on the response status codes.
func printDirs(dirs map[string]Asset, ignore []string, outputFile string, mutex *sync.Mutex, plain bool) {
	mutex.Lock()
	for domain, asset := range dirs {
		if !asset.Printed {
			dir := Asset{
				Value:   asset.Value,
				Printed: true,
			}
			dirs[domain] = dir
			var resp = asset.Value
			if !plain {
				if string(resp[0]) == "2" || string(resp[0]) == "3" {
					fmt.Fprint(os.Stdout, "\r \r")
					fmt.Printf("[+]FOUND: %s ", domain)
					color.Green("%s\n", resp)
					if outputFile != "" {
						appendWhere(domain, fmt.Sprint(resp), outputFile)
					}
				} else if (resp[:3] != "404") || string(resp[0]) == "5" {
					fmt.Fprint(os.Stdout, "\r \r")
					fmt.Printf("[+]FOUND: %s ", domain)
					color.Red("%s\n", resp)

					if outputFile != "" {
						appendWhere(domain, fmt.Sprint(resp), outputFile)
					}
				}
			} else {
				if resp[:3] != "404" {
					fmt.Printf("%s\n", domain)
					if outputFile != "" {
						appendWhere(domain, fmt.Sprint(resp), outputFile)
					}
				}
			}
		}
	}
	mutex.Unlock()
}

//cleanURL takes as input a string and it tries to
//remove the fragment and the query
//Example: https://example.com/directory1/?id=abcdef&path=ok#fragment1
//Output: https://example.com/directory1/
func cleanURL(input string) string {
	u, err := url.Parse(input)
	if err != nil {
		return input
	}
	return u.Scheme + "://" + u.Host + u.Path
}

//spawnCrawler spawn a crawler that search for
//links with this characteristic:
//- only http, https or ftp protocols allowed
func spawnCrawler(target string, ignore []string, dirs map[string]Asset, subs map[string]Asset, outputFile string, mutex *sync.Mutex, what string, plain bool) {
	ignoreBool := len(ignore) != 0
	c := colly.NewCollector()
	if what == "dir" {
		c = colly.NewCollector(
			colly.URLFilters(
				regexp.MustCompile("(http://|https://|ftp://)" + "(www.)?" + target + "*"),
			),
		)
	} else {
		c = colly.NewCollector(
			colly.URLFilters(
				regexp.MustCompile("(http://|https://|ftp://)" + "+." + target),
			),
		)
	}
	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if e.Attr("href") != "" {
			url := cleanURL(e.Attr("href"))
			if what == "dir" {
				if !presentDirs(url, dirs) && url != target {

					e.Request.Visit(url)
				}
			} else {
				if !presentSubs(url, subs) && url != target {

					e.Request.Visit(url)
				}
			}
		}
	})
	c.OnRequest(func(r *colly.Request) {
		var status = httpGet(r.URL.String())
		if ignoreBool {
			statusArray := strings.Split(status, " ")
			statusInt, err := strconv.Atoi(statusArray[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not get response status %s\n", status)
				os.Exit(1)
			}
			if !ignoreResponse(statusInt, ignore) {
				if what == "dir" {
					addDirs(r.URL.String(), status, dirs, mutex)
					printDirs(dirs, ignore, outputFile, mutex, plain)
				} else {
					addSubs(r.URL.String(), status, subs, mutex)
					printSubs(subs, ignore, outputFile, mutex, plain)
				}
			}
		} else {
			if what == "dir" {
				addDirs(r.URL.String(), status, dirs, mutex)
				printDirs(dirs, ignore, outputFile, mutex, plain)
			} else {
				addSubs(r.URL.String(), status, subs, mutex)
				printSubs(subs, ignore, outputFile, mutex, plain)
			}
		}
	})
	c.Visit("http://" + target)
}

//httpGet
func httpGet(input string) string {
	resp, err := http.Get(input)
	if err != nil {
		return "ERROR"
	}
	defer resp.Body.Close()
	return resp.Status
}

//addSubs adds the target found to the subs map
func addSubs(target string, value string, subs map[string]Asset, mutex *sync.Mutex) {
	sub := Asset{
		Value:   value,
		Printed: false,
	}
	mutex.Lock()
	if !presentSubs(target, subs) {
		subs[target] = sub
	}
	mutex.Unlock()
}

//addDirs adds the target found to the dirs map
func addDirs(target string, value string, dirs map[string]Asset, mutex *sync.Mutex) {
	dir := Asset{
		Value:   value,
		Printed: false,
	}
	mutex.Lock()
	if !presentDirs(target, dirs) {
		dirs[target] = dir
	}
	mutex.Unlock()
}

//presentSubs
func presentSubs(input string, subs map[string]Asset) bool {
	_, ok := subs[input]
	return ok
}

//presentDirs
func presentDirs(input string, dirs map[string]Asset) bool {
	_, ok := dirs[input]
	return ok
}
