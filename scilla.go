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
	"github.com/asaskevich/govalidator"
	"github.com/fatih/color"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

//intro prints the banner
func intro() {
	banner1 := "*****************************************\n"
	banner2 := "*                 _ _ _                 *\n"
	banner3 := "*        ___  ___(_) | | __ _           *\n"
	banner4 := "*       / __|/ __| | | |/ _` |          *\n"
	banner5 := "*       \\__ \\ (__| | | | (_| |          *\n"
	banner6 := "*       |___/\\___|_|_|_|\\__,_|          *\n"
	banner7 := "*                                       *\n"
	banner8 := "* https://github.com/edoardottt/scilla  *\n"
	banner9 := "* https://www.edoardoottavianelli.it    *\n"
	banner10 := "*                                       *\n"
	banner11 := "*****************************************\n"
	banner := banner1 + banner2 + banner3 + banner4 + banner5 + banner6 + banner7 + banner8 + banner9 + banner10 + banner11
	fmt.Println(banner)
}

//help prints in stdout scilla usage
func help() {
	fmt.Println("Information Gathering tool - DNS/subdomain/port enumeration")
	fmt.Println("")
	fmt.Println("usage: scilla [subcommand] { options }")
	fmt.Println("")
	fmt.Println("	Available subcommands:")
	fmt.Println("		- dns { -target <target (URL)> REQUIRED}")
	fmt.Println("		- subdomain { -target <target (URL)> REQUIRED}")
	fmt.Println("		- port { [-p <start-end>] -target <target (URL/IP)> REQUIRED}")
	fmt.Println("		- report { -target <target (URL/IP)> REQUIRED}")
	fmt.Println("		- help")
	fmt.Println("	Examples:")
	fmt.Println("		- scilla subdomain -target target.domain")
	fmt.Println("		- scilla port -p -450 -target target.domain")
}

//main
func main() {
	intro()
	input := readArgs()
	execute(input)
}

//execute reads inputs and start the correct procedure
func execute(input Input) {
	if input.ReportTarget != "" {
		fmt.Println("=============== REPORT ===============")
		target := cleanProtocol(input.SubdomainTarget)
		fmt.Printf("target: %s\n", target)
		fmt.Println("=============== SUBDOMAINS ===============")
		strings1 := createSubdomains(target)
		asyncGet(strings1)
		fmt.Println("=============== PORT SCANNING ===============")
		fmt.Printf("target: %s\n", target)
		asyncPort(input.StartPort, input.EndPort, target)
		fmt.Println("=============== DNS SCANNING ===============")
		fmt.Printf("target: %s\n", target)
	}
	if input.DnsTarget != "" {
		target := cleanProtocol(input.SubdomainTarget)
		fmt.Println("=============== DNS SCANNING ===============")
		fmt.Printf("target: %s\n", target)
	}
	if input.SubdomainTarget != "" {
		target := cleanProtocol(input.SubdomainTarget)
		fmt.Println("=============== SUBDOMAINS ===============")
		fmt.Printf("target: %s\n", target)
		strings1 := createSubdomains(target)
		asyncGet(strings1)
	}
	if input.PortTarget != "" {
		target := input.PortTarget
		if isUrl(target) {
			target = cleanProtocol(input.PortTarget)
		}
		fmt.Println("=============== PORT SCANNING ===============")
		fmt.Printf("target: %s\n", target)
		asyncPort(input.StartPort, input.EndPort, target)
	}

}

//cleanProtocol remove from the url the protocol scheme
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
	return target
}

//Input contains the input parameters
type Input struct {
	ReportTarget    string
	DnsTarget       string
	SubdomainTarget string
	PortTarget      string
	StartPort       int
	EndPort         int
}

//readArgs reads arguments/options from stdin
// Subcommands:
// 		report     ==> Full report
// 		dns        ==> Dns records enumeration
// 		subdomains ==> SubDomains enumeration
// 		port	   ==> port enumeration
// 		help       ==> doc
func readArgs() Input {
	reportCommand := flag.NewFlagSet("report", flag.ExitOnError)
	dnsCommand := flag.NewFlagSet("dns", flag.ExitOnError)
	subdomainCommand := flag.NewFlagSet("subdomain", flag.ExitOnError)
	portCommand := flag.NewFlagSet("port", flag.ExitOnError)
	helpCommand := flag.NewFlagSet("help", flag.ExitOnError)

	// report subcommand flag pointers
	reportTargetPtr := reportCommand.String("target", "", "Target {URL/IP} (Required)")

	// dns subcommand flag pointers
	dnsTargetPtr := dnsCommand.String("target", "", "Target {URL/IP} (Required)")

	// subdomains subcommand flag pointers
	subdomainTargetPtr := subdomainCommand.String("target", "", "Target {URL/IP} (Required)")

	// port subcommand flag pointers
	portTargetPtr := portCommand.String("target", "", "Target {URL/IP} (Required)")
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
	// FlagSet.Parse() requires a set of arguments to parse as input
	// os.Args[2:] will be all arguments starting after the subcommand at os.Args[1]
	switch os.Args[1] {
	case "report":
		reportCommand.Parse(os.Args[2:])
	case "dns":
		dnsCommand.Parse(os.Args[2:])
	case "subdomain":
		subdomainCommand.Parse(os.Args[2:])
	case "port":
		portCommand.Parse(os.Args[2:])
	case "help":
		helpCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Check which subcommand was Parsed using the FlagSet.Parsed() function. Handle each case accordingly.
	// FlagSet.Parse() will evaluate to false if no flags were parsed (i.e. the user did not provide any flags)
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
	}

	// Check which subcommand was Parsed using the FlagSet.Parsed() function. Handle each case accordingly.
	// FlagSet.Parse() will evaluate to false if no flags were parsed (i.e. the user did not provide any flags)
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
	}

	// Check which subcommand was Parsed using the FlagSet.Parsed() function. Handle each case accordingly.
	// FlagSet.Parse() will evaluate to false if no flags were parsed (i.e. the user did not provide any flags)
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
	}

	// Check which subcommand was Parsed using the FlagSet.Parsed() function. Handle each case accordingly.
	// FlagSet.Parse() will evaluate to false if no flags were parsed (i.e. the user did not provide any flags)
	if portCommand.Parsed() {

		// Required Flags
		if *portTargetPtr == "" {
			portCommand.PrintDefaults()
			os.Exit(1)
		}

		// If there's ports range, define it as inputs for the struct
		if *portsPtr != "" {
			portsRange := string(*portsPtr)
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
				maybeStart, err := strconv.Atoi(portsRange[:len(portsRange)])
				if err != nil {
					fmt.Println("The inputted port range is not valid.")
					os.Exit(1)
				}
				if maybeStart > 0 && maybeStart < EndPort {
					StartPort = maybeStart
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
		}

		//Verify good inputs
		if !isUrl(*portTargetPtr) {
			fmt.Println("The inputted target is not valid.")
			os.Exit(1)
		}
	}

	// Check which subcommand was Parsed using the FlagSet.Parsed() function. Handle each case accordingly.
	// FlagSet.Parse() will evaluate to false if no flags were parsed (i.e. the user did not provide any flags)
	if helpCommand.Parsed() {
		// Print help
		help()
		os.Exit(0)
	}

	result := Input{*reportTargetPtr, *dnsTargetPtr, *subdomainTargetPtr, *portTargetPtr, StartPort, EndPort}
	return result
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

//get performs a HTTP GET request to the target
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

//readDict scan all the possible subdomains from file
func readDict(inputFile string) []string {
	// os.Open() opens specific file in
	// read-only mode and this return
	// a pointer of type os.
	file, err := os.Open(inputFile)

	if err != nil {
		log.Fatalf("failed to open")

	}

	// The bufio.NewScanner() function is called in which the
	// object os.File passed as its parameter and this returns a
	// object bufio.Scanner which is further used on the
	// bufio.Scanner.Split() method.
	scanner := bufio.NewScanner(file)

	// The bufio.ScanLines is used as an
	// input to the method bufio.Scanner.Split()
	// and then the scanning forwards to each
	// new line using the bufio.Scanner.Scan()
	// method.
	scanner.Split(bufio.ScanLines)
	var text []string

	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	// The method os.File.Close() is called
	// on the os.File object to close the file
	file.Close()
	return text
}

//createSubdomains returns a list of subdomains
func createSubdomains(url string) []string {
	subs := readDict("lists/subdomains/subdomains.txt")
	result := []string{}
	for _, sub := range subs {
		path := buildUrl(sub, url)
		result = append(result, path)
	}
	return result
}

//HttpResp is a struct representing the fundamental data of HTTP response
type HttpResp struct {
	Id     string
	Target string
	Resp   *http.Response
	Err    error
}

//asyncGet performs concurrent requests to the specified
//urls and prints the results
func asyncGet(urls []string) {

	limiter := make(chan string, 200) // Limits simultaneous requests

	wg := sync.WaitGroup{} // Needed to not prematurely exit before all requests have been finished

	for i, domain := range urls {
		wg.Add(1)
		limiter <- domain

		go func(i int, domain string) {
			defer func() { <-limiter }()
			defer wg.Done()

			resp, err := http.Get(domain)
			if err != nil {
				return
			}

			fmt.Printf("[+]FOUND: %s: ", domain)
			if string(resp.Status[0]) == "2" {
				color.Green("%s\n", resp.Status)
			} else {
				color.Red("%s", resp.Status)
			}
		}(i, domain)
	}

	wg.Wait()
}

//isOpenPort scan if a port is open
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
func asyncPort(StartingPort int, EndingPort int, host string) {

	limiter := make(chan string, 1000) // Limits simultaneous requests

	wg := sync.WaitGroup{} // Needed to not prematurely exit before all requests have been finished

	for port := StartingPort; port <= EndingPort; port++ {
		wg.Add(1)
		portStr := fmt.Sprint(port)
		limiter <- portStr

		go func(portStr string, host string) {
			defer func() { <-limiter }()
			defer wg.Done()

			resp := isOpenPort(host, portStr)
			if resp {
				fmt.Printf("[+]FOUND: %s:", host)
				color.Green("%s\n", portStr)
			}
		}(portStr, host)
	}

	wg.Wait()
}
