// Copyright 2015 Andrew E. Bruno, 2019 Alexnder (Sasha) Favorov.
// All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package nwwreject
import "log"
import "fmt"

var Version string = "0.0.4"

var (
	Up   byte = 1
	Left byte = 2
	NW   byte = 3
	Here byte = 4
	Stop byte = 5
)

var f []int //not to allocate each call, it is the distance matrix
var pointer []byte //not to allocate each call, it is the distance matrix

func Init_distance_matrix(maxlen int){
	f=make([]int,maxlen*maxlen)
}

func Init_pointer_matrix(maxlen int){
	pointer=make([]byte,maxlen*maxlen)
} //these two are not obligatory, but the are good style


func idx(i, j, bLen int) int {
	return (i * bLen) + j
}

func reverse(a []byte) {
	for i := 0; i < len(a)/2; i++ {
		j := len(a) - 1 - i
		a[i], a[j] = a[j], a[i]
	}
}

func fmtmati(mat []int,aLen int, bLen int) {
	for i := 0; i < aLen; i++ {
		for j := 0; j < bLen; j++ {
			fmt.Print(mat[idx(i, j, bLen)],"  ")
		}
		fmt.Println()
	}
	fmt.Println()
}

func fmtmatb(mat []byte,aLen int, bLen int) {
	for i := 0; i < aLen; i++ {
		for j := 0; j < bLen; j++ {
			fmt.Print(mat[idx(i, j, bLen)],"  ")
		}
		fmt.Println()
	}
	fmt.Println()
}


func Align(a, b string, mismatch, gap, threshold int) (alignA, alignB string, dist int, ok bool) {
	//this is the most full version, returns alignment if success

	aLen := len(a) + 1
	bLen := len(b) + 1

	maxLen := aLen
	if maxLen < bLen {
		maxLen = bLen
	}

	aBytes := make([]byte, 0, maxLen)
	bBytes := make([]byte, 0, maxLen)

	we_broke_at:=bLen
	give_up:=false
	//we_broke_at is a critical value. 
	//if we are on the right of we_broke_at (we are on the next line),
	//we do not check any directions other than left
	//if we do not open good cells on this line and we are under the we_broke_at
	//we do not test we_broke_at+1
	//we break completely by setting give_up flag to true


	start_next_at:=1
	//where from to start next line of alignment matrix
	first_good_prev:=0
	//usually, it is start_next_at-1
	first_stopped_head:=aLen+1
	//which line is first to get >threshold distance and 
	//stop as pointer at the head (0) position
	//init with never

	if len(f) < aLen*bLen {
		f = make([]int, aLen*bLen)
	}

	if len(pointer) < aLen*bLen {
		pointer = make([]byte, aLen*bLen)
	}

	pointer[0] = Here

	for i := 1; i < aLen; i++ { //vertical
		dist := gap * i
		f[idx(i, 0, bLen)] = dist
		if dist <= threshold {
			pointer[idx(i, 0, bLen)] = Up
		} else {
			pointer[idx(i, 0, bLen)] = Stop
			first_stopped_head=i
			break
		}
	}

	for j := 1; j < bLen; j++ { //horizontal
		dist := gap * j
		f[idx(0, j, bLen)] = dist
		if dist <= threshold {
			pointer[idx(0, j, bLen)] = Left
		} else {
			pointer[idx(0, j, bLen)] = Stop
			we_broke_at=j
			//where the Stop appeared in a finshed line
			break
		}
	}

	for i := 1; i < aLen; i++ {
		if give_up {break;} //chau
		nonstop_already_found := false
		for j := start_next_at; j < bLen; j++ {
			var min int
			if (j<=we_broke_at) && (j>first_good_prev) { //we can test all 3 direction	
				matchMismatch := mismatch
				if a[i-1] == b[j-1] {
					matchMismatch = 0
				}

				min = f[idx(i-1, j-1, bLen)] + matchMismatch
				vgap := f[idx(i, j-1, bLen)] + gap
				hgap := f[idx(i-1, j, bLen)] + gap

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
			} else if (j>first_good_prev) {//ok, we can test only left
				pointer[idx(i, j, bLen)] = Left
				min=f[idx(i, j-1, bLen)] + gap
			} else if (j<=we_broke_at) { //ok, maybe, up will help, left and NW are forbidden
				pointer[idx(i, j, bLen)] = Up
				min=f[idx(i-1, j, bLen)] + gap
			} else {
				pointer[idx(i, j, bLen)] = Stop
				min=threshold+1
			}

			f[idx(i, j, bLen)] = min

			//debug	
			//fmt.Println("i:",i," j:",j," min:",min," p:",pointer[idx(i, j, bLen)]," tr:",threshold," st:",start_next_at," fgp: ",first_good_prev," br:",we_broke_at," msaf:",nonstop_already_found) 

			if min<=threshold {
				//we are in good area
				if !nonstop_already_found {
					//good area just started, mark it
					nonstop_already_found=true
					start_next_at=j //makes no sense to start under stop
					if 1==j && first_stopped_head > i {
						first_good_prev=0
					} else {
						first_good_prev=j
					}
				}
				continue; //go on, next j
			}
			//if we are here, min > threshold 
			pointer[idx(i, j, bLen)] = Stop //the value is set already

			if j>= we_broke_at || j==bLen-1 {
				// we are under we_broke_at or we are at the end of the line
				//no chance to move right any more
				if i==aLen-1 {
					//if it is the last line, so we cannot rich the SE corner, so we break completely (give_up)
					give_up=true
					break
				}
				if nonstop_already_found {
					//ther was good are in the line, so we go on
					we_broke_at=j
					break
				}
				//if we are here, we did not find any good area on tis line, so we break completely (give_up)
				give_up=true
				break
			}
		} //j cycle
	} //i cycle

	//debug	
	//fmtmati(f,aLen,bLen)
	//fmtmatb(pointer,aLen,bLen)

	if (give_up) { //we gave up
		return "","",threshold+1,false
	}

	//debug	
	//log.Println("restoring..")

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
		} else {
			log.Fatalln("Unknown something is found on the alignment way. I am lost.")
		}

	}

	reverse(aBytes)
	reverse(bBytes)

	return string(aBytes), string(bBytes), dist, true
}

func Distance(a, b string, mismatch, gap, threshold int) (dist int, ok bool) {
	//this is the fast version, returns distance only if success -- 
	//no pointer matrix 

	aLen := len(a) + 1
	bLen := len(b) + 1

	maxLen := aLen
	if maxLen < bLen {
		maxLen = bLen
	}

	we_broke_at:=bLen
	give_up:=false
	//we_broke_at is a critical value. 
	//if we are on the right of we_broke_at (we are on the next line),
	//we do not check any directions other than left
	//if we do not open good cells on this line and we are under the we_broke_at
	//we do not test we_broke_at+1
	//we break completely by setting give_up flag to true


	start_next_at:=1
	//where from to start next line of alignment matrix
	first_good_prev:=0
	//usually, it is start_next_at-1, sometimes. start_next_at
	first_stopped_head:=aLen+1
	//which line is first to get >threshold distance and 


	if len(f) < aLen*bLen {
		f = make([]int, aLen*bLen)
	}

	for i := 1; i < aLen; i++ { //vertical
		dist := gap * i
		f[idx(i, 0, bLen)] = dist
		if dist > threshold {
			first_stopped_head=i
			break
		}
	}

	for j := 1; j < bLen; j++ { //horizontal
		dist := gap * j
		f[idx(0, j, bLen)] = dist
		if dist > threshold {
			we_broke_at=j
			//where the Stop appeared in a finshed line
			break
		}
	}

	for i := 1; i < aLen; i++ {
		if give_up {break;} //chau
		nonstop_already_found := false
		for j := start_next_at; j < bLen; j++ {
			var min int
			if (j<=we_broke_at) && (j>first_good_prev) { //we can test all 3 direction	
				matchMismatch := mismatch
				if a[i-1] == b[j-1] {
					matchMismatch = 0
				}

				min = f[idx(i-1, j-1, bLen)] + matchMismatch
				vgap := f[idx(i, j-1, bLen)] + gap
				hgap := f[idx(i-1, j, bLen)] + gap

				if hgap < min {
					min = hgap
				}
				if vgap < min {
					min = vgap
				}

				f[idx(i, j, bLen)] = min
			} else if (j>first_good_prev) {//ok, we can test only left
				min=f[idx(i, j-1, bLen)] + gap
			} else if (j<=we_broke_at) { //ok, maybe, up will help, left and NW are forbidden
				min=f[idx(i-1, j, bLen)] + gap
			} else {
				min=threshold+1
			}

			f[idx(i, j, bLen)] = min

			//debug	
			//log.Println("i:",i," j:",j," min:",min," tr:",threshold," st:",start_next_at," fgp: ",first_good_prev," br:",we_broke_at," msaf:",nonstop_already_found) 

			if min<=threshold {
				//we are in good area
				if !nonstop_already_found {
					//good area just started, mark it
					nonstop_already_found=true
					start_next_at=j //makes no sense to start under stop
					if 1==j && first_stopped_head > i {
						first_good_prev=0
					} else {
						first_good_prev=j
					}
				}
				continue; //go on, next j
			}
			//if we are here, min > threshold 

			if j>= we_broke_at || j==bLen-1 {
				// we are under we_broke_at or we are at the end of the line
				//no chance to move right any more
				if i==aLen-1 {
					//if it is the last line, so we cannot rich the SE corner, so we break completely (give_up)
					give_up=true
					break
				}
				if nonstop_already_found {
					//ther was good are in the line, so we go on
					we_broke_at=j
					break
				}
				//if we are here, we did not find any good area on tis line, so we break completely (give_up)
				give_up=true
				break
			}
		} //j cycle
	} //i cycle

	//debug	
	//fmtmati(f,aLen,bLen)

	if (give_up) { //we gave up
		return threshold+1,false
	}

	return f[idx(aLen-1, bLen-1, bLen)], true
}
