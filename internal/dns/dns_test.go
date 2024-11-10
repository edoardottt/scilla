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

package utils_test

import (
	"testing"

	dnsUtils "github.com/edoardottt/scilla/internal/dns"
	"github.com/stretchr/testify/require"
)

func TestSimpleDNSLookup(t *testing.T) {
	require.Equal(t, false, dnsUtils.SimpleDNSLookup("whtbiwvuwivtytyiethiubwixfniqf.wefbywf.com"))
}

func TestNewCustomResolver(t *testing.T) {
	r := dnsUtils.NewCustomResolver("8.8.8.8")
	require.NotNil(t, dnsUtils.CustomDNSLookup(r, "whtbiwvuwivtytyiethiubwixfniqf.wefbywf.com"))
}

func TestCustomDNSLookup(t *testing.T) {
	r := dnsUtils.NewCustomResolver("8.8.8.8")
	require.Equal(t, false, dnsUtils.CustomDNSLookup(r, "whtbiwvuwivtytyiethiubwixfniqf.wefbywf.com"))
}
