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

var (
	ErrInvalidArray = errors.New("invalid port array format")
	ErrInvalidRange = errors.New("invalid port range format")
)

const (
	maxPort             = 65535
	minPort             = 1
	portsArrayDelimiter = byte(',')
	portsRangeDelimiter = byte('-')
)

// CheckPortsArray checks the basic rules to
// be valid and then returns the ports array to scan.
// - remove duplicates.
// - check if they can be converted to integers.
// - check if they are in the port array (1 - 65535).
func CheckPortsArray(input string) ([]int, error) {
	sliceOfPorts := strings.Split(input, string(portsArrayDelimiter))
	sliceOfPorts = sliceUtils.RemoveDuplicateValues(sliceOfPorts)
	result := []int{}

	for _, elem := range sliceOfPorts {
		try, err := strconv.Atoi(elem)
		if err != nil || try < minPort || try > maxPort {
			return nil, fmt.Errorf("%w", ErrInvalidArray)
		}

		result = append(result, try)
	}

	return result, nil
}

// CheckPortsRange checks the basic rules to
// be valid and then returns the starting port and the ending port.
func CheckPortsRange(portsRange string, startPort int, endPort int) (int, int, error) {
	switch {
	case portsRange[0] == portsRangeDelimiter:
		{
			maybeEnd, err := strconv.Atoi(portsRange[1:])
			if err != nil {
				return 0, 0, fmt.Errorf("%w", ErrInvalidRange)
			}

			if maybeEnd >= 1 && maybeEnd <= endPort {
				endPort = maybeEnd
			}

			break
		}

	case portsRange[len(portsRange)-1] == portsRangeDelimiter:
		{
			// If ending port isn't specified
			maybeStart, err := strconv.Atoi(portsRange[:len(portsRange)-1])
			if err != nil {
				return 0, 0, fmt.Errorf("%w", ErrInvalidRange)
			}

			if maybeStart > 0 && maybeStart < endPort {
				startPort = maybeStart
			}

			break
		}

	case !strings.Contains(portsRange, string(portsRangeDelimiter)):
		{
			// If a single port is specified
			maybePort, err := strconv.Atoi(portsRange)
			if err != nil {
				return 0, 0, fmt.Errorf("%w", ErrInvalidRange)
			}

			if maybePort > 0 && maybePort < endPort {
				startPort = maybePort
				endPort = maybePort
			}

			break
		}

	default:
		{
			// If a range is specified
			sliceOfPorts := strings.Split(portsRange, string(portsRangeDelimiter))
			if len(sliceOfPorts) != 2 {
				return 0, 0, fmt.Errorf("%w", ErrInvalidRange)
			}

			maybeStart, err := strconv.Atoi(sliceOfPorts[0])
			if err != nil {
				return 0, 0, fmt.Errorf("%w", ErrInvalidRange)
			}

			maybeEnd, err := strconv.Atoi(sliceOfPorts[1])
			if err != nil {
				return 0, 0, fmt.Errorf("%w", ErrInvalidRange)
			}

			if maybeStart > maybeEnd || maybeStart < 1 || maybeEnd > endPort {
				return 0, 0, fmt.Errorf("%w", ErrInvalidRange)
			}

			startPort = maybeStart
			endPort = maybeEnd
		}
	}

	return startPort, endPort, nil
}

func PortsInputHelper(portsPtr *string, startPort, endPort int, portsArray []int,
	portArrayBool bool) (int, int, []int, bool) {
	var err error

	if *portsPtr != "" {
		if strings.Contains(*portsPtr, "-") && strings.Contains(*portsPtr, ",") {
			fmt.Println("you can specify a ports range or an array, not both")
			os.Exit(1)
		}

		switch {
		case strings.Contains(*portsPtr, "-"):
			{
				portsRange := *portsPtr

				startPort, endPort, err = CheckPortsRange(portsRange, startPort, endPort)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				portArrayBool = false
			}
		case strings.Contains(*portsPtr, ","):
			{
				portsArray, err = CheckPortsArray(*portsPtr)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				portArrayBool = true
			}
		default:
			{
				portsRange := *portsPtr

				startPort, endPort, err = CheckPortsRange(portsRange, startPort, endPort)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				portArrayBool = false
			}
		}
	}

	return startPort, endPort, portsArray, portArrayBool
}
