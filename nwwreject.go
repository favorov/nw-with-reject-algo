// Copyright 2015 Andrew E. Bruno, 2019 Alexnder (Sasha) Favorov.
// All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package nwwreject

import "log"

var (
	Up   byte = 1
	Left byte = 2
	NW   byte = 3
	Here byte = 4
	Stop byte = 5
)

func idx(i, j, bLen int) int {
	return (i * bLen) + j
}

func Align(a, b string, mismatch, gap, threshold int) (alignA, alignB string, dist int, ok bool) {

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
		dist := gap * i
		f[idx(i, 0, bLen)] = dist
		if dist <= threshold {
			pointer[idx(i, 0, bLen)] = Up
		} else {
			pointer[idx(i, 0, bLen)] = Stop
			break
		}
	}
	for j := 1; j < bLen; j++ {
		dist := gap * j
		f[idx(0, j, bLen)] = dist
		if dist <= threshold {
			pointer[idx(0, j, bLen)] = Left
		} else {
			pointer[idx(0, j, bLen)] = Stop
			we_broke_at:=j
			//where the Stop appeared in a finshed line
			break
		}
	}

	pointer[0] = Here

	start_next_at:=1
	//where from to start next line of alignment matrix

	//we_broke_at is a critical value. 
	//if we are on the right of (we are on the next line),
	//we do not check any directions other than left
	//if we do not open good cells on this line and we are under the we_broke_at
	//we do not test we_broke_at+1
	//we break completely
	for i := 1; i < aLen; i++ {
		for j := start_next_at; j < bLen; j++ {
			nonstop_already_found := false 
			//we count them; if we get break and there are 0 there, we return fail

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

			if min > threshold {
				pointer[idx(i, j, bLen)] = Stop
				f[idx(i, j, bLen)] = min
				if nonstop_already_found {
					//do not go right, it was ok and then it is bad again,
					//we left the good area
					we_broke_at=j
					break 
				} 
				else {
					continue //looking
				}
			}
			else { //good area!!
				nonstop_already_found=true
				start_next_at=j //makes no sense to start under stop
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

	for p := pointer[idx(i, j, bLen)]; p != Here; p = pointer[idx(i, j, bLen)] {
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
		} else if p == Stop {
			log.Fatalln("Stop is found on the alignment way. I am lost.")
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
