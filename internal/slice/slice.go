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

	@Author:      edoardottt, https://edoardottt.com

	@License: https://github.com/edoardottt/scilla/blob/main/LICENSE

*/

package utils

// RemoveDuplicateValues removes from a slice of string the
// duplicate values.
func RemoveDuplicateValues(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range strSlice {
		if ok := keys[entry]; !ok {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	return list
}

// Difference computes the difference between
// two slices of string (A - B).
func Difference(stringA, stringB []string) []string {
	mapDiff := make(map[string]bool)
	diff := []string{}

	for _, item := range stringB {
		mapDiff[item] = true
	}

	for _, item := range stringA {
		if _, ok := mapDiff[item]; !ok {
			diff = append(diff, item)
		}
	}

	return diff
}
