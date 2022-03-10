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

//CommonPorts >
func CommonPorts() []int {
	return []int{13, 20, 21, 22, 23, 25, 42, 50, 51, 53, 67, 68,
		69, 70, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 90, 102,
		107, 109, 110, 111, 113, 115, 118, 119, 123, 135, 136, 137,
		138, 139, 143, 156, 161, 162, 179, 194, 220, 300, 311, 389,
		443, 444, 445, 464, 500, 512, 513, 514, 515, 530, 543, 546,
		547, 556, 587, 591, 593, 631, 636, 660, 749, 802, 832, 853,
		873, 902, 981, 989, 990, 992, 993, 994, 995, 1000, 1010, 1025,
		1080, 1194, 1241, 1293, 1311, 1337, 1417, 1433, 1434, 1527, 1755,
		1812, 1813, 1880, 1883, 1935, 2000, 2020, 2049, 2052, 2053, 2078,
		2079, 2080, 2082, 2083, 2086, 2087, 2095, 2096, 2222, 2480, 2483,
		2484, 2638, 3000, 3004, 3030, 3128, 3268, 3283, 3306, 3333, 3389,
		3434, 4000, 4040, 4100, 4243, 4431, 4433, 4443, 4444, 4567, 4711,
		4712, 4993, 5000, 5050, 5104, 5108, 5280, 5432, 5555, 5800, 5938,
		6000, 6060, 6379, 6543, 6666, 7000, 7001, 7002, 7070, 7071, 7080,
		7081, 7396, 7443, 7474, 7547, 7777, 8000, 8001, 8002, 8003, 8004,
		8005, 8008, 8010, 8014, 8042, 8069, 8080, 8081, 8082, 8084, 8085,
		8087, 8088, 8089, 8090, 8091, 8099, 8118, 8123, 8172, 8188, 8200,
		8222, 8243, 8280, 8281, 8333, 8383, 8443, 8500, 8834, 8880, 8888,
		8983, 9000, 9001, 9003, 9009, 9043, 9050, 9060, 9080, 9090, 9091,
		9200, 9443, 9800, 9981, 9090, 9998, 9999, 10000, 10125, 10443, 12443,
		16080, 18091, 18092, 20000, 20720, 28017, 50000, 54321}
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
