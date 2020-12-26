/*
          _ _ _
 ___  ___(_) | | __ _
/ __|/ __| | | |/ _` |
\__ \ (__| | | | (_| |
|___/\___|_|_|_|\__,_|

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
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gookit/color"
)

//intro prints the banner when the program is started
func intro() {
	banner := `
	*****************************************
	*                 _ _ _                 *
	*        ___  ___(_) | | __ _           *
	*       / __|/ __| | | |/ _  |          *
	*       \__ \ (__| | | | (_| |          *
	*       |___/\___|_|_|_|\__,_|          *
	*                                       *
	* https://github.com/edoardottt/scilla  *
	* https://www.edoardoottavianelli.it    *
	*                                       *
	*****************************************`
	fmt.Println(banner)
}

//help prints in stdout scilla usage
func help() {
	fmt.Println(`
Information Gathering tool - DNS / Subdomain / Ports / Directories enumeration
usage: scilla [subcommand] { options }
Available subcommands:
  - dns { -target <target (URL)> REQUIRED}
  - subdomain { [-w wordlist] -target <target (URL)> REQUIRED}
  - port { [-p <start-end>] -target <target (URL/IP)> REQUIRED}
  - dir { [-w wordlist] -target <target (URL/IP)> REQUIRED}
  - report { [-p <start-end>] -target <target (URL/IP)> REQUIRED}
  - help
Examples:
  - scilla dns -target target.domain
  - scilla subdomain -target target.domain
  - scilla port -p -450 -target target.domain
  - scilla port -p 90- -target target.domain
  - scilla report -p 80 -target target.domain
  - scilla report -p 50-200 -target target.domain
  - scilla dir -w directs.txt -target target.domain`)
}

//Input struct contains the input parameters
type Input struct {
	ReportTarget    string
	ReportDirWord   string
	ReportDNSWord   string
	DNSTarget       string
	SubdomainTarget string
	SubdomainWord   string
	DirTarget       string
	DirWord         string
	PortTarget      string
	StartPort       int
	EndPort         int
}

var (
	subs, dirs             []string
	subsLength, dirsLength int

	green  = color.FgGreen.Render
	red    = color.FgRed.Render
	yellow = color.FgYellow.Render
	cyan   = color.FgCyan.Render
)

func main() {
	intro()

	// Loads the deafult wordlist files
	if runtime.GOOS == "windows" {
		subs = readDict("lists/subdomains.txt")
		dirs = readDict("lists/dirs.txt")
	} else { // linux
		subs = readDict("/usr/bin/lists/subdomains.txt")
		dirs = readDict("/usr/bin/lists/dirs.txt")
	}

	input := readArgs()
	execute(input)
}

//execute reads inputs and starts the correct procedure
func execute(input Input) {
	if input.ReportTarget != "" {
		fmt.Printf("\n=============== FULL REPORT ===============\n")
		target := input.ReportTarget

		fmt.Printf("\n=============== SUBDOMAINS ===============\n")
		domains := createSubdomains(input.ReportDNSWord, cleanProtocol(target))
		fmt.Printf("%v Target: %s\n\n", cyan("[#]"), target)
		asyncResolve(domains)

		fmt.Printf("\n=============== PORT SCANNING ===============\n")
		fmt.Printf("%v Target: %s\n\n", cyan("[#]"), target)
		asyncPort(input.StartPort, input.EndPort, cleanProtocol(target))

		fmt.Printf("\n=============== DNS SCANNING ===============\n")
		fmt.Printf("%v Target: %s\n\n", cyan("[#]"), target)
		lookupDNS(cleanProtocol(target))

		fmt.Printf("\n=============== DIRECTORIES ===============\n")
		urls := createUrls(input.ReportDirWord, target)
		fmt.Printf("%v Target: %s\n\n", cyan("[#]"), target)
		asyncReq(urls)

	}

	if input.DNSTarget != "" {
		target := cleanProtocol(input.DNSTarget)
		fmt.Printf("\n=============== DNS SCANNING ===============\n")
		fmt.Printf("%v Target: %s\n\n", cyan("[#]"), target)
		lookupDNS(target)
	}

	if input.SubdomainTarget != "" {
		target := cleanProtocol(input.SubdomainTarget)
		fmt.Printf("\n=============== SUBDOMAINS ===============\n")
		fmt.Printf("%v Target: %s\n\n", cyan("[#]"), target)
		domains := createSubdomains(input.SubdomainWord, target)
		asyncResolve(domains)
	}

	if input.DirTarget != "" {
		target := input.DirTarget
		fmt.Printf("\n=============== DIRECTORIES ===============\n")
		fmt.Printf("%v Target: %s\n\n", cyan("[#]"), target)
		urls := createUrls(input.DirWord, target)
		asyncReq(urls)
	}

	if input.PortTarget != "" {

		target := input.PortTarget
		fmt.Printf("\n=============== PORT SCANNING ===============\n")
		fmt.Printf("%v Target: %s\n\n", cyan("[#]"), target)
		asyncPort(input.StartPort, input.EndPort, cleanProtocol(target))
	}
}

//readArgs reads arguments/options from stdin
// Subcommands:
// 		report     ==> Full report
// 		dns        ==> Dns records enumeration
// 		subdomains ==> Subdomains enumeration
// 		port	   ==> Ports enumeration
//		dir		   ==> Directiories enumeration
// 		help       ==> Help
func readArgs() Input {
	reportCommand := flag.NewFlagSet("report", flag.ExitOnError)
	dnsCommand := flag.NewFlagSet("dns", flag.ExitOnError)
	subdomainCommand := flag.NewFlagSet("subdomain", flag.ExitOnError)
	portCommand := flag.NewFlagSet("port", flag.ExitOnError)
	dirCommand := flag.NewFlagSet("dir", flag.ExitOnError)
	helpCommand := flag.NewFlagSet("help", flag.ExitOnError)

	// report subcommand flag pointers
	reportTargetPtr := reportCommand.String("target", "", "Target {URL/IP} (Required)")
	portsReportPtr := reportCommand.String("p", "", "ports range <start-end>")
	wordlistsDirReportPtr := reportCommand.String("wdir", "", "wordlist to use for directory brute force (default enabled)")
	wordlistsDNSReportPtr := reportCommand.String("wdns", "", "wordlist to use for dns discovery (default enabled)")

	// dns subcommand flag pointers
	dnsTargetPtr := dnsCommand.String("target", "", "Target {URL/IP} (Required)")
	subdomainTargetPtr := subdomainCommand.String("target", "", "Target {URL/IP} (Required)")
	subdomainWordlistPtr := subdomainCommand.String("w", "", "wordlist to use (default enabled)")

	// dir subcommand flag pointers
	dirTargetPtr := dirCommand.String("target", "", "Target {URL/IP} (Required)")
	dirWordlistPtr := dirCommand.String("w", "", "wordlist to use (default enabled)")

	// port subcommand flag pointers
	portTargetPtr := portCommand.String("target", "", "Target {URL/IP} (Required)")
	portsPtr := portCommand.String("p", "", "ports range <start-end>")

	// Default ports
	StartPort := 1
	EndPort := 65535

	if len(os.Args) < 2 {
		fmt.Println("[ERROR] subcommand is required.")
		fmt.Println("	Type: scilla help")
		os.Exit(1)
	}

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
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	if reportCommand.Parsed() {
		if *reportTargetPtr == "" || !isURL(*reportTargetPtr) {
			reportCommand.PrintDefaults()
			fmt.Println("The inputted target is not valid or didn't provided")
			os.Exit(1)
		}

		if *portsReportPtr != "" {
			portsRange := string(*portsReportPtr)
			StartPort, EndPort = checkPortsRange(portsRange, StartPort, EndPort)
		}
	}

	if dnsCommand.Parsed() {
		if *dnsTargetPtr == "" {
			dnsCommand.PrintDefaults()
			fmt.Println("The inputted target is not valid or didn't provided")
			os.Exit(1)
		}
	}

	if subdomainCommand.Parsed() {
		if *subdomainTargetPtr == "" {
			subdomainCommand.PrintDefaults()
			fmt.Println("The inputted target is not valid or didn't provided")
			os.Exit(1)
		}
	}

	if portCommand.Parsed() {
		if *portTargetPtr == "" {
			portCommand.PrintDefaults()
			fmt.Println("The inputted target is not valid or didn't provided")
			os.Exit(1)
		}

		if *portsPtr != "" {
			portsRange := string(*portsPtr)
			StartPort, EndPort = checkPortsRange(portsRange, StartPort, EndPort)
		}
	}

	if dirCommand.Parsed() {
		// Required Flags
		if *dirTargetPtr == "" || !isURL(*dirTargetPtr) {
			dirCommand.PrintDefaults()
			fmt.Println("The inputted target is not valid or didn't provided.")
			os.Exit(1)
		}
	}

	if helpCommand.Parsed() {
		help()
		os.Exit(0)
	}

	result := Input{
		*reportTargetPtr,
		*wordlistsDirReportPtr,
		*wordlistsDNSReportPtr,
		*dnsTargetPtr,
		*subdomainTargetPtr,
		*subdomainWordlistPtr,
		*dirTargetPtr,
		*dirWordlistPtr,
		*portTargetPtr,
		StartPort,
		EndPort,
	}

	return result
}

//cleanProtocol remove from the url the protocol scheme
// (http - https - tls)
func cleanProtocol(target string) string {
	regex := regexp.MustCompile(`^((https?)|(tls))://`)
	return regex.ReplaceAllString(target, "")
}

//checkPortsRange checks the basic rules to
//be valid and then returns the starting port and the ending port.
func checkPortsRange(portsRange string, StartPort, EndPort int) (int, int) {
	delimiter := byte('-')
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
		maybeStart, err := strconv.Atoi(portsRange)
		if err != nil {
			fmt.Println("The inputted port range is not valid.")
			os.Exit(1)
		}
		if maybeStart > 0 && maybeStart < EndPort {
			StartPort = maybeStart
		}
	} else if !strings.Contains(portsRange, string(delimiter)) {
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

//isURL checks if the inputted Url is valid
func isURL(target string) bool {
	return regexp.MustCompile(`^https?://`).MatchString(target)
}

//readDict scan all the possible subdomains from file
func readDict(inputFile string) []string {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("failed to open %s ", inputFile)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	return text
}

//createSubdomains returns a list of subdomains
//from the default file lists/subdomains.txt.
func createSubdomains(filename string, url string) []string {
	if filename != "" {
		subs = readDict(filename)
	}

	result := make([]string, 0)
	for _, sub := range subs {
		result = append(result, sub+"."+url)
	}
	return result
}

//createUrls returns a list of directories
//from the default file lists/dirs.txt.
func createUrls(filename string, url string) []string {
	if filename != "" {
		dirs = readDict(filename)
	}
	result := make([]string, 0)
	for _, dir := range dirs {
		result = append(result, url+"/"+dir)
	}
	return result
}

func asyncResolve(hosts []string) {
	var wg sync.WaitGroup

	limiter := make(chan string, 80)
	nxdomain, err := net.LookupHost("nonexistingsubdomain.nonexistingdomain.nonexistingtld")
	if err != nil {
		log.Fatal(err)
	}

	for _, host := range hosts {
		wg.Add(1)
		limiter <- host
		go func(nxdomain, host string) {
			defer func() { <-limiter }()
			defer wg.Done()
			ips, _ := net.LookupHost(host)

			for _, ip := range ips {
				if ip != nxdomain {
					fmt.Printf("%v FOUND: %s ==> %v\n", green("[+]"), host, ip)
				}
			}
		}(nxdomain[0], host)
	}
	wg.Wait()
}

//isOpenPort scans if a port is open
func isOpenPort(host string, port, timeout int) bool {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", addr, time.Second*time.Duration(timeout))
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
func asyncPort(StartingPort, EndingPort int, host string) {
	limiter := make(chan int, 300) // Limits simultaneous requests
	var wg sync.WaitGroup

	for port := StartingPort; port <= EndingPort; port++ {
		wg.Add(1)
		limiter <- port

		go func(port int, host string) {
			defer func() { <-limiter }()
			defer wg.Done()

			resp := isOpenPort(host, port, 1)
			if resp {
				fmt.Printf("%v FOUND: %s %d\n", green("[+]"), host, port)
			}
		}(port, host)
	}
	wg.Wait()
}

//lookupDNS prints the DNS servers for the inputted domain
func lookupDNS(domain string) {

	// -- A RECORDS --
	ips, err := net.LookupIP(domain)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v Could not get A/AAAA Records: %v\n", red("[-]"), err)
	}
	for _, ip := range ips {
		fmt.Printf("%v FOUND %s IN A/AAAA: %s\n", green("[+]"), domain, ip.String())
	}

	// -- CNAME RECORD --
	cname, err := net.LookupCNAME(domain)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v Could not get CNAME Records: %v\n", red("[-]"), err)
	}
	fmt.Printf("%v FOUND %s IN CNAME: %s\n", green("[+]"), domain, cname)

	// -- NS RECORDS --
	nameserver, err := net.LookupNS(domain)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v Could not get NS Records: %v\n", red("[-]"), err)
	}
	for _, ns := range nameserver {
		fmt.Printf("%v FOUND %s IN NS: %s\n", green("[+]"), domain, ns.Host)
	}

	// -- MX RECORDS --
	mxrecords, err := net.LookupMX(domain)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v Could not get MX Records: %v\n", red("[-]"), err)
	}
	for _, mx := range mxrecords {
		fmt.Printf("%v FOUND %s IN MX: %s %v\n", green("[+]"), domain, mx.Host, mx.Pref)
	}

	// -- SRV SERVICE --
	_, srvs, err := net.LookupSRV("xmpp-server", "tcp", domain)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v Could not get SRV Records: %v\n", red("[-]"), err)
	}
	for _, srv := range srvs {
		fmt.Printf("%v FOUND %s IN SRV: %v:%v:%d:%d\n", green("[+]"), domain, srv.Target, srv.Port, srv.Priority, srv.Weight)
	}

	// -- TXT RECORDS --
	txtrecords, err := net.LookupTXT(domain)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v Could not get TXT Records: %v\n", red("[-]"), err)
	}

	for _, txt := range txtrecords {
		fmt.Printf("%v FOUND %s IN TXT: %s\n", green("[+]"), domain, txt)
	}
}

//asyncDir performs concurrent requests to the specified
//urls and prints the results
func asyncReq(urls []string) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: time.Second * 10,
	}

	limiter := make(chan string, 80) // Limits simultaneous requests

	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		limiter <- url

		go func(url string) {
			defer func() { <-limiter }()
			defer wg.Done()

			resp, err := client.Get(url)
			if err != nil {
				log.Println(err)
			}

			if resp.StatusCode == 200 {
				fmt.Printf("%v FOUND: %s >> %s\n", green("[+]"), url, resp.Status)
			} else if resp.StatusCode != 404 {
				fmt.Printf("%v FOUND: %s >> %s\n", yellow("[!]"), url, resp.Status)
			}
		}(url)
	}

	wg.Wait()
}
