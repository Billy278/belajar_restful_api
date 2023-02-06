package test

import (
	"belajar-golang-api/app"
	"belajar-golang-api/controller"
	"belajar-golang-api/helper"
	"belajar-golang-api/middleware"
	"belajar-golang-api/model/domain"
	"belajar-golang-api/repository"
	"belajar-golang-api/service"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func setupTestDB() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/belajar_golang_restfull_api_test")
	helper.PanicIfError(err)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(20)
	db.SetConnMaxIdleTime(10 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)
	return db

}

func setupRouter(db *sql.DB) http.Handler {
	//db := setupTestDB()
	validate := validator.New()
	categoryRepository := repository.NewCategoryRepository()
	categoryService := service.NewCategoryService(categoryRepository, db, validate)
	categoryController := controller.NewCategoryController(categoryService)
	router := app.NewRouter(categoryController)

	return middleware.NewAuthMiddlewate(router)
}

func truncateCategory(db *sql.DB) {
	db.Exec("TRUNCATE category")

}

func TestCreateCategorySuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)
	router := setupRouter(db)

	requestBody := strings.NewReader(`{"name" : "Gadget"}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/categories", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-key", "RAHASIA")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)
	fmt.Println(responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "Ok", responseBody["status"])
	assert.Equal(t, "Gadget", responseBody["data"].(map[string]interface{})["name"])

}
func TestCreateCategoryFailed(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)
	router := setupRouter(db)

	requestBody := strings.NewReader(`{"name" : ""}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/categories", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-key", "RAHASIA")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)
	//fmt.Println(responseBody)

	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Equal(t, "Bad Request", responseBody["status"])

}
func TestUpdateCategorySuccess(t *testing.T) {
	//ini bisa kita buat
	// db := setupTestDB()
	// //truncateCategory(db)
	// router := setupRouter(db)

	// requestBody := strings.NewReader(`{"name" : "Gadget"}`)
	// request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/categories/1", requestBody)
	// request.Header.Add("Content-Type", "application/json")
	// request.Header.Add("X-API-key", "RAHASIA")
	// recorder := httptest.NewRecorder()
	// router.ServeHTTP(recorder, request)

	// response := recorder.Result()
	// assert.Equal(t, 200, response.StatusCode)

	// body, _ := io.ReadAll(response.Body)
	// var responseBody map[string]interface{}
	// json.Unmarshal(body, &responseBody)
	// fmt.Println(responseBody)

	// assert.Equal(t, 200, int(responseBody["code"].(float64)))
	// assert.Equal(t, "Ok", responseBody["status"])
	// assert.Equal(t, "Gadget", responseBody["data"].(map[string]interface{})["name"])
	//====================================================================
	//tapi dalam tutorial dia menambahkan data dulu baru di update

	db := setupTestDB()
	truncateCategory(db)

	categoryRepository := repository.NewCategoryRepository()
	tx, _ := db.Begin()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Bola",
	})

	tx.Commit()
	router := setupRouter(db)

	requestBody := strings.NewReader(`{"name" : "Bola"}`)
	request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/categories/"+strconv.Itoa(category.Id), requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-key", "RAHASIA")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)
	//fmt.Println(responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "Ok", responseBody["status"])
	assert.Equal(t, category.Id, int(responseBody["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, "Bola", responseBody["data"].(map[string]interface{})["name"])

}

func TestUpdateCategoryFailed(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	categoryRepository := repository.NewCategoryRepository()
	tx, _ := db.Begin()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Bola",
	})

	tx.Commit()
	router := setupRouter(db)

	requestBody := strings.NewReader(`{"name" : ""}`)
	request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/categories/"+strconv.Itoa(category.Id), requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-key", "RAHASIA")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)
	//fmt.Println(responseBody)

	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Equal(t, "Bad Request", responseBody["status"])

}

func TestGetCategorySuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	categoryRepository := repository.NewCategoryRepository()
	tx, _ := db.Begin()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Buku",
	})

	tx.Commit()
	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories/"+strconv.Itoa(category.Id), nil)
	//request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-key", "RAHASIA")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)
	//fmt.Println(responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "Ok", responseBody["status"])
	assert.Equal(t, category.Id, int(responseBody["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, category.Name, responseBody["data"].(map[string]interface{})["name"])

}

func TestGetCategoryFailed(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)
	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories/404", nil)
	request.Header.Add("X-API-key", "RAHASIA")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, 404, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)
	//fmt.Println(responseBody)

	assert.Equal(t, 404, int(responseBody["code"].(float64)))
	assert.Equal(t, "Not Found", responseBody["status"])

}

func TestDeleteCategorySuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	categoryRepository := repository.NewCategoryRepository()
	tx, _ := db.Begin()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Baju",
	})

	tx.Commit()
	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/categories/"+strconv.Itoa(category.Id), nil)
	request.Header.Add("X-API-key", "RAHASIA")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)
	fmt.Println(responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "Ok", responseBody["status"])

}

func TestDeleteCategoryFailed(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)
	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/categories/123", nil)
	request.Header.Add("X-API-key", "RAHASIA")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, 404, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)
	//fmt.Println(responseBody)

	assert.Equal(t, 404, int(responseBody["code"].(float64)))
	assert.Equal(t, "Not Found", responseBody["status"])

}

func TestListCategorySuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	categoryRepository := repository.NewCategoryRepository()
	tx, _ := db.Begin()
	category1 := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Buku",
	})

	category2 := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Komputer",
	})
	tx.Commit()
	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories", nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-key", "RAHASIA")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)
	fmt.Println(responseBody)
	//jadi disini maksudnya data dari response body kita konversi ke
	//slice map[string]interface{}
	var categories = responseBody["data"].([]interface{})
	categoriesResponse1 := categories[0].(map[string]interface{})
	categoriesResponse2 := categories[1].(map[string]interface{})
	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "Ok", responseBody["status"])
	assert.Equal(t, category1.Id, int(categoriesResponse1["id"].(float64)))
	assert.Equal(t, category1.Name, categoriesResponse1["name"])
	assert.Equal(t, category2.Id, int(categoriesResponse2["id"].(float64)))
	assert.Equal(t, category2.Name, categoriesResponse2["name"])
}

func TestUnauthorized(t *testing.T) {

	db := setupTestDB()
	truncateCategory(db)

	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories", nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-key", "salah")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, 401, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)
	fmt.Println(responseBody)

	assert.Equal(t, 401, int(responseBody["code"].(float64)))
	assert.Equal(t, "UnAouthorized", responseBody["status"])

}
