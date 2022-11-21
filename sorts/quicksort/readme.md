# Concurrent Quicksort

This is an implementation of quicksort using multiple goroutines.

Concurrency does not necessarily speed up any cpu-bound algorithm except that it
distributes work effectively over available cores.
Nor is the number of cores consistent on similar hosts, for example ResourceQuota
constraints on kubernetes deployments can compress core resources such that NumCPU()
is not the optimal choice of procs for a single host. See 100 Mistakes in Golang for an explanation.

This implementation simply utilizes two concurrency strategies:
1) distribute work over available cores (usually a modest number)
2) call insertionSort on slice partitions less than size n, where n is chosen
   to utilize cpu cache as much as possible.

Only empirical analysis can establish one concurrent implementation over another.
The benchmark tests perform some very modest analyses, but not much.

TODO:
0) Genericize the sorting interface; []int was just for development
1) The stdlib implementation outperforms this one for certain inputs, seemingly when it
   can memoize them. It would be nice to know what it is doing internally.
