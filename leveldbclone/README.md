TODO


Optional:
-improve the conversion from []byte -> string for map storage



Benchmark Results

** In Memory Only **

    BenchmarkFillSeqTest-10               	 4407808	       306.6 ns/op	     190 B/op	       2 allocs/op
    BenchmarkFillRandom-10                	 5334057	       251.5 ns/op	     159 B/op	       2 allocs/op
    BenchmarkOverwrite-10                 	 5457484	       297.2 ns/op	       8 B/op	       2 allocs/op
    BenchmarkDeleteSeq-10                 	 8044210	       189.9 ns/op	       4 B/op	       1 allocs/op
    BenchmarkDeleteRandom-10              	 4782680	       271.0 ns/op	       4 B/op	       1 allocs/op
    BenchmarkReadSeq-10                   	11905912	       113.2 ns/op	       4 B/op	       1 allocs/op
    BenchmarkReadReverse-10               	13104931	       117.9 ns/op	       4 B/op	       1 allocs/op
    BenchmarkReadRandom-10                	 5662317	       278.2 ns/op	       4 B/op	       1 allocs/op
    BenchmarkRangeScanNoIteration-10      	   10000	   6333095 ns/op	 2768433 B/op	   10054 allocs/op
    BenchmarkRangeScanWithIteration-10    	   10000	   6390621 ns/op	 2764970 B/op	   10054 allocs/op
    BenchmarkRangeAndPut-10               	   10000	   2929551 ns/op	 1321614 B/op	    5041 allocs/op

** Linked List **

    BenchmarkFillSeqTest-10               	   59344	    113802 ns/op	      68 B/op	       2 allocs/op
    BenchmarkFillRandom-10                	   62784	    120517 ns/op	      68 B/op	       2 allocs/op
    BenchmarkOverwrite-10                 	   61996	    119040 ns/op	       4 B/op	       1 allocs/op
    BenchmarkDeleteSeq-10                 	   61534	    118228 ns/op	       4 B/op	       1 allocs/op
    BenchmarkDeleteRandom-10              	   62212	    120112 ns/op	       4 B/op	       1 allocs/op
    BenchmarkReadSeq-10                   	   61818	    119047 ns/op	       4 B/op	       1 allocs/op
    BenchmarkReadReverse-10               	   62358	    119808 ns/op	       4 B/op	       1 allocs/op
    BenchmarkReadRandom-10                	   62943	    120758 ns/op	       4 B/op	       1 allocs/op
    BenchmarkRangeScanNoIteration-10      	   97008	    113425 ns/op	      40 B/op	       3 allocs/op
    BenchmarkRangeScanWithIteration-10    	   31520	    117309 ns/op	      40 B/op	       3 allocs/op
    BenchmarkRangeAndPut-10               	   41062	    117173 ns/op	     108 B/op	       5 allocs/op

** Skip List **

    BenchmarkFillSeqTest-10               	 3230851	       403.7 ns/op	     484 B/op	       4 allocs/op
    BenchmarkFillRandom-10                	 3333127	       413.0 ns/op	     484 B/op	       4 allocs/op
    BenchmarkOverwrite-10                 	 1000000	      1848 ns/op	     196 B/op	       2 allocs/op
    BenchmarkDeleteSeq-10                 	29680746	        41.50 ns/op	       4 B/op	       1 allocs/op
    BenchmarkDeleteRandom-10              	24913774	        50.20 ns/op	       4 B/op	       1 allocs/op
    BenchmarkReadSeq-10                   	32148058	        49.15 ns/op	       4 B/op	       1 allocs/op
    BenchmarkReadReverse-10               	32404333	        40.00 ns/op	       4 B/op	       1 allocs/op
    BenchmarkReadRandom-10                	27844515	        42.82 ns/op	       4 B/op	       1 allocs/op
    BenchmarkRangeScanNoIteration-10      	 1000000	      1570 ns/op	      40 B/op	       3 allocs/op
    BenchmarkRangeScanWithIteration-10    	   34322	    144363 ns/op	      40 B/op	       3 allocs/op
    BenchmarkRangeAndPut-10               	 1000000	      1747 ns/op	     524 B/op	       7 allocs/op

** Skip List with Write Ahead Log **

    BenchmarkFillSeqTest-10               	     369	   3640800 ns/op	     663 B/op	      10 allocs/op
    BenchmarkFillRandom-10                	     319	   3714427 ns/op	     663 B/op	      10 allocs/op
    BenchmarkOverwrite-10                 	     361	   4040867 ns/op	     375 B/op	       8 allocs/op
    BenchmarkDeleteSeq-10                 	     334	   3566051 ns/op	     151 B/op	       6 allocs/op
    BenchmarkDeleteRandom-10              	     304	   3444294 ns/op	     151 B/op	       6 allocs/op
    BenchmarkReadSeq-10                   	46219464	        26.43 ns/op	       4 B/op	       1 allocs/op
    BenchmarkReadReverse-10               	48169924	        27.97 ns/op	       4 B/op	       1 allocs/op
    BenchmarkReadRandom-10                	35209245	        36.89 ns/op	       4 B/op	       1 allocs/op
    BenchmarkRangeScanNoIteration-10      	 8673613	       158.2 ns/op	      40 B/op	       3 allocs/op
    BenchmarkRangeScanWithIteration-10    	 9561294	       115.7 ns/op	      40 B/op	       3 allocs/op
    BenchmarkRangeAndPut-10               	       1	4568132584 ns/op	  703696 B/op	   13000 allocs/op