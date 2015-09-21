# Game of life, SIMD performance
[![GoDoc](https://godoc.org/github.com/barnex/life?status.svg)](https://godoc.org/github.com/barnex/life) [![Build Status](https://travis-ci.org/barnex/life.svg?branch=master)](https://travis-ci.org/barnex/life)

Fast implementation of Conway's game of life.

This implementation packs 16 cell states in a single 64-bit integer, using 4 bits of storage per cell. Most operations are done SIMD-style, operating on 16 nibbles at a time. 

Large boards (>1M cells) process at about 2 billion cells per second on a single core of my i7-3612QM CPU @ 2.10GHz.
That's about one clock cycle per cell!

## weblife
Shows life state in browser. E.g.:

![fig](img.png)

## evolution
Investigate evolution from random start state. Starting form different fill fractions, most boards evolve towards about 2.5% filled:

![fig](evolution.png)


