package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Builder struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty" json:"_id"`
	Parent      *primitive.ObjectID `bson:"parent" json:"parent"`
	Name        string              `bson:"name" json:"name"`
	Kind        string              `bson:"kind" json:"kind"`
	Icon        string              `bson:"icon" json:"icon"`
	Description string              `bson:"description" json:"description"`
	Schema      *BuilderSchema      `bson:"schema" json:"schema"`
	Sort        int64               `bson:"sort" json:"sort"`
	CreateTime  time.Time           `bson:"create_time" json:"create_time"`
	UpdateTime  time.Time           `bson:"update_time" json:"update_time"`
}

type BuilderSchema struct {
	Key    string               `bson:"key" json:"key"`
	Fields []BuilderSchemaField `bson:"fields" json:"fields"`
}

type BuilderSchemaField struct {
	Name        string       `bson:"name" json:"name"`
	Key         string       `bson:"key" json:"key"`
	Type        string       `bson:"type" json:"type"`
	Required    bool         `bson:"required" json:"required"`
	Private     bool         `bson:"private" json:"private"`
	Default     interface{}  `bson:"default" json:"default"`
	Placeholder string       `bson:"placeholder" json:"placeholder"`
	Description string       `bson:"description" json:"description"`
	Option      *FieldOption `bson:"option" json:"option"`
	Sort        int64        `bson:"sort" json:"sort"`
}

type FieldOption struct {
	// Type: number
	Max     int64 `bson:"max" json:"max"`
	Min     int64 `bson:"min" json:"min"`
	Decimal int64 `bson:"decimal" json:"decimal"`
	// Type: date,dates
	Time bool `bson:"time" json:"time"`
	// Type: radio,checkbox,select
	Enum []FieldOptionEnum `bson:"enum" json:"enum"`
	// Type: ref
	Ref     string   `bson:"ref" json:"ref"`
	RefKeys []string `bson:"ref_keys" json:"ref_keys"`
	// Type: manual
	Component string `bson:"component" json:"component"`
	// Type: other
	Multiple bool `bson:"multiple" json:"multiple"`
}

type FieldOptionEnum struct {
	Label string      `bson:"label" json:"label"`
	Value interface{} `bson:"value" json:"value"`
}

func SetBuilders(ctx context.Context, db *mongo.Database) (err error) {
	var ns []string
	if ns, err = db.ListCollectionNames(ctx, bson.M{"name": "builders"}); err != nil {
		return
	}
	var jsonSchema bson.D
	if err = LoadJsonSchema("builder", &jsonSchema); err != nil {
		return
	}
	if len(ns) == 0 {
		option := options.CreateCollection().SetValidator(jsonSchema)
		if err = db.CreateCollection(ctx, "builders", option); err != nil {
			return
		}
	} else {
		if err = db.RunCommand(ctx, bson.D{
			{"collMod", "builders"},
			{"validator", jsonSchema},
			{"validationLevel", "strict"},
		}).Err(); err != nil {
			return
		}
	}
	return
}
