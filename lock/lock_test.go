package lock

import (
    "testing"
    "time"
    "math/rand"
)

func TestLock(t *testing.T) {
    rand.Seed(time.Now().UnixNano())
    tests := 999
    max := 9999
    lk1 := MakeLock(tests)
    lk1.ToggleAllBusy()
    for i := 0; i < tests; i++ {
        go func(h int) {
            n := rand.Intn(max)
            lk2 := MakeLock(n)
            lk2.ToggleAllBusy()
            for j := 0; j < n; j++ {
                go func(k int){
                    lk2.ToggleFinished(k)
                }(j)
            }
            lk2.ConcurrentJoin()
            lk1.ToggleFinished(h)
        }(i)
    }
    lk1.ConcurrentJoin()
}

