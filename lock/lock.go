package lock

// TODO: Unit tests

type Lock struct {
    size int
    busy []bool
}

func (lk *Lock) init(n int) {
    lk.size = n
    lk.busy = make([]bool, n)
}

func MakeLock(n int) Lock {
    lk := Lock{}
    lk.init(n)
    return lk
}

// Sets all members to "busy":
func (lk *Lock) ToggleAllBusy() {
    for i:= range lk.busy {
        lk.busy[i] = true
    }
}

// Releases the concurrency lock for a given member.  
func (lk *Lock) ToggleFinished(n int) {
    lk.busy[n] = false
}

// Returns true if all members are unlocked.  
func (lk *Lock) AllFinished() bool {
    for i := range lk.busy {
        if lk.busy[i] {
            return false
        }
    }
    return true
}

/* Delays execution until all members of the Lock are
   no longer busy.  */
func (lk *Lock) ConcurrentJoin() {
    for ;; {
        if lk.AllFinished() {
            break
        }
    }
}

