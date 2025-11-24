package data

import (
	"context"
	"errors"
	"time"

	"github.com/skdebela/task_manager_api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	defaultTimeout = 10 * time.Second
)

type TaskService struct {
	client		*mongo.Client
	collection	*mongo.Collection
}

func NewTaskService(ctx context.Context, uri, dbName string) (*TaskService, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(ctx)
		return nil, err
	}

	col := client.Database(dbName).Collection("tasks")

	indexModel := mongo.IndexModel{
		Keys:		bson.D{{Key: "id", Value: 1}},
		Options:	options.Index().SetUnique(true),
	}

	if _, err := col.Indexes().CreateOne(ctx, indexModel); err != nil {
		return nil, err
	}

	return &TaskService{
		client:		client,
		collection: col,
	}, nil
}

func (s *TaskService) Disconnect(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	return s.client.Disconnect(ctx)
}


func (s *TaskService) GetTasks(ctx context.Context) ([]models.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	cur, err := s.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []models.Task
	for cur.Next(ctx) {
		var t models.Task
		if err := cur.Decode(&t); err != nil {
			return nil, err
		}
		out = append(out, t)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *TaskService) GetTaskByID(ctx context.Context, id string) (models.Task, bool, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var t models.Task
	filter := bson.D{{Key: "id", Value: id}}
	err := s.collection.FindOne(ctx, filter).Decode(&t)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Task{}, false, nil
		}
		return models.Task{}, false, err
	}
	return t, true, nil
}

func (s *TaskService) AddTask(ctx context.Context, task models.Task) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	_, err := s.collection.InsertOne(ctx, task)
	if mongo.IsDuplicateKeyError(err) {
		return errors.New("task with this id already exists")
	}
	return err
}

func (s *TaskService) UpdateTask(ctx context.Context, id string, updates bson.D) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	filter := bson.D{{Key: "id", Value: id}}
	res, err := s.collection.UpdateOne(ctx, filter, updates)
	if err != nil {
		return false, err
	}
	return res.MatchedCount > 0, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id string) (bool, error) {
	
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	filter := bson.D{{Key: "id", Value: id}}
	res, err := s.collection.DeleteOne(ctx, filter)
	if err != nil {
		return false, err
	}
	return res.DeletedCount > 0, nil
}