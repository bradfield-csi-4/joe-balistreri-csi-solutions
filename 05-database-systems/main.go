package main

import (
	"fmt"

	"github.com/jdbalistreri/bradfield-csi-solutions/05-database-systems/node"
)

func main() {
	// find movie with id 5000
	query1 := []node.QueryExpression{
		{Name: "PROJECTION", Args: []string{"title"}},
		{Name: "SELECTION", Args: []string{"movieId", "EQUALS", "5000"}},
		{Name: "SCAN", Args: []string{"movies"}},
	}

	// list the first 100 movies
	query2 := []node.QueryExpression{
		{Name: "PROJECTION", Args: []string{"movieId", "genres", "title"}},
		{Name: "LIMIT", Args: []string{"100"}},
		{Name: "SCAN", Args: []string{"movies"}},
	}

	// count how many distinct movies have been rated
	query3 := []node.QueryExpression{
		// {Name: "AGG", Args: []string{"COUNT"}},
		{Name: "DISTINCT", Args: []string{"movieId"}},
		{Name: "SCAN", Args: []string{"ratings"}},
	}

	// // list the top 10 highest rated movies
	// query4 := []node.QueryExpression{
	// 	{Name: "LIMIT", Args: []string{"10"}},
	// 	{Name: "SORT", Args: []string{"avg(rating)", "DESC"}},
	// 	{Name: "AGG", Args: []string{"COUNT", "AVG", "rating", "GROUP BY", "movieId"}},
	// 	{Name: "SCAN", Args: []string{"ratings"}},
	// }

	// // list the ratings of the top 10 most rated movies
	// query5 := []node.QueryExpression{
	// 	{Name: "LIMIT", Args: []string{"10"}},
	// 	{Name: "SORT", Args: []string{"count", "DESC"}},
	// 	{Name: "AGG", Args: []string{"COUNT", "AVG", "rating", "GROUP BY", "movieId"}},
	// 	{Name: "SCAN", Args: []string{"ratings"}},
	// }

	for i, query := range [][]node.QueryExpression{
		query1,
		query2,
		query3,
		// query4,
		// query5,
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
