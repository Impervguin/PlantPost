package postapi

import (
	"PlantSite/internal/api/post-api/mapper"
	"PlantSite/internal/api/post-api/request"
	"PlantSite/internal/models"
	"PlantSite/internal/models/post"
	postservice "PlantSite/internal/services/post-service"

	"github.com/gin-gonic/gin"
)

type PostRouter struct {
	post *postservice.PostService
}

func (r *PostRouter) Init(router *gin.RouterGroup, post *postservice.PostService) {
	r.post = post
	gr := router.Group("/post")
	gr.POST("/create", r.Create)
	gr.GET("/get/:id", r.Get)
	gr.DELETE("/delete/:id", r.Delete)
	gr.PUT("/text/:id", r.Update)
}

func (r *PostRouter) Create(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.CreatePostRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	if req.Tags == nil {
		req.Tags = make([]string, 0)
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	files := make([]models.FileData, 0)
	for _, file := range form.File["files"] {
		f, err := file.Open()
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		files = append(files, models.FileData{
			Name:        file.Filename,
			ContentType: file.Header.Get("Content-Type"),
			Reader:      f,
		})
	}

	content, err := post.NewContent(req.Content, post.ContentTypePlainText)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	_, err = r.post.CreatePost(ctx, postservice.CreatePostTextData{
		Title:   req.Title,
		Content: *content,
		Tags:    req.Tags,
	}, files)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{})
}

func (r *PostRouter) Get(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapPostGetRequest(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	post, err := r.post.GetPost(ctx, req.ID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	resp := mapper.MapGetPostResponse(post)
	if resp == nil {
		c.JSON(404, gin.H{"error": "post not found"})
		return
	}

	c.JSON(200, gin.H{"post": resp})
}

func (r *PostRouter) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapPostDeleteRequest(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	err = r.post.Delete(ctx, req.ID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{})
}

func (r *PostRouter) Update(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapPostUpdateRequest(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	newContent, err := post.NewContent(req.Content, post.ContentTypePlainText)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	_, err = r.post.UpdatePost(ctx, req.ID, postservice.UpdatePostTextData{
		Title:   req.Title,
		Content: *newContent,
		Tags:    req.Tags,
	})
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{})
}
