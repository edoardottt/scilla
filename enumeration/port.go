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

package enumeration

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/edoardottt/scilla/output"
	"github.com/edoardottt/scilla/utils"
	"github.com/fatih/color"
)

//CommonPorts is a slice of integers containing the common ports.
func CommonPorts() []int {
	return []int{13, 20, 21, 22, 23, 25, 42, 50, 51, 53, 67, 68,
		69, 70, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 90, 102,
		107, 109, 110, 111, 113, 115, 118, 119, 123, 135, 136, 137,
		138, 139, 143, 156, 161, 162, 177, 179, 194, 201, 220, 264,
		300, 311, 318, 381, 382, 383, 389, 411, 412, 443, 444, 445,
		464, 465, 497, 500, 512, 513, 514, 515, 520, 521, 530, 540,
		543, 546, 547, 554, 556, 560, 563, 587, 591, 593, 631, 636,
		639, 646, 660, 691, 749, 802, 832, 853, 860, 873, 902, 981,
		989, 990, 992, 993, 994, 995, 1000, 1010, 1025, 1026, 1027,
		1028, 1029, 1080, 1194, 1214, 1241, 1293, 1311, 1337, 1417,
		1433, 1434, 1512, 1527, 1589, 1701, 1723, 1725, 1741, 1755,
		1812, 1813, 1863, 1880, 1883, 1935, 1985, 2000, 2002, 2020,
		2049, 2052, 2053, 2078, 2079, 2080, 2081, 2082, 2083, 2086,
		2087, 2095, 2096, 2100, 2222, 2480, 2483, 2484, 2638, 2745,
		2967, 3000, 3004, 3030, 3050, 3074, 3124, 3127, 3128, 3222,
		3260, 3268, 3283, 3306, 3333, 3389, 3434, 3689, 3690, 3724,
		3784, 3785, 4000, 4040, 4100, 4243, 4333, 4431, 4433, 4443,
		4444, 4567, 4664, 4672, 4711, 4712, 4899, 4993, 5000, 5001,
		5002, 5003, 5004, 5005, 5050, 5060, 5104, 5108, 5190, 5222,
		5223, 5280, 5432, 5500, 5554, 5555, 5631, 5632, 5800, 5900,
		5938, 6000, 6001, 6060, 6112, 6129, 6257, 6346, 6347, 6379,
		6500, 6543, 6566, 6588, 6665, 6666, 6667, 6668, 6669, 6679,
		6697, 6699, 6881, 6891, 6901, 6970, 6999, 7000, 7001, 7002,
		7070, 7071, 7080, 7081, 7212, 7396, 7443, 7474, 7547, 7648,
		7649, 7777, 8000, 8001, 8002, 8003, 8004, 8005, 8008, 8010,
		8014, 8042, 8069, 8080, 8081, 8082, 8084, 8085, 8086, 8087,
		8088, 8089, 8090, 8091, 8099, 8118, 8123, 8172, 8188, 8200,
		8222, 8243, 8280, 8281, 8333, 8383, 8443, 8500, 8767, 8834,
		8880, 8888, 8983, 9000, 9001, 9003, 9009, 9043, 9050, 9060,
		9080, 9090, 9091, 9100, 9101, 9102, 9103, 9119, 9200, 9443,
		9800, 9898, 9981, 9988, 9998, 9999, 10000, 10113, 10114, 10115,
		10116, 10125, 10443, 11371, 12345, 12443, 13720, 13721, 14567,
		15118, 16080, 18091, 18092, 19226, 19638, 20000, 20720, 24800,
		25999, 28017, 30000, 40000, 50000, 54321}
}

//IsOpenPort checks if a port is open (listening)
func IsOpenPort(host string, port string, timeout int) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), time.Duration(timeout)*time.Second)
	if err != nil {
		return false
	}
	if conn != nil {
		defer conn.Close()
		return true
	}
	return false
}

//AsyncPort performs concurrent requests to the specified
//ports range and, if someone is open it prints the results
func AsyncPort(portsArray []int, portsArrayBool bool, StartingPort int,
	EndingPort int, host string, outputFile string, common bool,
	commonPorts []int, plain bool, timeout int) {
	var count int = 0
	var total int = (EndingPort - StartingPort) + 1
	if portsArrayBool {
		total = len(portsArray)
	}
	if common {
		total = len(commonPorts)
	}
	limiter := make(chan string, 200) // Limits simultaneous requests
	wg := sync.WaitGroup{}            // Needed to not prematurely exit before all requests have been finished
	if outputFile != "" {
		if outputFile[len(outputFile)-4:] == "html" {
			output.HeaderHTML("PORT SCANNING", outputFile)
		}
	}
	ports := []int{}
	if !common {
		if portsArrayBool {
			ports = portsArray
		} else {
			for port := StartingPort; port <= EndingPort; port++ {
				ports = append(ports, port)
			}
		}
	} else {
		ports = commonPorts
	}
	for _, port := range ports {
		wg.Add(1)
		portStr := fmt.Sprint(port)
		limiter <- portStr
		if !plain && count%100 == 0 { // update counter
			fmt.Fprint(os.Stdout, "\r \r")
			fmt.Printf("%0.2f%% : %d / %d", utils.Percentage(count, total), count, total)
		}
		go func(portStr string, host string) {
			defer func() { <-limiter }()
			defer wg.Done()
			resp := IsOpenPort(host, portStr, timeout)
			count++
			if resp {
				if !plain {
					fmt.Fprint(os.Stdout, "\r \r")
					fmt.Printf("[+]FOUND: %s ", host)
					color.Green("%s\n", portStr)
				} else {
					fmt.Printf("%s:%s\n", host, portStr)
				}
				if outputFile != "" {
					output.AppendWhere("http://"+host+":"+portStr, "", "PORT", "", outputFile)
				}
			}
		}(portStr, host)
	}
	wg.Wait()
	fmt.Fprint(os.Stdout, "\r \r")
	if outputFile != "" {
		if outputFile[len(outputFile)-4:] == "html" {
			output.FooterHTML(outputFile)
		}
	}
}
