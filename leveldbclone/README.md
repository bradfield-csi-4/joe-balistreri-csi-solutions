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


** 1/12/23 - In Memory Only ** (with ResetTimer fix)
    goos: darwin
    goarch: amd64
    pkg: github.com/jdbalistreri/bradfield-csi-solutions/leveldbclone/db
    cpu: VirtualApple @ 2.50GHz
    BenchmarkFillSeqTest-10     	 3123774	       332.7 ns/op	     141 B/op	       2 allocs/op
    BenchmarkFillRandom-10      	 3494415	       396.2 ns/op	     237 B/op	       2 allocs/op
    BenchmarkOverwrite-10       	16265038	        75.26 ns/op	       8 B/op	       2 allocs/op
    BenchmarkDeleteSeq-10       	81230418	        14.42 ns/op	       4 B/op	       1 allocs/op
    BenchmarkDeleteRandom-10    	18007585	        71.74 ns/op	       4 B/op	       1 allocs/op
    BenchmarkReadSeq-10         	27856068	        45.19 ns/op	       4 B/op	       1 allocs/op
    BenchmarkReadReverse-10     	27066804	        44.45 ns/op	       4 B/op	       1 allocs/op
    BenchmarkReadRandom-10      	18062416	        70.52 ns/op	       4 B/op	       1 allocs/op
    PASS
    coverage: 22.5% of statements
ok  	github.com/jdbalistreri/bradfield-csi-solutions/leveldbclone/db	12.437s