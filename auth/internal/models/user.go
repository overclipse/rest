package models

import "golang.org/x/crypto/bcrypt"

// User представляет модель пользователя в системе
// Он содержит ID, имя пользователя и пароль
// Используется для хранения и управления данными пользователей
// Важно отметить, что пароль должен храниться в зашифрованном виде
type User struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"password,omitempty" gorm:"not null"`
}

// bcryptCost определяет стоимость хеширования пароля
// Чем выше значение, тем больше времени требуется для хеширования
// Это влияет на безопасность, но также увеличивает время обработки запросов
// Рекомендуется использовать значение от 10 до 14 для bcrypt
const bcryptCost = 12

// Метод для хеширования пароля
func (u *User) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// Метод для проверки пароля
func (u *User) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}