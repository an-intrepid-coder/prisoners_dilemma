module github.com/prisoners_dilemma/cas

go 1.17

replace github.com/prisoners_dilemma/lock => ../lock

require (
	github.com/prisoners_dilemma/lock v0.0.0-00010101000000-000000000000
	github.com/prisoners_dilemma/util v0.0.0-00010101000000-000000000000
)

replace github.com/prisoners_dilemma/util => ../util
