package models

type Task struct {
	ID int `json:"id" gorm:"primaryKey"`
	Title string `json:"title"`
 	Status string `json:"status"`
}

type DependencyStatus struct {
	Status string `json:"status"`
	Error string `json:"error,omitempty"`
}

type HealthStatus struct {
	Status string `json:"status"`
	Uptime string `json:"uptime"`
	Dependencies map[string]DependencyStatus `json:"dependencies"`
}