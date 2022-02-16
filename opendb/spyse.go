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

const searchMethodResultsLimit = 10000
const defaultScrollResultsLimit = 20000

//SpyseSubdomains appends to the subdomains in the list
//the subdomains found with the Spyse service.
func SpyseSubdomains(target string, accessToken string) []string {

	var result []string

	client, err := spyse.NewClient(accessToken, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	svc := spyse.NewDomainService(client)

	//Dot before the domain name is important because search fetch any domains ending with ".$target"
	var searchDomain = "." + target
	var subdomainsSearchParams spyse.QueryBuilder
	var ctx = context.Background()

	subdomainsSearchParams.AppendParam(spyse.QueryParam{
		Name:     svc.Params().Name.Name,
		Operator: svc.Params().Name.Operator.EndsWith,
		Value:    searchDomain,
	})

	totalResults, err := svc.SearchCount(ctx, subdomainsSearchParams.Query)
	if err != nil {
		log.Fatal(err.Error())
	}

	if totalResults == 0 {
		return result
	}

	// The default "Search" method returns only first 10 000 subdomains.
	// To obtain more than 10 000 subdomains the "Scroll" method should be used.
	// Note: The "Scroll" method is only available for "PRO" customers, so we need to check
	// quota.IsScrollSearchEnabled parameter.
	if totalResults > searchMethodResultsLimit && client.Account().IsScrollSearchEnabled {
		var scrollID string
		var scrollResults *spyse.DomainScrollResponse

		for {
			if scrollResults, err = svc.ScrollSearch(ctx, subdomainsSearchParams.Query, scrollID); err != nil {
				if len(result) > 0 {
					spyseErr, ok := err.(*spyse.ErrResponse)
					if ok && spyseErr.Err.Code == spyse.CodeRequestsLimitReached {
						break
					}
				}
				log.Fatal(err.Error())
			}
			if len(scrollResults.Items) > 0 {
				scrollID = scrollResults.SearchID

				for _, domain := range scrollResults.Items {
					result = append(result, domain.Name)
				}
				// The default "Scroll" limit, to avoid results that more than 20k
				// If a limit of the number of requested subdomains will be added this block should be changed
				if len(result) > defaultScrollResultsLimit {
					break
				}
			}
		}
	} else {
		var limit = 100
		var searchResults []spyse.Domain

		for offset := 0; int64(offset) < totalResults && int64(offset) < searchMethodResultsLimit; offset += limit {
			if searchResults, err = svc.Search(ctx, subdomainsSearchParams.Query, limit, offset); err != nil {
				if len(result) > 0 {
					spyseErr, ok := err.(*spyse.ErrResponse)
					if ok && spyseErr.Err.Code == spyse.CodeRequestsLimitReached {
						break
					}
				}
				log.Fatal(err.Error())
			}

			for _, domain := range searchResults {
				result = append(result, domain.Name)
			}
		}
	}
	return result
}
