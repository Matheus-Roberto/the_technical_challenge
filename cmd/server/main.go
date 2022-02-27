package main

import (
	"context"
	"fmt"
	"log"
	"net"

	//"github.com/Matheus-Roberto/the_technical_challenge_klever/pb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

var mongoCtx context.Context

var db *mongo.Client
var postDB *mongo.Collection


func main() {
	grpcServer := grpc.NewServer()

	fmt.Println("Starting server on port :50051...")
	listener, err := net.Listen("tcp", "localhost:50051")

	if err != nil {
		log.Fatalf("Unable to listen on port :50051: %v", err)
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

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	fmt.Println("Server succesfully started on port :50051")

	fmt.Println("\nStopping the server...")
	grpcServer.Stop()
	listener.Close()
	fmt.Println("Closing MongoDB connection")
	db.Disconnect(mongoCtx)
	fmt.Println("Done.")
}


