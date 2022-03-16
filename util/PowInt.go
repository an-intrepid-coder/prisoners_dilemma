package util

// Function for small powers of 2 using ints:
func Pow2Int(a int) int {
    b := a
    if b > 64 {
        b = 64
    }
    m := 1
    for i := 0; i < b; i++ {
        m *= 2
    }
    return m
}

