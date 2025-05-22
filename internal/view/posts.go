package view

import (
	postsquery "PlantSite/internal/api-utils/query-filters/posts-query"
	"PlantSite/internal/models/post"
	"PlantSite/internal/models/post/parser"
	"PlantSite/internal/models/search"
	searchservice "PlantSite/internal/services/search-service"
	"PlantSite/internal/view/components"
	"PlantSite/internal/view/gintemplrenderer"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *ViewRouter) handlePostsWithPlant(c *gin.Context, psts []*searchservice.SearchPost) (map[uuid.UUID]*searchservice.SearchPlant, error) {
	plantMap := make(map[uuid.UUID]*searchservice.SearchPlant)
	plants := make(map[uuid.UUID]struct{})
	for _, pst := range psts {
		if post.CheckContentWithPlant(&pst.Content) {
			parser, err := parser.GetParser(&pst.Content, r.plntGet)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return nil, err
			}
			content, err := post.NewContentWithPlant(pst.Content.Text, pst.Content.ContentType, parser)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return nil, err
			}
			plantIDs := content.PlantIDs()
			for _, plntID := range plantIDs {
				if _, ok := plants[plntID]; !ok {
					plants[plntID] = struct{}{}
				}
			}
		}
	}
	srch := search.NewPlantSearch()
	plantIDs := make([]uuid.UUID, 0, len(plants))
	for plntID := range plants {
		plantIDs = append(plantIDs, plntID)
	}
	srch.AddFilter(search.NewPlantIDsFilter(plantIDs))

	plnts, err := r.srch.SearchPlants(c.Request.Context(), srch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil, err
	}

	for _, plnt := range plnts {
		plantMap[plnt.ID] = plnt
	}
	return plantMap, nil
}

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

	plantMap, err := r.handlePostsWithPlant(c, posts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, plnt := range plantMap {
		plnt.MainPhoto.URL = r.plantMedia.GetUrl(plnt.MainPhoto.URL)
	}

	rend := gintemplrenderer.New(c.Request.Context(), http.StatusOK, components.Posts(user, posts, tags, authors, plantMap))
	c.Render(http.StatusOK, rend)
}

type postView struct {
	ID string `uri:"id" binding:"required"`
}

func (r *ViewRouter) handlePostWithPlant(c *gin.Context, pst *searchservice.GetPost) (map[uuid.UUID]*searchservice.SearchPlant, error) {
	parser, err := parser.GetParser(&pst.Content, r.plntGet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil, err
	}
	content, err := post.NewContentWithPlant(pst.Content.Text, pst.Content.ContentType, parser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil, err
	}
	plantIDs := content.PlantIDs()
	srch := search.NewPlantSearch()
	srch.AddFilter(search.NewPlantIDsFilter(plantIDs))

	plants, err := r.srch.SearchPlants(c.Request.Context(), srch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil, err
	}

	plantMap := make(map[uuid.UUID]*searchservice.SearchPlant)
	for _, plnt := range plants {
		plantMap[plnt.ID] = plnt
	}

	return plantMap, nil
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

	pst, err := r.srch.GetPost(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i, _ := range pst.Photos {
		pst.Photos[i].File.URL = r.postMedia.GetUrl(pst.Photos[i].File.URL)
	}

	plantMap := make(map[uuid.UUID]*searchservice.SearchPlant)
	if post.CheckContentWithPlant(&pst.Content) {
		plantMap, err = r.handlePostWithPlant(c, pst)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	for _, plnt := range plantMap {
		plnt.MainPhoto.URL = r.plantMedia.GetUrl(plnt.MainPhoto.URL)
	}

	rend := gintemplrenderer.New(c.Request.Context(), http.StatusOK, components.PostView(user, pst, plantMap))
	c.Render(http.StatusOK, rend)
}

func (r *ViewRouter) CreatePostHandler(c *gin.Context) {
	ctx := c.Request.Context()
	user := r.auth.UserFromContext(ctx)
	if !user.HasAuthorRights() {
		c.Redirect(http.StatusFound, "/view/plants")
		return
	}

	rend := gintemplrenderer.New(c.Request.Context(), http.StatusOK, components.PostCreate(user))
	c.Render(http.StatusOK, rend)
}

func (r *ViewRouter) UpdatePostHandler(c *gin.Context) {
	ctx := c.Request.Context()
	user := r.auth.UserFromContext(ctx)
	if !user.HasAuthorRights() {
		c.Redirect(http.StatusFound, "/view/plants")
		return
	}

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

	rend := gintemplrenderer.New(c.Request.Context(), http.StatusOK, components.PostUpdate(user, post))
	c.Render(http.StatusOK, rend)
}
