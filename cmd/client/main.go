package main

import (
	"log"
	"net/http"

	"github.com/Matheus-Roberto/the_technical_challenge_klever/pb"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial("localhost:5000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewPostServiceClient(conn)

	request := gin.Default()

	request.POST("create", func(context *gin.Context) {
		post := pb.Post{}

		if err := context.ShouldBindJSON(&post); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Missing values"})
			return
		}

		req := &pb.CreatePostRequest{
			Post: &post,
		}
		res, err := client.CreatePost(context, req)
		if err != nil {
			context.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.JSON(http.StatusOK, res.Post)
	})

	request.GET("/post/:id", func(context *gin.Context) {
		id := context.Param("id")

		req := &pb.GetPostRequest{Id: id}
		res, err := client.GetPost(context, req)
		if err != nil {
			context.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.JSON(http.StatusOK, res.Post)
	})

	request.GET("/all/post/", func(context *gin.Context) {

		req := &pb.ListPostRequest{}
		res, err := client.ListPost(context, req)
		if err != nil {
			context.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		context.JSON(http.StatusOK, res)
	})

	request.DELETE("/post/:id", func(context *gin.Context) {
		id := context.Param("id")

		req := &pb.DeletePostRequest{Id: id}
		res, err := client.DeletePost(context, req)
		if err != nil {
			context.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.JSON(http.StatusOK, gin.H{
			"Id:":           id,
			"Post deleted:": res.Success,
		})

	})

	request.PUT("/update/post/", func(context *gin.Context) {
		post := pb.Post{}

		if err := context.ShouldBindJSON(&post); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Missing values"})
			return
		}

		req := &pb.UpdatePostRequest{
			Post: &post,
		}
		res, err := client.UpdatePost(context, req)
		if err != nil {
			context.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.JSON(http.StatusOK, res.Post)
	})

	request.PUT("/post/upvote", func(context *gin.Context) {
		post := pb.Post{}
		if err := context.ShouldBindJSON(&post); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Missing values"})
			return
		}
		req := &pb.UpVoteRequest{
			Post: &post,
		}

		res, err := client.UpVote(context, req)
		if err != nil {
			context.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.JSON(http.StatusOK, res.Post)
	})

	request.PUT("/post/downvote", func(c *gin.Context) {
		post := pb.Post{}

		if err := c.ShouldBindJSON(&post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Missing values"})
			return
		}

		req := &pb.DownVoteRequest{
			Post: &post,
		}

		res, err := client.DownVote(c, req)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, res.Post)
	})

	if err := request.Run(":8052"); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
