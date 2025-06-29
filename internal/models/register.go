package models

type RegisterRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Message string `json:"message"`
	User    struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"user"`
}
