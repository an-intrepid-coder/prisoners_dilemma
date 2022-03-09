package util

// Returns a / b * 100 and checks for Divide-by-Zero:
func Percent(a float64, b float64) float64 {
    if b == 0 {
        return a // TODO: throw an error here
    }
    return a / b * 100.0
}

