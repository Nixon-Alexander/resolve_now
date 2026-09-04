package model

import "time"

type CompanyStatus string

const (
	COMPANY_ACTIVE    CompanyStatus = "active"
	COMPANY_SUSPENDED CompanyStatus = "suspended"
)

type AddCompanyRequest struct {
	Name   string        `json:"name"`
	Status CompanyStatus `json:"status"`
}

type GetAllCompanyResponse struct {
	Id        string        `json:"id"`
	Name      string        `json:"name"`
	Status    CompanyStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
