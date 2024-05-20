package repository

type Filter struct {
	Limit     *int
	Offset    *int
	OrderBy   *string
	Direction *string
}
