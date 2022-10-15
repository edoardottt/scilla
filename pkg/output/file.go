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

	fileUtils "github.com/edoardottt/scilla/internal/file"
)

// CreateOutputFolder creates the output folder.
func CreateOutputFolder(path string) {
	// Create a folder/directory at a full qualified path
	if strings.Trim(path, " ") != "" {
		err := os.MkdirAll(path, fileUtils.Permission0755)
		if err != nil {
			fmt.Println("Can't create output folder.")
			os.Exit(1)
		}
	}
}

// CreateOutputFile creates the output file (txt/json/html).
func CreateOutputFile(path string) string {
	dir, file := filepath.Split(path)
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		sepPresent := strings.Contains(path, string(os.PathSeparator))
		if _, err := os.Stat(dir); os.IsNotExist(err) && sepPresent {
			CreateOutputFolder(dir)
		}
		// If the file doesn't exist, create it.
		fOpen, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, fileUtils.Permission0644)
		if err != nil {
			fmt.Println("Can't create output file.")
			os.Exit(1)
		}

		fOpen.Close()
	} else {
		// The file already exists, check what the user want.
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("The output file %s already esists, do you want to overwrite? (Y/n): ", file)
		text, _ := reader.ReadString('\n')
		answer := strings.ToLower(text)
		answer = strings.TrimSpace(answer)

		if answer == "y" || answer == "yes" || answer == "" {
			fOpen, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, fileUtils.Permission0644)
			if err != nil {
				fmt.Println("Can't create output file.")
				os.Exit(1)
			}
			err = fOpen.Truncate(0)
			if err != nil {
				fmt.Println("Can't create output file.")
				os.Exit(1)
			}
			fOpen.Close()
		} else {
			os.Exit(1)
		}
	}

	return path
}

// AppendWhere checks which format the output should be (html, json or txt).
func AppendWhere(what string, status string, key string, record string, format string, outputFile string) {
	switch {
	case format == "html":
		{
			AppendOutputToHTML(what, status, outputFile)
		}
	case format == "json":
		{
			AppendOutputToJSON(what, key, record, outputFile)
		}
	default:
		{
			AppendOutputToTxt(what, outputFile)
		}
	}
}

// AppendExtension appends to the path the given extension.
func AppendExtension(path, extension string) string {
	if len(path) < len(extension)+1 {
		return path + "." + extension
	}

	if path[len(path)-len(extension)-1:] != "."+extension {
		return path + "." + extension
	}

	return path
}
