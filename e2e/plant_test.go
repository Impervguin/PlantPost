package e2e_test

import (
	"PlantSite/e2e/containers/appcnt"
	"PlantSite/e2e/containers/miniocnt"
	pgcnt "PlantSite/e2e/containers/pgcnt"
	"PlantSite/internal/models/plant"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
)

type PlantTestSuite struct {
	suite.Suite
	appCnt      testcontainers.Container
	pgcnt       testcontainers.Container
	miniocnt    testcontainers.Container
	appConfig   *appcnt.AppConfig
	pgConfig    *pgcnt.PostgresConfig
	minioConfig *miniocnt.MinioConfig
	network     *testcontainers.DockerNetwork
}

const configPath = "./config"

func (s *PlantTestSuite) setupAppConfig() {
	s.appConfig.DbHost = &s.pgConfig.Host
	s.appConfig.DbPort = &s.pgConfig.Port
	s.appConfig.DbName = &s.pgConfig.Database
	s.appConfig.DbUser = &s.pgConfig.User
	s.appConfig.DbPassword = &s.pgConfig.Password

	s.appConfig.MinioHost = &s.minioConfig.Host
	s.appConfig.MinioPort = &s.minioConfig.Port
	s.appConfig.MinioUser = &s.minioConfig.User
	s.appConfig.MinioPassword = &s.minioConfig.Password
	s.appConfig.MinioPostBucket = &s.minioConfig.Buckets[0]
	s.appConfig.MinioPlantBucket = &s.minioConfig.Buckets[1]
}

func (s *PlantTestSuite) BeforeAll(t provider.T) {
	ctx := context.Background()
	var err error

	s.appConfig, err = appcnt.GetConfig(configPath)
	require.NoError(t, err)

	if s.appConfig.External {
		s.pgConfig, err = pgcnt.GetConfig(configPath)
		require.NoError(t, err)

		s.minioConfig, err = miniocnt.GetConfig(configPath)
		require.NoError(t, err)

		s.pgConfig.OuterHost = &s.pgConfig.Host
		s.pgConfig.OuterPort = &s.pgConfig.Port

		s.minioConfig.OuterHost = &s.minioConfig.Host
		s.minioConfig.OuterPort = &s.minioConfig.Port

		s.appConfig.OuterHost = &s.appConfig.Host
		s.appConfig.OuterPort = &s.appConfig.Port

		miniocnt.Migrate(context.Background(), s.minioConfig)
	} else {
		s.network, err = network.New(ctx)
		require.NoError(t, err)

		// Create test containers
		s.pgcnt, s.pgConfig, err = pgcnt.NewTestPostgres(ctx, configPath, s.network.Name)
		require.NoError(t, err)

		s.miniocnt, s.minioConfig, err = miniocnt.NewTestMinio(ctx, configPath, s.network.Name)
		require.NoError(t, err)

		err := pgcnt.Migrate(context.Background(), s.pgConfig)
		if err != nil {
			panic(err)
		}
		err = miniocnt.Migrate(context.Background(), s.minioConfig)
		if err != nil {
			panic(err)
		}
	}
	s.setupAppConfig()
}

func (s *PlantTestSuite) BeforeEach(t provider.T) {
	t.Epic("Plant E2E")
	t.Feature("Plant")
	var err error
	if !s.appConfig.External {
		if s.appCnt == nil {
			s.appCnt, s.appConfig, err = appcnt.NewTestApp(s.appConfig, s.network.Name)
			require.NoError(t, err)
		} else {
			s.appCnt.Start(context.Background())
		}
	}
}

func (s *PlantTestSuite) AfterEach(t provider.T) {
	// err := pgcnt.TruncateTables(context.Background(), s.pgConfig)
	// require.NoError(t, err)
	// err = miniocnt.CleanUpBucket(context.Background(), s.minioConfig)
	// require.NoError(t, err)
}

func (s *PlantTestSuite) AfterAll(t provider.T) {
	ctx := context.Background()
	if !s.appConfig.External {
		s.network.Remove(ctx)
	}
	if s.pgcnt != nil {
		s.pgcnt.Terminate(ctx)
	}
	if s.miniocnt != nil {
		s.miniocnt.Terminate(ctx)
	}
	if s.appCnt != nil {
		s.appCnt.Terminate(ctx)
	}
}

func TestPlantSuite(t *testing.T) {
	suite.RunSuite(t, new(PlantTestSuite))
}

type ConiferousSpecificationDataBuilder struct {
	HeightM         float64
	DiameterM       float64
	Flowering       plant.FloweringPeriod
	SoilAcidity     plant.SoilAcidity
	SoilMoisture    plant.SoilMoisture
	Light           plant.LightRelation
	SoilType        plant.Soil
	WinterHardiness plant.WinterHardiness
}

func (b *ConiferousSpecificationDataBuilder) Build() map[string]interface{} {
	return map[string]interface{}{
		"height_m":         b.HeightM,
		"diameter_m":       b.DiameterM,
		"flowering":        b.Flowering,
		"soil_acidity":     b.SoilAcidity,
		"soil_moisture":    b.SoilMoisture,
		"light_relation":   b.Light,
		"soil_type":        b.SoilType,
		"winter_hardiness": b.WinterHardiness,
	}
}

func NewConiferousSpecificationDataBuilder() *ConiferousSpecificationDataBuilder {
	return &ConiferousSpecificationDataBuilder{
		HeightM:         1.5,
		DiameterM:       0.5,
		Flowering:       plant.Spring,
		SoilAcidity:     10,
		SoilMoisture:    plant.MediumMoisture,
		Light:           plant.HalfShadow,
		SoilType:        plant.MediumSoil,
		WinterHardiness: plant.WinterHardiness(10),
	}
}

func (b *ConiferousSpecificationDataBuilder) WithHeightM(heightM float64) *ConiferousSpecificationDataBuilder {
	b.HeightM = heightM
	return b
}

func (b *ConiferousSpecificationDataBuilder) WithDiameterM(diameterM float64) *ConiferousSpecificationDataBuilder {
	b.DiameterM = diameterM
	return b
}

func (b *ConiferousSpecificationDataBuilder) WithFlowering(flowering plant.FloweringPeriod) *ConiferousSpecificationDataBuilder {
	b.Flowering = flowering
	return b
}

func (b *ConiferousSpecificationDataBuilder) WithSoilAcidity(soilAcidity plant.SoilAcidity) *ConiferousSpecificationDataBuilder {
	b.SoilAcidity = soilAcidity
	return b
}

func (b *ConiferousSpecificationDataBuilder) WithSoilMoisture(soilMoisture plant.SoilMoisture) *ConiferousSpecificationDataBuilder {
	b.SoilMoisture = soilMoisture
	return b
}

func (b *ConiferousSpecificationDataBuilder) WithLight(light plant.LightRelation) *ConiferousSpecificationDataBuilder {
	b.Light = light
	return b
}

func (b *ConiferousSpecificationDataBuilder) WithSoilType(soilType plant.Soil) *ConiferousSpecificationDataBuilder {
	b.SoilType = soilType
	return b
}

func (b *ConiferousSpecificationDataBuilder) WithWinterHardiness(winterHardiness plant.WinterHardiness) *ConiferousSpecificationDataBuilder { //nolint:golint
	b.WinterHardiness = winterHardiness
	return b
}

type DeciduousSpecificationDataBuilder struct {
	HeightM         float64
	DiameterM       float64
	Flowering       plant.FloweringPeriod
	SoilAcidity     plant.SoilAcidity
	SoilMoisture    plant.SoilMoisture
	Light           plant.LightRelation
	SoilType        plant.Soil
	WinterHardiness plant.WinterHardiness
	FloweringPeriod plant.FloweringPeriod
}

func (b *DeciduousSpecificationDataBuilder) Build() map[string]interface{} {
	return map[string]interface{}{
		"height_m":         b.HeightM,
		"diameter_m":       b.DiameterM,
		"flowering":        b.Flowering,
		"soil_acidity":     b.SoilAcidity,
		"soil_moisture":    b.SoilMoisture,
		"light_relation":   b.Light,
		"soil_type":        b.SoilType,
		"winter_hardiness": b.WinterHardiness,
		"flowering_period": b.FloweringPeriod,
	}
}

func NewDeciduousSpecificationDataBuilder() *DeciduousSpecificationDataBuilder {
	return &DeciduousSpecificationDataBuilder{
		HeightM:         1.5,
		DiameterM:       0.5,
		Flowering:       plant.Spring,
		SoilAcidity:     10,
		SoilMoisture:    plant.MediumMoisture,
		Light:           plant.HalfShadow,
		SoilType:        plant.MediumSoil,
		WinterHardiness: plant.WinterHardiness(10),
		FloweringPeriod: plant.Spring,
	}
}

func (b *DeciduousSpecificationDataBuilder) WithHeightM(heightM float64) *DeciduousSpecificationDataBuilder {
	b.HeightM = heightM
	return b
}

func (b *DeciduousSpecificationDataBuilder) WithDiameterM(diameterM float64) *DeciduousSpecificationDataBuilder {
	b.DiameterM = diameterM
	return b
}

func (b *DeciduousSpecificationDataBuilder) WithFlowering(flowering plant.FloweringPeriod) *DeciduousSpecificationDataBuilder {
	b.Flowering = flowering
	return b
}

func (b *DeciduousSpecificationDataBuilder) WithSoilAcidity(soilAcidity plant.SoilAcidity) *DeciduousSpecificationDataBuilder {
	b.SoilAcidity = soilAcidity
	return b
}

func (b *DeciduousSpecificationDataBuilder) WithSoilMoisture(soilMoisture plant.SoilMoisture) *DeciduousSpecificationDataBuilder {
	b.SoilMoisture = soilMoisture
	return b
}

func (b *DeciduousSpecificationDataBuilder) WithLight(light plant.LightRelation) *DeciduousSpecificationDataBuilder {
	b.Light = light
	return b
}

func (b *DeciduousSpecificationDataBuilder) WithSoilType(soilType plant.Soil) *DeciduousSpecificationDataBuilder {
	b.SoilType = soilType
	return b
}

func (b *DeciduousSpecificationDataBuilder) WithWinterHardiness(winterHardiness plant.WinterHardiness) *DeciduousSpecificationDataBuilder {
	b.WinterHardiness = winterHardiness
	return b
}

func (b *DeciduousSpecificationDataBuilder) WithFloweringPeriod(floweringPeriod plant.FloweringPeriod) *DeciduousSpecificationDataBuilder {
	b.FloweringPeriod = floweringPeriod
	return b
}

type PlantSpecificationDataBuilder interface {
	Build() map[string]interface{}
}

type PlantRequestDataBuilder struct {
	Name          string
	LatinName     string
	Description   string
	MainPhotoName string
	MainPhoto     []byte
	Category      string
	Specification PlantSpecificationDataBuilder
}

type PlantRequest struct {
	buf         *bytes.Buffer
	contentType string
}

func (b *PlantRequestDataBuilder) Build() *PlantRequest {
	buf := &bytes.Buffer{}

	writer := multipart.NewWriter(buf)
	defer writer.Close()

	part, err := writer.CreateFormFile("file", b.MainPhotoName)
	if err != nil {
		panic(err)
	}
	_, err = part.Write(b.MainPhoto)
	if err != nil {
		panic(err)
	}
	err = writer.WriteField("name", b.Name)
	if err != nil {
		panic(err)
	}
	err = writer.WriteField("latin_name", b.LatinName)
	if err != nil {
		panic(err)
	}
	err = writer.WriteField("description", b.Description)
	if err != nil {
		panic(err)
	}
	err = writer.WriteField("category", b.Category)
	if err != nil {
		panic(err)
	}

	jsonBytes, err := json.Marshal(b.Specification.Build())
	if err != nil {
		panic(err)
	}

	writer.WriteField("specification", string(jsonBytes))

	return &PlantRequest{
		buf:         buf,
		contentType: writer.FormDataContentType(),
	}
}

func NewPlantRequestDataBuilder() *PlantRequestDataBuilder {
	uuu := uuid.New()
	return &PlantRequestDataBuilder{
		Name:          uuu.String()[:24],
		LatinName:     uuu.String() + " Plantus",
		Description:   uuu.String() + " Plantus desc",
		MainPhotoName: uuu.String() + "_photo.jpg",
		MainPhoto:     []byte("test photo content"),
		Category:      "coniferous",
		Specification: NewConiferousSpecificationDataBuilder(),
	}
}

func (b *PlantRequestDataBuilder) WithName(name string) *PlantRequestDataBuilder {
	b.Name = name
	return b
}

func (b *PlantRequestDataBuilder) WithLatinName(latinName string) *PlantRequestDataBuilder {
	b.LatinName = latinName
	return b
}

func (b *PlantRequestDataBuilder) WithDescription(description string) *PlantRequestDataBuilder {
	b.Description = description
	return b
}

func (b *PlantRequestDataBuilder) WithMainPhoto(mainPhoto []byte) *PlantRequestDataBuilder {
	b.MainPhoto = mainPhoto
	return b
}

func (b *PlantRequestDataBuilder) WithCategory(category string) *PlantRequestDataBuilder {
	b.Category = category
	return b
}

func (b *PlantRequestDataBuilder) WithSpecification(specification PlantSpecificationDataBuilder) *PlantRequestDataBuilder {
	b.Specification = specification
	return b
}

func (s *PlantTestSuite) TestE2EPlant(t provider.T) {
	t.Tags("e2e", "plant")
	t.Description("Test plant creation and retrieval functionality")

	addr := fmt.Sprintf("http://%s:%d/api", *s.appConfig.OuterHost, *s.appConfig.OuterPort)

	t.Run("Complex searching through plants and creating albums", func(t provider.T) {
		jar, err := cookiejar.New(nil)
		require.NoError(t, err)
		client := &http.Client{
			Jar: jar,
		}

		t.WithNewStep("Login as admin", func(pctx provider.StepCtx) {
			resp, err := client.Post(addr+"/auth/login", "application/json", strings.NewReader(fmt.Sprintf(`{"username":"%s","password":"%s"},`, s.appConfig.AdminUser, s.appConfig.AdminPassword)))
			require.Equal(t, http.StatusOK, resp.StatusCode)
			require.NoError(t, err)
		})
		plantUUID := uuid.New()
		kustUUID := uuid.New()

		t.WithNewStep("Fill storage with plants", func(pctx provider.StepCtx) {
			plantsData := []*PlantRequest{
				NewPlantRequestDataBuilder().
					WithName("Test " + plantUUID.String() + "1").
					WithLatinName("Testus Plantus 1").
					WithDescription("Testus plantus description 1").
					WithMainPhoto([]byte("test photo content 1")).
					WithCategory("coniferous").
					WithSpecification(NewConiferousSpecificationDataBuilder().
						WithHeightM(1.5).
						WithDiameterM(0.5).
						WithFlowering(plant.Spring).
						WithSoilAcidity(10).
						WithSoilMoisture(plant.MediumMoisture).
						WithLight(plant.HalfShadow).
						WithSoilType(plant.MediumSoil).
						WithWinterHardiness(plant.WinterHardiness(10))).
					Build(),
				NewPlantRequestDataBuilder().
					WithName("Test " + plantUUID.String() + "2").
					WithLatinName("Testus Plantus 2").
					WithDescription("Testus plantus description 2").
					WithMainPhoto([]byte("test photo content 2")).
					WithCategory("coniferous").
					WithSpecification(NewConiferousSpecificationDataBuilder().
						WithHeightM(4).
						WithDiameterM(0.5).
						WithFlowering(plant.Autumn).
						WithSoilAcidity(10).
						WithSoilMoisture(plant.HighMoisture).
						WithLight(plant.Shadow).
						WithSoilType(plant.MediumSoil).
						WithWinterHardiness(plant.WinterHardiness(10))).
					Build(),
				NewPlantRequestDataBuilder().
					WithName("Test " + plantUUID.String() + "3").
					WithLatinName("Testus Plantus 3").
					WithDescription("Testus plantus description 3").
					WithMainPhoto([]byte("test photo content 3")).
					WithCategory("deciduous").
					WithSpecification(NewDeciduousSpecificationDataBuilder().
						WithHeightM(25).
						WithDiameterM(1).
						WithFlowering(plant.Spring).
						WithSoilAcidity(10).
						WithSoilMoisture(plant.MediumMoisture).
						WithLight(plant.HalfShadow).
						WithSoilType(plant.MediumSoil).
						WithWinterHardiness(plant.WinterHardiness(10)).
						WithFloweringPeriod(plant.Spring)).
					Build(),
				NewPlantRequestDataBuilder().
					WithName("Test " + kustUUID.String() + "1").
					WithLatinName("Testus Kustus 1").
					WithDescription("Testus kustus description 1").
					WithMainPhoto([]byte("test kust photo content 1")).
					WithCategory("coniferous").
					WithSpecification(NewConiferousSpecificationDataBuilder().
						WithHeightM(1.5).
						WithDiameterM(15).
						WithFlowering(plant.Spring).
						WithSoilAcidity(10).
						WithSoilMoisture(plant.MediumMoisture).
						WithLight(plant.HalfShadow).
						WithSoilType(plant.MediumSoil).
						WithWinterHardiness(plant.WinterHardiness(10))).
					Build(),
				NewPlantRequestDataBuilder().
					WithName("Test " + kustUUID.String() + "2").
					WithLatinName("Testus Kustus 2").
					WithDescription("Testus kustus description 2").
					WithMainPhoto([]byte("test kust photo content 2")).
					WithCategory("coniferous").
					WithSpecification(NewConiferousSpecificationDataBuilder().
						WithHeightM(4).
						WithDiameterM(15).
						WithFlowering(plant.Autumn).
						WithSoilAcidity(10).
						WithSoilMoisture(plant.HighMoisture).
						WithLight(plant.Shadow).
						WithSoilType(plant.MediumSoil).
						WithWinterHardiness(plant.WinterHardiness(10))).
					Build(),
				NewPlantRequestDataBuilder().
					WithName("Test " + kustUUID.String() + "3").
					WithLatinName("Testus Kustus 3").
					WithDescription("Testus kustus description 3").
					WithMainPhoto([]byte("test kust photo content 3")).
					WithCategory("deciduous").
					WithSpecification(NewDeciduousSpecificationDataBuilder().
						WithHeightM(25).
						WithDiameterM(7.2).
						WithFlowering(plant.Spring).
						WithSoilAcidity(20).
						WithSoilMoisture(plant.MediumMoisture).
						WithLight(plant.Light).
						WithSoilType(plant.MediumSoil).
						WithWinterHardiness(plant.WinterHardiness(10)).
						WithFloweringPeriod(plant.Summer)).
					Build(),
			}
			for _, plantData := range plantsData {
				resp, err := client.Post(addr+"/plant/create", plantData.contentType, plantData.buf)
				require.Equal(t, http.StatusOK, resp.StatusCode)
				require.NoError(t, err)
			}
		})

		t.WithNewStep("Logout", func(pctx provider.StepCtx) {
			resp, err := client.Post(addr+"/auth/logout", "application/json", strings.NewReader(""))
			require.Equal(t, http.StatusOK, resp.StatusCode)
			require.NoError(t, err)
		})

		t.WithNewStep("Search for decidious kusts", func(pctx provider.StepCtx) {
			expectedNames := []string{"Test " + kustUUID.String() + "1", "Test " + kustUUID.String() + "2", "Test " + kustUUID.String() + "3"}
			resp, err := client.Get(addr + "/search/plants?name=" + kustUUID.String())
			require.Equal(t, http.StatusOK, resp.StatusCode)
			require.NoError(t, err)
			var searchResp map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&searchResp)
			require.NoError(t, err)
			require.Equal(t, 3, len(searchResp["plants"].([]interface{})))
			for _, plant := range searchResp["plants"].([]interface{}) {
				plantMap := plant.(map[string]interface{})
				require.Contains(t, expectedNames, plantMap["name"].(string))
			}
		})

		t.WithNewStep("Search for plants with summer flowering period", func(pctx provider.StepCtx) {
			expectedNames := []string{"Test " + kustUUID.String() + "3"}
			resp, err := client.Get(addr + "/search/plants?flowering_period=summer&name=" + kustUUID.String())
			require.Equal(t, http.StatusOK, resp.StatusCode)
			require.NoError(t, err)
			var searchResp map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&searchResp)
			require.NoError(t, err)
			require.Equal(t, 1, len(searchResp["plants"].([]interface{})))
			for _, plant := range searchResp["plants"].([]interface{}) {
				plantMap := plant.(map[string]interface{})
				require.Contains(t, expectedNames, plantMap["name"].(string))
			}
		})

		t.WithNewStep("Search for plants with height between 1.5 and 5 and diameter between 0.5 and 1.5", func(pctx provider.StepCtx) {
			expectedNames := []string{"Test " + plantUUID.String() + "1", "Test " + plantUUID.String() + "2"}
			resp, err := client.Get(addr + "/search/plants?height=1.5-5&diameter=0.5-1.5&name=" + plantUUID.String())
			require.Equal(t, http.StatusOK, resp.StatusCode)
			require.NoError(t, err)
			var searchResp map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&searchResp)
			require.NoError(t, err)
			require.Equal(t, 2, len(searchResp["plants"].([]interface{})))
			for _, plant := range searchResp["plants"].([]interface{}) {
				plantMap := plant.(map[string]interface{})
				require.Contains(t, expectedNames, plantMap["name"].(string))
			}
		})
		testUserUUID := uuid.New()

		t.WithNewStep("Register user", func(pctx provider.StepCtx) {
			resp, err := client.Post(addr+"/auth/register", "application/json", strings.NewReader(fmt.Sprintf(`{"username":"%s","password":"%s", "email":"%s"},`, testUserUUID.String()[:24], testUserUUID.String(), testUserUUID.String()+"@test.com")))
			require.Equal(t, http.StatusOK, resp.StatusCode)
			require.NoError(t, err)
		})

		t.WithNewStep("Login as test user", func(pctx provider.StepCtx) {
			resp, err := client.Post(addr+"/auth/login", "application/json", strings.NewReader(fmt.Sprintf(`{"username":"%s","password":"%s"},`, testUserUUID.String()[:24], testUserUUID.String())))
			require.Equal(t, http.StatusOK, resp.StatusCode)
			require.NoError(t, err)
		})

		var KustIds []string
		t.WithNewStep("Get Kust IDs", func(pctx provider.StepCtx) {
			resp, err := client.Get(addr + "/search/plants?name=" + kustUUID.String())
			require.Equal(t, http.StatusOK, resp.StatusCode)
			require.NoError(t, err)
			var searchResp map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&searchResp)
			require.NoError(t, err)
			require.Equal(t, 3, len(searchResp["plants"].([]interface{})))
			for _, plant := range searchResp["plants"].([]interface{}) {
				plantMap := plant.(map[string]interface{})
				KustIds = append(KustIds, plantMap["id"].(string))
			}
		})

		testAlbumUUID := uuid.New()

		t.WithNewStep("Create album", func(pctx provider.StepCtx) {
			var body map[string]interface{} = map[string]interface{}{
				"name":        "Test " + testAlbumUUID.String(),
				"description": "Test album description",
				"plant_ids":   KustIds,
			}
			jsonBytes, err := json.Marshal(body)
			require.NoError(t, err)
			resp, err := client.Post(addr+"/album/create", "application/json", strings.NewReader(string(jsonBytes)))
			require.Equal(t, http.StatusOK, resp.StatusCode)
			require.NoError(t, err)
		})

		t.WithNewStep("List albums", func(pctx provider.StepCtx) {
			resp, err := client.Get(addr + "/album/list")
			require.Equal(t, http.StatusOK, resp.StatusCode)
			require.NoError(t, err)
			var albumsResp map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&albumsResp)
			require.NoError(t, err)
		})

		t.WithNewStep("Logout", func(pctx provider.StepCtx) {
			resp, err := client.Post(addr+"/auth/logout", "application/json", strings.NewReader(""))
			require.Equal(t, http.StatusOK, resp.StatusCode)
			require.NoError(t, err)
		})

		t.WithNewStep("Try List albums without authorization", func(pctx provider.StepCtx) {
			resp, err := client.Get(addr + "/album/list")
			require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
			require.NoError(t, err)
		})
	})
}
