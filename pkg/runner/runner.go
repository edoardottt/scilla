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

package runner

import (
	"fmt"
	"sync"

	ipUtils "github.com/edoardottt/scilla/internal/ip"
	sliceUtils "github.com/edoardottt/scilla/internal/slice"
	urlUtils "github.com/edoardottt/scilla/internal/url"
	"github.com/edoardottt/scilla/pkg/crawler"
	"github.com/edoardottt/scilla/pkg/enumeration"
	"github.com/edoardottt/scilla/pkg/input"
	"github.com/edoardottt/scilla/pkg/opendb"
	"github.com/edoardottt/scilla/pkg/output"
)

const (
	httpProtocol = "http"
)

// Runner.
type Runner struct {
	Input input.Input
}

// New returns a new Runner with the provided input.
func New() *Runner {
	input := input.ReadArgs()
	return &Runner{Input: input}
}

// Execute reads the input and starts the correct procedure.
func (r *Runner) Execute(dirs, subs map[string]output.Asset) {
	var (
		mutex           = &sync.Mutex{}
		commandProvided = false
	)

	// :::::::: REPORT SUBCOMMAND HANDLER ::::::::
	if r.Input.ReportTarget != "" {
		commandProvided = true

		ReportSubcommandHandler(r.Input, mutex, dirs, subs)
	}

	// :::::::: DNS SUBCOMMAND HANDLER ::::::::
	if r.Input.DNSTarget != "" {
		commandProvided = true

		DNSSubcommandHandler(r.Input)
	}

	// :::::::: SUBDOMAIN SUBCOMMAND HANDLER ::::::::
	if r.Input.SubdomainTarget != "" {
		commandProvided = true

		SubdomainSubcommandHandler(r.Input, mutex, dirs, subs)
	}

	// :::::::: DIRECTORIES SUBCOMMAND HANDLER ::::::::
	if r.Input.DirTarget != "" {
		commandProvided = true

		DirSubcommandHandler(r.Input, mutex, dirs, subs)
	}

	// :::::::: PORT SUBCOMMAND HANDLER ::::::::
	if r.Input.PortTarget != "" {
		commandProvided = true

		PortSubcommandHandler(r.Input, enumeration.CommonPorts())
	}

	if !commandProvided {
		output.Help()
	}
}

// ReportSubcommandHandler.
func ReportSubcommandHandler(userInput input.Input, mutex *sync.Mutex,
	dirs map[string]output.Asset, subs map[string]output.Asset) {
	output.Intro()

	target := userInput.ReportTarget
	if target[len(target)-1] == byte('/') {
		target = target[:len(target)-1]
	}

	var protocolTemp string
	// if there isn't a scheme use http.
	if !urlUtils.ProtocolExists(target) {
		protocolTemp = httpProtocol
	} else {
		protocolTemp = urlUtils.RetrieveProtocol(target)
	}

	var targetIP string

	fmt.Printf("target: %s\n", target)
	fmt.Println("================ FULL REPORT ========================")

	// - json output -
	var outputFileJSON string
	if userInput.ReportOutputJSON != "" {
		outputFileJSON = output.CreateOutputFile(userInput.ReportOutputJSON)
	}

	// - html output -
	var outputFileHTML string
	if userInput.ReportOutputHTML != "" {
		outputFileHTML = output.CreateOutputFile(userInput.ReportOutputHTML)
		output.BannerHTML(userInput.ReportTarget, outputFileHTML)
	}

	// - txt output -
	var outputFileTXT string
	if userInput.ReportOutputTXT != "" {
		outputFileTXT = output.CreateOutputFile(userInput.ReportOutputTXT)
	}

	fmt.Println("================ SCANNING SUBDOMAINS ================")

	var subdomains []string
	// change from ip to Hostname
	if ipUtils.IsIP(target) {
		targetIP = target
		target = ipUtils.IPToHostname(targetIP)
	}

	target = urlUtils.CleanProtocol(target)

	if outputFileHTML != "" {
		output.HeaderHTML("SUBDOMAINS ENUMERATION", outputFileHTML)
	}

	if userInput.ReportCrawlerSub {
		go crawler.SpawnCrawler(urlUtils.CleanProtocol(target), protocolTemp,
			userInput.ReportIgnoreSub, dirs, subs, outputFileJSON, outputFileHTML,
			outputFileTXT, mutex, "sub", false, userInput.ReportUserAgent, userInput.ReportRandomUserAgent)
	}

	subdomains = input.CreateSubdomains(userInput.ReportWordSub, protocolTemp, urlUtils.CleanProtocol(target))

	if userInput.ReportSubdomainDB {
		crtsh := opendb.CrtshSubdomains(urlUtils.CleanProtocol(target), false)
		subdomains = opendb.AppendDBSubdomains(crtsh, subdomains)
		threatcrowd := opendb.ThreatcrowdSubdomains(urlUtils.CleanProtocol(target), false)
		subdomains = opendb.AppendDBSubdomains(threatcrowd, subdomains)
		hackerTarget := opendb.HackerTargetSubdomains(urlUtils.CleanProtocol(target), false)
		subdomains = opendb.AppendDBSubdomains(hackerTarget, subdomains)
		anubis := opendb.AnubisSubdomains(urlUtils.CleanProtocol(target), false)
		subdomains = opendb.AppendDBSubdomains(anubis, subdomains)
		threatminer := opendb.ThreatMinerSubdomains(urlUtils.CleanProtocol(target), false)
		subdomains = opendb.AppendDBSubdomains(threatminer, subdomains)

		// Service Not Working
		// bufferOverrun := opendb.BufferOverrunSubdomains(urlUtils.CleanProtocol(target), false)
		// subdomains = opendb.AppendDBSubdomains(bufferOverrun, subdomains)

		// Service not working
		// sonar := opendb.SonarSubdomains(urlUtils.CleanProtocol(target), false)
		// subdomains = opendb.AppendDBSubdomains(sonar, subdomains)

		if userInput.ReportVirusTotal {
			vtSubs := opendb.VirusTotalSubdomains(urlUtils.CleanProtocol(target), input.GetKey("virustotal"), false)
			subdomains = opendb.AppendDBSubdomains(vtSubs, subdomains)
		}

		subdomains = opendb.ShuffleSubdomains(subdomains)
	}

	// be sure to not scan duplicate values
	subdomains = sliceUtils.RemoveDuplicateValues(urlUtils.CleanSubdomainsOk(urlUtils.CleanProtocol(target), subdomains))

	enumeration.AsyncGet(protocolTemp, subdomains, userInput.ReportIgnoreSub, outputFileJSON,
		outputFileHTML, outputFileTXT, subs, mutex, false, userInput.ReportUserAgent, userInput.ReportRandomUserAgent,
		userInput.ReportAlive, userInput.ReportDNS)

	if outputFileHTML != "" {
		output.FooterHTML(outputFileHTML)
	}

	if targetIP != "" {
		target = targetIP
	}

	fmt.Println("================ SCANNING PORTS =====================")

	enumeration.AsyncPort(userInput.PortsArray, userInput.PortArrayBool, userInput.StartPort, userInput.EndPort,
		urlUtils.CleanProtocol(target), outputFileJSON, outputFileHTML, outputFileTXT,
		userInput.ReportCommon, enumeration.CommonPorts(), false, userInput.ReportTimeoutPort)

	fmt.Println("================ SCANNING DNS =======================")
	enumeration.LookupDNS(urlUtils.CleanProtocol(target), outputFileJSON, outputFileHTML, outputFileTXT, false)

	fmt.Println("================ SCANNING DIRECTORIES ===============")

	var urls = input.CreateUrls(userInput.ReportWordDir, protocolTemp, urlUtils.CleanProtocol(target))

	if outputFileHTML != "" {
		output.HeaderHTML("DIRECTORIES ENUMERATION", outputFileHTML)
	}

	if userInput.ReportCrawlerDir {
		go crawler.SpawnCrawler(urlUtils.CleanProtocol(target), protocolTemp,
			userInput.ReportIgnoreDir, dirs, subs, outputFileJSON, outputFileHTML, outputFileTXT,
			mutex, "dir", false, userInput.ReportUserAgent, userInput.ReportRandomUserAgent)
	}

	enumeration.AsyncDir(urls, userInput.ReportIgnoreDir, outputFileJSON, outputFileHTML, outputFileTXT,
		dirs, mutex, false, userInput.ReportRedirect, userInput.ReportUserAgent, userInput.ReportRandomUserAgent)

	if outputFileHTML != "" {
		output.FooterHTML(outputFileHTML)
		output.BannerFooterHTML(outputFileHTML)
	}
}

// DNSSubcommandHandler.
func DNSSubcommandHandler(userInput input.Input) {
	if !userInput.DNSPlain {
		output.Intro()
	}

	target := urlUtils.CleanProtocol(userInput.DNSTarget)
	// change from ip to Hostname
	if ipUtils.IsIP(target) {
		target = ipUtils.IPToHostname(target)
	}

	if target[len(target)-1] == byte('/') {
		target = target[:len(target)-1]
	}

	if !userInput.DNSPlain {
		fmt.Printf("target: %s\n", target)
		fmt.Println("================ SCANNING DNS =======================")
	}

	// - json output -
	var outputFileJSON string
	if userInput.DNSOutputJSON != "" {
		outputFileJSON = output.CreateOutputFile(userInput.DNSOutputJSON)
	}

	// - html output -
	var outputFileHTML string
	if userInput.DNSOutputHTML != "" {
		outputFileHTML = output.CreateOutputFile(userInput.DNSOutputHTML)
		output.BannerHTML(userInput.DNSTarget, outputFileHTML)
	}

	// - txt output -
	var outputFileTXT string
	if userInput.DNSOutputTXT != "" {
		outputFileTXT = output.CreateOutputFile(userInput.DNSOutputTXT)
	}

	enumeration.LookupDNS(target, outputFileJSON, outputFileHTML, outputFileTXT, userInput.DNSPlain)

	if userInput.DNSOutputHTML != "" {
		output.BannerFooterHTML(outputFileHTML)
	}
}

// SubdomainSubcommandHandler.
func SubdomainSubcommandHandler(userInput input.Input, mutex *sync.Mutex,
	dirs map[string]output.Asset, subs map[string]output.Asset) {
	if !userInput.SubdomainPlain {
		output.Intro()
	}

	target := userInput.SubdomainTarget

	var protocolTemp string
	// if there isn't a scheme use http.
	if !urlUtils.ProtocolExists(target) {
		protocolTemp = httpProtocol
	} else {
		protocolTemp = urlUtils.RetrieveProtocol(target)
	}

	// change from ip to Hostname
	if ipUtils.IsIP(target) {
		target = ipUtils.IPToHostname(target)
	}

	if target[len(target)-1] == byte('/') {
		target = target[:len(target)-1]
	}

	if !userInput.SubdomainPlain {
		fmt.Printf("target: %s\n", target)
		fmt.Println("================ SCANNING SUBDOMAINS ================")
	}

	// - json output -
	var outputFileJSON string
	if userInput.SubdomainOutputJSON != "" {
		outputFileJSON = output.CreateOutputFile(userInput.SubdomainOutputJSON)
	}

	// - html output -
	var outputFileHTML string
	if userInput.SubdomainOutputHTML != "" {
		outputFileHTML = output.CreateOutputFile(userInput.SubdomainOutputHTML)
		output.BannerHTML(userInput.SubdomainTarget, outputFileHTML)
	}

	// - txt output -
	var outputFileTXT string
	if userInput.SubdomainOutputTXT != "" {
		outputFileTXT = output.CreateOutputFile(userInput.SubdomainOutputTXT)
	}

	var subdomains []string
	if !userInput.SubdomainNoCheck {
		subdomains = input.CreateSubdomains(userInput.SubdomainWord, protocolTemp, urlUtils.CleanProtocol(target))
	}

	if userInput.SubdomainDB {
		crtsh := opendb.CrtshSubdomains(urlUtils.CleanProtocol(target), userInput.SubdomainPlain)
		subdomains = opendb.AppendDBSubdomains(crtsh, subdomains)
		threatcrowd := opendb.ThreatcrowdSubdomains(urlUtils.CleanProtocol(target), userInput.SubdomainPlain)
		subdomains = opendb.AppendDBSubdomains(threatcrowd, subdomains)
		hackerTarget := opendb.HackerTargetSubdomains(urlUtils.CleanProtocol(target), userInput.SubdomainPlain)
		subdomains = opendb.AppendDBSubdomains(hackerTarget, subdomains)
		anubis := opendb.AnubisSubdomains(urlUtils.CleanProtocol(target), false)
		subdomains = opendb.AppendDBSubdomains(anubis, subdomains)
		threatminer := opendb.ThreatMinerSubdomains(urlUtils.CleanProtocol(target), false)
		subdomains = opendb.AppendDBSubdomains(threatminer, subdomains)

		// Service Not Working
		// bufferOverrun := opendb.BufferOverrunSubdomains(urlUtils.CleanProtocol(target), userInput.SubdomainPlain)
		// subdomains = opendb.AppendDBSubdomains(bufferOverrun, subdomains)

		// Service not working
		// sonar := opendb.SonarSubdomains(urlUtils.CleanProtocol(target), userInput.SubdomainPlain)
		// subdomains = opendb.AppendDBSubdomains(sonar, subdomains)

		if userInput.SubdomainVirusTotal {
			vtSubs := opendb.VirusTotalSubdomains(urlUtils.CleanProtocol(target), input.GetKey("virustotal"),
				userInput.SubdomainPlain)
			subdomains = opendb.AppendDBSubdomains(vtSubs, subdomains)
		}

		if userInput.SubdomainBuiltWith {
			builtWithSubs := opendb.BuiltWithSubdomains(urlUtils.CleanProtocol(target), input.GetKey("builtwith"),
				userInput.SubdomainPlain)
			subdomains = opendb.AppendDBSubdomains(builtWithSubs, subdomains)
		}
	}

	if outputFileHTML != "" {
		output.HeaderHTML("SUBDOMAINS ENUMERATION", outputFileHTML)
	}

	if userInput.SubdomainCrawler && !userInput.SubdomainNoCheck {
		go crawler.SpawnCrawler(urlUtils.CleanProtocol(target), protocolTemp,
			userInput.SubdomainIgnore, dirs, subs, outputFileJSON, outputFileHTML, outputFileTXT,
			mutex, "sub", userInput.SubdomainPlain, userInput.SubdomainUserAgent, userInput.SubdomainRandomUserAgent)
	}

	// be sure to not scan duplicate values
	subdomains = sliceUtils.RemoveDuplicateValues(urlUtils.CleanSubdomainsOk(urlUtils.CleanProtocol(target), subdomains))

	if !userInput.SubdomainNoCheck {
		enumeration.AsyncGet(protocolTemp, subdomains, userInput.SubdomainIgnore, outputFileJSON, outputFileHTML,
			outputFileTXT, subs, mutex, userInput.SubdomainPlain, userInput.SubdomainUserAgent,
			userInput.SubdomainRandomUserAgent, userInput.SubdomainAlive, userInput.SubdomainDNS)
	} else {
		for _, elem := range subdomains {
			fmt.Println(elem)
			if outputFileJSON != "" {
				output.AppendOutputToJSON(elem, "SUB", "", outputFileJSON)
			}
			if outputFileHTML != "" {
				output.AppendOutputToHTML(elem, "", outputFileHTML)
			}
			if outputFileTXT != "" {
				output.AppendOutputToTxt(elem, outputFileTXT)
			}
		}
	}

	if outputFileHTML != "" {
		output.FooterHTML(outputFileHTML)
		output.BannerFooterHTML(outputFileHTML)
	}
}

// DirSubcommandHandler.
func DirSubcommandHandler(userInput input.Input, mutex *sync.Mutex,
	dirs map[string]output.Asset, subs map[string]output.Asset) {
	if !userInput.DirPlain {
		output.Intro()
	}

	target := userInput.DirTarget

	var protocolTemp string

	// if there isn't a scheme use http.
	if !urlUtils.ProtocolExists(target) {
		protocolTemp = httpProtocol
	} else {
		protocolTemp = urlUtils.RetrieveProtocol(target)
	}

	if target[len(target)-1] == byte('/') {
		target = target[:len(target)-1]
	}

	if !userInput.DirPlain {
		fmt.Printf("target: %s\n", target)
		fmt.Println("================ SCANNING DIRECTORIES ===============")
	}

	target = urlUtils.CleanProtocol(target)

	// - json output -
	var outputFileJSON string
	if userInput.DirOutputJSON != "" {
		outputFileJSON = output.CreateOutputFile(userInput.DirOutputJSON)
	}

	// - html output -
	var outputFileHTML string
	if userInput.DirOutputHTML != "" {
		outputFileHTML = output.CreateOutputFile(userInput.DirOutputHTML)
		output.BannerHTML(userInput.DirTarget, outputFileHTML)
	}

	// - txt output -
	var outputFileTXT string
	if userInput.DirOutputTXT != "" {
		outputFileTXT = output.CreateOutputFile(userInput.DirOutputTXT)
	}

	var urls = input.CreateUrls(userInput.DirWord, protocolTemp, target)

	if outputFileHTML != "" {
		output.HeaderHTML("DIRECTORIES ENUMERATION", outputFileHTML)
	}

	if userInput.DirCrawler {
		go crawler.SpawnCrawler(urlUtils.CleanProtocol(target), protocolTemp,
			userInput.DirIgnore, dirs, subs, outputFileJSON, outputFileHTML, outputFileTXT,
			mutex, "dir", userInput.DirPlain, userInput.DirUserAgent, userInput.DirRandomUserAgent)
	}

	enumeration.AsyncDir(urls, userInput.DirIgnore, outputFileJSON, outputFileHTML, outputFileTXT,
		dirs, mutex, userInput.DirPlain, userInput.DirRedirect, userInput.DirUserAgent, userInput.DirRandomUserAgent)

	if outputFileHTML != "" {
		output.FooterHTML(outputFileHTML)
		output.BannerFooterHTML(outputFileHTML)
	}
}

// PortSubcommandHandler.
func PortSubcommandHandler(userInput input.Input, common []int) {
	if !userInput.PortPlain {
		output.Intro()
	}

	target := userInput.PortTarget
	if urlUtils.IsURL(target) {
		target = urlUtils.CleanProtocol(userInput.PortTarget)
	}

	if target[len(target)-1] == byte('/') {
		target = target[:len(target)-1]
	}

	// - json output -
	var outputFileJSON string
	if userInput.PortOutputJSON != "" {
		outputFileJSON = output.CreateOutputFile(userInput.PortOutputJSON)
	}

	// - html output -
	var outputFileHTML string
	if userInput.PortOutputHTML != "" {
		outputFileHTML = output.CreateOutputFile(userInput.PortOutputHTML)
		output.BannerHTML(userInput.PortTarget, outputFileHTML)
	}

	// - txt output -
	var outputFileTXT string
	if userInput.PortOutputTXT != "" {
		outputFileTXT = output.CreateOutputFile(userInput.PortOutputTXT)
	}

	if !userInput.PortPlain {
		fmt.Printf("target: %s\n", target)
		fmt.Println("================ SCANNING PORTS =====================")
	}

	enumeration.AsyncPort(userInput.PortsArray, userInput.PortArrayBool, userInput.StartPort, userInput.EndPort,
		target, outputFileJSON, outputFileHTML, outputFileTXT,
		userInput.PortCommon, common, userInput.PortPlain, userInput.PortTimeout)

	if userInput.PortOutputHTML != "" {
		output.BannerFooterHTML(outputFileHTML)
	}
}
