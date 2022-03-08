package cas

import (
    "math/rand"
    "sort"
)

type CohortMetadata struct { 
    Size int
    Generation int
    Fitness float64
} 

/* The Cohort is mostly a slice of Agents, but it also
   tracks metadata and handles concurrent generational
   processing as well as post-generation evolution. It is
   the current highest level of the cas API, but soon I will
   have a few more levels so that Cohorts which are testing
   different parameters can be compared directly, and more.  */
type Cohort struct {
    size int
    members []Agent
    busy []bool
    generation int
    fitness float64
    Metadata CohortMetadata
}

func (c *Cohort) SetFitness(n float64) {
    c.fitness = n
}

func (c *Cohort) Fitness() float64 {
    return c.fitness
}

func (c *Cohort) Generation() int {
    return c.generation
}

func (c *Cohort) init(n int) {
    m := n
    if m % 2 != 0 {
        m++
    }
    c.size = m
    c.members = make([]Agent, m)
    c.busy = make([]bool, m)
    for i := 0; i < m; i++ {
        c.members[i] = MakeAgent()
        c.busy[i] = false
    }
    c.generation = 0
    c.fitness = 0.0
}

func MakeCohort(n int) Cohort {
    c := Cohort{}
    c.init(n)
    return c
}

func (c *Cohort) SortByResources() {
    sort.Slice(c.members, func(i, j int) bool { // descending order
        return c.members[i].Resources() > c.members[j].Resources()
    }) 
}

/* To evolve the Cohort, I followed the steps suggested by John Holland's
   paper. The size of the generation never changes. The next generation is
   first filled by the offspring of fit parents, then by the fit parents
   themselves, and then by however many of the rest can fit. The Cohort is
   shuffled at the end of every Evolve() just to be safe.  */
func (c *Cohort) Evolve(n int, g int) {

    // Sort generation in descending order by resources
    c.SortByResources()

    // Next generation:
    s := make([]Agent, c.size, c.size)

    r := 0
    for ;; {
        a := c.members[r]
        if a.Resources() < n  || r >= c.size {
            break
        }
        s[r] = a
        s[r].TakeResources(n)
        s[r].Metadata.Resources = s[r].Resources()
        r++
    }
    // r is now the index in s after the last reproducing agent was
    // inserted
    
    q, h := r - 1, r
    if q % 2 != 0 {
        q--
    }
    // q is now the index in s of the last reproducing agent
    // h is now the index in c.members which is the next agent

    for i := 0; i < q; i += 2 {
        a, b := &s[i], &s[i + 1]
        p := a.Combine(b)
        for j := range p {
            p[j].Metadata.Generation = g
            if r < c.size {
                s[r] = p[j]
                r++
            }
        }
    }
    // r is now the index in s after the last inserted offspring

    for ; r < c.size; r++ { 
        s[r] = c.members[h]
        h++
    }
    // s should now be full of the next generation

    c.members = s
    rand.Shuffle(c.size, func(i, j int) {
        c.members[i], c.members[j] = c.members[j], c.members[i]
    })
    c.Metadata = CohortMetadata{c.size, c.generation, c.fitness}
    c.generation++ 
}

func (c *Cohort) Member(i int) *Agent {
    return &c.members[i]
}

func (c *Cohort) Size() int {
    return c.size
}

/* "sets" the concurrent "lock". ConcurrentJoin() won't
   finish while any members are busy.  */
func (c *Cohort) ToggleAllBusy() {
    for i := range c.busy {
        c.busy[i] = true
    }
}

// Releases the concurrency lock for a given member.  
func (c *Cohort) ToggleFinished(n int) {
    c.busy[n] = false
}

// Returns true if all members are unlocked.  
func (c *Cohort) AllFinished() bool {
    for i := range c.busy {
        if c.busy[i] {
            return false
        }
    }
    return true
}

/* Delays execution until all members of the Cohort are
   no longer busy.  */
func (c *Cohort) ConcurrentJoin() {
    for ;; {
        if c.AllFinished() {
            break
        }
    }
}

