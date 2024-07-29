package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Polygon struct {
	ID         primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	Type       string                 `json:"type" bson:"type"`
	Geometry   Geometry               `json:"geometry" bson:"geometry"`
	Properties map[string]interface{} `json:"properties" bson:"properties"`
}

type Geometry struct {
	Type        string          `json:"type" bson:"type"`
	Coordinates [][][][]float64 `json:"coordinates" bson:"coordinates"`
}
