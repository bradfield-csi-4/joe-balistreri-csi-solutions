package node

import (
	"testing"

	"github.com/smartystreets/assertions/should"
)

func TestAggregator(t *testing.T) {
	var db = []string{
		"id,rating,genre,length",
		"1,3.0,horror,120",
		"2,4.0,horror,95",
		"3,5.0,anime,85",
		"4,4.0,anime,91",
		"5,2.0,disco,30",
	}

	// var groupByField = "genre"

	t.Run("count works without a groupby", func(t *testing.T) {
		n := NewTestScanNode(db)
		a := NewAggregatorNode(n, AggOptions{
			Aggregators: []Aggregator{NewCountAggregator(nil)},
		})
		res := readAll(a)
		So(t, should.Equal(len(res), 2))
		So(t, should.Equal(len(res[0]), 1))
		So(t, should.Equal(len(res[1]), 1))
		So(t, should.Equal(res[0], "count"))
		So(t, should.Equal(res[1], "5"))
	})
	// t.Run("avg works without a groupby", func(t *testing.T) {
	// 	n := NewTestScanNode(db)
	// 	a := NewAggregatorNode(n, AggOptions{
	// 		Aggregators: []Aggregator{NewAvgAggregator("rating", nil, false)},
	// 	})
	// 	res := readAll(a)
	// 	So(t, should.Equal(len(res), 1))
	// 	So(t, should.Equal(res[0]["avg(rating)"], "3.6"))
	// })
	// t.Run("sum works without a groupby", func(t *testing.T) {
	// 	n := NewTestScanNode(db)
	// 	a := NewAggregatorNode(n, AggOptions{
	// 		Aggregators: []Aggregator{NewAvgAggregator("rating", nil, true)},
	// 	})
	// 	res := readAll(a)
	// 	So(t, should.Equal(len(res), 1))
	// 	So(t, should.Equal(res[0]["sum(rating)"], "18"))
	// })
	// t.Run("count works with a groupby", func(t *testing.T) {
	// 	n := NewTestScanNode(db)
	// 	a := NewAggregatorNode(n, AggOptions{
	// 		Aggregators: []Aggregator{NewCountAggregator(&groupByField)},
	// 	})
	// 	res := readAll(a)
	// 	So(t, should.Equal(len(res), 3))
	// 	sort.SliceStable(res, func(i, j int) bool {
	// 		return res[i]["genre"] < res[j]["genre"]
	// 	})
	// 	So(t, should.Equal(res[0]["genre"], "anime"))
	// 	So(t, should.Equal(res[0]["count"], "2"))
	// 	So(t, should.Equal(res[1]["genre"], "disco"))
	// 	So(t, should.Equal(res[1]["count"], "1"))
	// 	So(t, should.Equal(res[2]["genre"], "horror"))
	// 	So(t, should.Equal(res[2]["count"], "2"))
	// })
	// t.Run("avg works with a groupby", func(t *testing.T) {
	// 	n := NewTestScanNode(db)
	// 	a := NewAggregatorNode(n, AggOptions{
	// 		Aggregators: []Aggregator{NewAvgAggregator("rating", &groupByField, false)},
	// 	})
	// 	res := readAll(a)
	// 	So(t, should.Equal(len(res), 3))
	// 	sort.SliceStable(res, func(i, j int) bool {
	// 		return res[i]["genre"] < res[j]["genre"]
	// 	})
	// 	So(t, should.Equal(res[0]["genre"], "anime"))
	// 	So(t, should.Equal(res[0]["avg(rating)"], "4.5"))
	// 	So(t, should.Equal(res[1]["genre"], "disco"))
	// 	So(t, should.Equal(res[1]["avg(rating)"], "2"))
	// 	So(t, should.Equal(res[2]["genre"], "horror"))
	// 	So(t, should.Equal(res[2]["avg(rating)"], "3.5"))
	// })
	// t.Run("sum works with a groupby", func(t *testing.T) {
	// 	n := NewTestScanNode(db)
	// 	a := NewAggregatorNode(n, AggOptions{
	// 		Aggregators: []Aggregator{NewAvgAggregator("rating", &groupByField, true)},
	// 	})
	// 	res := readAll(a)
	// 	So(t, should.Equal(len(res), 3))
	// 	sort.SliceStable(res, func(i, j int) bool {
	// 		return res[i]["genre"] < res[j]["genre"]
	// 	})
	// 	So(t, should.Equal(res[0]["genre"], "anime"))
	// 	So(t, should.Equal(res[0]["sum(rating)"], "9"))
	// 	So(t, should.Equal(res[1]["genre"], "disco"))
	// 	So(t, should.Equal(res[1]["sum(rating)"], "2"))
	// 	So(t, should.Equal(res[2]["genre"], "horror"))
	// 	So(t, should.Equal(res[2]["sum(rating)"], "7"))
	// })
	// t.Run("multiple aggregators work at the same time with groupby", func(t *testing.T) {
	// 	n := NewTestScanNode(db)
	// 	a := NewAggregatorNode(n, AggOptions{
	// 		Aggregators: []Aggregator{
	// 			NewCountAggregator(&groupByField),
	// 			NewAvgAggregator("rating", &groupByField, false),
	// 			NewAvgAggregator("length", &groupByField, false),
	// 			NewAvgAggregator("rating", &groupByField, true),
	// 		},
	// 	})
	// 	res := readAll(a)
	// 	So(t, should.Equal(len(res), 3))
	// 	sort.SliceStable(res, func(i, j int) bool {
	// 		return res[i]["genre"] < res[j]["genre"]
	// 	})
	// 	So(t, should.Equal(res[0]["genre"], "anime"))
	// 	So(t, should.Equal(res[0]["sum(rating)"], "9"))
	// 	So(t, should.Equal(res[0]["avg(rating)"], "4.5"))
	// 	So(t, should.Equal(res[0]["avg(length)"], "88"))
	// 	So(t, should.Equal(res[0]["count"], "2"))
	// 	So(t, should.Equal(res[1]["genre"], "disco"))
	// 	So(t, should.Equal(res[1]["sum(rating)"], "2"))
	// 	So(t, should.Equal(res[1]["avg(rating)"], "2"))
	// 	So(t, should.Equal(res[1]["avg(length)"], "30"))
	// 	So(t, should.Equal(res[1]["count"], "1"))
	// 	So(t, should.Equal(res[2]["genre"], "horror"))
	// 	So(t, should.Equal(res[2]["sum(rating)"], "7"))
	// 	So(t, should.Equal(res[2]["avg(rating)"], "3.5"))
	// 	So(t, should.Equal(res[2]["avg(length)"], "107.5"))
	// 	So(t, should.Equal(res[2]["count"], "2"))
	// })
	// t.Run("multiple aggregators work at the same time without groupby", func(t *testing.T) {
	// 	n := NewTestScanNode(db)
	// 	a := NewAggregatorNode(n, AggOptions{
	// 		Aggregators: []Aggregator{
	// 			NewCountAggregator(nil),
	// 			NewAvgAggregator("rating", nil, false),
	// 			NewAvgAggregator("length", nil, false),
	// 			NewAvgAggregator("rating", nil, true),
	// 		},
	// 	})
	// 	res := readAll(a)
	// 	So(t, should.Equal(len(res), 1))
	// 	So(t, should.Equal(res[0]["sum(rating)"], "18"))
	// 	So(t, should.Equal(res[0]["avg(rating)"], "3.6"))
	// 	So(t, should.Equal(res[0]["avg(length)"], "84.2"))
	// 	So(t, should.Equal(res[0]["count"], "5"))
	// })
}

func readAll(n ExecutionNode) (result []Row) {
	for curr := n.Next(); curr != nil; curr = n.Next() {
		result = append(result, curr)
	}
	return result
}

func So(t *testing.T, s string) {
	t.Helper()
	if s != "" {
		t.Fatal(s)
	}
}
