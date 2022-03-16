package cas

import (
    "math/rand"

    "github.com/prisoners_dilemma/util"
)

type Classifier struct {
    rule []int
    depth int
}

func (c *Classifier) Depth() int {
    return c.depth
}

func (c *Classifier) Rule() []int {
    return c.rule
}

func (c *Classifier) init(d int) { 
    b := util.Pow2Int(d * 2)       
    c.rule = make([]int, b)       
    for i := 0; i < b; i++ {
        c.rule[i] = rand.Intn(2)
    }
    c.depth = d
}

func MakeClassifier(d int) Classifier {
    c := Classifier{}
    c.init(d)
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
func (c *Classifier) Combine(d *Classifier, freq int) []Classifier {
    a, b := MakeClassifier(c.Depth()), MakeClassifier(c.Depth())

    // A random pivot is chosen:
    l := util.Pow2Int(c.Depth() * 2)
    p := rand.Intn(l)
    
    // Points before the pivot are overlaid on to the 
    // new Rules:
    for i := 0; i < p; i++ {
        a.rule[i] = d.rule[i]
        b.rule[i] = c.rule[i]
    }
    // Points from the pivot onward are swapped
    // from parent to offspring:
    for i := p; i < l; i++ {
        a.rule[i] = c.rule[i]
        b.rule[i] = d.rule[i]
    }
    // Mutation chance is applied:
    f := func(n int) int { 
        if rand.Intn(freq) == 0 {
            return (n + 1) % 2 
        }
        return n
    } 
    for i := 0; i < l; i++ { 
        a.rule[i] = f(a.rule[i])
        b.rule[i] = f(b.rule[i])
    }
    return []Classifier{a, b}
}

