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

package opendb

import (
	"context"
	"log"

	spyse "github.com/spyse-com/go-spyse/pkg"
)

//SpyseSubdomains appends to the subdomains in the list
//thr subdomains found with the Spyse service.
func SpyseSubdomains(target string, accessToken string) []string {

	var result []string

	client, err := spyse.NewClient(accessToken, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	svc := spyse.NewDomainService(client)

	//Dot before the domain name is important because search fetch any domains that end with ".$target"
	var searchDomain = "." + target
	var subdomainsSearchParams spyse.QueryBuilder

	subdomainsSearchParams.AppendParam(spyse.QueryParam{
		Name:     svc.Params().Name.Name,
		Operator: svc.Params().Name.Operator.EndsWith,
		Value:    searchDomain,
	})

	countResults, err := svc.SearchCount(context.Background(), subdomainsSearchParams.Query)
	if err != nil {
		log.Fatal(err.Error())
	}

	var limit = 100
	var offset = 0
	var searchResults []spyse.Domain
	var domain spyse.Domain
	for ; int64(offset) < countResults; offset += limit {
		//Notice that you can fetch only the first 10000 (can depend on your subscription plan) results using the Search method
		searchResults, err = svc.Search(context.Background(), subdomainsSearchParams.Query, limit, offset)
		if err != nil {
			log.Fatal(err.Error())
		}
		for _, domain = range searchResults {
			result = append(result, domain.Name)
		}
	}
	return result
}
