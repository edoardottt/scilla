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

package output

import (
	"log"
	"os"

	"github.com/edoardottt/scilla/utils"
)

// BannerHTML writes in the input file the HTML banner
func BannerHTML(target string, filename string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, utils.Permission0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	_, err = file.WriteString(`<html><body>
	<div style='background-color:#4adeff;color:white'>
	<h1>Scilla - Information Gathering Tool</h1>
	<ul>
	<li><a href='https://github.com/edoardottt/scilla'>github.com/edoardottt/scilla</a></li>
	<li>edoardottt, <a href='https://www.edoardoottavianelli.it'>edoardoottavianelli.it</a></li>
	<li>Released under <a href='http://www.gnu.org/licenses/gpl-3.0.html'>GPLv3 License</a></li>
	</ul></div>`)

	if err != nil {
		log.Printf(err.Error())

		return
	}

	_, err = file.WriteString("<h4>target: " + target + "</h4>")
	if err != nil {
		log.Printf(err.Error())

		return
	}

	file.Close()
}

// AppendOutputToHTML appends a (html) row in the HTML output file
func AppendOutputToHTML(output string, status string, filename string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, utils.Permission0644)
	if err != nil {
		log.Fatal(err)
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

	if _, err := file.WriteString("<li><a target='_blank' href='" + output + "'>" +
		utils.CleanProtocol(output) +
		"</a> " + statusColor + "</li>"); err != nil {
		log.Printf(err.Error())

		return
	}

	file.Close()
}

// HeaderHTML writes in the (html) output file the header (directories, dns ...)
func HeaderHTML(header string, filename string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, utils.Permission0644)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := file.WriteString("<h3>" + header + "</h3><ul>"); err != nil {
		log.Fatal(err)
	}

	file.Close()
}

// FooterHTML closes the HTML list
func FooterHTML(filename string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, utils.Permission0644)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	if _, err := file.WriteString("</ul>"); err != nil {
		log.Printf(err.Error())

		return
	}

	file.Close()
}

// BannerFooterHTML writes in the (html) output file the HTML footer
func BannerFooterHTML(filename string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, utils.Permission0644)

	if err != nil {
		log.Fatal(err)

		return
	}

	_, err = file.WriteString(`<div style='background-color:#4adeff;color:white'>
	<ul><li><a href='https://github.com/edoardottt/scilla'>Contribute to scilla</a></li>
	<li>Released under <a href='http://www.gnu.org/licenses/gpl-3.0.html'>GPLv3 License</a></li></ul></div>`)
	if err != nil {
		log.Fatal(err)
	}

	file.Close()
}
