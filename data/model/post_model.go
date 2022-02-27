package model

import(
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostServiceServer struct{
}

type Post struct{
	Id			primitive.ObjectID	`json:"_id" bson:"_id"`
	Title		string				`json:"title" bson:"title"`
	Content		string				`json:"content" bson:"content"`
	User		string				`json:"user" bson:"user"`
	Votes		int64				`json:"votes" bson:"votes"`
}

type Posts []Post