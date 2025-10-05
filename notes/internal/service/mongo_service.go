package service

import (
	"context"
	"fmt"
	"notes/internal/config"
	"notes/internal/database"
	"notes/internal/errors"
	"notes/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoService - реализация интерфейса Service для работы с базой данных MongoDB
type MongoService struct {
	db         *mongo.Client     // Указатель на клиент MongoDB для выполнения операций с базой данных
	collection *mongo.Collection // Коллекция заметок в MongoDB
}

// Проверка, что MongoService реализует интерфейс Service
// Это позволяет гарантировать, что MongoService соответствует контракту интерфейса Service
var _ Service = (*MongoService)(nil)

// NewService - конструктор для создания нового экземпляра MongoService
// Он принимает конфигурацию, используется для подключения к MongoDB
func NewService(cfg *config.Config) (Service, error) {
	db, err := database.NewDatabase(cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrDatabaseConnection, err)
	}

	// Получаем коллекцию заметок
	collection := db.Database(cfg.DB_NAME).Collection(cfg.DB_COLLECTION)

	return &MongoService{
		db:         db,
		collection: collection,
	}, nil
}

// Create создает новую заметку в базе данных
func (m *MongoService) Create(ctx context.Context, note models.Note) (*models.Note, error) {
	// Создаем новый ObjectId для заметки
	result, err := m.collection.InsertOne(ctx, bson.M{
		"name":      note.Name,
		"content":   note.Content,
		"author_id": note.AuthorID,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrNoteCreation, err)
	}

	// Получаем ID созданной заметки
	insertedID := result.InsertedID.(primitive.ObjectID)
	note.ID = insertedID.Hex() // Преобразуем ObjectID в строку

	return &note, nil
}

// GetByID получает заметку по идентификатору
func (m *MongoService) GetByID(ctx context.Context, id string) (*models.Note, error) {

	// Преобразуем строку ID в ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrInvalidNoteID, err)
	}

	var note models.Note
	err = m.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&note)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%w: заметка с ID %s не найдена", errors.ErrNoteNotFound, id)
		}
		return nil, fmt.Errorf("%w: %v", errors.ErrDatabaseOperation, err)
	}

	note.ID = objectID.Hex()

	return &note, nil
}

// GetAll получает все заметки из базы данных для конкретного автора
func (m *MongoService) GetAll(ctx context.Context, authorId int) ([]models.Note, error) {
	// Запрос к базе данных
	filter := bson.M{"author_id": authorId}

	cursor, err := m.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrDatabaseOperation, err)
	}
	defer cursor.Close(ctx)

	var notes []models.Note
	for cursor.Next(ctx) {
		var note models.Note
		if err := cursor.Decode(&note); err != nil {
			return nil, fmt.Errorf("%w: %v", errors.ErrDecodeNote, err)
		}
		notes = append(notes, note)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrIterationNotes, err)
	}

	return notes, nil
}

// Update обновляет существующую заметку
func (m *MongoService) Update(ctx context.Context, note models.Note) (*models.Note, error) {
	// Преобразуем строку ID в ObjectID
	objectID, err := primitive.ObjectIDFromHex(note.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrInvalidNoteID, err)
	}

	// Сначала получаем существующую заметку для получения AuthorID
	var existingNote models.Note
	err = m.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&existingNote)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%w: заметка с ID %s не найдена", errors.ErrNoteNotFound, note.ID)
		}
		return nil, fmt.Errorf("%w: %v", errors.ErrDatabaseOperation, err)
	}

	update := bson.M{
		"$set": bson.M{
			"name":    note.Name,
			"content": note.Content,
		},
	}

	result, err := m.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrNoteUpdate, err)
	}

	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("%w: заметка с ID %s не найдена", errors.ErrNoteNotFound, note.ID)
	}
	// Устанавливаем AuthorID из существующей заметки
	note.AuthorID = existingNote.AuthorID

	return &note, nil
}

// Delete удаляет заметку по идентификатору
func (m *MongoService) Delete(ctx context.Context, id string) error {
	// Преобразуем строку ID в ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrInvalidNoteID, err)
	}

	// Сначала получаем заметку для получения AuthorID перед удалением
	var existingNote models.Note
	err = m.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&existingNote)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("%w: заметка с ID %s не найдена", errors.ErrNoteNotFound, id)
		}
		return fmt.Errorf("%w: %v", errors.ErrDatabaseOperation, err)
	}

	result, err := m.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrNoteDeletion, err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("%w: заметка с ID %s не найдена", errors.ErrNoteNotFound, id)
	}

	return nil
}

// Close закрывает соединение с базой данных
func (m *MongoService) Close() error {
	return database.CloseDB(m.db, &config.Config{Timeout: 10})
}
