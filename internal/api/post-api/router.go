package postapi

import (
	"PlantSite/internal/api/post-api/mapper"
	"PlantSite/internal/api/post-api/request"
	_ "PlantSite/internal/api/post-api/response"
	"PlantSite/internal/models"
	"PlantSite/internal/models/auth"
	"PlantSite/internal/models/post"
	postservice "PlantSite/internal/services/post-service"
	"errors"
	"net/http"

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

// @Summary Create a new post
// @Description Creates a new post with text content and optional images
// @Tags post
// @Accept mpfd
// @Param title formData string true "Post title"
// @Param content formData string true "Post content"
// @Param tags formData []string false "List of tags"
// @Param files formData []file false "Attached files"
// @Success 200  "Post created successfully"
// @Failure 400  "Bad Request - Invalid input or missing required fields"
// @Failure 401  "Unauthorized - Not authorized to create post"
// @Failure 403  "Forbidden - Does not have author rights to create post"
// @Failure 500 "Internal Server Error - Failed to create post"
// @Router /post/create [post]
func (r *PostRouter) Create(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.CreatePostRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	if req.Tags == nil {
		req.Tags = make([]string, 0)
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	files := make([]models.FileData, 0)
	for _, file := range form.File["files"] {
		f, err := file.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	_, err = r.post.CreatePost(ctx, postservice.CreatePostTextData{
		Title:   req.Title,
		Content: *content,
		Tags:    req.Tags,
	}, files)
	if errors.Is(err, auth.ErrNotAuthorized) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Error(err)
		return
	} else if errors.Is(err, auth.ErrNoAuthorRights) {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		c.Error(err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// Get Post Handler
// @Summary Get post
// @Description Gets a post by ID
// @Tags post
// @Produce json
// @Param id path string true "Post ID"
// @Success 200  {object} response.GetPostResponse "Post fetch successfully"
// @Failure 400  "Bad Request - Invalid input"
// @Failure 401  "Unauthorized - Not authorized to get post"
// @Failure 403  "Forbidden - Does not have author rights to get post"
// @Failure 500 "Internal Server Error - Failed to get post"
// @Router /post/get/{id} [get]
func (r *PostRouter) Get(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapPostGetRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	post, err := r.post.GetPost(ctx, req.ID)
	if errors.Is(err, auth.ErrNotAuthorized) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Error(err)
		return
	} else if errors.Is(err, auth.ErrNoAuthorRights) {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		c.Error(err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	resp := mapper.MapGetPostResponse(post)
	if resp == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "post not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": resp})
}

// Delete Post Handler
// @Summary Delete post
// @Description Deletes a post by ID
// @Tags post
// @Param id path string true "Post ID"
// @Success 200  "Post deleted successfully"
// @Failure 400  "Bad Request - Invalid input"
// @Failure 401  "Unauthorized - Not authorized to delete post"
// @Failure 403  "Forbidden - Does not have author rights to delete post"
// @Failure 500 "Internal Server Error - Failed to delete post"
// @Router /post/delete/{id} [delete]
func (r *PostRouter) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapPostDeleteRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	err = r.post.Delete(ctx, req.ID)
	if errors.Is(err, auth.ErrNotAuthorized) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Error(err)
		return
	} else if errors.Is(err, auth.ErrNoAuthorRights) {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		c.Error(err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// Update Post Handler
// @Summary Update post text data
// @Description Updates the text data of a post
// @Tags post
// @Accept json
// @Param id path string true "Post ID"
// @Param request body mapper.UpdatePostRequestBody true "Update post request body"
// @Success 200  "Post updated successfully"
// @Failure 400  "Bad Request - Invalid input or missing required fields"
// @Failure 401  "Unauthorized - Not authorized to update post"
// @Failure 403  "Forbidden - Does not have author rights to update post"
// @Failure 500 "Internal Server Error - Failed to update post"
// @Router /post/text/{id} [put]
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	_, err = r.post.UpdatePost(ctx, req.ID, postservice.UpdatePostTextData{
		Title:   req.Title,
		Content: *newContent,
		Tags:    req.Tags,
	})
	if errors.Is(err, auth.ErrNotAuthorized) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Error(err)
		return
	} else if errors.Is(err, auth.ErrNoAuthorRights) {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		c.Error(err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
