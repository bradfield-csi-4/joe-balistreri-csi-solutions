TODO


Optional:
-improve the conversion from []byte -> string for map storage



Benchmark Results

** 1/12/23 - In Memory Only **

    goos: darwin
    goarch: amd64
    pkg: github.com/jdbalistreri/bradfield-csi-solutions/leveldbclone/db
    cpu: VirtualApple @ 2.50GHz
    BenchmarkFillSeqTest-10     	 3759535	       425.8 ns/op	     221 B/op	       2 allocs/op
    BenchmarkFillRandom-10      	 4216839	       391.8 ns/op	     198 B/op	       2 allocs/op
    BenchmarkOverwrite-10       	16444137	        76.42 ns/op	       8 B/op	       2 allocs/op
    BenchmarkDeleteSeq-10       	81909630	        14.32 ns/op	       4 B/op	       1 allocs/op
    BenchmarkDeleteRandom-10    	18122713	        68.95 ns/op	       4 B/op	       1 allocs/op
    BenchmarkReadSeq-10         	28138106	        43.92 ns/op	       4 B/op	       1 allocs/op
    BenchmarkReadReverse-10     	26745714	        43.14 ns/op	       4 B/op	       1 allocs/op
    BenchmarkReadRandom-10      	18807057	        67.90 ns/op	       4 B/op	       1 allocs/op
    PASS
    coverage: 22.5% of statements
    ok  	github.com/jdbalistreri/bradfield-csi-solutions/leveldbclone/db	13.013s