// Copyright 2015 Andrew E. Bruno, 2019 Alexnder (Sasha) Favorov. 
// All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package nwwreject

import (
	"testing"
)

func TestAlign(t *testing.T) {
	seqs := [][]string{
		[]string{"CGAGAGA", "GAGAGA", "CGAGAGA", "-GAGAGA"}}

	for _, a := range seqs {
		aln1, aln2, dist, ok := Align(a[0], a[1], 1, 1, 1)
		if aln1 != a[2] || aln2 != a[3] || dist !=1 || !ok{
			t.Errorf("Align(%s, %s)\n***GOT***\n%s\n%s\n%d\n%t\n***WANT***\n%s\n%s\n%d\n%t", a[0], a[1], 1, true, aln1, aln2, a[2], a[3], dist, ok)
		}
	}

}

func BenchmarkAlign(b *testing.B) {
	seq1 := "GGAATTAATCCAGGTAATGGACCCCAAGAT"
	seq2 := "GCCAGGATTCCCAGATATGGCCAAGGTTCC"

	for i := 0; i < b.N; i++ {
		Align(seq1, seq2, 1, 1)
	}
}
