package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Matheus-Roberto/the_technical_challenge_klever/pb"
	"google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial("localhost:5000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewPostServiceClient(conn)

	postCreate := pb.Post{
		Title:   "title_default",
		Content: "content_default",
		User:    "user_default",
	}

	reqCreate := &pb.CreatePostRequest{
		Post: &postCreate,
	}

	resCreate, err := client.CreatePost(context.Background(), reqCreate)
	if err != nil {
		fmt.Println("could not create:")

	} else {
		fmt.Println("post created")
		fmt.Println(resCreate.Post)
	}

	reqGetId := &pb.GetPostRequest{Id: resCreate.Post.Id}
	resGetId, err := client.GetPost(context.Background(), reqGetId)
	if err != nil {
		fmt.Println("could not get by id:")

	} else {
		fmt.Println("getted by id ")
		fmt.Println(resGetId.Post)
	}
	postUpdate := pb.Post{
		Id:      resCreate.Post.Id,
		Title:   "title",
		Content: "content",
		User:    "user",
	}
	reqUpdate := &pb.UpdatePostRequest{
		Post: &postUpdate,
	}
	resUpdate, err := client.UpdatePost(context.Background(), reqUpdate)
	if err != nil {
		fmt.Println("could not updated")

	} else {
		fmt.Println("post updated")
		fmt.Println(resUpdate.Post)
	}

	postvote := pb.Post{
		Id:      resCreate.Post.Id,
	}
	reqUpvote := &pb.UpVoteRequest{
		Post: &postvote,
	}
	resUpvote, err := client.UpVote(context.Background(), reqUpvote)
	if err != nil {
		fmt.Println("could not upVoted")

	} else {
		fmt.Println("post upVoted")
		fmt.Println(resUpvote.Post)
	}

	reqDownVote := &pb.DownVoteRequest{
		Post: &postvote,
	}
	resDownVote, err := client.DownVote(context.Background(), reqDownVote)
	if err != nil {
		fmt.Println("could not DownVoted")

	} else {
		fmt.Println("post DownVoted")
		fmt.Println(resDownVote.Post)
	}

	reqDelete := &pb.DeletePostRequest{
		Id: resCreate.Post.Id,
	}
	resDelete, err := client.DeletePost(context.Background(), reqDelete)
	if err != nil {
		fmt.Println("could not deleted")

	} else {
		fmt.Println("post deleted")
		fmt.Println(resDelete.Success)
	}

	req := &pb.ListPostRequest{}
		res, err := client.ListPost(context.Background(), req)
		if err != nil {
			fmt.Println("could not get")
		}else{
			fmt.Println("posts getted")
			fmt.Println(res.Post)
		}
}

