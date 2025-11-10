package payload

type CreateMenuCategoryReq struct {
	Name string `json:"name" validate:"required"`
}
