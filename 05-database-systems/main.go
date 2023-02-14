package main

import (
	"fmt"

	"github.com/jdbalistreri/bradfield-csi-solutions/05-database-systems/node"
)

func main() {
	query := []node.QueryExpression{
		{Name: "PROJECTION", Args: []string{"name"}},
		{Name: "SELECTION", Args: []string{"id", "EQUALS", "5000"}},
		{Name: "SCAN", Args: []string{"movies"}},
	}

	var rootNode node.ExecutionNode
	for i := len(query) - 1; i > 0; i-- {
		rootNode = node.ParseNode(query[i], rootNode)
	}

	for {
		next := rootNode.Next()
		if next == nil {
			break
		}
		fmt.Println(next)
	}
}
