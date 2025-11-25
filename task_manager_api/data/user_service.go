package data

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/skdebela/task_manager_api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

const defaultUserServiceTimeout = 10 * time.Second

var (
	once     sync.Once
	firstUserAdmin bool = true // Track if first user has been created
	mu       sync.Mutex
)

type UserService struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func NewUserService(ctx context.Context, uri, dbName string) (*UserService, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultUserServiceTimeout)
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

	collection := client.Database(dbName).Collection("users")

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err = collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		return nil, err
	}

	return &UserService{
		client:     client,
		collection: collection,
	}, nil
}

func (s *UserService) Disconnect(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, defaultUserServiceTimeout)
	defer cancel()
	return s.client.Disconnect(ctx)
}

func (s *UserService) GetUsers(ctx context.Context) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultUserServiceTimeout)
	defer cancel()

	cur, err := s.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var users []models.User
	for cur.Next(ctx) {
		var u models.User
		if err := cur.Decode(&u); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, cur.Err()
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (models.User, bool, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultUserServiceTimeout)
	defer cancel()

	var user models.User
	filter := bson.M{"email": email}
	err := s.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.User{}, false, nil
		}
		return models.User{}, false, err
	}
	return user, true, nil
}

func (s *UserService) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultUserServiceTimeout)
	defer cancel()

	mu.Lock()
	if firstUserAdmin {
		user.Role = models.RoleAdmin
		firstUserAdmin = false
	} else {
		user.Role = models.RoleUser
	}
	mu.Unlock()

	hashed, err := hashPassword(user.Password)
	if err != nil {
		return models.User{}, err
	}
	user.Password = hashed
	user.ID = primitive.NewObjectID().Hex() 

	_, err = s.collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return models.User{}, errors.New("email already exists")
		}
		return models.User{}, err
	}

	user.Password = ""
	return user, nil
}

func (s *UserService) PromoteToAdmin(ctx context.Context, id string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultUserServiceTimeout)
	defer cancel()

	filter := bson.M{"_id": id} 
	update := bson.M{"$set": bson.M{"role": models.RoleAdmin}}

	res, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}

	return res.ModifiedCount > 0, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (models.User, bool, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultUserServiceTimeout)
	defer cancel()

	var user models.User
	filter := bson.M{"_id": id}
	err := s.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.User{}, false, nil
		}
		return models.User{}, false, err
	}
	return user, true, nil
}