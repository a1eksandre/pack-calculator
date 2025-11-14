package calculator

import (
	"errors"
	"fmt"
	"math"
	"sort"
)

// CalculatePacks determins how many packs are required for a given number of items.
//   - Only whole packs are used.
//   - At least `items` shipped.
//   - Minimise the total items shipped.
//   - Fewest packs possible for the minimised total items.
func CalculatePacks(items int, packSizes []int) (PackResult, error) {
	if items <= 0 {
		return PackResult{}, fmt.Errorf("items must be > 0, got %d", items)
	}
	if len(packSizes) == 0 {
		return PackResult{}, errors.New("at least one pack size is required")
	}

	// Remove duplicates and check values are positive.
	sizes := make([]int, 0, len(packSizes))
	seen := make(map[int]struct{})
	for _, s := range packSizes {
		if s <= 0 {
			return PackResult{}, fmt.Errorf("pack size must be > 0, got %d", s)
		}
		if _, found := seen[s]; found {
			continue
		}
		seen[s] = struct{}{}
		sizes = append(sizes, s)
	}
	if len(sizes) == 0 {
		return PackResult{}, errors.New("no valid pack sizes after cleaning input")
	}

	sort.Ints(sizes)

	maxSize := sizes[len(sizes)-1]
	limit := items + maxSize - 1

	// Guard against blowing up memory.
	if limit <= 0 || limit > 10_000_000 {
		return PackResult{}, fmt.Errorf("items and pack sizes lead to too large limit: %d", limit)
	}

	const INF = math.MaxInt32 / 2

	curr := make([]int, limit+1) // curr[i] = min packs to get exactly i items
	prev := make([]int, limit+1)

	for i := 0; i <= limit; i++ {
		curr[i] = INF
		prev[i] = -1
	}
	curr[0] = 0

	// For each pack size try to extend all totals.
	for _, s := range sizes {
		for total := s; total <= limit; total++ {
			if curr[total-s]+1 < curr[total] {
				curr[total] = curr[total-s] + 1
				prev[total] = s
			}
		}
	}

	// Pick the best total:
	//   1) smallest total >= items
	//   2) if tied, fewest packs.
	bestTotal := -1
	bestPacks := INF

	for total := items; total <= limit; total++ {
		if curr[total] == INF {
			continue
		}
		if bestTotal == -1 ||
			total < bestTotal ||
			(total == bestTotal && curr[total] < bestPacks) {

			bestTotal = total
			bestPacks = curr[total]
		}
	}

	if bestTotal == -1 {
		return PackResult{}, fmt.Errorf("no combination of packs can cover %d items", items)
	}

	// Rebuild the counts from the prev[] array.
	counts := make(map[int]int)
	cur := bestTotal

	for cur > 0 {
		s := prev[cur]
		if s == -1 {
			return PackResult{}, errors.New("internal error while rebuilding solution")
		}
		counts[s]++
		cur -= s
	}

	sol := PackResult{
		Packs:      counts,
		TotalItems: bestTotal,
		ExtraItems: bestTotal - items,
	}
	return sol, nil
}
