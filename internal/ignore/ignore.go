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

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	sliceUtils "github.com/edoardottt/scilla/internal/slice"
)

const (
	statusCodeLength = 3
)

var (
	ErrWrongStatusCodeLength  = errors.New("invalid status code: It should consist of three digits")
	ErrInvalidStatusCode      = errors.New("invalid status code: 100 <= code <= 599")
	ErrInvalidStatusCodeClass = errors.New("invalid status code: You can use * only in 1**,2**,3**,4**,5**")
)

// CheckIgnore checks the inputted status code(s) to be ignored.
// It can be a list e.g. 301,302,400,404,500
// It can be a 'class' of codes e.g. 3**.
func CheckIgnore(input string) ([]string, error) {
	result := []string{}
	temp := strings.Split(input, ",")
	temp = sliceUtils.RemoveDuplicateValues(temp)

	for _, elem := range temp {
		elem := strings.TrimSpace(elem)
		if len(elem) != statusCodeLength {
			return nil, fmt.Errorf("%w", ErrWrongStatusCodeLength)
		}

		if ignoreInt, err := strconv.Atoi(elem); err == nil {
			// if it is a valid status code without * (e.g. 404)
			if 100 <= ignoreInt && ignoreInt <= 599 {
				result = append(result, elem)
			} else {
				return nil, fmt.Errorf("%w", ErrInvalidStatusCode)
			}
		} else if strings.Contains(elem, "*") {
			// if it is a valid status code without * (e.g. 4**)
			if ignoreClassOk(elem) {
				result = append(result, elem)
			} else {
				return nil, fmt.Errorf("%w", ErrInvalidStatusCodeClass)
			}
		}
	}

	result = sliceUtils.RemoveDuplicateValues(result)
	result = deleteUnusefulIgnoreresponses(result)

	return result, nil
}

// deleteUnusefulIgnoreresponses removes from to-be-ignored arrays
// the responses already included with * as classes.
func deleteUnusefulIgnoreresponses(input []string) []string {
	var result []string

	toberemoved := []string{}
	classes := []string{}

	for _, elem := range input {
		if strings.Contains(elem, "*") {
			classes = append(classes, elem)
		}
	}

	for _, class := range classes {
		for _, elem := range input {
			if class[0] == elem[0] && elem[1] != '*' {
				toberemoved = append(toberemoved, elem)
			}
		}
	}

	result = sliceUtils.Difference(input, toberemoved)

	return result
}

// ignoreClassOk states if the class of ignored status codes
// is correct or not (4**,2**...)
func ignoreClassOk(input string) bool {
	if strings.Contains(input, "*") {
		if _, err := strconv.Atoi(string(input[0])); err == nil {
			i, err := strconv.Atoi(string(input[0]))
			if err != nil {
				// handle error
				fmt.Println(err)
				os.Exit(1)
			}

			if i >= 1 && i <= 5 {
				if input[1] == byte('*') && input[2] == byte('*') {
					return true
				}
			}
		}
	}

	return false
}

// IgnoreResponse returns a boolean if the response
// should be ignored or not.
func IgnoreResponse(response int, ignore []string) bool {
	responseString := strconv.Itoa(response)
	// if I don't have to ignore responses, just return true
	if len(ignore) == 0 {
		return false
	}

	for _, ignorePort := range ignore {
		if strings.Contains(ignorePort, "*") {
			if responseString[0] == ignorePort[0] {
				return true
			}
		}

		if responseString == ignorePort {
			return true
		}
	}

	return false
}
