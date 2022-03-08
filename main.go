package main

import (
    "fmt"
    "math/rand"
    "os"
    "strings"
    "time"

    "github.com/prisoners_dilemma/cas"
    "github.com/prisoners_dilemma/queue"
)

const (
    DECISION_DEPTH = 3 
    COHORT_SIZE = 100
    NUM_ROUNDS = 100   
    REPRODUCTION_THRESHOLD = 3
    GENERATION_CAP = 10000 
    FITNESS_GOAL = 99.0

    USE_SYSTEM_TIME = -1
    SQUELCH_NOTIFICATIONS = -1

    COOPERATE = 0
    DEFECT = 1

    PUNISHMENT = 2
    REWARD = 1
    SUCKERS = 3
    TEMPTATION = 0
)

func main() {
    // format: ./prisoners_dilemma -seed=42 -notifications=1 (etc...)
    args := map[string]int {
        /* NOTE: When an additional scope is testing multiple Cohorts then all of 
           these parameters can be in flux at once. For now only one Cohort is being 
           tested at a time using these args. Aside from the decision depth (for 
           now), they are all tunable via command line. The defaults currently
           work pretty fast, but it would be a neat additional project to use
           another Genetic Algorithm to determine the best parameters for this 
           one. That is probably a near-term step so that this framework can be
           used for examining more complex problems.  */
        "-decisionDepth=": DECISION_DEPTH,
        // TODO: Allow for arbitrary decision depth.
        "-cohortSize=": COHORT_SIZE,
        "-numRounds=": NUM_ROUNDS,
        "-rThreshold=": REPRODUCTION_THRESHOLD,
        "-genCap=": GENERATION_CAP,
        /* Fitness is in this case defined by the percentage of a Cohort which can
           "beat" a randomly generated Prisoner's Dilemma Classifier Rule.  */
        "-fitGoal=": FITNESS_GOAL,
        "-seed=": USE_SYSTEM_TIME,
        "-notifications=": SQUELCH_NOTIFICATIONS,
        // TODO: Allow adjustable mutation frequency.
    }
    // Collect and parse the args from the command line (if any):
    for i := range os.Args { 
        s := strings.SplitAfter(os.Args[i], "=")
        _, ok := args[s[0]]
        if ok {
            var v int
            fmt.Sscanf(s[1], "%d", &v)
            args[s[0]] = v
        }
    }

    // Seed the PRNG:
    s := args["-seed="]
    if s == USE_SYSTEM_TIME {
        rand.Seed(time.Now().UnixNano())
    } else {
        rand.Seed(int64(s))
        /* TODO: Currently, it won't reproduce exactly the same results for 
           the same seed. There is some research to do in order to fix that.
           I was hoping it would lead to exact reproducible results, but runs
           still vary in detail even if the seed is the same. It is a long-term
           goal to ensure reproducibility to that degree.  */
    }

    // Make a test Cohort: 
    c := cas.MakeCohort(args["-cohortSize="])

    fmt.Println("Evolving...")

    // Process/Evolve Loop:
    for ;; {
        if args["-notifications="] != SQUELCH_NOTIFICATIONS {
            fmt.Printf("Beginning generation #%d\n", c.Generation())
            fmt.Printf("\tCohort Fitness: %.02f\n", c.Fitness())
        }

        // Process the generation:
        pdGeneration(&c, args["-numRounds="], DECISION_DEPTH)

        // Evolve the Cohort:
        c.Evolve(args["-rThreshold="], c.Generation() + 1)

        if args["-notifications="] != SQUELCH_NOTIFICATIONS {
            fmt.Printf("\tCohort Fitness: %.02f\n", c.Fitness())
        }

        // End simulation if goal reached or cap hit:
        if c.Generation() >= args["-genCap="] || c.Fitness() >= float64(args["-fitGoal="]) {
            break
        }
    }

    // Print the results for the Cohort:
    fmt.Println("Cohort evolution complete!")
    fmt.Printf("\tCohort has %d members.\n", c.Metadata.Size)
    fmt.Printf("\tCohort ran for %d generations.\n", c.Metadata.Generation)
    fmt.Printf("\tCohort fitness score is: %f\n", c.Metadata.Fitness)

    // Find the best member of the Cohort:
    fmt.Println("Finding champion...")
    v := pdChamp(&c, args["-numRounds="], DECISION_DEPTH)

    // TODO: Progress notifications for pdChamp() stage

    // Print final results:
    fmt.Println("Champion found!")
    fmt.Printf("\tChampion ID #%d\n", v.Metadata.Id)
    fmt.Printf("\tChampion Generation: %d\n", v.Metadata.Generation)
    fmt.Printf("\tChampion Resources: %d\n", v.Metadata.Resources)
    fmt.Printf("\tChampion Wins/Losses: %d / %d (%.02f)\n", v.Metadata.Wins, v.Metadata.Losses, v.Metadata.WinRate)
    fmt.Print("\tChampion Classifier Rule: ")
    x := v.Rule()
    for i := range x {
        fmt.Print(x[i])
    }
    fmt.Print("\n")
}

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
    for i := 0; i < DECISION_DEPTH; i++ {
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
        }
    } else {
        w = b
        if counts {
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
        u.Metadata.WinRate = float64(y) / float64(y + z) * 100.0
    }
    f(a)
    f(b)
    return w
}

func pdGeneration(c *cas.Cohort, rounds int, depth int) {
    c.ToggleAllBusy()

    numFit := 0
    // Send each individual game to a goroutine for concurrent processing:
    for i := 0; i < c.Size(); i++ { 
        go func(j int) { 
            a, b := c.Member(j), cas.MakeAgent()
            w := pdGame(a, &b, rounds, true, depth) 
            if *w == *a {
                numFit++
            }
            c.ToggleFinished(j)
        }(i) 
    } 
    c.ConcurrentJoin()

    // Calculate fitness for current generation:
    c.SetFitness(float64(numFit) / float64(c.Size()) * 100.0)
}

/* To find the champ, each member of the Cohort plays each other member of
   the Cohort (including themselves), and the winner is the one with the
   most wins.  */
func pdChamp(c *cas.Cohort, rounds int, depth int) cas.Agent {
    r := make([]int, c.Size()) 

    c.ToggleAllBusy()

    for i := range r {
        go func(j int) {
            a := c.Member(i)
            for j := range r {
                b := c.Member(j)
                w := pdGame(a, b, rounds, false, depth)
                if *w == *a {
                    r[i]++
                }
            }
            c.ToggleFinished(j)
        }(i)
    }

    c.ConcurrentJoin()

    var v *cas.Agent
    m := 0
    for i := range r {
        if r[i] > m {
            m = r[i]
            v = c.Member(i)
        }
    }
    return *v
}

