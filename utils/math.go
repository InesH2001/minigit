package utils

func Max(a, b, c int) int {
    if a > b && a > c {
        return a
    }
    if b > c {
        return b
    }
    return c
}