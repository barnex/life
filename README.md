# Game of life, high performance
[![GoDoc](https://godoc.org/github.com/barnex/life?status.svg)](https://godoc.org/github.com/barnex/life) 

Accelerated implementation of Conway's game of life. Large boards (>1M cells) process at about 2 billion cells per second (!) on a single core of my i7-3612QM CPU @ 2.10GHz.
That's about one clock cycle per cell.

## Universe representation

We pack 16 cell states in a single 64-bit integer, using 4 bits of storage per cell. Dead cell: `0000`, live cell: `0001`. Using 4 bits per cell gives us just enough headroom to do all computations on entire words at once.

|0000|0001|0000|0000|0000|0001|0000|0000|
|----|----|----|----|----|----|----|----|
|0000|0001|0001|0000|0000|0001|0001|0000|
|0000|0000|0001|0000|0000|0000|0001|0000|

## Counting neighbors



Most operations are done SIMD-style, operating on 16 nibbles at a time. 


## cmd/weblife
Command weblife shows life state in browser. E.g.:

![fig](img.png)

## cmd/evolution
Investigate evolution from random start state. Starting form different fill fractions, most boards evolve towards about 2.5% filled:

![fig](evolution.png)


