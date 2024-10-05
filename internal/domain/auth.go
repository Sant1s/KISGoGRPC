package domain

type RegisterUserRequest struct {
	Login        string `db:"nickname" json:"login"`
	PasswordHash string `db:"password_hash" json:"password_hash"`
	Permission   string `db:"permission" json:"permission"`
}

type LoginUserRequest struct {
	Login        string `db:"nickname" json:"login"`
	PasswordHash string `db:"password_hash" json:"password_hash"`
}

type RegisterUserResponse struct {
	Code   int64  `json:"code"`
	Output string `json:"output"`
}

type LoginUserResponse struct {
	Code   int64  `json:"code"`
	Output string `json:"output"`
}
