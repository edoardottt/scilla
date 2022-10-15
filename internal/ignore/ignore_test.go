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

package utils_test

import (
	"testing"

	ignoreUtils "github.com/edoardottt/scilla/internal/ignore"
	"github.com/stretchr/testify/require"
)

func TestCheckIgnore(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
		erro  error
	}{
		{
			name:  "single value",
			input: "301",
			want:  []string{"301"},
			erro:  nil,
		},
		{
			name:  "empty value",
			input: "",
			want:  nil,
			erro:  ignoreUtils.ErrWrongStatusCodeLength,
		},
		{
			name:  "withous duplicates",
			input: "301,302,400",
			want:  []string{"301", "302", "400"},
			erro:  nil,
		},
		{
			name:  "has duplicates",
			input: "301,301,302,301",
			want:  []string{"301", "302"},
			erro:  nil,
		},
		{
			name:  "class",
			input: "3**",
			want:  []string{"3**"},
			erro:  nil,
		},
		{
			name:  "multiple classes",
			input: "3**,4**",
			want:  []string{"3**", "4**"},
			erro:  nil,
		},
		{
			name:  "mixed",
			input: "301,4**",
			want:  []string{"301", "4**"},
			erro:  nil,
		},
		{
			name:  "mixed with duplicates",
			input: "301,4**,4**",
			want:  []string{"301", "4**"},
			erro:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ignoreUtils.CheckIgnore(tt.input)
			require.Equal(t, tt.erro, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestIgnoreResponse(t *testing.T) {
	tests := []struct {
		name     string
		response int
		ignore   []string
		want     bool
	}{
		{
			name:     "single value",
			response: 200,
			ignore:   []string{"301"},
			want:     false,
		},
		{
			name:     "empty value",
			response: 0,
			ignore:   []string{"301"},
			want:     false,
		},
		{
			name:     "empty ignore",
			response: 200,
			ignore:   []string{},
			want:     false,
		},
		{
			name:     "invalid response",
			response: 44545345,
			ignore:   []string{},
			want:     false,
		},
		{
			name:     "withous duplicates",
			response: 301,
			ignore:   []string{"301", "302", "400"},
			want:     true,
		},
		{
			name:     "has duplicates",
			response: 301,
			ignore:   []string{"301", "301", "302"},
			want:     true,
		},
		{
			name:     "class",
			response: 301,
			ignore:   []string{"3**"},
			want:     true,
		},
		{
			name:     "multiple classes",
			response: 401,
			ignore:   []string{"3**", "4**"},
			want:     true,
		},
		{
			name:     "mixed",
			response: 401,
			ignore:   []string{"301", "4**"},
			want:     true,
		},
		{
			name:     "mixed with duplicates",
			response: 500,
			ignore:   []string{"301", "4**", "5**", "5**"},
			want:     true,
		},
		{
			name:     "don't ignore",
			response: 500,
			ignore:   []string{"301", "4**", "302", "429"},
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ignoreUtils.IgnoreResponse(tt.response, tt.ignore)
			require.Equal(t, tt.want, got)
		})
	}
}
