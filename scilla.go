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
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/fatih/color"
)

//intro prints the banner when the program is started
func intro() {
	banner1 := " *****************************************\n"
	banner2 := " *                 _ _ _                 *\n"
	banner3 := " *        ___  ___(_) | | __ _           *\n"
	banner4 := " *       / __|/ __| | | |/ _` |          *\n"
	banner5 := " *       \\__ \\ (__| | | | (_| |          *\n"
	banner6 := " *       |___/\\___|_|_|_|\\__,_|          *\n"
	banner7 := " *                                       *\n"
	banner8 := " * https://github.com/edoardottt/scilla  *\n"
	banner9 := " * https://www.edoardoottavianelli.it    *\n"
	banner10 := " *                                       *\n"
	banner11 := " *****************************************\n"
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
	fmt.Println("		- subdomain [-w wordlist] [-o output-format] -target <target (URL)> REQUIRED")
	fmt.Println("		- port [-p <start-end>] [-o output-format] -target <target (URL/IP)> REQUIRED")
	fmt.Println("		- dir [-w wordlist] [-o output-format] -target <target (URL/IP)> REQUIRED")
	fmt.Println("		- report [-p <start-end>] [-w wordlist] [-o output-format] -target <target (URL/IP)> REQUIRED")
	fmt.Println("		- help")
	fmt.Println("		- examples")
	fmt.Println()
}

//examples prints some examples
func examples() {
	fmt.Println("	Examples:")
	fmt.Println("		- scilla dns -target target.domain")
	fmt.Println("		- scilla dns -target -o txt target.domain")
	fmt.Println()
	fmt.Println("		- scilla subdomain -target target.domain")
	fmt.Println("		- scilla subdomain -w wordlist.txt -target target.domain")
	fmt.Println("		- scilla subdomain -o txt -target target.domain")
	fmt.Println()
	fmt.Println("		- scilla port -p -450 -target target.domain")
	fmt.Println("		- scilla port -p 90- -target target.domain")
	fmt.Println("		- scilla port -p 10-1000 -target target.domain")
	fmt.Println("		- scilla port -o txt -target target.domain")
	fmt.Println()
	fmt.Println("		- scilla dir -target target.domain")
	fmt.Println("		- scilla dir -o txt -target target.domain")
	fmt.Println("		- scilla dir -w wordlist.txt -target target.domain")
	fmt.Println()
	fmt.Println("		- scilla report -p 80 -target target.domain")
	fmt.Println("		- scilla report -p 50-200 -target target.domain")
	fmt.Println("		- scilla report -wd dirs.txt -target target.domain")
	fmt.Println("		- scilla report -ws subdomains.txt -target target.domain")
	fmt.Println("")
}

//main function
func main() {
	intro()
	input := readArgs()
	execute(input)
}

//execute reads inputs and starts the correct procedure
func execute(input Input) {

	if input.ReportTarget != "" {

		fmt.Println("=============== FULL REPORT ===============")
		target := cleanProtocol(input.ReportTarget)
		outputFile := ""
		if input.ReportOutput != "" {
			outputFile = createOutputFile(input.ReportTarget, input.ReportOutput)
		}
		fmt.Printf("====== Target: %s ======\n", target)

		fmt.Println("=============== SUBDOMAINS ===============")
		var strings1 []string
		strings1 = createSubdomains(input.ReportWordSub, target)
		fmt.Printf("target: %s\n", target)
		asyncGet(strings1, outputFile)

		fmt.Println("=============== PORT SCANNING ===============")
		fmt.Printf("target: %s\n", target)
		asyncPort(input.StartPort, input.EndPort, target, outputFile)

		fmt.Println("=============== DNS SCANNING ===============")
		fmt.Printf("target: %s\n", target)
		lookupDNS(target, outputFile)

		fmt.Println("=============== DIRECTORIES ===============")
		fmt.Printf("target: %s\n", target)
		var strings2 []string
		strings2 = createUrls(input.ReportWordDir, target)
		asyncDir(strings2, outputFile)

	}
	if input.DnsTarget != "" {

		target := cleanProtocol(input.DnsTarget)
		fmt.Println("=============== DNS SCANNING ===============")
		outputFile := ""
		if input.DnsOutput != "" {
			outputFile = createOutputFile(input.DnsTarget, input.DnsOutput)
		}
		fmt.Printf("target: %s\n", target)
		lookupDNS(target, outputFile)

	}
	if input.SubdomainTarget != "" {

		target := cleanProtocol(input.SubdomainTarget)
		fmt.Println("=============== SUBDOMAINS ===============")
		outputFile := ""
		if input.SubdomainOutput != "" {
			outputFile = createOutputFile(input.SubdomainTarget, input.SubdomainOutput)
		}
		fmt.Printf("target: %s\n", target)
		var strings1 []string
		strings1 = createSubdomains(input.SubdomainWord, target)
		asyncGet(strings1, outputFile)

	}
	if input.DirTarget != "" {

		target := cleanProtocol(input.DirTarget)
		fmt.Println("=============== DIRECTORIES ===============")
		outputFile := ""
		if input.DirOutput != "" {
			outputFile = createOutputFile(input.DirTarget, input.DirOutput)
		}
		fmt.Printf("target: %s\n", target)
		var strings2 []string
		strings2 = createUrls(input.DirWord, target)
		asyncDir(strings2, outputFile)
	}
	if input.PortTarget != "" {

		target := input.PortTarget
		if isUrl(target) {
			target = cleanProtocol(input.PortTarget)
		}
		outputFile := ""
		if input.PortOutput != "" {
			outputFile = createOutputFile(input.PortTarget, input.PortOutput)
		}
		fmt.Println("=============== PORT SCANNING ===============")
		fmt.Printf("target: %s\n", target)
		asyncPort(input.StartPort, input.EndPort, target, outputFile)

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

	acceptedOutput := [1]string{"txt"}
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
	ReportTarget    string
	ReportWordDir   string
	ReportWordSub   string
	ReportOutput    string
	DnsTarget       string
	DnsOutput       string
	SubdomainTarget string
	SubdomainWord   string
	SubdomainOutput string
	DirTarget       string
	DirWord         string
	DirOutput       string
	PortTarget      string
	PortOutput      string
	StartPort       int
	EndPort         int
}

//readArgs reads arguments/options from stdin
// Subcommands:
// 		report     ==> Full report
// 		dns        ==> Dns records enumeration
// 		subdomains ==> Subdomains enumeration
// 		port	   ==> ports enumeration
//		dir		   ==> directiories enumeration
// 		help       ==> doc
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
	reportOutputPtr := reportCommand.String("o", "", "output format (txt)")

	// dns subcommand flag pointers
	dnsTargetPtr := dnsCommand.String("target", "", "Target {URL/IP} (Required)")

	// report subcommand flag pointers
	dnsOutputPtr := dnsCommand.String("o", "", "output format (txt)")

	// subdomains subcommand flag pointers
	subdomainTargetPtr := subdomainCommand.String("target", "", "Target {URL} (Required)")

	// subdomains subcommand wordlist
	subdomainWordlistPtr := subdomainCommand.String("w", "", "wordlist to use (default enabled)")

	// report subcommand flag pointers
	subdomainOutputPtr := subdomainCommand.String("o", "", "output format (txt)")

	// dir subcommand flag pointers
	dirTargetPtr := dirCommand.String("target", "", "Target {URL/IP} (Required)")

	// dor subcommand wordlist
	dirWordlistPtr := dirCommand.String("w", "", "wordlist to use (default enabled)")

	// report subcommand flag pointers
	dirOutputPtr := dirCommand.String("o", "", "output format (txt)")

	// port subcommand flag pointers
	portTargetPtr := portCommand.String("target", "", "Target {URL/IP} (Required)")

	// report subcommand flag pointers
	portOutputPtr := portCommand.String("o", "", "output format (txt)")

	portsPtr := portCommand.String("p", "", "ports range <start-end>")
	// Default ports
	StartPort := 1
	EndPort := 65535

	// Verify that a subcommand has been provided
	// os.Arg[0] is the main command
	// os.Arg[1] will be the subcommand
	if len(os.Args) < 2 {
		fmt.Println("[ERROR] subcommand is required.")
		fmt.Println("	Type: scilla help")
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
		if !isUrl(*reportTargetPtr) {
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
	}

	// DNS subcommand
	if dnsCommand.Parsed() {
		// Required Flags
		if *dnsTargetPtr == "" {
			dnsCommand.PrintDefaults()
			os.Exit(1)
		}

		//Verify good inputs
		if !isUrl(*dnsTargetPtr) {
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
		if !isUrl(*subdomainTargetPtr) {
			fmt.Println("The inputted target is not valid.")
			os.Exit(1)
		}
		if !outputFormatIsOk(*subdomainOutputPtr) {
			fmt.Println("The output format is not valid.")
			os.Exit(1)
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
		if !isUrl(*portTargetPtr) {
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
		if !isUrl(*dirTargetPtr) {
			fmt.Println("The inputted target is not valid.")
			os.Exit(1)
		}
		if !outputFormatIsOk(*dirOutputPtr) {
			fmt.Println("The output format is not valid.")
			os.Exit(1)
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
		// Print help
		examples()
		os.Exit(0)
	}

	result := Input{
		*reportTargetPtr,
		*reportWordlistDirPtr,
		*reportWordlistSubdomainPtr,
		*reportOutputPtr,
		*dnsTargetPtr,
		*dnsOutputPtr,
		*subdomainTargetPtr,
		*subdomainWordlistPtr,
		*subdomainOutputPtr,
		*dirTargetPtr,
		*dirWordlistPtr,
		*dirOutputPtr,
		*portTargetPtr,
		*portOutputPtr,
		StartPort,
		EndPort,
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
		maybeStart, err := strconv.Atoi(portsRange)
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

//isIp checks if the inputted Ip is valid
func isIp(str string) bool {
	return govalidator.IsIPv4(str) || govalidator.IsIPv6(str)
}

//isUrl checks if the inputted Url is valid
func isUrl(str string) bool {
	target := cleanProtocol(str)
	str = "http://" + target
	u, err := url.Parse(str)
	return err == nil && u.Host != ""
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
func buildUrl(subdomain string, domain string) string {
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
		path := buildUrl(sub, url)
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

//HttpResp is a struct representing the fundamental data of an HTTP response
type HttpResp struct {
	Id     string
	Target string
	Resp   *http.Response
	Err    error
}

//appendOutputToFile
func appendOutputToFile(output string, filename string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	if _, err := file.WriteString(output + "\n"); err != nil {
		log.Fatal(err)
	}
}

//asyncGet performs concurrent requests to the specified
//urls and prints the results
func asyncGet(urls []string, outputFile string) {

	var count int = 0
	var total int = len(urls)
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	limiter := make(chan string, 80) // Limits simultaneous requests
	wg := sync.WaitGroup{}           // Needed to not prematurely exit before all requests have been finished

	output := outputFile != ""

	for i, domain := range urls {
		wg.Add(1)
		limiter <- domain

		go func(i int, domain string) {
			defer func() { <-limiter }()
			defer wg.Done()

			resp, err := client.Get(domain)
			count++
			if count%100 == 0 { // update counter
				fmt.Fprint(os.Stdout, "\r \r")
				fmt.Printf("%d / %d", count, total)
			}
			if err != nil {
				return
			}

			fmt.Fprint(os.Stdout, "\r \r")
			subDomainFound := cleanProtocol(domain)
			fmt.Printf("[+]FOUND: %s ", subDomainFound)
			if string(resp.Status[0]) == "2" {
				if output {
					appendOutputToFile(domain, outputFile)
				}
				color.Green("%s\n", resp.Status)
			} else {
				color.Red("%s\n", resp.Status)
			}
		}(i, domain)
	}
	wg.Wait()
}

//isOpenPort scans if a port is open
func isOpenPort(host string, port string) bool {
	timeout := time.Second
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
	limiter := make(chan string, 300) // Limits simultaneous requests
	wg := sync.WaitGroup{}            // Needed to not prematurely exit before all requests have been finished

	output := outputFile != ""

	for port := StartingPort; port <= EndingPort; port++ {
		wg.Add(1)
		portStr := fmt.Sprint(port)
		limiter <- portStr

		go func(portStr string, host string) {
			defer func() { <-limiter }()
			defer wg.Done()

			resp := isOpenPort(host, portStr)
			count++
			if count%100 == 0 { // update counter
				fmt.Fprint(os.Stdout, "\r \r")
				fmt.Printf("%d / %d", count, total)
			}
			if resp {
				fmt.Fprint(os.Stdout, "\r \r")
				fmt.Printf("[+]FOUND: %s ", host)
				color.Green("%s\n", portStr)
				if output {
					appendOutputToFile(portStr, outputFile)
				}
			}
		}(portStr, host)
	}
	wg.Wait()
}

//lookupDNS prints the DNS servers for the inputted domain
func lookupDNS(domain string, outputFile string) {

	output := outputFile != ""

	// -- A RECORDS --
	ips, err := net.LookupIP(domain)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
	}
	for _, ip := range ips {
		fmt.Printf("[+]FOUND %s IN A: ", domain)
		color.Green("%s\n", ip.String())
		if output {
			appendOutputToFile(ip.String(), outputFile)
		}
	}

	// -- CNAME RECORD --
	cname, err := net.LookupCNAME(domain)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get CNAME: %v\n", err)
	}
	fmt.Printf("[+]FOUND %s IN CNAME: ", domain)
	color.Green("%s\n", cname)
	if output {
		appendOutputToFile(cname, outputFile)
	}

	// -- NS RECORDS --
	nameserver, err := net.LookupNS(domain)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get NSs: %v\n", err)
	}
	for _, ns := range nameserver {
		fmt.Printf("[+]FOUND %s IN NS: ", domain)
		color.Green("%s\n", ns)
		if output {
			appendOutputToFile(ns.Host, outputFile)
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
		if output {
			appendOutputToFile(mx.Host, outputFile)
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
		if output {
			appendOutputToFile(srv.Target, outputFile)
		}
	}

	// -- TXT RECORDS --
	txtrecords, _ := net.LookupTXT(domain)

	for _, txt := range txtrecords {
		fmt.Printf("[+]FOUND %s IN TXT: ", domain)
		color.Green("%s\n", txt)
		if output {
			appendOutputToFile(txt, outputFile)
		}
	}
}

//asyncDir performs concurrent requests to the specified
//urls and prints the results
func asyncDir(urls []string, outputFile string) {

	var count int = 0
	var total int = len(urls)
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	limiter := make(chan string, 80) // Limits simultaneous requests
	wg := sync.WaitGroup{}           // Needed to not prematurely exit before all requests have been finished

	output := outputFile != ""

	for i, domain := range urls {
		wg.Add(1)
		limiter <- domain

		go func(i int, domain string) {
			defer func() { <-limiter }()
			defer wg.Done()

			resp, err := client.Get(domain)
			count++
			if count%100 == 0 { // update counter
				fmt.Fprint(os.Stdout, "\r \r")
				fmt.Printf("%d / %d", count, total)
			}
			if err != nil {
				return
			}

			if string(resp.Status[0]) == "2" || string(resp.Status[0]) == "3" {
				fmt.Fprint(os.Stdout, "\r \r")
				fmt.Printf("[+]FOUND: %s ", domain)
				color.Green("%s\n", resp.Status)
				if output {
					appendOutputToFile(domain, outputFile)
				}
			} else if (resp.StatusCode != 404) || string(resp.Status[0]) == "5" {
				fmt.Fprint(os.Stdout, "\r \r")
				fmt.Printf("[+]FOUND: %s ", domain)
				color.Red("%s\n", resp.Status)
				if output {
					appendOutputToFile(domain, outputFile)
				}
			}
		}(i, domain)
	}
	wg.Wait()
}
