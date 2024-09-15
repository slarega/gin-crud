package models

type User struct {
	GUID     int    `gorm:"type:int;primary_key"`
	Email    string `gorm:"type:varchar(100);unique"`
	Password string `gorm:"type:varchar(100)"`
}

type UserForm struct {
	Email    string `validate:"max=100" json:"email" example:"user@mail.com"`
	Password string `validate:"max=15,min=5" json:"password" example:"123456"`
}

func (f *UserForm) ToModel() User {
	return User{
		Email:    f.Email,
		Password: f.Password,
	}
}
