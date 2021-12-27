/*
 * Copyright (c) 2021 Andreas Signer <asigner@gmail.com>
 *
 * This file is part of gosfxr.
 *
 * gosfxr is free software: you can redistribute it and/or
 * modify it under the terms of the GNU General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * gosfxr is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with gosfxr.  If not, see <http://www.gnu.org/licenses/>.
 */
package ui

import (
	"testing"
)

func Test_fixExtensions(t *testing.T) {
	type args struct {
		filename string
		ext      string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Extensions is added",
			args: args{
				filename: "foo",
				ext:      ".bar",
			},
			want: "foo.bar",
		},
		{
			name: "Extensions is not overridden",
			args: args{
				filename: "foo.bar",
				ext:      ".baz",
			},
			want: "foo.bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fixExtensions(tt.args.filename, tt.args.ext); got != tt.want {
				t.Errorf("fixExtensions() = %v, want %v", got, tt.want)
			}
		})
	}
}
