package lib

import (
	"math"
	"slices"
	"strconv"
	"strings"
)

type CountedItem struct {
	Name  string
	Count int
}

func (c CountedItem) String() string {
	return c.Name + " (" + strconv.Itoa(c.Count) + ")"
}

func AllocateN(entries []CountedItem, n int) []string {
	total := len(entries)
	if n <= 0 || total == 0 {
		return nil
	}

	sum := n
	for _, e := range entries {
		sum += e.Count
	}

	average := int(math.Ceil(float64(sum) / float64(total)))
	ls := make([]string, 0, n)

	slices.SortStableFunc(entries, func(a, b CountedItem) int { // allocation to start with items having the lowest counts
		if a.Count == b.Count {
			return strings.Compare(a.Name, b.Name)
		}

		return a.Count - b.Count
	})

	for _, e := range entries {
		if e.Count > average { // skip the ones with more than average, they get nothing
			continue
		}

		w := average - e.Count // how many more to add to reach average

		if w > n { // do not exceed what we need
			w = n
		}

		n -= w

		for i := 0; i < w; i++ {
			ls = append(ls, e.Name)
		}

		if n == 0 { // we do not need more
			return ls
		}
	}

	return ls
}
