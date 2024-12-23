package day23

import (
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type node struct {
	name       string
	neighbours map[string]bool
}

type graph map[string]*node

func loadInput(file string) (graph, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	g := make(graph)
	for _, l := range strings.Split(string(content), "\n") {
		parts := strings.Split(strings.TrimSpace(string(l)), "-")
		if g[parts[0]] == nil {
			g[parts[0]] = &node{parts[0], map[string]bool{parts[1]: true}}
		} else {
			g[parts[0]].neighbours[parts[1]] = true
		}
		if g[parts[1]] == nil {
			g[parts[1]] = &node{parts[1], map[string]bool{parts[0]: true}}
		} else {
			g[parts[1]].neighbours[parts[0]] = true
		}
	}
	return g, nil
}

func findClusters(g graph) (map[string]bool, int) {
	clusters := make(map[string]bool)
	hc := 0
	for _, nd := range g {
		for nb := range nd.neighbours {
			for nb2 := range g[nb].neighbours {
				if nb2 == nd.name {
					continue
				}
				if g[nb2].neighbours[nd.name] {
					c := []string{nd.name, nb, nb2}
					sort.Strings(c)
					cls := strings.Join(c, ",")
					if !clusters[cls] {
						clusters[cls] = true
						if nd.name[0] == 't' || nb[0] == 't' || nb2[0] == 't' {
							hc++
						}
					}
				}
			}
		}
	}
	return clusters, hc
}

func TestDay23Part1(t *testing.T) {
	g, err := loadInput("input_test.txt")
	require.NoError(t, err)
	cls, hc := findClusters(g)
	require.Equal(t, 12, len(cls))
	require.Equal(t, 7, hc)

	g, err = loadInput("input.txt")
	require.NoError(t, err)
	cls, hc = findClusters(g)
	assert.Equal(t, 11011, len(cls))
	assert.Equal(t, 1043, hc)
}

func canAdd(n string, cluster map[string]bool, g graph) bool {
	nd := g[n]
	for n := range cluster {
		if nd.neighbours[n] == false {
			return false
		}
	}
	return true
}

func extendCluster(nd string, cluster map[string]bool, g graph) map[string]bool {
	for nb := range g[nd].neighbours {
		if canAdd(nb, cluster, g) {
			cluster[nb] = true
			extendCluster(nb, cluster, g)
		}
	}
	return cluster
}

func findLargestCluster(g graph) string {
	nodes := make([]string, 0, len(g))
	for n := range g {
		nodes = append(nodes, n)
	}
	sort.Strings(nodes)
	processed := make(map[string]bool)
	cluster := make(map[string]bool)
	for _, n := range nodes {
		if processed[n] {
			continue
		}
		nc := extendCluster(n, map[string]bool{n: true}, g)
		if len(nc) > len(cluster) {
			cluster = nc
			for n := range cluster {
				processed[n] = true
			}
		}
	}

	players := make([]string, 0, len(cluster))
	for n := range cluster {
		players = append(players, n)
	}
	sort.Strings(players)
	return strings.Join(players, ",")
}

func TestDay23Part2(t *testing.T) {
	g, err := loadInput("input_test.txt")
	require.NoError(t, err)
	cls := findLargestCluster(g)
	require.Equal(t, "co,de,ka,ta", cls)

	g, err = loadInput("input.txt")
	require.NoError(t, err)
	cls = findLargestCluster(g)
	assert.Equal(t, "ai,bk,dc,dx,fo,gx,hk,kd,os,uz,xn,yk,zs", cls)
}
