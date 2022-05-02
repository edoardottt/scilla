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
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//CreateOutputFolder creates the output folder
func CreateOutputFolder(path string) {
	//Create a folder/directory at a full qualified path
	err := os.MkdirAll(path, 0755)
	if err != nil {
		fmt.Println("Can't create output folder.")
		os.Exit(1)
	}
}

//CreateOutputFile creates the output file (txt/json/html)
func CreateOutputFile(path string) string {
	dir, _ := filepath.Split("/some/path/to/remove/file.name")

	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			CreateOutputFolder(dir)
		}
		// If the file doesn't exist, create it.
		f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
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
			f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
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
	return path
}

//AppendWhere checks which format the output should be (html, json or txt)
func AppendWhere(what string, status string, key string, record string, outputFile string) {
	if outputFile[len(outputFile)-4:] == "html" {
		AppendOutputToHTML(what, status, outputFile)
	} else if outputFile[len(outputFile)-4:] == "json" {
		AppendOutputToJSON(what, key, record, outputFile)
	} else {
		AppendOutputToTxt(what, outputFile)
	}
}
