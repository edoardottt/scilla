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
	"flag"
	"fmt"
	"log"
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
	banner1 := " =========================================\n"
	banner2 := " |                 _ _ _                 |\n"
	banner3 := " |        ___  ___(_) | | __ _           |\n"
	banner4 := " |       / __|/ __| | | |/ _` |          |\n"
	banner5 := " |       \\__ \\ (__| | | | (_| |          |\n"
	banner6 := " |       |___/\\___|_|_|_|\\__,_|          |\n"
	banner7 := " |                                       |\n"
	banner8 := " | https://github.com/edoardottt/scilla  |\n"
	banner9 := " | https://www.edoardoottavianelli.it    |\n"
	banner10 := " | Released under GPLv3 license          |\n"
	banner11 := " =========================================\n"
	banner := banner1 + banner2 + banner3 + banner4 + banner5 + banner6 + banner7 + banner8 + banner9 + banner10 + banner11
	fmt.Println(banner)
}

//help prints in stdout scilla usage
func help() {
	fmt.Println("Information Gathering tool - DNS / Subdomain / Ports / Directories enumeration")
	fmt.Println("")
	fmt.Println("usage: scilla subcommand { options }")
	fmt.Println("")
	fmt.Println("	Available subcommands:")
	fmt.Println("		- dns -target [-o output-format] <target (URL)> REQUIRED")
	fmt.Println("		- port [-p <start-end>] [-o output-format] -target <target (URL/IP)> REQUIRED")
	fmt.Println("		- subdomain [-w wordlist]")
	fmt.Println("					[-o output-format]")
	fmt.Println("					[-i ignore status codes]")
	fmt.Println("					[-c use also a web crawler (SLOWER)]")
	fmt.Println("					-target <target (URL)> REQUIRED")
	fmt.Println("		- dir [-w wordlist]")
	fmt.Println("			  [-o output-format]")
	fmt.Println("			  [-i ignore status codes]")
	fmt.Println("			  [-c use also a web crawler (SLOWER)]")
	fmt.Println("			  -target <target (URL)> REQUIRED")
	fmt.Println("		- report [-p <start-end>]")
	fmt.Println("				 [-ws subdomains wordlist]")
	fmt.Println("			 	 [-wd directories wordlist]")
	fmt.Println("			 	 [-o output-format]")
	fmt.Println("			 	 [-id ignore status codes in directories scanning]")
	fmt.Println("			 	 [-is ignore status codes in subdomains scanning]")
	fmt.Println("			 	 [-cd use also a web crawler for directories scanning (SLOWER)]")
	fmt.Println("			 	 [-cs use also a web crawler for subdomains scanning (SLOWER)]")
	fmt.Println("			 	 -target <target (URL/IP)> REQUIRED")
	fmt.Println("		- help")
	fmt.Println("		- examples")
	fmt.Println()
}

//examples prints some examples
func examples() {
	fmt.Println("	Examples:")
	fmt.Println("		- scilla dns -target target.domain")
	fmt.Println("		- scilla dns -target -o txt target.domain")
	fmt.Println("		- scilla dns -target -o html target.domain")
	fmt.Println()
	fmt.Println("		- scilla subdomain -target target.domain")
	fmt.Println("		- scilla subdomain -w wordlist.txt -target target.domain")
	fmt.Println("		- scilla subdomain -o txt -target target.domain")
	fmt.Println("		- scilla subdomain -o html -target target.domain")
	fmt.Println("		- scilla subdomain -i 400 -target target.domain")
	fmt.Println("		- scilla subdomain -i 4** -target target.domain")
	fmt.Println("		- scilla subdomain -c -target target.domain")
	fmt.Println()
	fmt.Println("		- scilla port -p -450 -target target.domain")
	fmt.Println("		- scilla port -p 90- -target target.domain")
	fmt.Println("		- scilla port -p 10-1000 -target target.domain")
	fmt.Println("		- scilla port -o txt -target target.domain")
	fmt.Println("		- scilla port -o html -target target.domain")
	fmt.Println()
	fmt.Println("		- scilla dir -target target.domain")
	fmt.Println("		- scilla dir -o txt -target target.domain")
	fmt.Println("		- scilla dir -o html -target target.domain")
	fmt.Println("		- scilla dir -w wordlist.txt -target target.domain")
	fmt.Println("		- scilla dir -i 500,401 -target target.domain")
	fmt.Println("		- scilla dir -i 5**,401 -target target.domain")
	fmt.Println("		- scilla dir -c -target target.domain")
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
	fmt.Println("")
}

//main function
func main() {
	intro()
	input := readArgs()
	// common assets found (only subdomain and dir)
	subs := make(map[string]Asset)
	dirs := make(map[string]Asset)
	execute(input, subs, dirs)
}

//Asset gives information about the asset found
type Asset struct {
	Value   string
	Printed bool
}

//execute reads inputs and starts the correct procedure
func execute(input Input, subs map[string]Asset, dirs map[string]Asset) {

	var mutex = &sync.Mutex{}
	if input.ReportTarget != "" {

		target := cleanProtocol(input.ReportTarget)
		var targetIP string
		fmt.Printf("target: %s\n", target)
		fmt.Println("=============== FULL REPORT ===============")
		outputFile := ""
		if input.ReportOutput != "" {
			outputFile = createOutputFile(input.ReportTarget, input.ReportOutput)
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
			go spawnCrawler(target, input.ReportIgnoreSub, dirs, subs, outputFile, mutex, "sub")
		}
		strings1 = createSubdomains(input.ReportWordSub, target)
		asyncGet(strings1, input.ReportIgnoreSub, outputFile, subs, mutex)
		if outputFile != "" {
			if outputFile[len(outputFile)-4:] == "html" {
				footerHTML(outputFile)
			}
		}

		if targetIP != "" {
			target = targetIP
		}
		fmt.Println("=============== PORT SCANNING ===============")
		asyncPort(input.StartPort, input.EndPort, target, outputFile)

		fmt.Println("=============== DNS SCANNING ===============")
		lookupDNS(target, outputFile)

		fmt.Println("=============== DIRECTORIES SCANNING ===============")
		var strings2 []string
		strings2 = createUrls(input.ReportWordDir, target)
		if outputFile != "" {
			if outputFile[len(outputFile)-4:] == "html" {
				headerHTML("DIRECTORY SCANNING", outputFile)
			}
		}
		if input.ReportCrawlerDir {
			go spawnCrawler(target, input.ReportIgnoreDir, dirs, subs, outputFile, mutex, "dir")
		}
		asyncDir(strings2, input.ReportIgnoreDir, outputFile, dirs, mutex)
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
		target := cleanProtocol(input.DNSTarget)
		// change from ip to Hostname
		if isIP(target) {
			target = ipToHostname(target)
		}
		fmt.Printf("target: %s\n", target)
		fmt.Println("=============== DNS SCANNING ===============")
		outputFile := ""
		if input.DNSOutput != "" {
			outputFile = createOutputFile(input.DNSTarget, input.DNSOutput)

			if outputFile[len(outputFile)-4:] == "html" {
				bannerHTML(input.DNSTarget, outputFile)
			}
		}
		lookupDNS(target, outputFile)
		if input.DNSOutput != "" {
			if outputFile[len(outputFile)-4:] == "html" {
				bannerFooterHTML(outputFile)
			}
		}
	}

	if input.SubdomainTarget != "" {

		target := cleanProtocol(input.SubdomainTarget)
		// change from ip to Hostname
		if isIP(target) {
			target = ipToHostname(target)
		}
		fmt.Printf("target: %s\n", target)
		fmt.Println("=============== SUBDOMAINS SCANNING ===============")
		outputFile := ""
		if input.SubdomainOutput != "" {
			outputFile = createOutputFile(input.SubdomainTarget, input.SubdomainOutput)
			if outputFile[len(outputFile)-4:] == "html" {
				bannerHTML(input.SubdomainTarget, outputFile)
			}
		}
		var strings1 []string
		strings1 = createSubdomains(input.SubdomainWord, target)
		if outputFile != "" {
			if outputFile[len(outputFile)-4:] == "html" {
				headerHTML("SUBDOMAIN SCANNING", outputFile)
			}
		}
		if input.SubdomainCrawler {
			go spawnCrawler(target, input.SubdomainIgnore, dirs, subs, outputFile, mutex, "sub")
		}
		asyncGet(strings1, input.SubdomainIgnore, outputFile, subs, mutex)
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

		target := cleanProtocol(input.DirTarget)
		fmt.Printf("target: %s\n", target)
		fmt.Println("=============== DIRECTORIES SCANNING ===============")
		outputFile := ""
		if input.DirOutput != "" {
			outputFile = createOutputFile(input.DirTarget, input.DirOutput)

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
			spawnCrawler(target, input.ReportIgnoreDir, dirs, subs, outputFile, mutex, "dir")
		}
		asyncDir(strings2, input.DirIgnore, outputFile, dirs, mutex)
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
		target := input.PortTarget
		if isURL(target) {
			target = cleanProtocol(input.PortTarget)
		}
		outputFile := ""
		if input.PortOutput != "" {
			outputFile = createOutputFile(input.PortTarget, input.PortOutput)
			if outputFile[len(outputFile)-4:] == "html" {
				bannerHTML(input.PortTarget, outputFile)
			}
		}
		fmt.Printf("target: %s\n", target)
		fmt.Println("=============== PORT SCANNING ===============")
		asyncPort(input.StartPort, input.EndPort, target, outputFile)

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

// output formats accepted
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
	ReportTarget     string
	ReportWordDir    string
	ReportWordSub    string
	ReportOutput     string
	ReportIgnoreDir  []string
	ReportIgnoreSub  []string
	ReportCrawlerDir bool
	ReportCrawlerSub bool
	DNSTarget        string
	DNSOutput        string
	SubdomainTarget  string
	SubdomainWord    string
	SubdomainOutput  string
	SubdomainIgnore  []string
	SubdomainCrawler bool
	DirTarget        string
	DirWord          string
	DirOutput        string
	DirIgnore        []string
	DirCrawler       bool
	PortTarget       string
	PortOutput       string
	StartPort        int
	EndPort          int
}

//readArgs reads arguments/options from stdin
// Subcommands:
// 		report     ==> Full report
// 		dns        ==> Dns records enumeration
// 		subdomains ==> Subdomains enumeration
// 		port	   ==> ports enumeration
//		dir		   ==> directiories enumeration
// 		help       ==> doc
//		examples   ==> examples
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

	// dns subcommand flag pointers
	dnsTargetPtr := dnsCommand.String("target", "", "Target {URL/IP} (Required)")

	// dns subcommand flag pointers
	dnsOutputPtr := dnsCommand.String("o", "", "output format (txt/html)")

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

	// port subcommand flag pointers
	portTargetPtr := portCommand.String("target", "", "Target {URL/IP} (Required)")

	// report subcommand flag pointers
	portOutputPtr := portCommand.String("o", "", "output format (txt/html)")

	portsPtr := portCommand.String("p", "", "ports range <start-end>")
	// Default ports
	StartPort := 1
	EndPort := 65535

	// Verify that a subcommand has been provided
	// os.Arg[0] is the main command
	// os.Arg[1] will be the subcommand
	if len(os.Args) < 2 {
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
		helpCommand.Parse(os.Args[2:])
	case "examples":
		examplesCommand.Parse(os.Args[2:])
	default:
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
		if *reportPortsPtr != "" {
			portsRange := string(*reportPortsPtr)
			StartPort, EndPort = checkPortsRange(portsRange, StartPort, EndPort)
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
		if *portsPtr != "" {
			portsRange := string(*portsPtr)
			StartPort, EndPort = checkPortsRange(portsRange, StartPort, EndPort)
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
		*dnsTargetPtr,
		*dnsOutputPtr,
		*subdomainTargetPtr,
		*subdomainWordlistPtr,
		*subdomainOutputPtr,
		subdomainIgnore,
		*subdomainCrawlerPtr,
		*dirTargetPtr,
		*dirWordlistPtr,
		*dirOutputPtr,
		dirIgnore,
		*dirCrawlerPtr,
		*portTargetPtr,
		*portOutputPtr,
		StartPort,
		EndPort,
	}
	return result
}

//checkIgnore
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

//replaceBadCharacterOutput
func replaceBadCharacterOutput(input string) string {
	result := strings.ReplaceAll(input, "/", "-")
	return result
}

// Create Output Folder
func createOutputFolder() {
	//Create a folder/directory at a full qualified path
	err := os.Mkdir("output", 0755)
	if err != nil {
		fmt.Println("Can't create output folder.")
		os.Exit(1)
	}
}

// Create Output File
func createOutputFile(target string, format string) string {
	target = replaceBadCharacterOutput(target)
	filename := "output" + "/" + target + "." + format
	_, err := os.Stat(filename)

	if os.IsNotExist(err) {
		if _, err := os.Stat("output/"); os.IsNotExist(err) {
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
func appendDir(domain string, dir string) string {
	return "http://" + domain + "/" + dir + "/"
}

//readDict scan all the possible subdomains from file
func readDict(inputFile string) []string {
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
			subs = readDict("lists/subdomains.txt")
		} else { // linux
			subs = readDict("/usr/bin/lists/subdomains.txt")
		}
	} else {
		subs = readDict(filename)
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
			dirs = readDict("lists/dirs.txt")
		} else { // linux
			dirs = readDict("/usr/bin/lists/dirs.txt")
		}
	} else {
		dirs = readDict(filename)
	}
	result := []string{}
	for _, dir := range dirs {
		path := appendDir(url, dir)
		result = append(result, path)
	}
	return result
}

//appendOutputToTxt
func appendOutputToTxt(output string, filename string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	if _, err := file.WriteString(output + "\n"); err != nil {
		log.Fatal(err)
	}
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
}

//headerHtml
func headerHTML(header string, filename string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	if _, err := file.WriteString("<h3>" + header + "</h3><ul>"); err != nil {
		log.Fatal(err)
	}
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
}

//bannerFooterHTML
func bannerFooterHTML(filename string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	file.WriteString("<div style='" + "background-color:#4adeff;color:white" + "'>")
	file.WriteString("<ul><li><a href='" + "https://github.com/edoardottt/scilla'" + ">Contribute to scilla</a></li>")
	file.WriteString("<li>Released under <a href='" + "http://www.gnu.org/licenses/gpl-3.0.html'" + ">GPLv3 License</a></li></ul></div>")
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
func asyncGet(urls []string, ignore []string, outputFile string, subs map[string]Asset, mutex *sync.Mutex) {
	ignoreBool := len(ignore) != 0
	var count int = 0
	var total int = len(urls)
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	limiter := make(chan string, 50) // Limits simultaneous requests
	wg := sync.WaitGroup{}           // Needed to not prematurely exit before all requests have been finished

	for i, domain := range urls {
		limiter <- domain
		wg.Add(1)
		if count%50 == 0 { // update counter
			fmt.Fprint(os.Stdout, "\r \r")
			printSubs(subs, ignore, outputFile, mutex)
		}
		if count%100 == 0 { // update counter
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
	printSubs(subs, ignore, outputFile, mutex)
	wg.Wait()
	printSubs(subs, ignore, outputFile, mutex)
}

//appendWhere
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
func asyncPort(StartingPort int, EndingPort int, host string, outputFile string) {
	var count int = 0
	var total int = EndingPort - StartingPort
	limiter := make(chan string, 200) // Limits simultaneous requests
	wg := sync.WaitGroup{}            // Needed to not prematurely exit before all requests have been finished
	if outputFile != "" {
		if outputFile[len(outputFile)-4:] == "html" {
			headerHTML("PORT SCANNING", outputFile)
		}
	}
	for port := StartingPort; port <= EndingPort; port++ {
		wg.Add(1)
		portStr := fmt.Sprint(port)
		limiter <- portStr
		if count%100 == 0 { // update counter
			fmt.Fprint(os.Stdout, "\r \r")
			fmt.Printf("%0.2f%% : %d / %d", percentage(count, total), count, total)
		}
		go func(portStr string, host string) {
			defer func() { <-limiter }()
			defer wg.Done()
			resp := isOpenPort(host, portStr)
			count++
			if resp {
				fmt.Fprint(os.Stdout, "\r \r")
				fmt.Printf("[+]FOUND: %s ", host)
				color.Green("%s\n", portStr)
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

//lookupDNS prints the DNS servers for the inputted domain
func lookupDNS(domain string, outputFile string) {
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
		fmt.Printf("[+]FOUND %s IN A: ", domain)
		color.Green("%s\n", ip.String())
		if outputFile != "" {
			appendWhere(ip.String(), "", outputFile)
		}
	}
	// -- CNAME RECORD --
	cname, err := net.LookupCNAME(domain)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get CNAME: %v\n", err)
	}
	fmt.Printf("[+]FOUND %s IN CNAME: ", domain)
	color.Green("%s\n", cname)
	if outputFile != "" {
		appendWhere(cname, "", outputFile)
	}
	// -- NS RECORDS --
	nameserver, err := net.LookupNS(domain)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get NSs: %v\n", err)
	}
	for _, ns := range nameserver {
		fmt.Printf("[+]FOUND %s IN NS: ", domain)
		color.Green("%s\n", ns.Host)
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
		fmt.Printf("[+]FOUND %s IN MX: ", domain)
		color.Green("%s %v\n", mx.Host, mx.Pref)
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
		fmt.Printf("[+]FOUND %s IN SRV: ", domain)
		color.Green("%v:%v:%d:%d\n", srv.Target, srv.Port, srv.Priority, srv.Weight)
		if outputFile != "" {
			appendWhere(srv.Target, "", outputFile)
		}
	}
	// -- TXT RECORDS --
	txtrecords, _ := net.LookupTXT(domain)
	for _, txt := range txtrecords {
		fmt.Printf("[+]FOUND %s IN TXT: ", domain)
		color.Green("%s\n", txt)
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
func asyncDir(urls []string, ignore []string, outputFile string, dirs map[string]Asset, mutex *sync.Mutex) {
	ignoreBool := len(ignore) != 0
	var count int = 0
	var total int = len(urls)
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	limiter := make(chan string, 50) // Limits simultaneous requests
	wg := sync.WaitGroup{}           // Needed to not prematurely exit before all requests have been finished
	for i, domain := range urls {
		limiter <- domain
		wg.Add(1)
		if count%50 == 0 { // update counter
			fmt.Fprint(os.Stdout, "\r \r")
			printDirs(dirs, ignore, outputFile, mutex)
		}
		if count%100 == 0 { // update counter
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
	printDirs(dirs, ignore, outputFile, mutex)
	wg.Wait()
	printDirs(dirs, ignore, outputFile, mutex)
}

//printSubs
func printSubs(subs map[string]Asset, ignore []string, outputFile string, mutex *sync.Mutex) {
	mutex.Lock()
	for domain, asset := range subs {
		if !asset.Printed {
			sub := Asset{
				Value:   asset.Value,
				Printed: true,
			}
			subs[domain] = sub
			var resp = asset.Value
			fmt.Fprint(os.Stdout, "\r \r")
			subDomainFound := cleanProtocol(domain)
			fmt.Printf("[+]FOUND: %s ", subDomainFound)
			if string(resp[0]) == "2" {
				if outputFile != "" {
					appendWhere(domain, fmt.Sprint(resp), outputFile)
				}
				color.Green("%s\n", resp)
			} else {
				if resp != "404" {
					if outputFile != "" {
						appendWhere(domain, fmt.Sprint(resp), outputFile)
					}
				}
				color.Red("%s\n", resp)
			}
		}
	}
	mutex.Unlock()
}

//printDirs
func printDirs(dirs map[string]Asset, ignore []string, outputFile string, mutex *sync.Mutex) {
	mutex.Lock()
	for domain, asset := range dirs {
		if !asset.Printed {
			dir := Asset{
				Value:   asset.Value,
				Printed: true,
			}
			dirs[domain] = dir
			var resp = asset.Value
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
		}
	}
	mutex.Unlock()
}

//spawnCrawler
func spawnCrawler(target string, ignore []string, dirs map[string]Asset, subs map[string]Asset, outputFile string, mutex *sync.Mutex, what string) {
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
			if what == "dir" {
				if !presentDirs(e.Attr("href"), dirs) && e.Attr("href") != target {

					e.Request.Visit(e.Attr("href"))
				}
			} else {
				if !presentSubs(e.Attr("href"), subs) && e.Attr("href") != target {

					e.Request.Visit(e.Attr("href"))
				}
			}
		}
	})
	c.OnRequest(func(r *colly.Request) {
		var status = httpGet(r.URL.String())
		if what == "dir" {
			addDirs(r.URL.String(), status, dirs, mutex)
			printDirs(dirs, ignore, outputFile, mutex)
		} else {
			addSubs(r.URL.String(), status, subs, mutex)
			printSubs(subs, ignore, outputFile, mutex)
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

//addSubs
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

//addDirs
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
