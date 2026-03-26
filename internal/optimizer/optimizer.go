package optimizer

import (
	"lem-in/internal/simulator"
	"lem-in/models"
	"sort"
)

type Optimizer struct {
	paths    [][]*models.Room
	antCount int
}

func NewOptimizer(paths [][]*models.Room, antCount int) *Optimizer {
	return &Optimizer{
		paths:    paths,
		antCount: antCount,
	}
}

// FindOptimalAssignment distributes ants across vertex-disjoint paths
// to minimise the number of turns.
//
// Formula: for a set of paths with lengths l1 ≤ l2 ≤ … ≤ lk,
// the total turns when path i gets ai ants is:
//
//	turns = max_i( li - 1 + ai )       (li-1 moves + ai-1 queuing delays + 1)
//
// To minimise the maximum, we raise the shorter paths' ant counts until
// their finish time equals the longer paths, then dump remaining ants on
// the shortest path.
func (o *Optimizer) FindOptimalAssignment() []simulator.PathAssignment {
	if len(o.paths) == 0 || o.antCount == 0 {
		return nil
	}

	// Sort paths by length ascending
	paths := make([][]*models.Room, len(o.paths))
	copy(paths, o.paths)
	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) < len(paths[j])
	})

	// Try every prefix of paths [0..k-1] and pick the best total turns
	bestTurns := int(^uint(0) >> 1) // MaxInt
	bestCounts := make([]int, len(paths))

	for k := 1; k <= len(paths); k++ {
		counts := o.distributeAnts(paths[:k])
		turns := o.calcTurns(paths[:k], counts)
		if turns < bestTurns {
			bestTurns = turns
			bestCounts = make([]int, len(paths))
			copy(bestCounts, counts)
		}
	}

	// Build assignment, skip paths with 0 ants
	result := make([]simulator.PathAssignment, 0, len(paths))
	for i, p := range paths {
		if bestCounts[i] > 0 {
			result = append(result, simulator.PathAssignment{
				Path:     p,
				AntCount: bestCounts[i],
			})
		}
	}
	return result
}

// distributeAnts computes the optimal ant counts for the given sorted paths.
// The idea: assign ants so that all paths finish at roughly the same turn.
// turns(path_i) = (len(path_i) - 1) + ant_count_i
// We want all of these equal to some T.
// ant_count_i = T - (len(path_i) - 1)
// Sum of ant_counts = n  =>  k*T - sum(len-1) = n
// T = (n + sum(len-1)) / k  (ceiling)
func (o *Optimizer) distributeAnts(paths [][]*models.Room) []int {
	k := len(paths)
	sumLenMinus1 := 0
	for _, p := range paths {
		sumLenMinus1 += len(p) - 1
	}

	// T = ceil((antCount + sumLenMinus1) / k)
	T := (o.antCount + sumLenMinus1 + k - 1) / k

	counts := make([]int, k)
	remaining := o.antCount
	for i, p := range paths {
		c := T - (len(p) - 1)
		if c < 0 {
			c = 0
		}
		if c > remaining {
			c = remaining
		}
		counts[i] = c
		remaining -= c
	}

	// If any ants are left over (due to rounding), assign to shortest path
	if remaining > 0 {
		counts[0] += remaining
	}

	return counts
}

func (o *Optimizer) calcTurns(paths [][]*models.Room, counts []int) int {
	max := 0
	for i, p := range paths {
		if counts[i] == 0 {
			continue
		}
		t := (len(p) - 1) + counts[i]
		if t > max {
			max = t
		}
	}
	return max
}