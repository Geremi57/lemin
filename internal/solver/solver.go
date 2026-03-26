package solver

import "lem-in/models"

// FindShortestPath finds the shortest path using BFS
func FindShortestPath(colony *models.Colony) []*models.Room {
	if colony.Start == nil || colony.End == nil {
		return nil
	}

	type state struct {
		room *models.Room
		path []*models.Room
	}

	visited := make(map[string]bool)
	visited[colony.Start.Name] = true
	queue := []state{{
		room: colony.Start,
		path: []*models.Room{colony.Start},
	}}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.room == colony.End {
			return current.path
		}

		for _, neighbor := range current.room.Connected {
			if visited[neighbor.Name] {
				continue
			}
			visited[neighbor.Name] = true
			newPath := make([]*models.Room, len(current.path)+1)
			copy(newPath, current.path)
			newPath[len(current.path)] = neighbor
			queue = append(queue, state{
				room: neighbor,
				path: newPath,
			})
		}
	}
	return nil
}

// FindDisjointPaths uses Edmonds-Karp (BFS max-flow) on a vertex-split graph
// to find the maximum set of vertex-disjoint paths from start to end.
//
// Each room r is split into two nodes:
//   rIn  = 2*i     (entry node)
//   rOut = 2*i + 1 (exit node)
//
// Internal edge rIn->rOut enforces at most one path per intermediate room.
// Start and end get capacity n (effectively unlimited).
// Tunnel u<->v becomes two directed edges: out(u)->in(v) and out(v)->in(u).
func FindDisjointPaths(colony *models.Colony) [][]*models.Room {
	if colony.Start == nil || colony.End == nil {
		return nil
	}

	// Index all rooms deterministically
	rooms := make([]*models.Room, 0, len(colony.Rooms))
	roomIndex := make(map[string]int)
	for _, r := range colony.Rooms {
		roomIndex[r.Name] = len(rooms)
		rooms = append(rooms, r)
	}

	n := len(rooms)
	size := 2 * n

	// ----- Flow graph -----
	type edge struct {
		to, cap, rev int
	}
	graph := make([][]edge, size)

	addEdge := func(u, v, cap int) {
		graph[u] = append(graph[u], edge{v, cap, len(graph[v])})
		graph[v] = append(graph[v], edge{u, 0, len(graph[u]) - 1})
	}

	startIdx := roomIndex[colony.Start.Name]
	endIdx := roomIndex[colony.End.Name]

	// Internal edges rIn -> rOut
	for _, r := range rooms {
		i := roomIndex[r.Name]
		cap := 1
		if r == colony.Start || r == colony.End {
			cap = n // unlimited
		}
		addEdge(2*i, 2*i+1, cap)
	}

	// Tunnel edges: add both directions, but deduplicate so we don't double-add
	addedTunnels := make(map[[2]int]bool)
	for _, r := range rooms {
		u := roomIndex[r.Name]
		for _, nb := range r.Connected {
			v := roomIndex[nb.Name]
			key := [2]int{u, v}
			revKey := [2]int{v, u}
			if !addedTunnels[key] && !addedTunnels[revKey] {
				// out(u) -> in(v)
				addEdge(2*u+1, 2*v, 1)
				// out(v) -> in(u)
				addEdge(2*v+1, 2*u, 1)
				addedTunnels[key] = true
			}
		}
	}

	source := 2 * startIdx
	sink := 2*endIdx + 1

	// ----- Edmonds-Karp BFS augmentation -----
	bfs := func() bool {
		prev := make([]int, size)
		prevEdge := make([]int, size)
		for i := range prev {
			prev[i] = -1
		}
		prev[source] = source
		queue := []int{source}
		for len(queue) > 0 && prev[sink] == -1 {
			u := queue[0]
			queue = queue[1:]
			for i, e := range graph[u] {
				if e.cap > 0 && prev[e.to] == -1 {
					prev[e.to] = u
					prevEdge[e.to] = i
					queue = append(queue, e.to)
				}
			}
		}
		if prev[sink] == -1 {
			return false
		}
		// Augment
		v := sink
		for v != source {
			u := prev[v]
			ei := prevEdge[v]
			graph[u][ei].cap--
			graph[v][graph[u][ei].rev].cap++
			v = u
		}
		return true
	}

	for bfs() {
		// exhaust all augmenting paths
	}

	// ----- Extract paths from flow -----
	// For each room u, collect all rooms v where the tunnel edge out(u)->in(v)
	// has flow (original cap was 1, now 0).
	// We use a slice per node so multiple paths from the same node are supported.
	flowNext := make(map[int][]int) // roomIndex -> list of next roomIndices with flow

	for _, r := range rooms {
		u := roomIndex[r.Name]
		for _, e := range graph[2*u+1] {
			// Tunnel edges land on in-nodes (even indices) of other rooms
			if e.to%2 == 0 {
				v := e.to / 2
				if v == u {
					continue // skip self-loops if any
				}
				// Flow used = original cap (1) - current cap = 0
				if e.cap == 0 {
					flowNext[u] = append(flowNext[u], v)
				}
			}
		}
	}

	// Trace paths greedily from startIdx, consuming each flow edge once
	var result [][]*models.Room

	for len(flowNext[startIdx]) > 0 {
		path := []*models.Room{colony.Start}
		cur := startIdx

		for cur != endIdx {
			nexts, ok := flowNext[cur]
			if !ok || len(nexts) == 0 {
				// Dead end — should not happen in a valid flow
				break
			}
			// Consume the first available next hop
			nxt := nexts[0]
			if len(nexts) == 1 {
				delete(flowNext, cur)
			} else {
				flowNext[cur] = nexts[1:]
			}
			cur = nxt
			path = append(path, rooms[cur])
		}

		if cur == endIdx {
			result = append(result, path)
		}
	}

	// Fallback to single shortest path if extraction yielded nothing
	if len(result) == 0 {
		sp := FindShortestPath(colony)
		if sp != nil {
			return [][]*models.Room{sp}
		}
	}

	return result
}

func GetShortestPathLength(colony *models.Colony) int {
	path := FindShortestPath(colony)
	if path == nil {
		return 0
	}
	return len(path)
}