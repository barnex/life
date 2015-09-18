# life
The classic Conway's game of life.

## weblife
Shows life state in browser. E.g.:

![fig](img.png)

## evolution
Investigate evolution from random start state. Starting form different fill fractions, most boards evolve towards about 2.5% filled:

![fig](evolution.png)


## performance
Dumb algorithm but reasonably optimized with Go's pprof tool. Large boards (>1M cells) process at >1 billion cells per second on my core i7-3612QM CPU @ 2.10GHz.
