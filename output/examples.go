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

import "fmt"

//Examples prints some examples
func Examples() {
	fmt.Println("	Examples:")
	fmt.Println("		- scilla dns -target target.domain")
	fmt.Println("		- scilla dns -target -o txt target.domain")
	fmt.Println("		- scilla dns -target -o html target.domain")
	fmt.Println("		- scilla dns -target -o json target.domain")
	fmt.Println("		- scilla dns -target -plain target.domain")
	fmt.Println()
	fmt.Println("		- scilla subdomain -target target.domain")
	fmt.Println("		- scilla subdomain -w wordlist.txt -target target.domain")
	fmt.Println("		- scilla subdomain -o txt -target target.domain")
	fmt.Println("		- scilla subdomain -o html -target target.domain")
	fmt.Println("		- scilla subdomain -o json -target target.domain")
	fmt.Println("		- scilla subdomain -i 400 -target target.domain")
	fmt.Println("		- scilla subdomain -i 4** -target target.domain")
	fmt.Println("		- scilla subdomain -c -target target.domain")
	fmt.Println("		- scilla subdomain -db -target target.domain")
	fmt.Println("		- scilla subdomain -plain -target target.domain")
	fmt.Println("		- scilla subdomain -db -no-check -target target.domain")
	fmt.Println("		- scilla subdomain -db -spyse -target target.domain")
	fmt.Println("		- scilla subdomain -db -vt -target target.domain")
	fmt.Println()
	fmt.Println("		- scilla port -p -450 -target target.domain")
	fmt.Println("		- scilla port -p 90- -target target.domain")
	fmt.Println("		- scilla port -p 10-1000 -target target.domain")
	fmt.Println("		- scilla port -o txt -target target.domain")
	fmt.Println("		- scilla port -o html -target target.domain")
	fmt.Println("		- scilla port -o json -target target.domain")
	fmt.Println("		- scilla port -p 21,25,80 -target target.domain")
	fmt.Println("		- scilla port -common -target target.domain")
	fmt.Println("		- scilla port -plain -target target.domain")
	fmt.Println("		- scilla port -t 2 -target target.domain")
	fmt.Println()
	fmt.Println("		- scilla dir -target target.domain")
	fmt.Println("		- scilla dir -o txt -target target.domain")
	fmt.Println("		- scilla dir -o html -target target.domain")
	fmt.Println("		- scilla dir -o json -target target.domain")
	fmt.Println("		- scilla dir -w wordlist.txt -target target.domain")
	fmt.Println("		- scilla dir -i 500,401 -target target.domain")
	fmt.Println("		- scilla dir -i 5**,401 -target target.domain")
	fmt.Println("		- scilla dir -c -target target.domain")
	fmt.Println("		- scilla dir -plain -target target.domain")
	fmt.Println("		- scilla dir -nr -target target.domain")
	fmt.Println()
	fmt.Println("		- scilla report -p 80 -target target.domain")
	fmt.Println("		- scilla report -o txt -target target.domain")
	fmt.Println("		- scilla report -o html -target target.domain")
	fmt.Println("		- scilla report -o json -target target.domain")
	fmt.Println("		- scilla report -p 50-200 -target target.domain")
	fmt.Println("		- scilla report -wd dirs.txt -target target.domain")
	fmt.Println("		- scilla report -ws subdomains.txt -target target.domain")
	fmt.Println("		- scilla report -id 500,501,502 -target target.domain")
	fmt.Println("		- scilla report -is 500,501,502 -target target.domain")
	fmt.Println("		- scilla report -id 5**,4** -target target.domain")
	fmt.Println("		- scilla report -is 5**,4** -target target.domain")
	fmt.Println("		- scilla report -cd -target target.domain")
	fmt.Println("		- scilla report -cs -target target.domain")
	fmt.Println("		- scilla report -db -target target.domain")
	fmt.Println("		- scilla report -p 21,25,80 -target target.domain")
	fmt.Println("		- scilla report -common -target target.domain")
	fmt.Println("		- scilla report -nr -target target.domain")
	fmt.Println("		- scilla report -db -spyse -target target.domain")
	fmt.Println("		- scilla report -db -vt -target target.domain")
	fmt.Println("")
}
