package models

type Pair struct {
	Key string
	Value interface{}
	Join Join
}

type Join struct {
	Source string
	Target string
}