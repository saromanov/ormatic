package models

type Insert struct {
	TableName string
	Pairs []Pair
}