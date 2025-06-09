package dto

type ReportRequest struct {
	StartDate  string `json:"start_date" binding:"required"`
	EndDate    string `json:"end_date" binding:"required"`
	AdminEmail string `json:"admin_email" binding:"required,email"`
}
