package main

const (
    /* TODO: Arbitrary decision depth. Also, using the
       methods discussed below to compare Cohorts of
       different depths.  */
    DECISION_DEPTH = 3 

    COHORT_SIZE = 300
    NUM_ROUNDS = 100   
    REPRODUCTION_THRESHOLD = 10
    GENERATION_CAP = 10000 
    FITNESS_GOAL = 95
    GAMES_PER_GENERATION = 10 
    RANDOM_SAMPLE_SIZE = 1000000 
    MUTATION_FREQUENCY = 10000

    USE_SYSTEM_TIME = -1
    SQUELCH_NOTIFICATIONS = -1

    COOPERATE = 0
    DEFECT = 1

    PUNISHMENT = 2
    REWARD = 1
    SUCKERS = 3
    TEMPTATION = 0
)

