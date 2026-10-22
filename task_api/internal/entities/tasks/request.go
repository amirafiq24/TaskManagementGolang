package tasks

type CreateTask struct {
	Title  string `json:"title"`
	Status string `json:"status"`
}

type UpdateTask struct {
	Title  *string `json:"title"`
	Status *string `json:"status"`
}
