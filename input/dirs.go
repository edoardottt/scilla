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

package input

import (
	"bufio"
	"log"
	"os"
	"runtime"

	"github.com/edoardottt/scilla/utils"
)

//ReadDictDirs reads all the possible dirs from input file
func ReadDictDirs(inputFile string) []string {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("failed to open %s ", inputFile)
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string
	var dir = ""
	for scanner.Scan() {
		dir = scanner.Text()
		if len(dir) > 0 {
			if string(dir[len(dir)-1:]) == "/" {
				dir = dir[:len(dir)-1]
			}
			text = append(text, dir)
		}
	}
	file.Close()
	text = utils.RemoveDuplicateValues(text)
	return text
}

//CreateUrls returns a list of directories
//from the default file lists/dirs.txt.
func CreateUrls(filename string, scheme string, url string) []string {
	var dirs []string
	if filename == "" {
		if runtime.GOOS == "windows" {
			dirs = ReadDictDirs("lists/dirs.txt")
		} else { // linux
			dirs = ReadDictDirs("/usr/bin/lists/dirs.txt")
		}
	} else {
		dirs = ReadDictDirs(filename)
	}
	result := []string{}
	for _, dir := range dirs {
		path, path2 := utils.AppendDir(scheme, url, dir)
		result = append(result, path)
		result = append(result, path2)
	}
	return result
}
