package web

//disni requset untuk create apa aja sih
//cukup name karena id outo increment
type CategoryCreateRequest struct {
	Name string `validate:"required,max=200,min=1" json:"name"`
}
