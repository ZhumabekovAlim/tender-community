package models

type DateRangeRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	UserId    int    `json:"user_id"`
}

type DateRangeRequestCompany struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	UserId    int    `json:"user_id"`
	CompanyId int    `json:"company_id"`
}
