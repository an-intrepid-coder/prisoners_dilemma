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
        "-decisionDepth=": DECISION_DEPTH,
        "-cohortSize=": COHORT_SIZE,
        "-numRounds=": NUM_ROUNDS,
        "-rThreshold=": REPRODUCTION_THRESHOLD,
        "-genCap=": GENERATION_CAP,
        "-fitGoal=": FITNESS_GOAL,
        "-seed=": USE_SYSTEM_TIME,
        "-notifications=": SQUELCH_NOTIFICATIONS,
        "-mutationFrequency=": MUTATION_FREQUENCY,
        "-controlSampleSize=": RANDOM_SAMPLE_SIZE,
        "-gamesPerGen=": GAMES_PER_GENERATION,
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

    depth := args["-decisionDepth="]
    if depth > DEPTH_CAP {
        depth = DEPTH_CAP
    }

    fmt.Println("... computing ...")

    // Discover a rule for Prisoner's Dilemma:
    r := DiscoverPdRule(seed,
                        args["-cohortSize="],
                        squelch,
                        args["-numRounds="],
                        depth,
                        args["-rThreshold="],
                        args["-genCap="],
                        args["-fitGoal="],
                        args["-mutationFrequency="],
                        args["-controlSampleSize="],
                        args["-gamesPerGen="])

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
    fmt.Printf("\tSize of search space: 2^%d 'bits'\n", util.Pow2Int(r.DecisionDepth * 2))
    fmt.Printf("\tNumber of rounds used: %d rounds/game\n", r.NumRounds)
    fmt.Printf("\tResource threshold used: %d\n", r.ResourceThreshold)
    fmt.Printf("\tCohort fitness goal used: %d percent\n", r.FitnessGoal)
    fmt.Printf("\tSeed used: %x\n", r.Seed)
    fmt.Printf("\tMutation frequency used: %.02f percent\n", util.Percent(1.0, float64(r.MutationFrequency)))
    fmt.Printf("\tControl Sample Size: %d random Agents\n", r.ControlSampleSize)
    fmt.Printf("\tGames per Agent per Generation: %d\n", r.GamesPerGen)
}

