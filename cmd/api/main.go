package main

import (
	authapi "PlantSite/internal/api/auth-api"
	"PlantSite/internal/api/middleware"
	postapi "PlantSite/internal/api/post-api"
	minioclient "PlantSite/internal/infra/minio-client"
	sessionstorage "PlantSite/internal/infra/session-storage"
	authrepo "PlantSite/internal/repositories/authrepo"
	filestorage "PlantSite/internal/repositories/pgminio/file-storage"
	authstorage "PlantSite/internal/repositories/postgres/auth-storage"
	poststorage "PlantSite/internal/repositories/postgres/post-storage"
	authservice "PlantSite/internal/services/auth-service"
	postservice "PlantSite/internal/services/post-service"
	"PlantSite/internal/utils/bcrypthasher"
	"PlantSite/internal/utils/logs"
	"context"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()
	engine := gin.New()
	engine.Use(middleware.RequestIDMiddleware())
	apiGroup := engine.Group(GetApiUrlPrefix())

	// ------------- LOGGING -------------
	logg := logs.InitTwoPlaceLogger(
		&logs.TwoPlaceConfig{
			Type:           logs.TypeDev,
			FileLevel:      logs.LevelInfo,
			ConsoleLevel:   logs.LevelDebug,
			LogFileFactory: logs.NewEveryDayFileFactory("/logs/", ".json"),
		},
	)
	apiGroup.Use(middleware.LogMiddleware(logg))

	// ------------- AUTH STORAGE -------------
	sessStorage := sessionstorage.NewMapSessionStorage()
	sqpgx := GetSqpgx(context.Background())
	hasher := bcrypthasher.NewBcryptHasher(GetHashCost())
	authRepo, err := authstorage.NewPostgresAuthRepository(ctx, sqpgx)
	if err != nil {
		panic(err)
	}

	adminsMap := GetAdminsMap(hasher)
	storageWithAdmins := authrepo.NewWithAdminRepository(adminsMap, authRepo)
	logg.Info("admins map initialized")

	// ------------- AUTH -------------
	authService := authservice.NewAuthService(sessStorage, storageWithAdmins, hasher)

	apiGroup.Use(middleware.AuthMiddleware(authService))

	authRouter := authapi.AuthRouter{}
	authRouter.Init(apiGroup, authService)

	// ------------- POST STORAGE -------------
	postMinioCl, err := minioclient.NewMinioClient(GetPostMinioConfig())
	if err != nil {
		panic(err)
	}

	postFStorage, err := filestorage.NewPgMinioStorage(ctx, sqpgx, postMinioCl)
	if err != nil {
		panic(err)
	}

	postRepo, err := poststorage.NewPostgresPostRepository(ctx, sqpgx)
	if err != nil {
		panic(err)
	}

	// ------------- POSTS -------------
	postservice := postservice.NewPostService(postRepo, postFStorage, authService)

	postRouter := postapi.PostRouter{}
	postRouter.Init(apiGroup, postservice)

	engine.Run(":8080")
}
