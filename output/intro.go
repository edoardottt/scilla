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

import (
	"fmt"

	"github.com/fatih/color"
)

//Intro prints the banner when the program starts
func Intro() {
	banner1 := "                  _ _ _\n"
	banner2 := "         ___  ___(_) | | __ _\n"
	banner3 := "        / __|/ __| | | |/ _` |\n"
	banner4 := "        \\__ \\ (__| | | | (_| |\n"
	banner5 := "        |___/\\___|_|_|_|\\__,_| v1.2.1\n"
	banner6 := " > github.com/edoardottt/scilla\n"
	banner7 := " > edoardoottavianelli.it"
	bannerPart1 := banner1 + banner2 + banner3 + banner4 + banner5
	bannerPart2 := banner6 + banner7
	color.Cyan("%s\n", bannerPart1)
	fmt.Println(bannerPart2)
	fmt.Println("========================================")
}
