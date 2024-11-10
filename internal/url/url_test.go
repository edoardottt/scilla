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

	urlUtils "github.com/edoardottt/scilla/internal/url"
	"github.com/stretchr/testify/require"
)

func TestProtocolExists(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "single value",
			input: "1",
			want:  false,
		},
		{
			name:  "good value",
			input: "http://ciao.com",
			want:  true,
		},
		{
			name:  "malformed input",
			input: "http/http/ciao.com/ciao",
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := urlUtils.ProtocolExists(tt.input)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestCleanProtocol(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "single value",
			input: "1",
			want:  "1",
		},
		{
			name:  "good value1",
			input: "http://ciao.com",
			want:  "ciao.com",
		},
		{
			name:  "good value2",
			input: "http://ciao.com/ciao?id=1&ip=1.1.1.1",
			want:  "ciao.com/ciao?id=1&ip=1.1.1.1",
		},
		{
			name:  "malformed input",
			input: "http/http/ciao.com/ciao",
			want:  "http/http/ciao.com/ciao",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := urlUtils.CleanProtocol(tt.input)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestCleanURL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "good value1",
			input: "http://ciao.com",
			want:  "http://ciao.com",
		},
		{
			name:  "good value2",
			input: "http://ciao.com/ciao?id=1&ip=1.1.1.1",
			want:  "http://ciao.com/ciao",
		},
		{
			name:  "malformed input",
			input: "http/http/ciao.com/ciao",
			want:  "http://http/http/ciao.com/ciao",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := urlUtils.CleanURL(tt.input)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestIsURL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "good value1",
			input: "http://ciao.com",
			want:  true,
		},
		{
			name:  "good value2",
			input: "http://ciao.com/ciao?id=1&ip=1.1.1.1",
			want:  true,
		},
		{
			name:  "malformed input1",
			input: "\\",
			want:  false,
		},
		{
			name:  "malformed input2",
			input: "#ciao.com",
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := urlUtils.IsURL(tt.input)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestBuildURL(t *testing.T) {
	tests := []struct {
		name      string
		scheme    string
		subdomain string
		domain    string
		want      string
	}{
		{
			name:      "good value1",
			scheme:    "http",
			subdomain: "sub",
			domain:    "ciao.com",
			want:      "http://sub.ciao.com",
		},
		{
			name:      "good value2",
			scheme:    "https",
			subdomain: "verylongsub.sub.sub",
			domain:    "ciao.com",
			want:      "https://verylongsub.sub.sub.ciao.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := urlUtils.BuildURL(tt.scheme, tt.subdomain, tt.domain)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestAppendDir(t *testing.T) {
	tests := []struct {
		name   string
		scheme string
		domain string
		dir    string
		want1  string
		want2  string
	}{
		{
			name:   "good value",
			scheme: "http",
			domain: "ciao.com",
			dir:    "secret",
			want1:  "http://ciao.com/secret/",
			want2:  "http://ciao.com/secret",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, got2 := urlUtils.AppendDir(tt.scheme, tt.domain, tt.dir)
			require.Equal(t, tt.want1, got1)
			require.Equal(t, tt.want2, got2)
		})
	}
}

func TestCleanSubdomainsOk(t *testing.T) {
	tests := []struct {
		name   string
		target string
		input  []string
		want   []string
	}{
		{
			name:   "good value1",
			target: "ciao.com",
			input:  []string{"sub1.ciao.com", "sub1ciao.com", "ciao.com", ".ciao.com"},
			want:   []string{"sub1.ciao.com"},
		},
		{
			name:   "good value2",
			target: "ciao.com",
			input:  []string{"sub2.sub1.ciao.com", "sub1.ciao.com", "sub1ciao.com", "ciao.com", ".ciao.com"},
			want:   []string{"sub2.sub1.ciao.com", "sub1.ciao.com"},
		},
		{
			name:   "bad value",
			target: "ciao.com",
			input:  []string{"sub1ciao.com", "sub1ciao.com", "ciaocom", ".ciaocom."},
			want:   []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := urlUtils.CleanSubdomainsOk(tt.target, tt.input)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestRetrieveProtocol(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "single value",
			input: "1",
			want:  "",
		},
		{
			name:  "good value1",
			input: "http://ciao.com",
			want:  "http",
		},
		{
			name:  "good value2",
			input: "http://ciao.com/ciao?id=1&ip=1.1.1.1",
			want:  "http",
		},
		{
			name:  "malformed input",
			input: "http/http/ciao.com/ciao",
			want:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := urlUtils.RetrieveProtocol(tt.input)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestAbsoluteURL(t *testing.T) {
	tests := []struct {
		name   string
		scheme string
		target string
		path   string
		want   string
	}{
		{
			name:   "good value",
			scheme: "http",
			target: "ciao.com",
			path:   "secret",
			want:   "http://ciao.com/secret",
		},
		{
			name:   "good value",
			scheme: "https",
			target: "sub.ciao.com",
			path:   "secret?ciao=1",
			want:   "https://sub.ciao.com/secret?ciao=1",
		},
		{
			name:   "absolute path",
			scheme: "https",
			target: "sub.ciao.com",
			path:   "http://ciao.com/secret?ciao=1",
			want:   "http://ciao.com/secret?ciao=1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := urlUtils.AbsoluteURL(tt.scheme, tt.target, tt.path)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestRetrieveHost(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "single value",
			input: "1",
			want:  "",
		},
		{
			name:  "good value1",
			input: "http://ciao.com",
			want:  "ciao.com",
		},
		{
			name:  "good value2",
			input: "http://ciao.com/ciao?id=1&ip=1.1.1.1",
			want:  "ciao.com",
		},
		{
			name:  "malformed input",
			input: "http/http/ciao.com/ciao",
			want:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := urlUtils.RetrieveHost(tt.input)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestGetRootHost(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "single value",
			input: "1",
			want:  "",
		},
		{
			name:  "good value1",
			input: "http://ciao.com",
			want:  "ciao.com",
		},
		{
			name:  "good value2",
			input: "http://sub1.ciao.com",
			want:  "ciao.com",
		},
		{
			name:  "good value3",
			input: "http://sub2.sub1.ciao.com",
			want:  "ciao.com",
		},
		{
			name:  "good value2",
			input: "http://ciao.com/ciao?id=1&ip=1.1.1.1",
			want:  "ciao.com",
		},
		{
			name:  "malformed input",
			input: "http/http/ciao.com/ciao",
			want:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := urlUtils.GetRootHost(tt.input)
			require.Equal(t, tt.want, got)
		})
	}
}
