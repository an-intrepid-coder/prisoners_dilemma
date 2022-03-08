package cas

import (
     "math/rand"
)

const (
    /* NOTE: Right now the Decision Depth is hard-coded
       at 3 turns, which takes 64 Classifier "bits" (but 
       not really bits as each one here is an integer word 
       in size despite being a 1 or 0).  */
    CLASSIFIER_BITS = 64

    /* In John Holland's paper, a mutation frequency of 
       1/10,000 was suggested as a good starting point.  */
    MUTATION_FREQ = 10000
)

type Classifier struct {
    rule []int
}

func (c *Classifier) Rule() []int {
    return c.rule
}

func (c *Classifier) init() {
    c.rule = make([]int, CLASSIFIER_BITS)
    for i := 0; i < CLASSIFIER_BITS; i++ {
        c.rule[i] = rand.Intn(2)
    }
}

func MakeClassifier() Classifier {
    c := Classifier{}
    c.init()
    return c
}

/* CalcMove iterates over the input slice of 1s and 0s as if it 
   were computing the conversion of a number from base-2 to base-10. 
   It doesn't really matter whether it reads the slice from left-to-
   right or vice versa, as it leads to a unique mapping either way
   as long as it is consistent. The slice's elements are converted 
   to a base-10 index, and the the value of the Classifier Rule at 
   that index is the response to the input.  */
func (c *Classifier) CalcMove(s []int) int {
    r := 0
    for i, v := range s {
        if v == 1 {
            x := 1
            for j := 0; j < i; j++ {
                x *= 2
            }
            r += x
        }
    }
    return c.rule[r]
}

/* As suggested in John Holland's paper, this combines two
   Classifiers by performing "Genetic Crossover" on their
   rules.  */
func (c *Classifier) Combine(d *Classifier) []Classifier {
    a, b := MakeClassifier(), MakeClassifier()

    // A random pivot is chosen:
    p := rand.Intn(CLASSIFIER_BITS)
    
    // Points before the pivot are overlaid on to the 
    // new Rules:
    for i := 0; i < p; i++ {
        a.rule[i] = d.rule[i]
        b.rule[i] = c.rule[i]
    }
    // Points from the pivot onward are swapped
    // from parent to offspring:
    for i := p; i < CLASSIFIER_BITS; i++ {
        a.rule[i] = c.rule[i]
        b.rule[i] = d.rule[i]
    }
    // Mutation chance is applied:
    f := func(n int) int { 
        if rand.Intn(MUTATION_FREQ) == 0 {
            return (n + 1) % 2 
        }
        return n
    } 
    for i := 0; i < CLASSIFIER_BITS; i++ { 
        a.rule[i] = f(a.rule[i])
        b.rule[i] = f(b.rule[i])
    }
    return []Classifier{a, b}
}

