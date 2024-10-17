// Copyright 2024 The Ebitengine Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vector

import (
	"testing"
)

func TestIsPointCloseToSegment(t *testing.T) {
	testCases := []struct {
		p     Point
		p0    Point
		p1    Point
		allow float32
		want  bool
	}{
		{
			p:     Point{0.5, 0.5},
			p0:    Point{0, 0},
			p1:    Point{1, 0},
			allow: 1,
			want:  true,
		},
		{
			p:     Point{0.5, 1.5},
			p0:    Point{0, 0},
			p1:    Point{1, 0},
			allow: 1,
			want:  false,
		},
		{
			p:     Point{0.5, 0.5},
			p0:    Point{0, 0},
			p1:    Point{1, 1},
			allow: 0,
			want:  true,
		},
		{
			p:     Point{0, 1},
			p0:    Point{0, 0},
			p1:    Point{1, 1},
			allow: 0.7,
			want:  false,
		},
		{
			p:     Point{0, 1},
			p0:    Point{0, 0},
			p1:    Point{1, 1},
			allow: 0.8,
			want:  true,
		},
		{
			// p0 and p1 are the same.
			p:     Point{0, 1},
			p0:    Point{0.5, 0.5},
			p1:    Point{0.5, 0.5},
			allow: 0.7,
			want:  false,
		},
		{
			// p0 and p1 are the same.
			p:     Point{0, 1},
			p0:    Point{0.5, 0.5},
			p1:    Point{0.5, 0.5},
			allow: 0.8,
			want:  true,
		},
	}
	for _, tc := range testCases {
		if got := IsPointCloseToSegment(tc.p, tc.p0, tc.p1, tc.allow); got != tc.want {
			t.Errorf("got: %v, want: %v", got, tc.want)
		}
	}
}
