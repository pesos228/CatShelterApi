package domain

type BaseModel struct {
	Id string `gorm:"type:uuid;primary_key"`
}
