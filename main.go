package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/yazu-codes/scanme-files.git/src/api"
	"github.com/yazu-codes/scanme-files.git/src/api/handlers"
	"github.com/yazu-codes/scanme-files.git/src/db"
	"github.com/yazu-codes/scanme-files.git/src/repositories"
	"github.com/yazu-codes/scanme-files.git/src/services"
	"github.com/yazu-codes/scanme-files.git/src/storage"
	"github.com/yazu-codes/scanme-files.git/src/util"
)

func main() {
	ctx := context.Background()

	var config *util.ConfigReader = util.NewConfigReader()
	config.Setup()

	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)
	logger = logger.With(slog.String("component", "image_service"))

	fmt.Println(config.MinIO.AccessKey)
	fmt.Println(config.MinIO.SecretKey)
	fmt.Println("Use SSL:", config.MinIO.UseSSL)

	client, err := minio.New(
		config.MinIO.Endpoint,
		&minio.Options{
			Creds: credentials.NewStaticV4(
				config.MinIO.AccessKey,
				config.MinIO.SecretKey,
				"",
			),
			Secure: config.MinIO.UseSSL,
		},
	)
	if err != nil {
		fmt.Println("PROBLEMMM")
		logger.Error(err.Error())
	}

	fmt.Println("ENDPOINT: ", client.EndpointURL().String())

	storage, err := storage.NewMinIOStorage(logger, client, config.MinIO.Bucket, ctx)
	if err != nil {
		panic(err)
		logger.Error(err.Error())
	}

	fmt.Println("Will be printing endpoint url")
	fmt.Println(storage.Client.EndpointURL().String())
	fmt.Println("Printed endpoint url")

	db := db.NewDatabase(config.Db.DSN(), logger)
	db.Connect()
	db.AutoMigrate()

	server := api.NewServer(config.Server.ConstructUrl(), logger)
	server.SetupDefaultConfig()

	imagesRepo := repositories.NewImages(db, logger)
	imagesService := services.NewImages(imagesRepo, logger, storage)

	imagesHandler := handlers.NewImages(imagesService, logger)

	server.Router.POST("/images", imagesHandler.AddImage)
	server.Router.POST("/images/bulk", imagesHandler.AddImages)
	server.Router.GET("/images", imagesHandler.GetImages)
	server.Router.GET("/images/:menu_id", imagesHandler.GetImagesByMenuId)
	server.Router.DELETE("/images/:id", imagesHandler.DeleteImageById)

	server.Run()

	// https://items.tapmy.menu/menu-items/images/1008e2b8-c943-4314-9282-65acfbbc4abc.jpg
}
