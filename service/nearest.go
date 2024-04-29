package service

import (
	"fmt"
	"math"
	"slices"
	"sort"

	h "github.com/Vorto-interview/helper"
	"github.com/Vorto-interview/models"
)

type NearestNeighborSolver struct {
	loads []models.Load
}

func NewNearestNeighborSolver(loads []models.Load) NearestNeighborSolver {
	solver := NearestNeighborSolver{
		loads: loads,
	}

	// Add our depot as a special load that starts and ends at the depot
	solver.loads = append([]models.Load{{Number: 0, Pickup: models.StartPoint(), Dropoff: models.Point{}}}, solver.loads...)

	//fmt.Printf("Nearest Neighbor Solver built with %v loads\n", len(loads))

	return solver
}

func (n NearestNeighborSolver) PlanRoutes() ([]h.Route, error) {
	for _, l := range n.loads {
		if models.StartPoint().DistanceTo(l.Pickup) > models.DriverMaxTime/2 ||
			models.StartPoint().DistanceTo(l.Dropoff) > models.DriverMaxTime/2 {
			return []h.Route{}, fmt.Errorf("load %v is too far away for any driver to complete in a single shift", l.Number)
		}
	}
	neighbors := n.getNeighborMap()

	roughMinTotal := n.estimateMinimumTime(neighbors)

	//Max shift length
	minDrivers := int(math.Ceil(roughMinTotal / models.DriverMaxTime))

	resultRoutes, totalCost := n.planRoutesForDrivers(minDrivers, neighbors)

	if len(resultRoutes) != minDrivers {

		// Adding new Driver Recalculate and see if starting with extra drivers yields cost improvement
		newRoutes, newTotalCost := n.planRoutesForDrivers(len(resultRoutes), neighbors)

		if newTotalCost < totalCost {
			resultRoutes = newRoutes
		}
	}

	// Prune out our depot 'loads' from our route
	for i, r := range resultRoutes {
		resultRoutes[i] = slices.Delete(r, 0, 1)
	}

	return resultRoutes, nil
}

func (n NearestNeighborSolver) planRoutesForDrivers(startingDrivers int, neighbors map[int][]models.Load) ([]h.Route, float64) {
	routes := make([]h.Route, startingDrivers)
	for i := range routes {
		routes[i] = append(routes[i], models.Load{})
	}

	remainingLoads := make(map[int]models.Load)
	for _, l := range n.loads {
		if l.Number == 0 {
			// Don't need to track the starting point
			continue
		}
		remainingLoads[l.Number] = l
	}

	var driverFound bool
	var driver int
	var nextLoad models.Load
	for len(remainingLoads) > 0 {
		driverFound = false
		for i, r := range routes {
			for _, l := range neighbors[r[len(r)-1].Number] {
				if _, ok := remainingLoads[l.Number]; ok {
					// Check driver capacity
					if r.CompletionTimeWithLoad(l) <= models.DriverMaxTime {
						driverFound = true
						driver = i
						nextLoad = l

						// As our neighbors are already sorted by nearest, we minimize deadhead by
						// matching here and moving on
						break
					}
				}
			}
		}

		if driverFound {
			routes[driver] = append(routes[driver], nextLoad)
			delete(remainingLoads, nextLoad.Number)
		} else {
			routes = append(routes, h.Route{models.Load{}})
		}
	}

	return routes, h.GetTotalCostOfRoutes(routes)
}

func (n NearestNeighborSolver) getNeighborMap() map[int][]models.Load {
	neighbors := make(map[int][]models.Load)

	for _, l := range n.loads {
		closest := make([]models.Load, 0)
		for _, o := range n.loads {
			if o == l {
				continue
			}
			closest = append(closest, o)
		}
		sort.Slice(closest, func(a, b int) bool {
			nextA := closest[a]
			nextB := closest[b]
			return l.Dropoff.DistanceTo(nextA.Pickup) < l.Dropoff.DistanceTo(nextB.Pickup)
		})
		neighbors[l.Number] = closest
	}

	return neighbors
}

// estimateMinimumTime does a rough esimate of the lower bound for time needed to process all loads
// by traversing the entire set following closest neighbors
func (n NearestNeighborSolver) estimateMinimumTime(neighbors map[int][]models.Load) float64 {
	roughMinTotal := 0.0
	visited := []int{0}
	currentLoadIndex := 0
	var current models.Load
	for len(visited) < len(n.loads) {
		current = n.loads[currentLoadIndex]

		for _, l := range neighbors[currentLoadIndex] {
			if !slices.Contains(visited, l.Number) {
				roughMinTotal += current.Dropoff.DistanceTo(l.Pickup)
				roughMinTotal += l.Cost()
				currentLoadIndex = l.Number
				visited = append(visited, currentLoadIndex)
				break
			}
		}
	}

	// Add our final return to depot
	roughMinTotal += n.loads[currentLoadIndex].Dropoff.DistanceTo(models.StartPoint())

	return roughMinTotal
}
