package models

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	ID       string `bson:"_id,omitempty" json:"id"`
	Username string `bson:"username" json:"username"`
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"-"`
	Role     Role   `bson:"role" json:"role"`
}