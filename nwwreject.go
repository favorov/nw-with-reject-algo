// Copyright 2015 Andrew E. Bruno, 2019 Alexnder (Sasha) Favorov. 
// All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package nwwrej

var (
	Up   byte = 1
	Left byte = 2
	NW   byte = 3
	None byte = 4
)

func idx(i, j, bLen int) int {
	return (i * bLen) + j
}

func Align(a, b string, mismatch, gap, threshold int) (alignA, alignB string, dist int, ok bool) {
	
	if(threshold>0){threshold++} //to do something

	aLen := len(a) + 1
	bLen := len(b) + 1

	maxLen := aLen
	if maxLen < bLen {
		maxLen = bLen
	}

	aBytes := make([]byte, 0, maxLen)
	bBytes := make([]byte, 0, maxLen)

	f := make([]int, aLen*bLen)
	pointer := make([]byte, aLen*bLen)

	for i := 1; i < aLen; i++ {
		f[idx(i, 0, bLen)] = gap * i
		pointer[idx(i, 0, bLen)] = Up
	}
	for j := 1; j < bLen; j++ {
		f[idx(0, j, bLen)] = gap * j
		pointer[idx(0, j, bLen)] = Left
	}

	pointer[0] = None

	for i := 1; i < aLen; i++ {
		for j := 1; j < bLen; j++ {
			matchMismatch := mismatch
			if a[i-1] == b[j-1] {
				matchMismatch = 0 
			}

			min := f[idx(i-1, j-1, bLen)] + matchMismatch
			hgap := f[idx(i-1, j, bLen)] + gap
			vgap := f[idx(i, j-1, bLen)] + gap

			if hgap < min {
				min = hgap
			}
			if vgap < min {
				min = vgap
			}

			p := NW
			if min == hgap {
				p = Up
			} else if min == vgap {
				p = Left
			}

			pointer[idx(i, j, bLen)] = p
			f[idx(i, j, bLen)] = min
		}
	}

	i := aLen - 1
	j := bLen - 1

	dist = f[idx(i, j, bLen)]

	for p := pointer[idx(i, j, bLen)]; p != None; p = pointer[idx(i, j, bLen)] {
		if p == NW {
			aBytes = append(aBytes, a[i-1])
			bBytes = append(bBytes, b[j-1])
			i--
			j--
		} else if p == Up {
			aBytes = append(aBytes, a[i-1])
			bBytes = append(bBytes, '-')
			i--
		} else if p == Left {
			aBytes = append(aBytes, '-')
			bBytes = append(bBytes, b[j-1])
			j--
		}
	}

	reverse(aBytes)
	reverse(bBytes)

	return string(aBytes), string(bBytes), dist, true 
}

func reverse(a []byte) {
	for i := 0; i < len(a)/2; i++ {
		j := len(a) - 1 - i
		a[i], a[j] = a[j], a[i]
	}
}
