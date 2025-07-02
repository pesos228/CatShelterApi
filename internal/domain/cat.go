package domain

type Cat struct {
	BaseModel
	Name   string
	Age    int16
	UserId string `gorm:"type:uuid;not null"`
}
