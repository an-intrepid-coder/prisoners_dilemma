package queue

import (
    "testing"
    "time"
    "math/rand"
)

func TestQueue(t *testing.T) {
    rand.Seed(time.Now().UnixNano())
    queueSize := 9
    max := 99999
    q := MakeQueue(queueSize)

    if !q.Empty() {
        t.Fatalf("queue not empty when initialized")
    }

    for i := 0; i < queueSize; i++ {
        if q.Full() {
            t.Fatalf("queue prematurely filled")
        }
        q.Insert(rand.Intn(max))
    }

    if !q.Full() {
        t.Fatalf("queue not reporting full after fill step")
    }

    for ; !q.Empty() ; {
        q.Del()
    }

    if !q.Empty() {
        t.Fatalf("queue not empty after drain step")
    }
}

