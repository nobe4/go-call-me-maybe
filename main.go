package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type CallMap map[string][]string
type Node struct {
	Name  string
	Depth int
}

func main() {
	// Generate the graph file by running:
	// callgragh github.com/cli/cli/cmd/gh
	cm := parseFile("./graph")

	// TODO: make this an argument
	mustMatch := regexp.MustCompile(".*github.com/cli/cli.*")

	// TODO: make this an argument
	maxDepth := 2

	// TODO: make this an argument
	nodes := match(cm, ".*ghinstance.IsLocal$")

	seen := make(map[string]bool)

	for {
		log.Printf("nodes: %v\n", nodes)
		log.Printf("count: %d\n", len(nodes))

		if len(nodes) == 0 {
			break
		}

		node := nodes[0]
		log.Printf("node: %s\n", node.Name)

		if _, ok := seen[node.Name]; !ok {
			seen[node.Name] = true

			for _, child := range cm[node.Name] {
				log.Printf("child: %s\n", child)
				if node.Depth >= maxDepth {
					continue
				}

				if !mustMatch.MatchString(child) {
					continue
				}

				fmt.Printf("%s --> %s\n", node.Name, child)
				nodes = append(nodes, Node{Name: child, Depth: node.Depth + 1})
			}
		}

		if len(nodes) == 1 {
			break
		}

		nodes = nodes[1:]
	}
}

func parseFile(path string) CallMap {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	cm := CallMap{}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "\t")
		caller, callee := parts[0], parts[2]

		if _, ok := cm[callee]; !ok {
			cm[callee] = []string{}
		}
		cm[callee] = append(cm[callee], caller)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return cm
}

func match(cm CallMap, target string) []Node {
	keys := []Node{}

	re := regexp.MustCompile(target)
	for key := range cm {
		if re.MatchString(key) {
			keys = append(keys, Node{Name: key, Depth: 0})
		}
	}

	return keys
}
