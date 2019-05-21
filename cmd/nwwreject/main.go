// Copyright 2015 Andrew E. Bruno, 2019 Alexnder (Sasha) Favorov. 
// All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"github.com/favorov/nwwreject"
)

var seq1 = flag.String("seq1", "", "first sequence")
var seq2 = flag.String("seq2", "", "second sequence")
var mismatch = flag.Int("mismatch", 1, "mismatch score")
var gap = flag.Int("gap", 1, "gap penalty")
var threshold = flag.Int("threshold", math.MaxInt32, "threshold to reject")

func main() {
	flag.Parse()
	if *seq1 == "" || *seq2 == "" {
		log.Fatal("Please provide 2 sequences to align. See nwwreject --help")
	}

	aln1, aln2, dist ,ok := nwwreject.Align(*seq1, *seq2, *mismatch, *gap,*threshold)
	if ok {
		fmt.Printf("%s\n%s\nDistance: %d\n", aln1, aln2, dist)
	} else {
		fmt.Printf("Sequences differ too much.\n")
	}
}
