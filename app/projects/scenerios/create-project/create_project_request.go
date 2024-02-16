package create_project

type CreateProjectRequest struct {
	Name string `json:"name" bson:"name" validate:"required,min=3"`
}
