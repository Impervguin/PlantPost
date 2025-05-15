package plantapi

import (
	"PlantSite/internal/api/plant-api/mapper"
	_ "PlantSite/internal/api/plant-api/request"
	_ "PlantSite/internal/api/plant-api/response"
	_ "PlantSite/internal/api/plant-api/spec"
	"PlantSite/internal/models"
	"PlantSite/internal/models/auth"
	plantservice "PlantSite/internal/services/plant-service"

	"errors"
	"fmt"
	"net/http"

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

// Create plant handler
// @Summary Create plant
// @Description Creates a new plant with the provided name, latin name, description, category and specification
// @Tags plant
// @Accept mpfd
// @Param name formData string true "plant name"
// @Param latin_name formData string true "plant latin name"
// @Param description formData string true "plant description"
// @Param category formData string true "plant category"
// @Param file formData file true "plant main image"
// @Param specification body spec.UnionSpecification false "plant specification"
// @Param specification formData string true "plant specification"
// @Success 200  "Plant created successfully"
// @Failure 400  "Bad Request - Invalid input or missing required fields"
// @Failure 401  "Unauthorized - Not authorized to create plant"
// @Failure 403  "Forbidden - Does not have author rights to create plant"
// @Failure 500 "Internal Server Error - Failed to create plant"
// @Router /plant/create [post]
func (r *PlantRouter) Create(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapCreatePlantRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	spec, err := req.Spec.ToDomain()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("can't convert specification to domain: %w", err).Error()})
		c.Error(err)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	if len(form.File["file"]) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no main photo provided"})
		return
	}

	if len(form.File["file"]) > 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only one main photo allowed"})
		return
	}

	file := form.File["file"][0]
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

// @Summary Get plant
// @Description Gets a plant by ID
// @Tags plant
// @Produce json
// @Param id path string true "Plant ID"
// @Success 200  {object} response.GetPlantResponse "Plant fetch successfully"
// @Failure 400  "Bad Request - Invalid input"
// @Failure 401  "Unauthorized - Not authorized to get plant"
// @Failure 403  "Forbidden - Does not have author rights to get plant"
// @Failure 500 "Internal Server Error - Failed to get plant"
// @Router /plant/get/{id} [get]
func (r *PlantRouter) Get(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapGetPlantRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	pl, err := r.plant.GetPlant(ctx, req.ID)
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
	resp, err := mapper.MapGetPlantResponse(pl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"plant": resp})
}

// @Summary Update plant specification
// @Description Updates the specification of a plant
// @Tags plant
// @Accept json
// @Param id path string true "Plant ID"
// @Param specification body spec.UnionSpecification false "plant specification"
// @Param specification formData string true "plant specification"
// @Success 200  "Plant specification updated successfully"
// @Failure 400  "Bad Request - Invalid input or missing required fields"
// @Failure 401  "Unauthorized - Not authorized to update plant specification"
// @Failure 403  "Forbidden - Does not have author rights to update plant specification"
// @Failure 500 "Internal Server Error - Failed to update plant specification"
// @Router /plant/specification/{id} [put]
func (r *PlantRouter) UpdateSpecification(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapUpdatePlantSpecRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	

	spec, err := req.Spec.ToDomain()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("can't convert specification to domain: %w", err).Error()})
		c.Error(err)
		return
	}

	err = r.plant.UpdatePlantSpec(ctx, req.ID, spec)
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

// @Summary Delete plant
// @Description Deletes a plant by ID
// @Tags plant
// @Param id path string true "Plant ID"
// @Success 200  "Plant deleted successfully"
// @Failure 400  "Bad Request - Invalid input"
// @Failure 401  "Unauthorized - Not authorized to delete plant"
// @Failure 403  "Forbidden - Does not have author rights to delete plant"
// @Failure 500 "Internal Server Error - Failed to delete plant"// @Param specification body spec.UnionSpecification false "plant specification"
// @Router /plant/delete/{id} [delete]
func (r *PlantRouter) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapDeletePlantRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	err = r.plant.DeletePlant(ctx, req.ID)
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

// @Summary Upload plant photo
// @Description Uploads a plant photo
// @Tags plant
// @Accept mpfd
// @Param id path string true "Plant ID"
// @Param file formData file true "Plant photo"
// @Param description formData string true "Plant photo description"
// @Success 200  "Plant photo uploaded successfully"
// @Failure 400  "Bad Request - Invalid input or missing required fields"
// @Failure 401  "Unauthorized - Not authorized to upload plant photo"
// @Failure 403  "Forbidden - Does not have author rights to upload plant photo"
// @Failure 500 "Internal Server Error - Failed to upload plant photo"
// @Router /plant/upload/{id} [post]
func (r *PlantRouter) UploadPhoto(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapUploadPlantPhotoRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	if len(form.File["file"]) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no main photo provided"})
		return
	}

	if len(form.File["file"]) > 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only one main photo allowed"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	data := models.FileData{
		Name:        file.Filename,
		ContentType: file.Header.Get("Content-Type"),
		Reader:      f,
	}

	err = r.plant.UploadPlantPhoto(ctx, req.ID, data, req.Description)
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
