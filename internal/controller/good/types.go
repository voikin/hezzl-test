package good

type ResponseRemove struct {
	Id        int  `json:"id"`
	ProjectId int  `json:"projectId"`
	Removed   bool `json:"removed"`
}

type RequestCreate struct {
	Name string `json:"name" binding:"required"`
}

type RequestUpdate struct {
	RequestCreate
	Description string `json:"description,omitempty"`
}

type RequestReprioritize struct {
	NewPriority int `json:"newPriority"`
}
