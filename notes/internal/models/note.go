package models

// Note - структура для представления заметки
// Содержит поля ID, Name, Content и AuthorID
// Используется для сериализации в JSON и BSON
// Используется в качестве модели для работы с базой данных
type Note struct {
	//  ID - уникальный идентификатор заметки
	ID string `json:"id,omitempty" bson:"id,omitempty" `
	//  Name - название заметки
	Name string `json:"name,omitempty" bson:"name,omitempty"`
	//  Content - содержимое заметки
	Content string `json:"content,omitempty" bson:"content,omitempty"`
	//  AuthorID - идентификатор автора заметки
	AuthorID int `json:"author_id,omitempty" bson:"author_id,omitempty"`
}