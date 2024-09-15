package response

type UserResponse struct {
	GUID     int    `json:"GUID"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

type UserOkResponse struct {
	Data UserResponse `json:"data"`
}

type UsersOkResponse struct {
	Data []UserResponse `json:"data"`
}
