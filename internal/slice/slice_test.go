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

	sliceUtils "github.com/edoardottt/scilla/internal/slice"
	"github.com/stretchr/testify/require"
)

func TestRemoveDuplicateValues(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{
			name:  "single value",
			input: []string{"301"},
			want:  []string{"301"},
		},
		{
			name:  "empty value",
			input: []string{""},
			want:  []string{""},
		},
		{
			name:  "empty slice",
			input: []string{},
			want:  []string{},
		},
		{
			name:  "no duplicates",
			input: []string{"1", "2", "3", "4"},
			want:  []string{"1", "2", "3", "4"},
		},
		{
			name:  "with duplicates",
			input: []string{"1", "2", "3", "4", "4", "3", "2", "1"},
			want:  []string{"1", "2", "3", "4"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sliceUtils.RemoveDuplicateValues(tt.input)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestDifference(t *testing.T) {
	tests := []struct {
		name   string
		inputA []string
		inputB []string
		want   []string
	}{
		{
			name:   "single value",
			inputA: []string{"301"},
			inputB: []string{},
			want:   []string{"301"},
		},
		{
			name:   "empty value",
			inputA: []string{""},
			inputB: []string{""},
			want:   []string{},
		},
		{
			name:   "empty slice",
			inputA: []string{},
			inputB: []string{},
			want:   []string{},
		},
		{
			name:   "no duplicates",
			inputA: []string{"1", "2", "3", "4"},
			inputB: []string{"30"},
			want:   []string{"1", "2", "3", "4"},
		},
		{
			name:   "with duplicates",
			inputA: []string{"1", "2", "3", "4", "4", "3", "2", "1"},
			inputB: []string{"301"},
			want:   []string{"1", "2", "3", "4", "4", "3", "2", "1"},
		},
		{
			name:   "difference",
			inputA: []string{"1", "2", "3", "4"},
			inputB: []string{"1", "2"},
			want:   []string{"3", "4"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sliceUtils.Difference(tt.inputA, tt.inputB)
			require.Equal(t, tt.want, got)
		})
	}
}
