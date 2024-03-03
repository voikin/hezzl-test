package project

type RequestCreate struct {
	Name string `json:"name" binding:"required"`
}

type RequestUpdate struct {
	RequestCreate
}
