package dto

type UserResponse struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Email   string  `json:"email"`
	Balance float64 `json:"balance"`
}

type UserDetailResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
