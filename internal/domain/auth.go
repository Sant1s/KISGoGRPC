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
	Id string `db:"id" json:"id"`
}

type LoginUserResponse struct {
	Id         string `db:"id" json:"id"`
	Permission string `db:"permission" json:"permission"`
}
