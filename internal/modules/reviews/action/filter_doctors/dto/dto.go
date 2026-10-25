package dto

type ReviewsSort int32

const (
	ReviewsSortDefault ReviewsSort = iota
	ReviewsSortDescending
	ReviewsSortAscending
)
