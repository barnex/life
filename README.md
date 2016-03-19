# Game of life, high performance
[![GoDoc](https://godoc.org/github.com/barnex/life?status.svg)](https://godoc.org/github.com/barnex/life) 

Accelerated implementation of Conway's game of life. Large boards (>1M cells) process at about 2 billion cells per second (!) on a single core of my i7-3612QM CPU @ 2.10GHz.
That's about one clock cycle per cell.

## Universe representation

We pack 16 cell states in a single 64-bit integer, using 4 bits of storage per cell. Dead cell: `0000`, live cell: `0001`. Using 4 bits per cell gives us just enough headroom to do all computations on entire words at once. E.g.: 3x16 cells are stored in these 3 64-bit words:

|    |    |    |    |    |    |    |    |
|----|----|----|----|----|----|----|----|
|0000|0001|0000|0000|0000|0001|0000|0001|
|0000|0001|0001|0000|0000|0001|0001|0001|
|0000|0000|0001|0000|0000|0000|0001|0001|

## Counting neighbors

By using 4 bits per cell, we can count neighbors without unpacking. We calculate the neighbors for one entire row of cells at a time. First we, make a partial sum adding up the row, the row above and the row below. E.g.: for the center row:

|    |    |    |    |    |    |    |    |
|----|----|----|----|----|----|----|----|
|0000|0001|0000|0000|0000|0001|0000|0001|
|0000|0001|0001|0000|0000|0001|0001|0001|
|0000|0000|0001|0000|0000|0000|0001|0001|

`+ =`

|0000|0010|0010|0000|0000|0010|0010|0011|
|----|----|----|----|----|----|----|----|
|    |    |    |    |    |    |    |    |

This requires only 2 64-bit additions.

Second, for each nibble in this partial sum, we need to add its left and right neighbor. We do this by adding the row with itself shifted by 4 bits to the left and 4 bits to the right.

|    |    |    |    |    |    |    |    |    |    |
|----|----|----|----|----|----|----|----|----|----|
|0000|0010|0010|0000|0000|0010|0010|0011|<<<<|    |
|    |0000|0010|0010|0000|0000|0010|0010|0011|    |
|    |>>>>|0000|0010|0010|0000|0000|0010|0010|0011|
| += |----|----|----|----|----|----|----|----|----|
|    |0000|0010|0100|0100|0010|0010|0100|0111|    |

Of course, when shifting, we need to fill in the 4 bits on the left or right side with the corresponding bits of the neighboring word. We marked those with `<<<<` for brevity.

We now have the number of neighbors for 16 cells, using only a few additions and shifts. A more straightforward implementation would have taken about 100 loads, 100 additions and 100 stores. Note that we include the central cell itself in the number of neighbors.

## Finding the next state

The liveness of a cell and its number of neighbors determine the cell's next state. 

|cell|neighbors|lookup-code|next state|
|----|---------|-----------|----------|
|0000|     0000|       0000|      0000|
|0000|     0001|       0001|      0000|
|0000|     0010|       0010|      0000|
|0000|     0011|       0011|      0001|
| ...|      ...|        ...|       ...|
|0001|     0000|       1000|      0000|
|0001|     0001|       1001|      0000|
|0001|     0010|       1010|      0000|
|0001|     0011|       1011|      0001|



## cmd/weblife
Command weblife shows life state in browser. E.g.:

![fig](img.png)

## cmd/evolution
Investigate evolution from random start state. Starting form different fill fractions, most boards evolve towards about 2.5% filled:

![fig](evolution.png)


