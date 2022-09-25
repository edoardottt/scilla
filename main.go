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

package main

import (
	"fmt"
	"sync"

	"github.com/edoardottt/scilla/crawler"
	"github.com/edoardottt/scilla/enumeration"
	"github.com/edoardottt/scilla/input"
	"github.com/edoardottt/scilla/opendb"
	"github.com/edoardottt/scilla/output"
	"github.com/edoardottt/scilla/utils"
)

const (
	httpProtocol = "http"
)

// main function.
func main() {
	input := input.ReadArgs()
	// common assets found (only subdomain and dir)
	subs := make(map[string]output.Asset)
	dirs := make(map[string]output.Asset)
	execute(input, subs, dirs, enumeration.CommonPorts())
}

// execute reads inputs and starts the correct procedure.
func execute(userInput input.Input, subs map[string]output.Asset, dirs map[string]output.Asset, common []int) {
	var mutex = &sync.Mutex{}

	var commandProvided = false

	// :::::::: REPORT SUBCOMMAND HANDLER ::::::::
	if userInput.ReportTarget != "" {
		commandProvided = true

		ReportSubcommandHandler(userInput, mutex, dirs, subs)
	}

	// :::::::: DNS SUBCOMMAND HANDLER ::::::::
	if userInput.DNSTarget != "" {
		commandProvided = true

		DNSSubcommandHandler(userInput)
	}

	// :::::::: SUBDOMAIN SUBCOMMAND HANDLER ::::::::
	if userInput.SubdomainTarget != "" {
		commandProvided = true

		SubdomainSubcommandHandler(userInput, mutex, dirs, subs)
	}

	// :::::::: DIRECTORIES SUBCOMMAND HANDLER ::::::::
	if userInput.DirTarget != "" {
		commandProvided = true

		DirSubcommandHandler(userInput, mutex, dirs, subs)
	}

	// :::::::: PORT SUBCOMMAND HANDLER ::::::::
	if userInput.PortTarget != "" {
		commandProvided = true

		PortSubcommandHandler(userInput, common)
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
	if !utils.ProtocolExists(target) {
		protocolTemp = httpProtocol
	} else {
		protocolTemp = utils.RetrieveProtocol(target)
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

	var strings1 []string
	// change from ip to Hostname
	if utils.IsIP(target) {
		targetIP = target
		target = utils.IPToHostname(targetIP)
	}

	target = utils.CleanProtocol(target)

	if outputFileHTML != "" {
		output.HeaderHTML("SUBDOMAINS ENUMERATION", outputFileHTML)
	}

	if userInput.ReportCrawlerSub {
		go crawler.SpawnCrawler(utils.CleanProtocol(target), protocolTemp,
			userInput.ReportIgnoreSub, dirs, subs, outputFileJSON, outputFileHTML,
			outputFileTXT, mutex, "sub", false)
	}

	strings1 = input.CreateSubdomains(userInput.ReportWordSub, protocolTemp, utils.CleanProtocol(target))

	if userInput.ReportSubdomainDB {
		if userInput.ReportVirusTotal {
			_ = input.GetVirusTotalKey()
		}

		sonar := opendb.SonarSubdomains(utils.CleanProtocol(target))
		strings1 = opendb.AppendDBSubdomains(sonar, strings1)
		crtsh := opendb.CrtshSubdomains(utils.CleanProtocol(target))
		strings1 = opendb.AppendDBSubdomains(crtsh, strings1)
		threatcrowd := opendb.ThreatcrowdSubdomains(utils.CleanProtocol(target))
		strings1 = opendb.AppendDBSubdomains(threatcrowd, strings1)
		hackerTarget := opendb.HackerTargetSubdomains(utils.CleanProtocol(target))
		strings1 = opendb.AppendDBSubdomains(hackerTarget, strings1)
		bufferOverrun := opendb.BufferOverrunSubdomains(utils.CleanProtocol(target))
		strings1 = opendb.AppendDBSubdomains(bufferOverrun, strings1)

		if userInput.ReportVirusTotal {
			vtSubs := opendb.VirusTotalSubdomains(utils.CleanProtocol(target), input.GetVirusTotalKey())
			strings1 = opendb.AppendDBSubdomains(vtSubs, strings1)
		}
	}

	// be sure to not scan duplicate values
	strings1 = utils.RemoveDuplicateValues(utils.CleanSubdomainsOk(utils.CleanProtocol(target), strings1))
	enumeration.AsyncGet(protocolTemp, strings1, userInput.ReportIgnoreSub, outputFileJSON,
		outputFileHTML, outputFileTXT, subs, mutex, false, userInput.ReportUserAgent, userInput.ReportRandomUserAgent)

	if outputFileHTML != "" {
		output.FooterHTML(outputFileHTML)
	}

	if targetIP != "" {
		target = targetIP
	}

	fmt.Println("================ SCANNING PORTS =====================")

	enumeration.AsyncPort(userInput.PortsArray, userInput.PortArrayBool, userInput.StartPort, userInput.EndPort,
		utils.CleanProtocol(target), outputFileJSON, outputFileHTML, outputFileTXT,
		userInput.ReportCommon, enumeration.CommonPorts(), false, userInput.ReportTimeoutPort)

	fmt.Println("================ SCANNING DNS =======================")
	enumeration.LookupDNS(utils.CleanProtocol(target), outputFileJSON, outputFileHTML, outputFileTXT, false)

	fmt.Println("================ SCANNING DIRECTORIES ===============")

	var strings2 = input.CreateUrls(userInput.ReportWordDir, protocolTemp, utils.CleanProtocol(target))

	if outputFileHTML != "" {
		output.HeaderHTML("DIRECTORIES ENUMERATION", outputFileHTML)
	}

	if userInput.ReportCrawlerDir {
		go crawler.SpawnCrawler(utils.CleanProtocol(target), protocolTemp,
			userInput.ReportIgnoreDir, dirs, subs, outputFileJSON, outputFileHTML, outputFileTXT, mutex, "dir", false)
	}

	enumeration.AsyncDir(strings2, userInput.ReportIgnoreDir, outputFileJSON, outputFileHTML, outputFileTXT,
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

	target := utils.CleanProtocol(userInput.DNSTarget)
	// change from ip to Hostname
	if utils.IsIP(target) {
		target = utils.IPToHostname(target)
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
	if !utils.ProtocolExists(target) {
		protocolTemp = httpProtocol
	} else {
		protocolTemp = utils.RetrieveProtocol(target)
	}

	// change from ip to Hostname
	if utils.IsIP(target) {
		target = utils.IPToHostname(target)
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

	var strings1 []string
	if !userInput.SubdomainNoCheck {
		strings1 = input.CreateSubdomains(userInput.SubdomainWord, protocolTemp, utils.CleanProtocol(target))
	}

	if userInput.SubdomainDB {
		sonar := opendb.SonarSubdomains(utils.CleanProtocol(target))
		strings1 = opendb.AppendDBSubdomains(sonar, strings1)
		crtsh := opendb.CrtshSubdomains(utils.CleanProtocol(target))
		strings1 = opendb.AppendDBSubdomains(crtsh, strings1)
		threatcrowd := opendb.ThreatcrowdSubdomains(utils.CleanProtocol(target))
		strings1 = opendb.AppendDBSubdomains(threatcrowd, strings1)
		hackerTarget := opendb.HackerTargetSubdomains(utils.CleanProtocol(target))
		strings1 = opendb.AppendDBSubdomains(hackerTarget, strings1)
		bufferOverrun := opendb.BufferOverrunSubdomains(utils.CleanProtocol(target))
		strings1 = opendb.AppendDBSubdomains(bufferOverrun, strings1)

		if userInput.SubdomainVirusTotal {
			vtSubs := opendb.VirusTotalSubdomains(utils.CleanProtocol(target), input.GetVirusTotalKey())
			strings1 = opendb.AppendDBSubdomains(vtSubs, strings1)
		}
	}

	if outputFileHTML != "" {
		output.HeaderHTML("SUBDOMAINS ENUMERATION", outputFileHTML)
	}

	if userInput.SubdomainCrawler && !userInput.SubdomainNoCheck {
		go crawler.SpawnCrawler(utils.CleanProtocol(target), protocolTemp,
			userInput.SubdomainIgnore, dirs, subs, outputFileJSON, outputFileHTML, outputFileTXT,
			mutex, "sub", userInput.SubdomainPlain)
	}

	// be sure to not scan duplicate values
	strings1 = utils.RemoveDuplicateValues(utils.CleanSubdomainsOk(utils.CleanProtocol(target), strings1))
	if !userInput.SubdomainNoCheck {
		enumeration.AsyncGet(protocolTemp, strings1, userInput.SubdomainIgnore, outputFileJSON, outputFileHTML, outputFileTXT,
			subs, mutex, userInput.SubdomainPlain, userInput.SubdomainUserAgent, userInput.SubdomainRandomUserAgent)
	} else {
		for _, elem := range strings1 {
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
	if !utils.ProtocolExists(target) {
		protocolTemp = httpProtocol
	} else {
		protocolTemp = utils.RetrieveProtocol(target)
	}

	if target[len(target)-1] == byte('/') {
		target = target[:len(target)-1]
	}

	if !userInput.DirPlain {
		fmt.Printf("target: %s\n", target)
		fmt.Println("================ SCANNING DIRECTORIES ===============")
	}

	target = utils.CleanProtocol(target)

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

	var strings2 = input.CreateUrls(userInput.DirWord, protocolTemp, target)

	if outputFileHTML != "" {
		output.HeaderHTML("DIRECTORIES ENUMERATION", outputFileHTML)
	}

	if userInput.DirCrawler {
		go crawler.SpawnCrawler(utils.CleanProtocol(target), protocolTemp,
			userInput.DirIgnore, dirs, subs, outputFileJSON, outputFileHTML, outputFileTXT,
			mutex, "dir", userInput.DirPlain)
	}

	enumeration.AsyncDir(strings2, userInput.DirIgnore, outputFileJSON, outputFileHTML, outputFileTXT,
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
	if utils.IsURL(target) {
		target = utils.CleanProtocol(userInput.PortTarget)
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
