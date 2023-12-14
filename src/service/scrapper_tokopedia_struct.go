package service

type List struct {
	Page  int    `json:"page"`
	Total uint64 `json:"total"`
}
type Result struct {
	Status  bool   `json:"status"`
	File    string `json:"file"`
	Message string `json:"message"`
}

type Recommendation struct {
	TID            string    `json:"tID"`
	Source         string    `json:"source"`
	Title          string    `json:"title"`
	ForeignTitle   string    `json:"foreignTitle"`
	WidgetURL      string    `json:"widgetUrl"`
	Recommendation []Product `json:"recommendation"`
}

type AceProduct struct {
	Count    int64     `json:"count"`
	Products []Product `json:"products"`
}

type Product struct {
	Type      string `json:"type"`
	Generated bool   `json:"generated"`
	ID        string `json:"id"`
	Typename  string `json:"typename"`
}

type DetailProduct struct {
	ID                 int64         `json:"id"`
	URL                string        `json:"url"`
	ImageURL           string        `json:"image_url"`
	ImageURL700        string        `json:"image_url_700"`
	CategoryID         int64         `json:"category_id"`
	GaKey              string        `json:"ga_key"`
	CountReview        int64         `json:"count_review"`
	DiscountPercentage int64         `json:"discount_percentage"`
	IsPreorder         bool          `json:"is_preorder"`
	Name               string        `json:"name"`
	Price              string        `json:"price"`
	PriceInt           int64         `json:"price_int"`
	OriginalPrice      string        `json:"original_price"`
	Rating             int64         `json:"rating"`
	Wishlist           bool          `json:"wishlist"`
	Labels             []interface{} `json:"labels"`
	Badges             []Shop        `json:"badges"`
	Shop               Shop          `json:"shop"`
	LabelGroups        []Shop        `json:"label_groups"`
	Typename           string        `json:"__typename"`
}

type Shop struct {
	Type      string `json:"type"`
	Generated bool   `json:"generated"`
	ID        string `json:"id"`
	Typename  string `json:"typename"`
}
