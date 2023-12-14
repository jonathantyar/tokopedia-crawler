package model

type Product struct {
	BaseEntitySF
	Name        string  `json:"name" gorm:"column:name;type:varchar(255)" export:"title:Name;"`
	Description string  `json:"description" gorm:"column:description;type:text" export:"title:Description"`
	ImageLink   string  `json:"image_link" gorm:"column:image_link;type:varchar(255)" export:"title:ImageLink"`
	Price       float64 `json:"price" gorm:"type:numeric(64,2)" export:"title:Price;"`
	Rating      float64 `json:"rating" gorm:"column:rating" export:"title:Rating;"`
	Merchant    string  `json:"merchant" gorm:"column:merchant;type:varchar(255)" export:"title:Merchant;"`
}
