package dto

type InteractionType uint8

const (
	ImpressionType InteractionType = 0
	ClickType      InteractionType = 1
)

type CustomToken struct {
	Interaction       InteractionType `json:"interaction"`
	AdID              uint            `json:"ad_id"`
	PublisherUsername string          `json:"publisher_username"`
	RedirectPath      string          `json:"redirect_path"`
	CreatedAt         int64           `json:"created_at"`
}
