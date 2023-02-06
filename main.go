package main

import (
	"belajar-golang-api/app"
	"belajar-golang-api/controller"
	"belajar-golang-api/helper"
	"belajar-golang-api/middleware"
	"belajar-golang-api/repository"
	"belajar-golang-api/service"
	"net/http"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := app.NewDB()
	validate := validator.New()
	categoryRepository := repository.NewCategoryRepository()
	categoryService := service.NewCategoryService(categoryRepository, db, validate)
	categoryController := controller.NewCategoryController(categoryService)
	// router := httprouter.New()

	// router.GET("/api/categories", categoryController.FindAll)
	// router.GET("/api/categories/:categoryId", categoryController.FindById)
	// router.POST("/api/categories", categoryController.Create)
	// router.PUT("/api/categories/:categoryId", categoryController.Update)
	// router.DELETE("/api/categories/:categoryId", categoryController.Delete)

	//untuk menghandler error agar program tidak berhenti

	// router.PanicHandler = exception.ErrorHandler

	router := app.NewRouter(categoryController)
	server := http.Server{
		Addr:    "localhost:3000",
		Handler: middleware.NewAuthMiddlewate(router),
	}
	err := server.ListenAndServe()
	helper.PanicIfError(err)

}
