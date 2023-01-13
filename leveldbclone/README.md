TODO


Optional:
-improve the conversion from []byte -> string for map storage



Benchmark Results

** 1/13/23 - In Memory Only - Starting Point **
    goos: darwin
    goarch: amd64
    pkg: github.com/jdbalistreri/bradfield-csi-solutions/leveldbclone/db
    cpu: VirtualApple @ 2.50GHz
    BenchmarkFillSeqTest-10               	 2947156	       363.9 ns/op	     147 B/op	       2 allocs/op
    BenchmarkFillRandom-10                	 3016215	       350.6 ns/op	     145 B/op	       2 allocs/op
    BenchmarkOverwrite-10                 	 4755320	       349.2 ns/op	       8 B/op	       2 allocs/op
    BenchmarkDeleteSeq-10                 	 6849459	       305.5 ns/op	       4 B/op	       1 allocs/op
    BenchmarkDeleteRandom-10              	 5082489	       290.4 ns/op	       4 B/op	       1 allocs/op
    BenchmarkReadSeq-10                   	 7175680	       266.1 ns/op	       4 B/op	       1 allocs/op
    BenchmarkReadReverse-10               	 6662030	       209.7 ns/op	       4 B/op	       1 allocs/op
    BenchmarkReadRandom-10                	 4509734	       307.1 ns/op	       4 B/op	       1 allocs/op
    BenchmarkRangeScanNoIteration-10      	   10000	  10242039 ns/op	 2765429 B/op	   10054 allocs/op
    BenchmarkRangeScanWithIteration-10    	   10000	  10361215 ns/op	 2759549 B/op	   10054 allocs/op
    BenchmarkRangeAndPut-10               	   10000	   4566075 ns/op	 1331528 B/op	    5042 allocs/op