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

package output

import "strings"

//Asset gives information about the asset found
type Asset struct {
	Value   string
	Printed bool
}

//OutputFormatIsOk checks if the specified output format is Ok
//(txt, html or json)
func OutputFormatIsOk(input string) bool {
	if input == "" {
		return true
	}
	acceptedOutput := [3]string{"txt", "html", "json"}
	input = strings.ToLower(input)
	for _, output := range acceptedOutput {
		if output == input {
			return true
		}
	}
	return false
}

//ReplaceBadCharacterOutput replaces slashes with dots
func ReplaceBadCharacterOutput(input string) string {
	result := strings.ReplaceAll(input, "/", "-")
	return result
}
