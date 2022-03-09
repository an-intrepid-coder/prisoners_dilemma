module github.com/prisoners_dilemma

go 1.17

replace github.com/prisoners_dilemma/cas => ./cas

replace github.com/prisoners_dilemma/queue => ./queue

require (
	github.com/prisoners_dilemma/cas v0.0.0-00010101000000-000000000000
	github.com/prisoners_dilemma/queue v0.0.0-00010101000000-000000000000
	github.com/prisoners_dilemma/util v0.0.0-00010101000000-000000000000
)

replace github.com/prisoners_dilemma/util => ./util
