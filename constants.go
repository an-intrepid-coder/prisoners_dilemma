package main

const (
    /* TODO: Arbitrary decision depth. Also, using the
       methods discussed below to compare Cohorts of
       different depths.  */
    DECISION_DEPTH = 3 

    /* TODO: Currently, the simulation reaches the fitness
       goal much faster when the Cohort Size is smaller. 
       This is because the fitness goal is defined as 
       the percentage of the Cohort which can beat a randomly
       generated Rule. This is arguably a very low bar, considering
       the non-viability of most Classifier Rules in the immense
       2^64 search space. The next step in the cas/ API is to 
       determine to what extent these are local maximums. This
       can be done by running multiple Cohorts in parallel and
       testing their "Champions" against each other. As noted in
       John Holland's books and papers, this is a good way to
       try and solve some problems within certain constraints.
       Taking it a level or two higher in scope allows for solving
       more complex problems. By raising it one scope level higher
       (which is the next major improvement) it will be possible
       to objectively determine which of the below parameters
       leads to the best Classifier Rule for the game of 
       Prisoner's Dilemma, in the shortest amount of time.
       This opens the door for using the framework on other parts 
       of Game Theory, for example, which will hopefully involve 
       raising the scope further.  */

    COHORT_SIZE = 100
    NUM_ROUNDS = 100   
    REPRODUCTION_THRESHOLD = 3
    GENERATION_CAP = 10000 
    FITNESS_GOAL = 99.0
    
    RANDOM_SAMPLE_SIZE = 1000000

    USE_SYSTEM_TIME = -1
    SQUELCH_NOTIFICATIONS = -1

    COOPERATE = 0
    DEFECT = 1

    PUNISHMENT = 2
    REWARD = 1
    SUCKERS = 3
    TEMPTATION = 0

    MUTATION_FREQUENCY = 10000
)

