package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type LogsetLogined struct {
	Timestamp time.Time             `bson:"timestamp" json:"timestamp"`
	Metadata  LogsetLoginedMetadata `bson:"metadata" json:"metadata"`
	UserAgent string                `bson:"user_agent"`
	Detail    interface{}           `bson:"detail"`
}

type LogsetLoginedMetadata struct {
	UserID   primitive.ObjectID `bson:"user_id"`
	ClientIP string             `bson:"client_ip"`
	Version  string             `bson:"version"`
	Source   string             `bson:"source" json:"source"`
}

func (x *LogsetLogined) SetDetail(v interface{}) {
	x.Detail = v
}

func (x *LogsetLogined) SetVersion(v string) {
	x.Metadata.Version = v
}

func NewLogsetLogined(uid primitive.ObjectID, ip string, source string, useragent string) *LogsetLogined {
	return &LogsetLogined{
		Timestamp: time.Now(),
		Metadata: LogsetLoginedMetadata{
			UserID:   uid,
			ClientIP: ip,
			Source:   source,
		},
		UserAgent: useragent,
	}
}