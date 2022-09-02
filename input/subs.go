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

package input

import (
	"bufio"
	"log"
	"os"
	"runtime"

	"github.com/edoardottt/scilla/utils"
)

// ReadDictSubs reads all the possible subdomains from file
func ReadDictSubs(inputFile string) []string {
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

// CreateSubdomains returns a list of subdomains
// from the default file lists/subdomains.txt.
func CreateSubdomains(filename string, scheme string, url string) []string {
	var subs []string

	if filename == "" {
		if runtime.GOOS == windows {
			subs = ReadDictSubs("lists/subdomains.txt")
		} else { // linux
			subs = ReadDictSubs("/usr/bin/lists/subdomains.txt")
		}
	} else {
		subs = ReadDictSubs(filename)
	}

	result := []string{}

	for _, sub := range subs {
		path := utils.BuildURL(scheme, sub, url)
		result = append(result, path)
	}

	return result
}
