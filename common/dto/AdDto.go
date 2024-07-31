package dto

type AdDTO struct {
	ID                uint
	Text              string
	ImagePath         string
	Bid               int64
	Website           string
	TotalCost         int64
	Impressions       int
	BalanceAdvertiser int64
	Score             float64
	Keywords          []string
}
