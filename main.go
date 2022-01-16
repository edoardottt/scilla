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
	"fmt"
	"sync"

	"github.com/edoardottt/scilla/crawler"
	"github.com/edoardottt/scilla/enumeration"
	"github.com/edoardottt/scilla/input"
	"github.com/edoardottt/scilla/opendb"
	"github.com/edoardottt/scilla/output"
	"github.com/edoardottt/scilla/utils"
)

//main function
func main() {
	input := input.ReadArgs()
	// common assets found (only subdomain and dir)
	subs := make(map[string]output.Asset)
	dirs := make(map[string]output.Asset)
	execute(input, subs, dirs, enumeration.CommonPorts())
}

//execute reads inputs and starts the correct procedure
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

//ReportSubcommandHandler >
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
	fmt.Println("=============== FULL REPORT ===============")
	outputFile := ""
	if userInput.ReportOutput != "" {
		outputFile = output.CreateOutputFile(userInput.ReportTarget, "report", userInput.ReportOutput)
		if outputFile[len(outputFile)-4:] == "html" {
			output.BannerHTML(userInput.ReportTarget, outputFile)
		}
	}
	fmt.Println("=============== SUBDOMAINS SCANNING ===============")
	var strings1 []string
	// change from ip to Hostname
	if utils.IsIP(target) {
		targetIP = target
		target = utils.IpToHostname(targetIP)
	}
	target = utils.CleanProtocol(target)
	if outputFile != "" {
		if outputFile[len(outputFile)-4:] == "html" {
			output.HeaderHTML("SUBDOMAIN SCANNING", outputFile)
		}
	}
	if userInput.ReportCrawlerSub {
		go crawler.SpawnCrawler(utils.CleanProtocol(target), protocolTemp,
			userInput.ReportIgnoreSub, dirs, subs, outputFile, mutex, "sub", false)
	}
	strings1 = input.CreateSubdomains(userInput.ReportWordSub, protocolTemp, utils.CleanProtocol(target))
	if userInput.ReportSubdomainDB {
		if userInput.ReportSpyse {
			_ = input.GetSpyseKey()
		}
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
		if userInput.ReportSpyse {
			spyseSubs := opendb.SpyseSubdomains(utils.CleanProtocol(target), input.GetSpyseKey())
			strings1 = opendb.AppendDBSubdomains(spyseSubs, strings1)
		}
		if userInput.ReportVirusTotal {
			vtSubs := opendb.VirusTotalSubdomains(utils.CleanProtocol(target), input.GetVirusTotalKey())
			strings1 = opendb.AppendDBSubdomains(vtSubs, strings1)
		}
	}
	// be sure to not scan duplicate values
	strings1 = utils.RemoveDuplicateValues(utils.CleanSubdomainsOk(utils.CleanProtocol(target), strings1))
	enumeration.AsyncGet(protocolTemp, strings1, userInput.ReportIgnoreSub, outputFile, subs, mutex, false)
	if outputFile != "" {
		if outputFile[len(outputFile)-4:] == "html" {
			output.FooterHTML(outputFile)
		}
	}

	if targetIP != "" {
		target = targetIP
	}
	fmt.Println("=============== PORT SCANNING ===============")

	enumeration.AsyncPort(userInput.PortsArray, userInput.PortArrayBool, userInput.StartPort, userInput.EndPort,
		utils.CleanProtocol(target), outputFile, userInput.ReportCommon, enumeration.CommonPorts(), false, userInput.ReportTimeoutPort)

	fmt.Println("=============== DNS SCANNING ===============")
	enumeration.LookupDNS(utils.CleanProtocol(target), outputFile, false)
	fmt.Println("=============== DIRECTORIES SCANNING ===============")
	var strings2 = input.CreateUrls(userInput.ReportWordDir, protocolTemp, utils.CleanProtocol(target))
	if outputFile != "" {
		if outputFile[len(outputFile)-4:] == "html" {
			output.HeaderHTML("DIRECTORY SCANNING", outputFile)
		}
	}
	if userInput.ReportCrawlerDir {
		go crawler.SpawnCrawler(utils.CleanProtocol(target), protocolTemp,
			userInput.ReportIgnoreDir, dirs, subs, outputFile, mutex, "dir", false)
	}
	enumeration.AsyncDir(strings2, userInput.ReportIgnoreDir, outputFile, dirs, mutex, false, userInput.ReportRedirect)
	if outputFile != "" {
		if outputFile[len(outputFile)-4:] == "html" {
			output.FooterHTML(outputFile)
		}
	}
	if userInput.ReportOutput != "" {
		if outputFile[len(outputFile)-4:] == "html" {
			output.BannerFooterHTML(outputFile)
		}
	}
}

//DNSSubcommandHandler >
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
		fmt.Println("=============== DNS SCANNING ===============")
	}
	outputFile := ""
	if userInput.DNSOutput != "" {
		outputFile = output.CreateOutputFile(userInput.DNSTarget, "dns", userInput.DNSOutput)

		if outputFile[len(outputFile)-4:] == "html" {
			output.BannerHTML(userInput.DNSTarget, outputFile)
		}
	}
	enumeration.LookupDNS(target, outputFile, userInput.DNSPlain)
	if userInput.DNSOutput != "" {
		if outputFile[len(outputFile)-4:] == "html" {
			output.BannerFooterHTML(outputFile)
		}
	}
}

//SubdomainSubcommandHandler >
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
		fmt.Println("=============== SUBDOMAINS SCANNING ===============")
	}
	outputFile := ""
	if userInput.SubdomainOutput != "" {
		outputFile = output.CreateOutputFile(userInput.SubdomainTarget, "subdomain", userInput.SubdomainOutput)
		if outputFile[len(outputFile)-4:] == "html" {
			output.BannerHTML(userInput.SubdomainTarget, outputFile)
		}
	}
	var strings1 []string
	if !userInput.SubdomainNoCheck {
		strings1 = input.CreateSubdomains(userInput.SubdomainWord, protocolTemp, utils.CleanProtocol(target))
	}
	if userInput.SubdomainDB {
		if userInput.SubdomainSpyse {
			_ = input.GetSpyseKey()
		}
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
		if userInput.SubdomainSpyse {
			spyseSubs := opendb.SpyseSubdomains(utils.CleanProtocol(target), input.GetSpyseKey())
			strings1 = opendb.AppendDBSubdomains(spyseSubs, strings1)
		}
		if userInput.SubdomainVirusTotal {
			vtSubs := opendb.VirusTotalSubdomains(utils.CleanProtocol(target), input.GetVirusTotalKey())
			strings1 = opendb.AppendDBSubdomains(vtSubs, strings1)
		}

	}
	if outputFile != "" {
		if outputFile[len(outputFile)-4:] == "html" {
			output.HeaderHTML("SUBDOMAIN SCANNING", outputFile)
		}
	}
	if userInput.SubdomainCrawler && !userInput.SubdomainNoCheck {
		go crawler.SpawnCrawler(utils.CleanProtocol(target), protocolTemp,
			userInput.SubdomainIgnore, dirs, subs, outputFile, mutex, "sub", userInput.SubdomainPlain)
	}
	// be sure to not scan duplicate values
	strings1 = utils.RemoveDuplicateValues(utils.CleanSubdomainsOk(utils.CleanProtocol(target), strings1))
	if !userInput.SubdomainNoCheck {
		enumeration.AsyncGet(protocolTemp, strings1, userInput.SubdomainIgnore, outputFile, subs, mutex, userInput.SubdomainPlain)
	} else {
		for _, elem := range strings1 {
			fmt.Println(elem)
			if userInput.SubdomainOutput == "txt" {
				output.AppendOutputToTxt(elem, outputFile)
			}
			if userInput.SubdomainOutput == "html" {
				output.AppendOutputToHTML(elem, "", outputFile)
			}
		}
	}
	if outputFile != "" {
		if outputFile[len(outputFile)-4:] == "html" {
			output.FooterHTML(outputFile)
		}
	}
	if userInput.SubdomainOutput != "" {
		if outputFile[len(outputFile)-4:] == "html" {
			output.BannerFooterHTML(outputFile)
		}
	}
}

//DirSubcommandHandler >
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
		fmt.Println("=============== DIRECTORIES SCANNING ===============")
	}
	target = utils.CleanProtocol(target)
	outputFile := ""
	if userInput.DirOutput != "" {
		outputFile = output.CreateOutputFile(userInput.DirTarget, "dir", userInput.DirOutput)

		if outputFile[len(outputFile)-4:] == "html" {
			output.BannerHTML(userInput.DirTarget, outputFile)
		}
	}
	var strings2 = input.CreateUrls(userInput.DirWord, protocolTemp, target)
	if outputFile != "" {
		if outputFile[len(outputFile)-4:] == "html" {
			output.HeaderHTML("DIRECTORY SCANNING", outputFile)
		}
	}
	if userInput.DirCrawler {
		go crawler.SpawnCrawler(utils.CleanProtocol(target), protocolTemp,
			userInput.DirIgnore, dirs, subs, outputFile, mutex, "dir", userInput.DirPlain)
	}
	enumeration.AsyncDir(strings2, userInput.DirIgnore, outputFile, dirs, mutex, userInput.DirPlain, userInput.DirRedirect)
	if outputFile != "" {
		if outputFile[len(outputFile)-4:] == "html" {
			output.FooterHTML(outputFile)
		}
	}
	if userInput.DirOutput != "" {
		if outputFile[len(outputFile)-4:] == "html" {
			output.BannerFooterHTML(outputFile)
		}
	}
}

//PortSubcommandHandler >
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
	outputFile := ""
	if userInput.PortOutput != "" {
		outputFile = output.CreateOutputFile(userInput.PortTarget, "port", userInput.PortOutput)
		if outputFile[len(outputFile)-4:] == "html" {
			output.BannerHTML(userInput.PortTarget, outputFile)
		}
	}
	if !userInput.PortPlain {
		fmt.Printf("target: %s\n", target)
		fmt.Println("=============== PORT SCANNING ===============")
	}
	enumeration.AsyncPort(userInput.PortsArray, userInput.PortArrayBool, userInput.StartPort, userInput.EndPort,
		target, outputFile, userInput.PortCommon, common, userInput.PortPlain, userInput.PortTimeout)

	if userInput.PortOutput != "" {
		if outputFile[len(outputFile)-4:] == "html" {
			output.BannerFooterHTML(outputFile)
		}
	}
}
