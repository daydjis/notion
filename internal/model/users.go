package model

type User struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	SecondName string `json:"second_name"`
	Age        string `json:"age"`
	Password   string `json:"password"`
}

type LoginInput struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type RegisterInput struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
type UserInput struct {
	Name     string `json:"name" binding:"required"`     // без этого JSON биндинг упадёт с ошибкой
	Password string `json:"password" binding:"required"` // пароль тоже обязателен в запросе
}
