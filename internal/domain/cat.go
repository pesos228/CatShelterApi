package domain

type Cat struct {
	BaseModel
	Name string
	Age  int16
	User *User
}
