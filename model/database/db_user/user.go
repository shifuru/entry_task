package db_user

import "time"

type Role string

const (
	Admin  Role = "admin"
	Normal Role = "user"
)

type User struct {
	Id        uint      `json:"id"  gorm:"primary_key;column:id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Role      Role      `json:"role"`
	TopicId   uint      `json:"topic_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (User) TableName() string {
	return "user_tab"
}
