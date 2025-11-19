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
	"fmt"
	"testing"

	transportUtils "github.com/edoardottt/scilla/internal/transport"
	"github.com/stretchr/testify/require"
)

func TestCheckPortsArray(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []int
		erro  error
	}{
		{
			name:  "single value",
			input: "1",
			want:  []int{1},
			erro:  nil,
		},
		{
			name:  "invalid value1",
			input: "1,123456",
			want:  nil,
			erro:  fmt.Errorf("%w", transportUtils.ErrInvalidArray),
		},
		{
			name:  "invalid value2",
			input: ",,",
			want:  nil,
			erro:  fmt.Errorf("%w", transportUtils.ErrInvalidArray),
		},
		{
			name:  "good array",
			input: "1,10",
			want:  []int{1, 10},
			erro:  nil,
		},
		{
			name:  "with duplicates",
			input: "1,1,10",
			want:  []int{1, 10},
			erro:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := transportUtils.CheckPortsArray(tt.input)
			require.Equal(t, tt.erro, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestCheckPortsRange(t *testing.T) {
	tests := []struct {
		name  string
		input string
		start int
		end   int
		erro  error
	}{
		{
			name:  "single value",
			input: "1",
			start: 1,
			end:   1,
			erro:  nil,
		},
		{
			name:  "invalid value",
			input: "1-123456",
			start: 0,
			end:   0,
			erro:  fmt.Errorf("%w", transportUtils.ErrInvalidRange),
		},
		{
			name:  "good array",
			input: "1-10",
			start: 1,
			end:   10,
			erro:  nil,
		},
		{
			name:  "end array",
			input: "-10",
			start: 1,
			end:   10,
			erro:  nil,
		},
		{
			name:  "start array",
			input: "1-",
			start: 1,
			end:   65535,
			erro:  nil,
		},
		{
			name:  "bad array1",
			input: "1--",
			start: 0,
			end:   0,
			erro:  fmt.Errorf("%w", transportUtils.ErrInvalidRange),
		},
		{
			name:  "bad array2",
			input: "-1-",
			start: 0,
			end:   0,
			erro:  fmt.Errorf("%w", transportUtils.ErrInvalidRange),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end, err := transportUtils.CheckPortsRange(tt.input, 1, 65535)
			require.Equal(t, tt.erro, err)
			require.Equal(t, tt.start, start)
			require.Equal(t, tt.end, end)
		})
	}
}
