package dto

type InteractionType uint8

const (
	ImpressionType InteractionType = 0
	ClickType      InteractionType = 1
)

type CustomToken struct {
	Interaction  InteractionType `json:"interaction"`
	AdID         uint            `json:"ad_id"`
	PublisherID  uint            `json:"publisher_id"`
	RedirectPath string          `json:"redirect_path"`
	Bid          int64           `json:"bid"`
	CreatedAt    int64           `json:"created_at"`
}
