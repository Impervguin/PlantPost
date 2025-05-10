package view

import (
	postsquery "PlantSite/internal/api-utils/query-filters/posts-query"
	"PlantSite/internal/view/components"
	"PlantSite/internal/view/gintemplrenderer"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *ViewRouter) PostsHandler(c *gin.Context) {
	ctx := c.Request.Context()
	user := r.auth.UserFromContext(ctx)

	postFilters, err := postsquery.ParseGinQueryPostSearch(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	posts, err := r.srch.SearchPosts(ctx, postFilters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, post := range posts {
		for i, _ := range post.Photos {
			post.Photos[i].File.URL = r.postMedia.GetUrl(post.Photos[i].File.URL)
		}
	}

	tags, err := r.srch.PostTags(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	authors, err := r.srch.PostAuthors(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rend := gintemplrenderer.New(c.Request.Context(), http.StatusOK, components.Posts(user, posts, tags, authors))
	c.Render(http.StatusOK, rend)
}

type postView struct {
	ID string `uri:"id" binding:"required"`
}

func (r *ViewRouter) PostViewHandler(c *gin.Context) {
	ctx := c.Request.Context()
	user := r.auth.UserFromContext(ctx)

	var postView postView

	if err := c.ShouldBindUri(&postView); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := uuid.Parse(postView.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, err := r.srch.GetPost(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i, _ := range post.Photos {
		post.Photos[i].File.URL = r.postMedia.GetUrl(post.Photos[i].File.URL)
	}

	rend := gintemplrenderer.New(c.Request.Context(), http.StatusOK, components.PostView(user, post))
	c.Render(http.StatusOK, rend)
}
