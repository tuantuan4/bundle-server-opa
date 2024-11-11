package models

type Auth struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Country   string `json:"country"`
	CreatedAt string `json:"created_at" gorm:"column:createdAt"` // Match the column name in the DB
	UpdatedAt string `json:"updated_at" gorm:"column:updatedAt"` // Match the column name in the DB
}
