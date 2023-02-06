package service

import (
	"belajar-golang-api/exception"
	"belajar-golang-api/helper"
	"belajar-golang-api/model/domain"
	"belajar-golang-api/model/web"
	"belajar-golang-api/repository"
	"context"
	"database/sql"

	"github.com/go-playground/validator/v10"
)

type CategoryServiceImpl struct {
	CategoryRepository repository.CategoryRepository
	DB                 *sql.DB
	Validate           *validator.Validate
}

func NewCategoryService(categoryRepository repository.CategoryRepository, DB *sql.DB, validate *validator.Validate) CategoryService {
	return &CategoryServiceImpl{
		CategoryRepository: categoryRepository,
		DB:                 DB,
		Validate:           validate,
	}
}

func (service *CategoryServiceImpl) Create(ctx context.Context, request web.CategoryCreateRequest) web.CategoryResponse {
	//jadi sebelum kita buat transaksinya
	//kita buat validasinya
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollBack(tx)
	category := domain.Category{
		Name: request.Name,
	}
	category = service.CategoryRepository.Save(ctx, tx, category)

	return helper.ToCategoryResponse(category)
	//	return web.CategoryResponse(category)
}
func (service *CategoryServiceImpl) Update(ctx context.Context, request web.CategoryUpdateRequest) web.CategoryResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollBack(tx)

	category, err := service.CategoryRepository.FindById(ctx, tx, request.Id)
	//helper.PanicIfError(err)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}
	category.Name = request.Name
	category = service.CategoryRepository.Update(ctx, tx, category)
	return helper.ToCategoryResponse(category)

}

func (service *CategoryServiceImpl) Delete(ctx context.Context, categoryId int) {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollBack(tx)

	category, err := service.CategoryRepository.FindById(ctx, tx, categoryId)
	//helper.PanicIfError(err)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	service.CategoryRepository.Delete(ctx, tx, category)

}

func (service *CategoryServiceImpl) FindById(ctx context.Context, categoryId int) web.CategoryResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollBack(tx)

	category, err := service.CategoryRepository.FindById(ctx, tx, categoryId)
	//helper.PanicIfError(err)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}
	return helper.ToCategoryResponse(category)
}

func (service *CategoryServiceImpl) FindAll(ctx context.Context) []web.CategoryResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollBack(tx)

	categories := service.CategoryRepository.FindAll(ctx, tx)

	return helper.ToCategoryResponses(categories)

}
