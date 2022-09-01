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

// main function
func main() {
	input := input.ReadArgs()
	// common assets found (only subdomain and dir)
	subs := make(map[string]output.Asset)
	dirs := make(map[string]output.Asset)
	execute(input, subs, dirs, enumeration.CommonPorts())
}

// execute reads inputs and starts the correct procedure
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

// ReportSubcommandHandler >
func ReportSubcommandHandler(userInput input.Input, mutex *sync.Mutex, dirs map[string]output.Asset, subs map[string]output.Asset) {
	output.Intro()
	target := userInput.ReportTarget
	if target[len(target)-1] == byte('/') {
		target = target[:len(target)-1]
	}
	var protocolTemp string
	// if there isn't a scheme use http.
	if !utils.ProtocolExists(target) {
		protocolTemp = "http"
	} else {
		protocolTemp = utils.RetrieveProtocol(target)
	}
	var targetIP string
	fmt.Printf("target: %s\n", target)
	fmt.Println("================ FULL REPORT ========================")

	// - json output -
	var outputFileJson string
	if userInput.ReportOutputJson != "" {
		outputFileJson = output.CreateOutputFile(userInput.ReportOutputJson)
	}
	// - html output -
	var outputFileHtml string
	if userInput.ReportOutputHtml != "" {
		outputFileHtml = output.CreateOutputFile(userInput.ReportOutputHtml)
		output.BannerHTML(userInput.ReportTarget, outputFileHtml)
	}
	// - txt output -
	var outputFileTxt string
	if userInput.ReportOutputTxt != "" {
		outputFileTxt = output.CreateOutputFile(userInput.ReportOutputTxt)
	}

	fmt.Println("================ SCANNING SUBDOMAINS ================")
	var strings1 []string
	// change from ip to Hostname
	if utils.IsIP(target) {
		targetIP = target
		target = utils.IpToHostname(targetIP)
	}
	target = utils.CleanProtocol(target)
	if outputFileHtml != "" {
		output.HeaderHTML("SUBDOMAINS ENUMERATION", outputFileHtml)
	}
	if userInput.ReportCrawlerSub {
		go crawler.SpawnCrawler(utils.CleanProtocol(target), protocolTemp,
			userInput.ReportIgnoreSub, dirs, subs, outputFileJson, outputFileHtml, outputFileTxt, mutex, "sub", false)
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
	enumeration.AsyncGet(protocolTemp, strings1, userInput.ReportIgnoreSub, outputFileJson,
		outputFileHtml, outputFileTxt, subs, mutex, false)
	if outputFileHtml != "" {
		output.FooterHTML(outputFileHtml)
	}

	if targetIP != "" {
		target = targetIP
	}
	fmt.Println("================ SCANNING PORTS =====================")

	enumeration.AsyncPort(userInput.PortsArray, userInput.PortArrayBool, userInput.StartPort, userInput.EndPort,
		utils.CleanProtocol(target), outputFileJson, outputFileHtml, outputFileTxt,
		userInput.ReportCommon, enumeration.CommonPorts(), false, userInput.ReportTimeoutPort)

	fmt.Println("================ SCANNING DNS =======================")
	enumeration.LookupDNS(utils.CleanProtocol(target), outputFileJson, outputFileHtml, outputFileTxt, false)

	fmt.Println("================ SCANNING DIRECTORIES ===============")
	var strings2 = input.CreateUrls(userInput.ReportWordDir, protocolTemp, utils.CleanProtocol(target))
	if outputFileHtml != "" {
		output.HeaderHTML("DIRECTORIES ENUMERATION", outputFileHtml)
	}
	if userInput.ReportCrawlerDir {
		go crawler.SpawnCrawler(utils.CleanProtocol(target), protocolTemp,
			userInput.ReportIgnoreDir, dirs, subs, outputFileJson, outputFileHtml, outputFileTxt, mutex, "dir", false)
	}
	enumeration.AsyncDir(strings2, userInput.ReportIgnoreDir, outputFileJson, outputFileHtml, outputFileTxt,
		dirs, mutex, false, userInput.ReportRedirect)
	if outputFileHtml != "" {
		output.FooterHTML(outputFileHtml)
		output.BannerFooterHTML(outputFileHtml)
	}
}

// DNSSubcommandHandler >
func DNSSubcommandHandler(userInput input.Input) {
	if !userInput.DNSPlain {
		output.Intro()
	}
	target := utils.CleanProtocol(userInput.DNSTarget)
	// change from ip to Hostname
	if utils.IsIP(target) {
		target = utils.IpToHostname(target)
	}
	if target[len(target)-1] == byte('/') {
		target = target[:len(target)-1]
	}
	if !userInput.DNSPlain {
		fmt.Printf("target: %s\n", target)
		fmt.Println("================ SCANNING DNS =======================")
	}
	// - json output -
	var outputFileJson string
	if userInput.DNSOutputJson != "" {
		outputFileJson = output.CreateOutputFile(userInput.DNSOutputJson)
	}
	// - html output -
	var outputFileHtml string
	if userInput.DNSOutputHtml != "" {
		outputFileHtml = output.CreateOutputFile(userInput.DNSOutputHtml)
		output.BannerHTML(userInput.DNSTarget, outputFileHtml)
	}
	// - txt output -
	var outputFileTxt string
	if userInput.DNSOutputTxt != "" {
		outputFileTxt = output.CreateOutputFile(userInput.DNSOutputTxt)
	}

	enumeration.LookupDNS(target, outputFileJson, outputFileHtml, outputFileTxt, userInput.DNSPlain)
	if userInput.DNSOutputHtml != "" {
		output.BannerFooterHTML(outputFileHtml)
	}
}

// SubdomainSubcommandHandler >
func SubdomainSubcommandHandler(userInput input.Input, mutex *sync.Mutex, dirs map[string]output.Asset, subs map[string]output.Asset) {
	if !userInput.SubdomainPlain {
		output.Intro()
	}

	target := userInput.SubdomainTarget
	var protocolTemp string
	// if there isn't a scheme use http.
	if !utils.ProtocolExists(target) {
		protocolTemp = "http"
	} else {
		protocolTemp = utils.RetrieveProtocol(target)
	}
	// change from ip to Hostname
	if utils.IsIP(target) {
		target = utils.IpToHostname(target)
	}
	if target[len(target)-1] == byte('/') {
		target = target[:len(target)-1]
	}
	if !userInput.SubdomainPlain {
		fmt.Printf("target: %s\n", target)
		fmt.Println("================ SCANNING SUBDOMAINS ================")
	}
	// - json output -
	var outputFileJson string
	if userInput.SubdomainOutputJson != "" {
		outputFileJson = output.CreateOutputFile(userInput.SubdomainOutputJson)
	}
	// - html output -
	var outputFileHtml string
	if userInput.SubdomainOutputHtml != "" {
		outputFileHtml = output.CreateOutputFile(userInput.SubdomainOutputHtml)
		output.BannerHTML(userInput.SubdomainTarget, outputFileHtml)
	}
	// - txt output -
	var outputFileTxt string
	if userInput.SubdomainOutputTxt != "" {
		outputFileTxt = output.CreateOutputFile(userInput.SubdomainOutputTxt)
	}
	var strings1 []string
	if !userInput.SubdomainNoCheck {
		strings1 = input.CreateSubdomains(userInput.SubdomainWord, protocolTemp, utils.CleanProtocol(target))
	}
	if userInput.SubdomainDB {
		if userInput.SubdomainVirusTotal {
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
		if userInput.SubdomainVirusTotal {
			vtSubs := opendb.VirusTotalSubdomains(utils.CleanProtocol(target), input.GetVirusTotalKey())
			strings1 = opendb.AppendDBSubdomains(vtSubs, strings1)
		}

	}
	if outputFileHtml != "" {
		output.HeaderHTML("SUBDOMAINS ENUMERATION", outputFileHtml)
	}
	if userInput.SubdomainCrawler && !userInput.SubdomainNoCheck {
		go crawler.SpawnCrawler(utils.CleanProtocol(target), protocolTemp,
			userInput.SubdomainIgnore, dirs, subs, outputFileJson, outputFileHtml, outputFileTxt,
			mutex, "sub", userInput.SubdomainPlain)
	}
	// be sure to not scan duplicate values
	strings1 = utils.RemoveDuplicateValues(utils.CleanSubdomainsOk(utils.CleanProtocol(target), strings1))
	if !userInput.SubdomainNoCheck {
		enumeration.AsyncGet(protocolTemp, strings1, userInput.SubdomainIgnore, outputFileJson, outputFileHtml, outputFileTxt,
			subs, mutex, userInput.SubdomainPlain)
	} else {
		for _, elem := range strings1 {
			fmt.Println(elem)
			if outputFileJson != "" {
				output.AppendOutputToJSON(elem, "SUB", "", outputFileJson)
			}
			if outputFileHtml != "" {
				output.AppendOutputToHTML(elem, "", outputFileHtml)
			}
			if outputFileTxt != "" {
				output.AppendOutputToTxt(elem, outputFileTxt)
			}
		}
	}
	if outputFileHtml != "" {
		output.FooterHTML(outputFileHtml)
		output.BannerFooterHTML(outputFileHtml)
	}
}

// DirSubcommandHandler >
func DirSubcommandHandler(userInput input.Input, mutex *sync.Mutex, dirs map[string]output.Asset, subs map[string]output.Asset) {
	if !userInput.DirPlain {
		output.Intro()
	}
	target := userInput.DirTarget
	var protocolTemp string
	// if there isn't a scheme use http.
	if !utils.ProtocolExists(target) {
		protocolTemp = "http"
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
	var outputFileJson string
	if userInput.DirOutputJson != "" {
		outputFileJson = output.CreateOutputFile(userInput.DirOutputJson)
	}
	// - html output -
	var outputFileHtml string
	if userInput.DirOutputHtml != "" {
		outputFileHtml = output.CreateOutputFile(userInput.DirOutputHtml)
		output.BannerHTML(userInput.DirTarget, outputFileHtml)
	}
	// - txt output -
	var outputFileTxt string
	if userInput.DirOutputTxt != "" {
		outputFileTxt = output.CreateOutputFile(userInput.DirOutputTxt)
	}
	var strings2 = input.CreateUrls(userInput.DirWord, protocolTemp, target)
	if outputFileHtml != "" {
		output.HeaderHTML("DIRECTORIES ENUMERATION", outputFileHtml)
	}
	if userInput.DirCrawler {
		go crawler.SpawnCrawler(utils.CleanProtocol(target), protocolTemp,
			userInput.DirIgnore, dirs, subs, outputFileJson, outputFileHtml, outputFileTxt,
			mutex, "dir", userInput.DirPlain)
	}
	enumeration.AsyncDir(strings2, userInput.DirIgnore, outputFileJson, outputFileHtml, outputFileTxt,
		dirs, mutex, userInput.DirPlain, userInput.DirRedirect)
	if outputFileHtml != "" {
		output.FooterHTML(outputFileHtml)
		output.BannerFooterHTML(outputFileHtml)
	}
}

// PortSubcommandHandler >
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
	var outputFileJson string
	if userInput.PortOutputJson != "" {
		outputFileJson = output.CreateOutputFile(userInput.PortOutputJson)
	}
	// - html output -
	var outputFileHtml string
	if userInput.PortOutputHtml != "" {
		outputFileHtml = output.CreateOutputFile(userInput.PortOutputHtml)
		output.BannerHTML(userInput.PortTarget, outputFileHtml)
	}
	// - txt output -
	var outputFileTxt string
	if userInput.PortOutputTxt != "" {
		outputFileTxt = output.CreateOutputFile(userInput.PortOutputTxt)
	}
	if !userInput.PortPlain {
		fmt.Printf("target: %s\n", target)
		fmt.Println("================ SCANNING PORTS =====================")
	}
	enumeration.AsyncPort(userInput.PortsArray, userInput.PortArrayBool, userInput.StartPort, userInput.EndPort,
		target, outputFileJson, outputFileHtml, outputFileTxt,
		userInput.PortCommon, common, userInput.PortPlain, userInput.PortTimeout)

	if userInput.PortOutputHtml != "" {
		output.BannerFooterHTML(outputFileHtml)
	}
}
