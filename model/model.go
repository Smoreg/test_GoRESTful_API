package model

import "gopkg.in/mgo.v2/bson"

type Meme struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `bson:"Name" json:"name"`
}