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

package opendb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	httpUtils "github.com/edoardottt/scilla/internal/http"
)

type builtWithResultWrapper struct {
	Results []struct {
		Result struct {
			SpendHistory []struct {
				D int64 `json:"D"`
				S int   `json:"S"`
			} `json:"SpendHistory"`
			IsDB  string `json:"IsDB"`
			Spend int    `json:"Spend"`
			Paths []struct {
				Technologies []struct {
					Name          string   `json:"Name"`
					Description   string   `json:"Description"`
					Link          string   `json:"Link"`
					Tag           string   `json:"Tag"`
					FirstDetected int64    `json:"FirstDetected"`
					LastDetected  int64    `json:"LastDetected"`
					IsPremium     string   `json:"IsPremium"`
					Categories    []string `json:"Categories,omitempty"`
					Parent        string   `json:"Parent,omitempty"`
				} `json:"Technologies"`
				FirstIndexed int64  `json:"FirstIndexed"`
				LastIndexed  int64  `json:"LastIndexed"`
				Domain       string `json:"Domain"`
				URL          string `json:"Url"`
				SubDomain    string `json:"SubDomain"`
			} `json:"Paths"`
		} `json:"Result"`
		Meta struct {
			Majestic    int      `json:"Majestic"`
			Umbrella    int      `json:"Umbrella"`
			Vertical    string   `json:"Vertical"`
			Social      []string `json:"Social"`
			CompanyName string   `json:"CompanyName"`
			Telephones  []string `json:"Telephones"`
			Emails      []string `json:"Emails"`
			City        string   `json:"City"`
			State       string   `json:"State"`
			Postcode    string   `json:"Postcode"`
			Country     string   `json:"Country"`
			Names       []struct {
				Name  string `json:"Name"`
				Type  int    `json:"Type"`
				Email string `json:"Email"`
			} `json:"Names"`
			ARank int `json:"ARank"`
			QRank int `json:"QRank"`
		} `json:"Meta"`
		Attributes struct {
			Employees    int `json:"Employees"`
			MJRank       int `json:"MJRank"`
			MJTLDRank    int `json:"MJTLDRank"`
			RefSN        int `json:"RefSN"`
			RefIP        int `json:"RefIP"`
			Followers    int `json:"Followers"`
			Sitemap      int `json:"Sitemap"`
			GTMTags      int `json:"GTMTags"`
			QubitTags    int `json:"QubitTags"`
			TealiumTags  int `json:"TealiumTags"`
			AdobeTags    int `json:"AdobeTags"`
			CDimensions  int `json:"CDimensions"`
			CGoals       int `json:"CGoals"`
			CMetrics     int `json:"CMetrics"`
			ProductCount int `json:"ProductCount"`
		} `json:"Attributes"`
		FirstIndexed int64  `json:"FirstIndexed"`
		LastIndexed  int64  `json:"LastIndexed"`
		Lookup       string `json:"Lookup"`
		SalesRevenue int    `json:"SalesRevenue"`
	} `json:"Results"`
	Errors []struct {
		Code    int    `json:"Code"`
		Message string `json:"Message"`
		Lookup  string `json:"Lookup"`
	} `json:"Errors"`
	Trust string `json:"Trust"`
}

// Builtwith retrieves from the url below some known subdomains.
func BuiltWithSubdomains(target, apikey string, plain bool) []string {
	result := []string{}

	if !plain {
		fmt.Println("Pulling data from BuiltWith")
	}

	client := http.Client{
		Timeout: httpUtils.Seconds30,
	}

	fetchURL := fmt.Sprintf(
		"https://api.builtwith.com/v21/api.json?&KEY=%s&LOOKUP=%s",
		apikey, target,
	)

	resp, err := client.Get(fetchURL)
	if err != nil {
		return result
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result
	}

	var bwrWrapper builtWithResultWrapper
	if err := json.Unmarshal(body, &bwrWrapper); err != nil {
		return result
	}

	if len(bwrWrapper.Results) == 0 {
		return result
	}

	for _, elem := range bwrWrapper.Results[0].Result.Paths {
		if elem.SubDomain == "" {
			continue
		}

		fullDomain := elem.SubDomain + "." + target
		result = append(result, fullDomain)
	}

	return result
}
