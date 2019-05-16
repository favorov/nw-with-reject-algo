===============================================================================
nwwreject - Needleman-Wunsch Alignment with rejecting in Go
===============================================================================

-------------------------------------------------------------------------------
About
-------------------------------------------------------------------------------

Here, we implement Needleman-Wunsch [1] algorithm with rejection to estimate the 
distance for almost-exact-match pairs of sequences. The implemntatin in Go

The Needleman-Wunsch global alignment algorith computes the alignment score and 
optimal global alignment. The modification we put here (NW with reject) rejects to proceed if all there is no way to biuld an alignment better than some threshold


It is based on: https://github.com/aebruno/nwalgo realisation of NW in Go.

-------------------------------------------------------------------------------
Install
-------------------------------------------------------------------------------

Fetch from github::

    $ go get github.com/favorov/nwwreject/...

-------------------------------------------------------------------------------
Usage
-------------------------------------------------------------------------------
Parameters: seq1, seq2, match, mismatch, gap, threshold.
match, mismatch and gap are integer gains to the alignment score.
threshold : if we cannot build anything not worse than threshold, 
we stop trying.

Returns: Distance returns score and bool status, e.g. whether nw succeded (if not, returned score value is the threshold), Align strings for alignment and then score and status.

Align 2 DNA sequences::

    $nwwreject -seq1 GAAAAAAT -seq2 GAAT 
    GAAAAAAT
    GAA----T
    Score: 0

    $nwwreject -seq1 GAAAAAAT -seq2 GAAT --threshold 1
		Sequences differ too much.

The package provide 2 functions: Distance and Align.
They do the same, but Distance does not return the alignment and thus work faster.


From code::

    package main

    import (
        "github.com/favorov/nwwreject"
    )

    func main() {
        score_1_d,ok_1_d := nwwreject.Distance("GAAAAAAT", "GAAT", 1, -1, -1, 1) //rejected, returns threshold, false
        score_0_d,ok_0_d := nwwreject.Distance("GAAAAAAT", "GAAT", 1, -1, -1, 0) //sucess, returns score, true
        aln1, aln2, score, ok_1_a := nwwreject.Align("GAAAAAAT", "GAAT", 1, -1, -1, 1)//rejected, returns "","", threshold, false
        aln1, aln2, score, ok_0_a := nwwreject.Align("GAAAAAAT", "GAAT", 1, -1, -1, 0) //returns aln1,aln2,score,true
    }

-------------------------------------------------------------------------------
References
-------------------------------------------------------------------------------

[1] http://en.wikipedia.org/wiki/Needleman-Wunsch_algorithm
