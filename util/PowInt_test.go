package util

import (
    "math"
    "math/rand"
    "testing"
)

func TestPowInt(t *testing.T) {
    n := 9999999
    max := 64
    for i := 0; i < n; i++ {
        a := rand.Intn(max)
        c := math.Pow(2, float64(a))
        d := Pow2Int(a)
        if int(c) != d {
            t.Fatalf("%x != %x\n", c, d)
        }
    }
}

