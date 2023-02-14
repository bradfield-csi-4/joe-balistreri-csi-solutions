package node

func ParseNode(q QueryExpression, nextNode ExecutionNode) ExecutionNode {
	return &ScanNode{input: []string{
		"movieId,title,genres",
		"1,Toy Story (1995),Adventure|Animation|Children|Comedy|Fantasy",
		"2,Jumanji (1995),Adventure|Children|Fantasy",
		"3,Grumpier Old Men (1995),Comedy|Romance",
		"4,Waiting to Exhale (1995),Comedy|Drama|Romance",
		"5,Father of the Bride Part II (1995),Comedy",
		"6,Heat (1995),Action|Crime|Thriller",
		"7,Sabrina (1995),Comedy|Romance",
		"8,Tom and Huck (1995),Adventure|Children",
		"9,Sudden Death (1995),Action",
	}}
}
