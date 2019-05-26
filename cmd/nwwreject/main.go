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
	//"github.com/favorov/nwwreject" 
	"./nwwreject" //debug , supposes symlink project/cmd/nwwreject/nwwreject/nwwreject.go ---> project/nwwreject.go 
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

	fmt.Println("nwwreject version",nwwreject.Version)
	fmt.Println("call Align",nwwreject.Version)
	aln1, aln2, dist ,ok := nwwreject.Align(*seq1, *seq2, *mismatch, *gap,*threshold)
	if ok {
		fmt.Printf("%s\n%s\nDistance: %d\n", aln1, aln2, dist)
	} else {
		fmt.Printf("Sequences differ too much.\n")
	}
	fmt.Println("call Distance",nwwreject.Version)
	distince ,diok := nwwreject.Distance(*seq1, *seq2, *mismatch, *gap,*threshold)
	if diok {
		fmt.Printf("Distance: %d\n", distince)
	} else {
		fmt.Printf("Sequences differ too much.\n")
	}

	seq3:=(*seq1)[:1] + "A" + (*seq1)[2:]
	fmt.Println("Mutate string 1 a bit (A to pos 2)")
	fmt.Println("call Align",nwwreject.Version)
	aln1, aln2, dist ,ok = nwwreject.Align(seq3, *seq2, *mismatch, *gap,*threshold)
	if ok {
		fmt.Printf("%s\n%s\nDistance: %d\n", aln1, aln2, dist)
	} else {
		fmt.Printf("Sequences differ too much.\n")
	}
	fmt.Println("call Distance",nwwreject.Version)
	distince ,diok = nwwreject.Distance(seq3, *seq2, *mismatch, *gap,*threshold)
	if diok {
		fmt.Printf("Distance: %d\n", distince)
	} else {
		fmt.Printf("Sequences differ too much.\n")
	}
}
