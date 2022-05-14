package gormutils

import "gorm.io/gorm"

// Predicate is a string that acts as a condition in the where clause
type Predicate string

const (
	EqualPredicate              = Predicate("=")
	NotEqualPredicate           = Predicate("<>")
	GreaterThanPredicate        = Predicate(">")
	GreaterThanOrEqualPredicate = Predicate(">=")
	SmallerThanPredicate        = Predicate("<")
	SmallerThanOrEqualPredicate = Predicate("<=")
	LikePredicate               = Predicate("LIKE")
	InPredicate                 = Predicate("IN")
	NotInPredicate              = Predicate("NOT IN")
)

type QueryBuilder struct {
	table string
	db    *gorm.DB
	order []string
	group []string
	where []struct {
		prefix string
		value  interface{}
	}
	limit  int
	offset int
}
