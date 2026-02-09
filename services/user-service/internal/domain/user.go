package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                   uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email                string    `gorm:"not null;uniqueIndex"`
	Firstname            string
	Lastname             string
	Username             string `gorm:"not null"`
	Profile              string
	Phonenumber          string
	Birthday             *time.Time `gorm:"type:date"`
	AutoShareAge         int        `gorm:"default:0"`
	IsAutoShareEnabled   bool       `gorm:"default:false"`
	IsAutoShareTriggered bool       `gorm:"default:false"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	JoinedGroups         []GroupMember `gorm:"foreignKey:UserID"`
	Friends              []User        `gorm:"many2many:friend_lists;joinForeignKey:UserID;joinReferences:FriendID"`
}

type FriendList struct {
	UserID        uuid.UUID `json:"user_id" gorm:"type:uuid;primaryKey;not null"`
	FriendID      uuid.UUID `json:"friend_id" gorm:"type:uuid;primaryKey;not null"`
	Status        string    `json:"status" gorm:"type:varchar(20);default:'PENDING';not null"`
	IsCloseFriend bool      `gorm:"default:false;not null"`
	User          User      `gorm:"foreignKey:UserID"`
	Friend        User      `gorm:"foreignKey:FriendID"`
}

type FriendLog struct {
	ID        uuid.UUID `gorm:"primaryKey;default:gen_random_uuid()"`
	OwnerID   uuid.UUID `gorm:"not null;index"`
	FriendID  uuid.UUID `gorm:"not null;index"`
	Metadata  string    `json:"metadata" gorm:"type:jsonb;default:'{}'"`
	LogType   string    `gorm:"not null"`
	Messages  string    `gorm:"not null"`
	CreatedBy uuid.UUID `gorm:"not null"`
	CreatedAt time.Time
}

type PrivateMessage struct {
	ID         uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	SenderID   uuid.UUID  `json:"sender_id" gorm:"type:uuid;not null;index"`
	ReceiverID uuid.UUID  `json:"receiver_id" gorm:"type:uuid;not null;index"`
	MsgType    string     `json:"msg_type" gorm:"type:varchar(20);default:'TEXT';not null"`
	Content    string     `json:"content" gorm:"type:text"`
	Metadata   string     `json:"metadata" gorm:"type:jsonb;default:'{}'"`
	CreatedAt  time.Time  `json:"created_at" gorm:"index"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at" gorm:"index"`
	Sender     *User      `json:"sender,omitempty" gorm:"foreignKey:SenderID"`
	Receiver   *User      `json:"receiver,omitempty" gorm:"foreignKey:ReceiverID"`
}
