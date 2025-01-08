package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllocateN(t *testing.T) {
	tests := []struct {
		name  string
		items []CountedItem
		n     int
		want  []string
	}{
		{
			name:  "no entries",
			items: []CountedItem{},
			n:     1,
			want:  nil,
		},
		{
			name: "n = 0",
			items: []CountedItem{
				{
					Name:  "a",
					Count: 1,
				},
				{
					Name:  "b",
					Count: 2,
				},
			},
			n:    0,
			want: nil,
		},
		{
			name: "two, one 1, need 1",
			items: []CountedItem{
				{
					Name:  "a",
					Count: 1,
				},
				{
					Name:  "b",
					Count: 0,
				},
			},
			n:    1,
			want: []string{"b"},
		},
		{
			name: "three, single 1, need 3",
			items: []CountedItem{
				{
					Name:  "a",
					Count: 1,
				},
				{
					Name:  "b",
					Count: 0,
				},
				{
					Name:  "c",
					Count: 0,
				},
			},
			n:    3,
			want: []string{"b", "b", "c"},
		},
		{
			name: "3 all zeroes, need 3",
			items: []CountedItem{
				{
					Name:  "a",
					Count: 0,
				},
				{
					Name:  "b",
					Count: 0,
				},
				{
					Name:  "c",
					Count: 0,
				},
			},
			n:    3,
			want: []string{"a", "b", "c"},
		},
		{
			name: "n less than number of entries",
			items: []CountedItem{
				{
					Name:  "a",
					Count: 2,
				},
				{
					Name:  "b",
					Count: 3,
				},
				{
					Name:  "c",
					Count: 1,
				},
			},
			n:    1,
			want: []string{"c"},
		},
		{
			name: "n greater than number of entries",
			items: []CountedItem{
				{
					Name:  "a",
					Count: 2,
				},
				{
					Name:  "b",
					Count: 1,
				},
				{
					Name:  "c",
					Count: 3,
				},
			},
			n:    4,
			want: []string{"b", "b", "b", "a"},
		},
		{
			name: "n greater than number of entries",
			items: []CountedItem{
				{
					Name:  "a",
					Count: 2,
				},
				{
					Name:  "b",
					Count: 1,
				},
				{
					Name:  "c",
					Count: 3,
				},
			},
			n:    5,
			want: []string{"b", "b", "b", "a", "a"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AllocateN(tt.items, tt.n)

			assert.Equal(t, tt.want, got)
		})
	}
}
