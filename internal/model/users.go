package model

// User представляет структуру пользователя для ответа (без пароля)
type User struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	SecondName *string `json:"second_name,omitempty"`
	Age        *int    `json:"age,omitempty"`
	Password   string  `json:"-"` // "-" предотвращает сериализацию в JSON
}

// LoginInput представляет данные для входа пользователя
type LoginInput struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterInput представляет данные для регистрации нового пользователя
type RegisterInput struct {
	Name       string  `json:"name" binding:"required"`
	SecondName *string `json:"second_name,omitempty"`
	Age        *int    `json:"age,omitempty"`
	Password   string  `json:"password" binding:"required"`
}

// UserInput представляет общие данные пользователя, например для обновления информации
type UserInput struct {
	Name       string  `json:"name" binding:"required"`
	SecondName *string `json:"second_name,omitempty"`
	Age        *int    `json:"age,omitempty"`
}
