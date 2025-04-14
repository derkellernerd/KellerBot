package model

import (
	"time"

	"gorm.io/gorm"
)

type ChatCommand struct {
	gorm.Model
	Command          string `gorm:"uniqueIndex"`
	Used             uint64
	TimeoutInSeconds uint64
	LastUsed         time.Time
	Action           string
}

func (c *ChatCommand) Increment() {
	c.Used += 1
	c.LastUsed = time.Now()
}
