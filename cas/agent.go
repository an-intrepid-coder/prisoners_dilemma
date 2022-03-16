package cas

var numAgents = 0

type AgentMetadata struct {
    Id int
    Resources int
    Generation int
    Wins int
    Losses int
    WinRate float64
}

type Agent struct { 
    id int
    // NOTE: More complex Agents might have multiple Classifiers.
    classifier *Classifier 
    // NOTE: More complex Agents might have multiple resource types.
    resources int
    /* NOTE: The Metadata, in this case, represents stuff which not only
       describes the Agent for reporting, but also might vary from experiment
       to experiment (beyond Prisoner's Dilemma).  */
    Metadata AgentMetadata
}

func (a *Agent) Depth() int {
    return a.classifier.Depth()
}

// Returns the classifier rule for this agent:
func (a *Agent) Rule() []int {
    return a.classifier.Rule()
}

func (a *Agent) init(d int) {
    a.id = numAgents
    numAgents++
    c := MakeClassifier(d)
    a.classifier = &c
    a.resources = 0
    a.Metadata = AgentMetadata{a.id, 0, 0, 0, 0, 0.0}
}

// Calculates the Agent's move based on the Classifier's logic:
func (a *Agent) CalcMove(s []int) int {
    return a.classifier.CalcMove(s)
}

/* Creates two new Agents with Classifier rules that are
   genetically crossed over reproductions of the parent
   Classifiers.  */
func (a *Agent) Combine(b *Agent, freq int) []Agent {
    c, d := MakeAgent(a.Depth()), MakeAgent(a.Depth())
    s := a.classifier.Combine(b.classifier, freq)
    c.classifier = &s[0]
    d.classifier = &s[1]
    return []Agent{c, d}
}

func MakeAgent(d int) Agent {
    a := Agent{}
    a.init(d)
    return a
}

func (a *Agent) Id() int {
    return a.id
}

func (a *Agent) Resources() int {
    return a.resources
}

func (a *Agent) AddResources(n int) {
    a.resources += n
}

func (a *Agent) TakeResources(n int) {
    a.resources -= n
    if a.resources < 0 {
        a.resources = 0
    }
}

