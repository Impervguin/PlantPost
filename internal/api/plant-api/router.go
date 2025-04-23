package plantapi

import (
	"PlantSite/internal/api/plant-api/mapper"
	"PlantSite/internal/models"
	plantservice "PlantSite/internal/services/plant-service"
	"fmt"

	"github.com/gin-gonic/gin"
)

type PlantRouter struct {
	plant *plantservice.PlantService
}

func (r *PlantRouter) Init(router *gin.RouterGroup, plantService *plantservice.PlantService) {
	r.plant = plantService
	gr := router.Group("/plant")
	gr.POST("/create", r.Create)
	gr.GET("/get/:id", r.Get)
	gr.PUT("/specification/:id", r.UpdateSpecification)
	gr.DELETE("/delete/:id", r.Delete)
	gr.POST("/upload/:id", r.UploadPhoto)
}

func (r *PlantRouter) Create(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapCreatePlantRequest(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	spec, err := req.Spec.ToDomain()
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("can't convert specification to domain: %w", err).Error()})
		c.Error(err)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	if len(form.File["file"]) == 0 {
		c.JSON(400, gin.H{"error": "no main photo provided"})
		return
	}

	if len(form.File["file"]) > 1 {
		c.JSON(400, gin.H{"error": "only one main photo allowed"})
		return
	}

	file := form.File["file"][0]
	f, err := file.Open()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	fileData := models.FileData{
		Name:        file.Filename,
		ContentType: file.Header.Get("Content-Type"),
		Reader:      f,
	}

	data := plantservice.CreatePlantData{
		Name:        req.Name,
		LatinName:   req.LatinName,
		Description: req.Description,
		Category:    req.Category,
		Spec:        spec,
	}

	err = r.plant.CreatePlant(ctx, data, fileData)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{})
}

func (r *PlantRouter) Get(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapGetPlantRequest(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	pl, err := r.plant.GetPlant(ctx, req.ID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	resp, err := mapper.MapGetPlantResponse(pl)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{"plant": resp})
}

func (r *PlantRouter) UpdateSpecification(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapUpdatePlantSpecRequest(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	spec, err := req.Spec.ToDomain()
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("can't convert specification to domain: %w", err).Error()})
		c.Error(err)
		return
	}

	err = r.plant.UpdatePlantSpec(ctx, req.ID, spec)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{})
}

func (r *PlantRouter) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapDeletePlantRequest(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	err = r.plant.DeletePlant(ctx, req.ID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{})
}

func (r *PlantRouter) UploadPhoto(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapUploadPlantPhotoRequest(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	if len(form.File["file"]) == 0 {
		c.JSON(400, gin.H{"error": "no main photo provided"})
		return
	}

	if len(form.File["file"]) > 1 {
		c.JSON(400, gin.H{"error": "only one main photo allowed"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	data := models.FileData{
		Name:        file.Filename,
		ContentType: file.Header.Get("Content-Type"),
		Reader:      f,
	}

	err = r.plant.UploadPlantPhoto(ctx, req.ID, data, req.Description)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{})
}
