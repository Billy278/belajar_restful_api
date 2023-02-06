package controller

import (
	"belajar-golang-api/helper"
	"belajar-golang-api/model/web"
	"belajar-golang-api/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// keanapa bentuk category service tidak pointer karena dia bentuk interface bukan struct
type CategoryControllerImpl struct {
	CategoryService service.CategoryService
}

func NewCategoryController(categoryService service.CategoryService) CategoryController {
	return &CategoryControllerImpl{
		CategoryService: categoryService,
	}

}

func (controller *CategoryControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	//jadi karena request kalo mau baca dari reader /artinya straming artinya
	//kita bisa langsung streaming encode dan decodenya tanpa perlu di conversi ke string
	decoder := json.NewDecoder(request.Body)

	categoryCreateRequest := web.CategoryCreateRequest{}
	err := decoder.Decode(&categoryCreateRequest)
	helper.PanicIfError(err)

	categoryResponse := controller.CategoryService.Create(request.Context(), categoryCreateRequest)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "Ok",
		Data:   categoryResponse,
	}

	writer.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	err = encoder.Encode(webResponse)
	helper.PanicIfError(err)

}

func (controller *CategoryControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	categoryUpdateRequest := web.CategoryUpdateRequest{}
	helper.ReadFromRequest(request, &categoryUpdateRequest)

	categoryId := params.ByName("categoryId")
	id, err := strconv.Atoi(categoryId)
	helper.PanicIfError(err)
	categoryUpdateRequest.Id = id
	categoryResponse := controller.CategoryService.Update(request.Context(), categoryUpdateRequest)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "Ok",
		Data:   categoryResponse,
	}

	helper.WriteToResponseBody(writer, &webResponse)

}

func (controller *CategoryControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	categoryId := params.ByName("categoryId")
	id, err := strconv.Atoi(categoryId)
	helper.PanicIfError(err)

	controller.CategoryService.Delete(request.Context(), id)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "Ok",
	}

	helper.WriteToResponseBody(writer, &webResponse)
}

func (controller *CategoryControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	categoryId := params.ByName("categoryId")
	id, err := strconv.Atoi(categoryId)
	helper.PanicIfError(err)

	categoryResponse := controller.CategoryService.FindById(request.Context(), id)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "Ok",
		Data:   categoryResponse,
	}

	helper.WriteToResponseBody(writer, &webResponse)
}

func (controller *CategoryControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	categoryResponses := controller.CategoryService.FindAll(request.Context())
	webResponse := web.WebResponse{
		Code:   200,
		Status: "Ok",
		Data:   categoryResponses,
	}

	helper.WriteToResponseBody(writer, &webResponse)
}
