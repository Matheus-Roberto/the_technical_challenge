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
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewPostServiceClient(conn)

	r := gin.Default()

	r.POST("create", func(c *gin.Context) {
		post := pb.Post{}

		if err := c.ShouldBindJSON(&post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Missing values"})
			return
		}

		req := &pb.CreatePostRequest{
			Post: &post,
		}
		res, err := client.CreatePost(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, res.Post)
	})

	r.GET("/post/:id", func(c *gin.Context) {
		id := c.Param("id")

		req := &pb.GetPostRequest{Id: id}
		res, err := client.GetPost(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, res.Post)
	})

	r.GET("/all/post/", func(c *gin.Context) {

		req := &pb.ListPostRequest{}
		res, err := client.ListPost(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, res)
	})

	r.DELETE("/post/:id", func(c *gin.Context) {
		id := c.Param("id")

		req := &pb.DeletePostRequest{Id: id}
		res, err := client.DeletePost(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"Id:":           id,
			"Post deleted:": res.Success,
		})

	})

	r.PUT("/update/post/", func(c *gin.Context) {
		post := pb.Post{}

		if err := c.ShouldBindJSON(&post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Missing values"})
			return
		}

		req := &pb.UpdatePostRequest{
			Post: &post,
		}
		res, err := client.UpdatePost(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, res.Post)
	})

	r.PUT("/post/upvote", func(c *gin.Context) {
		post := pb.Post{}
		if err := c.ShouldBindJSON(&post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Missing values"})
			return
		}
		req := &pb.UpVoteRequest{
			Post: &post,
		}

		res, err := client.UpVote(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, res.Post)
	})

	r.PUT("/post/downvote", func(c *gin.Context) {
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
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, res.Post)
	})

	if err := r.Run(":8052"); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
