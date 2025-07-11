package main

import (
	"PlantSite/internal/api-utils/urllib"
	albumapi "PlantSite/internal/api/album-api"
	authapi "PlantSite/internal/api/auth-api"
	"PlantSite/internal/api/middleware"
	plantapi "PlantSite/internal/api/plant-api"
	postapi "PlantSite/internal/api/post-api"
	searchapi "PlantSite/internal/api/search-api"
	minioclient "PlantSite/internal/infra/minio-client"
	filedir "PlantSite/internal/infra/os/file-dir"
	sessionstorage "PlantSite/internal/infra/session-storage"
	"PlantSite/internal/models"
	authrepo "PlantSite/internal/repositories/authrepo"
	miniofilestorage "PlantSite/internal/repositories/pgminio/file-storage"
	fsfilestorage "PlantSite/internal/repositories/pgos/file-storage"
	albumstorage "PlantSite/internal/repositories/postgres/album-storage"
	authstorage "PlantSite/internal/repositories/postgres/auth-storage"
	plantstorage "PlantSite/internal/repositories/postgres/plant-storage"
	poststorage "PlantSite/internal/repositories/postgres/post-storage"
	searchstorage "PlantSite/internal/repositories/postgres/search-storage"
	albumservice "PlantSite/internal/services/album-service"
	authservice "PlantSite/internal/services/auth-service"
	plantservice "PlantSite/internal/services/plant-service"
	postservice "PlantSite/internal/services/post-service"
	searchservice "PlantSite/internal/services/search-service"
	"PlantSite/internal/utils/bcrypthasher"
	"PlantSite/internal/utils/logs"
	"PlantSite/internal/view"
	"context"
	"fmt"

	docs "PlantSite/cmd/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	fmt.Println(GetPlantMinioConfig())
	ctx := context.Background()
	engine := gin.New()

	docs.SwaggerInfo.BasePath = "/api"
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	engine.Use(middleware.RequestIDMiddleware())
	apiGroup := engine.Group(GetApiUrlPrefix())

	// ------------- LOGGING -------------
	logconf := GetLoggerConfig()

	switch logconf.LogFileType {
	case LogFileTypeJson:
		break
	default:
		panic("unknown log file type")
	}

	var ff logs.LogFileFactory

	switch logconf.FileFactory {
	case EveryDayFileFactory:
		var err error
		ff, err = logs.NewEveryDayFileFactory(logconf.LogDir, ".json")
		if err != nil {
			panic(err)
		}
	default:
		panic("unknown file factory")
	}

	logg, err := logs.InitTwoPlaceLogger(
		&logs.TwoPlaceConfig{
			Type:           logconf.LogType,
			FileLevel:      logconf.LogFileLevel,
			ConsoleLevel:   logconf.LogConsoleLevel,
			LogFileFactory: ff,
		},
	)
	if err != nil {
		panic(fmt.Errorf("failed to init logger: %w", err))
	}
	apiGroup.Use(middleware.LogMiddleware(logg))

	sqpgx := GetSqpgx(context.Background())

	// ------------- MEDIA STORAGES -------------

	mediaType := GetMediaStorage()
	var postFStorage models.FileRepository
	var plantFStorage models.FileRepository

	switch mediaType {
	case MediaStorageFs:
		logg.Info("Choosed filesystem media storage")
		root := GetFsRoot()
		fClient, err := filedir.NewFileClient(root)
		if err != nil {
			panic(err)
		}
		postFStorage = fsfilestorage.NewPgOsFileStorage(GetFsBucket(FSPostBucketPrefix), fClient, sqpgx)
		plantFStorage = fsfilestorage.NewPgOsFileStorage(GetFsBucket(FSPlantBucketPrefix), fClient, sqpgx)
	case MediaStorageMinio:
		logg.Info("Choosed minio media storage")
		postMinioCl, err := minioclient.NewMinioClient(GetPostMinioConfig())
		if err != nil {
			panic(err)
		}
		postFStorage, err = miniofilestorage.NewPgMinioStorage(ctx, sqpgx, postMinioCl)
		if err != nil {
			panic(err)
		}

		plantMinioCl, err := minioclient.NewMinioClient(GetPlantMinioConfig())
		if err != nil {
			panic(err)
		}
		plantFStorage, err = miniofilestorage.NewPgMinioStorage(ctx, sqpgx, plantMinioCl)
		if err != nil {
			panic(err)
		}
	default:
		panic("unknown media storage")
	}

	// ------------- AUTH STORAGE -------------
	sessStorage := sessionstorage.NewMapSessionStorage()
	hasher := bcrypthasher.NewBcryptHasher(GetHashCost())
	authRepo, err := authstorage.NewPostgresAuthRepository(ctx, sqpgx)
	if err != nil {
		panic(err)
	}

	adminsMap := GetAdminsMap(hasher)
	storageWithAdmins := authrepo.NewWithAdminRepository(adminsMap, authRepo)
	logg.Info("admins map initialized")

	// ------------- AUTH -------------
	authservice.UpdateSessionExpireTime(GetSessionExpireTime())
	authService := authservice.NewAuthService(sessStorage, storageWithAdmins, hasher)

	apiGroup.Use(middleware.AuthMiddleware(authService))

	authRouter := authapi.AuthRouter{}
	authRouter.Init(apiGroup, authService)

	// ------------- SEARCH STORAGE -------------
	searchRepo, err := searchstorage.NewPostgresSearchRepository(ctx, sqpgx)
	if err != nil {
		panic(err)
	}

	plantGetter := searchstorage.NewSearchPlantGetter(searchRepo)

	postRepo, err := poststorage.NewPostgresPostRepository(ctx, sqpgx, plantGetter)
	if err != nil {
		panic(err)
	}

	// ------------- PLANT STORAGE -------------
	plantRepo, err := plantstorage.NewPostgresPlantRepository(ctx, sqpgx)
	if err != nil {
		panic(err)
	}

	plantCategoryRepo, err := plantstorage.NewPostgresPlantCategoryRepository(sqpgx)
	if err != nil {
		panic(err)
	}

	// ------------- PLANTS -------------
	plantService := plantservice.NewPlantService(plantRepo, plantCategoryRepo, plantFStorage, authService)

	plantRouter := plantapi.PlantRouter{}
	plantRouter.Init(apiGroup, plantService)

	// ------------- SEARCH -------------
	searchService := searchservice.NewSearchService(searchRepo, plantFStorage, postFStorage)

	searchRouter := searchapi.SearchRouter{}
	searchRouter.Init(apiGroup, searchService)

	// ------------- POSTS -------------
	postservice := postservice.NewPostService(postRepo, postFStorage, authService)

	postRouter := postapi.PostRouter{}
	postRouter.Init(apiGroup, postservice)

	// ------------- ALBUM STORAGE -------------
	albumRepo, err := albumstorage.NewPostgresAlbumRepository(ctx, sqpgx)
	if err != nil {
		panic(err)
	}

	// ------------- ALBUMS -------------
	albumService := albumservice.NewAlbumService(albumRepo, authService)

	albumRouter := albumapi.AlbumRouter{}
	albumRouter.Init(apiGroup, albumService)

	// ------------- VIEW -------------
	viewRouter := view.ViewRouter{}
	viewGroup := engine.Group("")
	viewGroup.Use(middleware.RequestIDMiddleware())
	viewGroup.Use(middleware.LogMiddleware(logg))
	viewGroup.Use(middleware.AuthMiddleware(authService))

	mediaStrategy := &urllib.StaticUrlStrategy{BaseUrl: GetMediaPath()}

	viewRouter.Init(viewGroup, GetStaticPath(), authService, searchService, albumService, plantGetter, mediaStrategy, mediaStrategy)

	engine.Run(fmt.Sprintf(":%d", GetApiPort()))
}
