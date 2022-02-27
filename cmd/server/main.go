package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/Matheus-Roberto/the_technical_challenge_klever/data/model"
	pb "github.com/Matheus-Roberto/the_technical_challenge_klever/pb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostServiceServer struct {}

var mongoCtx context.Context
var db *mongo.Client
var postDB *mongo.Collection

func (server *PostServiceServer) CreatePost(ctx context.Context, request *pb.CreatePostRequest) (*pb.CreatePostResponse, error) {

	post := request.GetPost()

	data := model.Post{
		Id:			primitive.NewObjectID(),
		Title:		post.GetTitle(),
		Content:	post.GetContent(),
		User:		post.User,
	}

	result, err := postDB.InsertOne(mongoCtx, data)

	if err != nil {

		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}

	oid := result.InsertedID.(primitive.ObjectID)
	post.Id = oid.Hex()
	post.Votes = 0

	return &pb.CreatePostResponse{Post: post}, nil
}

func (server *PostServiceServer) DeletePost(ctx context.Context, request *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {
	oid, err := primitive.ObjectIDFromHex(request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}

	filter := bson.M{"_id": oid}

	result := postDB.FindOneAndDelete(ctx, filter)

	decoded := model.Post{}
	err = result.Decode(&decoded)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find/delete post with id %s: %v", request.GetId(), err))
	}
	return &pb.DeletePostResponse{
		Success: true,
	}, nil
}

func (server *PostServiceServer) GetPost(ctx context.Context, request *pb.GetPostRequest) (*pb.GetPostResponse, error) {
	oid, err := primitive.ObjectIDFromHex(request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}
	result := postDB.FindOne(ctx, bson.M{"_id": oid})

	data := model.Post{}

	if err := result.Decode(&data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find post with Object Id %s: %v", request.GetId(), err))
	}

	response := &pb.GetPostResponse{
		Post: &pb.Post{
			Id:			oid.Hex(),
			Title:		data.Title,
			Content:	data.Content,
			User: 		data.User,
			Votes:   	data.Votes,
		},
	}
	return response, nil
}

func (server *PostServiceServer) ListPost(ctx context.Context, request *pb.ListPostRequest) (*pb.ListPostResponse, error) {

	filter := bson.M{}
	data := []*pb.Post{}

	cursor, _ := postDB.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var p model.Post
		err := cursor.Decode(&p)
		if err != nil {
			continue
		}
		data = append(data, &pb.Post{
			Id:      	p.Id.Hex(),
			Title:   	p.Title,
			Content: 	p.Content,
			User: 		p.User,
			Votes:   	p.Votes,
		})
	}
	return &pb.ListPostResponse{
		Post: data,
	}, nil
}

func (server *PostServiceServer) UpdatePost(ctx context.Context, request *pb.UpdatePostRequest) (*pb.UpdatePostResponse, error) {

	post := request.GetPost()

	oid, err := primitive.ObjectIDFromHex(post.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Could not convert the supplied post id to a MongoDB ObjectId: %v", err),
		)
	}

	update := bson.M{
		"title":   	post.GetTitle(),
		"content": 	post.GetContent(),
		"user":		post.GetUser(),
	}

	filter := bson.M{"_id": oid}

	result := postDB.FindOneAndUpdate(ctx, filter, bson.M{"$set": update}, options.FindOneAndUpdate().SetReturnDocument(1))

	decoded := model.Post{}
	err = result.Decode(&decoded)
	if err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Could not find post with supplied ID: %v", err),
		)
	}
	return &pb.UpdatePostResponse{
		Post: &pb.Post{
			Id:      	decoded.Id.Hex(),
			Title:   	decoded.Title,
			Content: 	decoded.Content,
			User:		decoded.User,
			Votes:   	decoded.Votes,
		},
	}, nil
}

func (server *PostServiceServer) UpVote(ctx context.Context, request *pb.UpVoteRequest) (*pb.UpVoteResponse, error) {
	post := request.GetPost()
	oid, err := primitive.ObjectIDFromHex(post.GetId())

	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Could not convert the supplied post id to a MongoDB ObjectId: %v", err),
		)
	}

	filter := bson.M{"_id": oid}

	result := postDB.FindOneAndUpdate(ctx, filter, bson.M{"$inc": bson.M{"votes": 1}}, options.FindOneAndUpdate().SetReturnDocument(1))

	decoded := model.Post{}
	err = result.Decode(&decoded)
	if err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Could not find post with supplied ID: %v", err),
		)
	}
	return &pb.UpVoteResponse{
		Post: &pb.Post{
			Id:      	decoded.Id.Hex(),
			Title:   	decoded.Title,
			Content: 	decoded.Content,
			User:		decoded.User,
			Votes:   	decoded.Votes,
		},
	}, nil
}

func (server *PostServiceServer) DownVote(ctx context.Context, request *pb.DownVoteRequest) (*pb.DownVoteResponse, error) {

	post := request.GetPost()

	oid, err := primitive.ObjectIDFromHex(post.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Could not convert the supplied post id to a MongoDB ObjectId: %v", err),
		)
	}

	filter := bson.M{"_id": oid}

	result := postDB.FindOneAndUpdate(ctx, filter, bson.M{"$inc": bson.M{"votes": -1}}, options.FindOneAndUpdate().SetReturnDocument(1))

	decoded := model.Post{}
	err = result.Decode(&decoded)
	if err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Could not find post with supplied ID: %v", err),
		)
	}
	return &pb.DownVoteResponse{
		Post: &pb.Post{
			Id:      	decoded.Id.Hex(),
			Title:   	decoded.Title,
			Content: 	decoded.Content,
			User:		decoded.User,
			Votes:   	decoded.Votes,
		},
	}, nil
}

func main() {
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	srv := &PostServiceServer{}

	pb.RegisterPostServiceServer(grpcServer, srv)

	fmt.Println("Starting server on port :5000...")
	listener, err := net.Listen("tcp", "localhost:5000")

	if err != nil {
		log.Fatalf("Unable to listen on port :5000: %v", err)
	}
	
	fmt.Println("Connecting to MongoDB...")
	mongoCtx = context.Background()
	db, err = mongo.Connect(mongoCtx, options.Client().ApplyURI("mongodb+srv://mroberto:Rob812345@cluster0.gwmb5.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping(mongoCtx, nil)
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %v\n", err)
	} else {
		fmt.Println("Connected to Mongodb")
	}

	postDB = db.Database("redditdb").Collection("post")

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	fmt.Println("Server succesfully started on port :5000")

	ch := make(chan os.Signal)

	signal.Notify(ch, os.Interrupt)

	<-ch

	fmt.Println("\nStopping the server...")
	grpcServer.Stop()
	listener.Close()
	fmt.Println("Closing MongoDB connection...")
	db.Disconnect(mongoCtx)
	fmt.Println("Done.")
}


