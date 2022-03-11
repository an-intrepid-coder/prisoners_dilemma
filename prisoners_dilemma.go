package main

import (
    "fmt"
    "math/rand"

    "github.com/prisoners_dilemma/cas"
    "github.com/prisoners_dilemma/lock"
    "github.com/prisoners_dilemma/queue"
    "github.com/prisoners_dilemma/util"
)

// This function will wrap the whole program thus far and enable 
// taking it a scope further in the future:

type DiscoverPdRuleMetadata struct {
    DecisionDepth int
    CohortSize int
    NumRounds int
    ResourceThreshold int
    GenerationCap int
    FitnessGoal int
    /* TODO: Currently, it won't reproduce exactly the same results for 
       the same seed. There is some research to do in order to fix that.
       I was hoping it would lead to exact reproducible results, but runs
       still vary in detail even if the seed is the same. It is a long-term
       goal to ensure reproducibility to that degree.  */
    Seed int64
    GenerationsUsed int
    Rule []int
    RuleWinPercent float64
    MutationFrequency int
}

func DiscoverPdRule(seed int64,
                    cohortSize int,
                    squelch bool,
                    numRounds int,
                    depth int,
                    rThreshold int,
                    genCap int,
                    fitGoal int) DiscoverPdRuleMetadata {
    // Seed the PRNG:
    if seed == USE_SYSTEM_TIME {
        rand.Seed(seed)
    } else {
        rand.Seed(seed)
    }

    // Make a Cohort: 
    c := cas.MakeCohort(cohortSize)

    if !squelch {
        fmt.Println("Discovering Prisoner's Dilemma Rule...")
    }

    // Process/Evolve Loop:
    for ;; {
        if !squelch {
            fmt.Printf("Beginning generation #%d\n", c.Generation())
        }

        // Process the generation:
        pdGeneration(&c, numRounds, depth)

        // Evolve the Cohort:
        c.Evolve(rThreshold, c.Generation() + 1)

        if !squelch {
            fmt.Printf("\tCohort Fitness: %.02f\n", c.Fitness())
        }

        // End simulation if goal reached or cap hit:
        if c.Generation() >= genCap || c.Fitness() >= float64(fitGoal) {
            break
        }
    }

    // Print the results for the Cohort:
    if !squelch {
        fmt.Println("Cohort evolution complete!")
        fmt.Printf("\tCohort has %d members.\n", c.Metadata.Size)
        fmt.Printf("\tCohort ran for %d generations.\n", c.Metadata.Generation)
        fmt.Printf("\tCohort fitness score is: %.02f\n", c.Metadata.Fitness)
    }

    // Find the best member of the Cohort:
    if !squelch {
        fmt.Println("Finding champion...")
    }
    v := pdChamp(&c, numRounds, depth)

    // TODO: Progress notifications for pdChamp() stage

    // Print initial results:
    if !squelch {
        fmt.Println("Champion found!")
        fmt.Printf("\tChampion ID #%d\n", v.Metadata.Id)
        fmt.Printf("\tChampion Generation: %d (age: %d)\n", v.Metadata.Generation, c.Generation() - v.Metadata.Generation - 1) 
        fmt.Printf("\tChampion Resources: %d\n", v.Metadata.Resources)
        fmt.Printf("\tChampion Wins/Losses: %d / %d (%.02f)\n", v.Metadata.Wins, v.Metadata.Losses, v.Metadata.WinRate)
        fmt.Print("\tChampion Classifier Rule: ")
        x := v.Rule()
        for i := range x {
            fmt.Print(x[i])
        }
        fmt.Print("\n")
    }

    // Print final results:
    if !squelch {
        fmt.Printf("Testing Champion against %d random samples...\n", RANDOM_SAMPLE_SIZE)
    }
    cr := pdTestAgentAgainstSamples(v, numRounds, depth, RANDOM_SAMPLE_SIZE)
    if !squelch {
        fmt.Printf("\tChampion win/loss percentage vs. random samples: %.02f\n", cr)
    }

    // Collect and return metadata:
    md := DiscoverPdRuleMetadata{}
    md.DecisionDepth = depth
    md.CohortSize = cohortSize
    md.NumRounds = numRounds
    md.ResourceThreshold = rThreshold
    md.GenerationCap = genCap
    md.FitnessGoal = fitGoal
    md.Seed = seed
    md.GenerationsUsed = c.Generation()
    md.Rule = v.Rule()
    md.RuleWinPercent = cr
    md.MutationFrequency = MUTATION_FREQUENCY
    return md
}

/* Tests an Agent against a given number of random Agents (preferably a very large
number) to get a good idea of what its general effectiveness is as a Prisoner's
Dilemma Classifier Rule. This is so that the "Champions" of multiple Cohorts can
be compared against each other.  */
func pdTestAgentAgainstSamples(a *cas.Agent, rounds int, depth int, samples int) float64 {
    cw, cl := 0, 0 
    lk := lock.MakeLock(samples)
    lk.ToggleAllBusy()
    for i := 0; i < samples; i++ {
        go func(k int) {
            b := cas.MakeAgent()
            w := pdGame(a, &b, rounds, false, DECISION_DEPTH)
            if w == a {
                cw++
            } else {
                cl++
            }
            lk.ToggleFinished(k)
        }(i)
    }
    lk.ConcurrentJoin()
    return util.Percent(float64(cw), float64(cw + cl))
}

// Plays a game of Prisoner's Dilemma and returns a pointer to the winner:
func pdGame(a *cas.Agent, b *cas.Agent, rounds int, counts bool, depth int) *cas.Agent { 
    // Random player goes first:              
    p := []*cas.Agent{a, b}
    t := rand.Intn(2)

    /* NOTE: Counts two kinds of scores: total cumulative
    "points" over the game, and number of rounds "won". In the
    spirit of the game, I am defining a "won" round as being one
    in which the player scores less than or equal to their 
    opponent, as opposed to strictly less. I am going to use the
    cumulative points as the default for now, but later on I will
    measure and contrast both during runtime. Note that, as in
    golf, the lower score is better. The Wikipedia article on 
    Prisoner's Dilemma uses negative scores, while some other
    people use positive ones. This has no real effect on the
    game as long as the comparisons are consistent. Although 
    a round can be "won" by both players in the event of
    mutual cooperation or defection, the game as a sequence
    of rounds can only go to one of the players.  */

    // Cumulative "points":
    sa, sb := 0, 0

    // Cumulative "wins":
    wa, wb := 0, 0

    // Queues are used to hold turn memory:
    qa, qb := queue.MakeQueue(depth), queue.MakeQueue(depth)
    for i := 0; i < depth; i++ {
        qa.Insert(0)
        qb.Insert(0)
    }

    // Players face off for n rounds:
    for i := 0; i < rounds; i++ {

        /* NOTE: One could also randomize the turn order each round.
           That could make a difference for some Classifiers. I will
           experiment with that down the road.  */

        // Player decision/score this round:
        var ra, rb int 

        // Each player takes a turn each round:
        for j := 0; j < 2; j++ {
            
            /* NOTE: The queues' contents are combined in to a slice of 1s
               and 0s which the Agents' Classifiers will treat as a 
               binary number. The default state {0, 0, 0, 0, 0, 0}
               is the same as if both players had cooperated for 3
               rounds in a row. This is just a starting point based on
               John Holland's paper. You could use many more rounds of 
               depth for this, up to the practical limits of computation.  */
            s := append(qa.Contents(), qb.Contents() ...)
           
            // Current player makes a decision to COOPERATE or DEFECT:
            if t == 0 {
                ra = p[t].CalcMove(s) 
                qa.Del()
                qa.Insert(ra)
            } else {
                rb = p[t].CalcMove(s)
                qb.Del()
                qb.Insert(rb)
            }
            t = (t + 1) % 2
        }

        // Tally wins/points: 
        if ra == COOPERATE && rb == COOPERATE {
            ra, rb = REWARD, REWARD
        } else if ra == COOPERATE && rb == DEFECT {
            ra, rb = SUCKERS, TEMPTATION
        } else if ra == DEFECT && rb == COOPERATE {
            ra, rb = TEMPTATION, SUCKERS
        } else {
            ra, rb = PUNISHMENT, PUNISHMENT
        }
        sa += ra
        sb += rb
        if ra <= rb {
            wa++
        }
        if rb <= ra {
            wb++
        }
    }
    /* Award resources to the one who got the least points. In the
       future, it may also award for the most "wins". In the unlikely
       event of a draw here neither get a reward.  */
    var w *cas.Agent
    if sa < sb {
        w = a
        if counts {
            a.Metadata.Wins++
            b.Metadata.Losses++
        }
    } else {
        w = b
        if counts {
            a.Metadata.Losses++
            b.Metadata.Wins++
        }
    }
    if counts {
        w.AddResources(1)
        w.Metadata.Resources++
    }
    // Update metadata:
    f := func(u *cas.Agent) {
        y, z := u.Metadata.Wins, u.Metadata.Losses
        u.Metadata.WinRate = util.Percent(float64(y), float64(y + z))
    }
    f(a)
    f(b)
    return w
}

func pdGeneration(c *cas.Cohort, rounds int, depth int) {
    c.Lock.ToggleAllBusy()
    numFit := 0
    // Send each individual game to a goroutine for concurrent processing:
    for i := 0; i < c.Size(); i++ { 
        go func(j int) { 
            a, b := c.Member(j), cas.MakeAgent()
            w := pdGame(a, &b, rounds, true, depth) 
            if *w == *a {
                numFit++
            }
            c.Lock.ToggleFinished(j)
        }(i) 
    } 
    c.Lock.ConcurrentJoin()
    // Calculate fitness for current generation:
    c.SetFitness(util.Percent(float64(numFit), float64(c.Size())))
}

/* To find the champ, each member of the Cohort plays each other member of
the Cohort (including themselves), and the winner is the one with the
most wins.  */
func pdChamp(c *cas.Cohort, rounds int, depth int) *cas.Agent {
    r := make([]int, c.Size()) 
    c.Lock.ToggleAllBusy()
    for i := range r {
        go func(k int) {
            a := c.Member(k)
            for j := range r {
                b := c.Member(j)
                w := pdGame(a, b, rounds, false, depth)
                if w.Id() == a.Id() {
                    r[k]++
                }
            }
            c.Lock.ToggleFinished(k)
        }(i)
    }
    c.Lock.ConcurrentJoin()
    var v *cas.Agent
    m := 0
    for i := range r {
        if r[i] > m {
            m = r[i]
            v = c.Member(i)
        }
    }
    return v
}

