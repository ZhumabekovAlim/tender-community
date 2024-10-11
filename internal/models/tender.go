package models

import (
	"errors"
	"time"
)

var ErrTenderNotFound = errors.New("tender not found")

type Tender struct {
	ID            int       `json:"id"`
	Type          string    `json:"type"`
	TenderNumber  *string   `json:"tender_number,omitempty"`
	UserID        int       `json:"user_id"`
	CompanyID     int       `json:"company_id"`
	Organization  string    `json:"organization,omitempty"`
	Total         float64   `json:"total"`
	Commission    float64   `json:"commission"`
	CompletedDate *string   `json:"completed_date,omitempty"`
	Date          time.Time `json:"date"`
	Status        int       `json:"status"`
	UserName      *string   `json:"username,omitempty"`
	CompanyName   *string   `json:"companyname,omitempty"`
}

type TenderDebt struct {
	GOIK float64 `json:"goik"`
	GOPP float64 `json:"gopp"`
}

type TenderCount struct {
	Total   int `json:"total"`
	Status0 int `json:"status_0"`
	Status1 int `json:"status_1"`
	Status2 int `json:"status_2"`
	Status3 int `json:"status_3"`
}
