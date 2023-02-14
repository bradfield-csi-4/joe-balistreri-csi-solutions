package main

import (
	"fmt"

	"github.com/jdbalistreri/bradfield-csi-solutions/05-database-systems/node"
)

func main() {
	// query1 := []node.QueryExpression{
	// 	{Name: "PROJECTION", Args: []string{"title"}},
	// 	{Name: "SELECTION", Args: []string{"movieId", "EQUALS", "5000"}},
	// 	{Name: "SCAN", Args: []string{"movies"}},
	// }

	// query2 := []node.QueryExpression{
	// 	{Name: "LIMIT", Args: []string{"100"}},
	// 	{Name: "SCAN", Args: []string{"movies"}},
	// }

	// query3 := []node.QueryExpression{
	// 	{Name: "COUNT", Args: []string{"GROUP BY", "movieId"}},
	// 	// {Name: "DISTINCT", Args: []string{"movieId"}},
	// 	{Name: "SCAN", Args: []string{"ratings"}},
	// }

	query4 := []node.QueryExpression{
		{Name: "LIMIT", Args: []string{"10"}},
		{Name: "SORT", Args: []string{"avg(rating)", "DESC"}},
		{Name: "AVG", Args: []string{"rating", "GROUP BY", "movieId"}},
		{Name: "SCAN", Args: []string{"ratings"}},
	}

	for i, query := range [][]node.QueryExpression{
		// query1,
		// query2,
		// query3,
		query4,
	} {
		fmt.Printf("Query %d\n", i+1)
		readAll(parseQuery(query))
		fmt.Printf("\n\n\n")
	}

}

func readAll(rootNode node.ExecutionNode) {
	for {
		next := rootNode.Next()
		if next == nil {
			break
		}
		fmt.Println(next)
	}
}

func parseQuery(query []node.QueryExpression) (rootNode node.ExecutionNode) {
	for i := len(query) - 1; i >= 0; i-- {
		rootNode = node.ParseNode(query[i], rootNode)
	}
	return rootNode
}
