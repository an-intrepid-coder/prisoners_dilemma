package main

import (
    "fmt"
    "os"
    "strings"
    "time"

    "github.com/prisoners_dilemma/util"
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
           "beat" a randomly generated Prisoner's Dilemma Classifier Rule. 
           TODO: Allow floating point input for fitGoal for finer results.  */
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

    var seed int64
    if args["-seed="] == USE_SYSTEM_TIME {
        seed = time.Now().UnixNano()
    } else {
        seed = int64(args["-seed="])
    }

    squelch := args["-notifications="] == SQUELCH_NOTIFICATIONS

    fmt.Println("... computing ...")

    // Discover a rule for Prisoner's Dilemma:
    r := DiscoverPdRule(seed,
                        args["-cohortSize="],
                        squelch,
                        args["-numRounds="],
                        DECISION_DEPTH,
                        args["-rThreshold="],
                        args["-genCap="],
                        args["-fitGoal="])

    // Results:
    fmt.Println("Rule discovered! Results:")
    fmt.Printf("\tRule: ")
    for i := range r.Rule {
        fmt.Print(r.Rule[i])
    }
    fmt.Printf("\n")
    fmt.Printf("\tRule effectiveness: %.02f percent\n", r.RuleWinPercent)
    fmt.Printf("\tIt took %d / %d generations.\n", r.GenerationsUsed, r.GenerationCap)
    fmt.Printf("\tDecision depth used: %d rounds\n", r.DecisionDepth)
    fmt.Printf("\tCohort size used: %d Agents\n", r.CohortSize)
    fmt.Printf("\tNumber of rounds used: %d rounds/game\n", r.NumRounds)
    fmt.Printf("\tResource threshold used: %d\n", r.ResourceThreshold)
    fmt.Printf("\tCohort fitness goal used: %d percent\n", r.FitnessGoal)
    fmt.Printf("\tSeed used: %x\n", r.Seed)
    fmt.Printf("\tMutation frequency used: %.02f percent \n", util.Percent(1.0, float64(r.MutationFrequency)))
}

