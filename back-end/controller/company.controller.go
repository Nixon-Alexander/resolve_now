package controller

import (
	"strconv"

	"github.com/Nixon-Alexander/resolve_now.git/config"
	"github.com/Nixon-Alexander/resolve_now.git/helper"
	"github.com/Nixon-Alexander/resolve_now.git/model"
	"github.com/Nixon-Alexander/resolve_now.git/repository"
	"github.com/gin-gonic/gin"
)

type CompanyController struct {
	repo *repository.CompanyRepository
	log  *config.Log
}

func InitCompanyController(repo *repository.CompanyRepository, log *config.Log) *CompanyController {
	return &CompanyController{
		repo: repo,
		log:  log,
	}
}

func (cc *CompanyController) AddCompany(c *gin.Context) {
	var addRequest model.AddCompanyRequest
	err := c.ShouldBindJSON(&addRequest)
	if err != nil {
		cc.log.Error("Add Company Bind JSON", "error", err.Error())
		helper.SendErrorJSON(c, 400, "Invalid data type", err)
		return
	}

	result, err := cc.repo.Add(c, &addRequest)
	if err != nil {
		cc.log.Error("Add Company Bind JSON", "error", err.Error())
		helper.SendErrorJSON(c, 400, "Error when add data: ", err)
		return
	}

	cc.log.Info("Add Company Request", "request", addRequest)
	helper.SendSuccessJSON(c, 202, "Succes add company", result)
}

func (cc *CompanyController) GetAllCompany(c *gin.Context) {
	rows, err := cc.repo.GetAll(c)
	if err != nil {
		cc.log.Error("GET ALL COMPANY ERROR", "error", err)
		helper.SendErrorJSON(c, 500, "Error when getting all company", err)
		return
	}

	var response []model.GetAllCompanyResponse
	for rows.Next() {
		var r model.GetAllCompanyResponse
		err = rows.Scan(&r.Id, &r.Name, &r.Status, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			cc.log.Error("ERROR WHEN MAPPING RESPONSE", "error", err)
			helper.SendErrorJSON(c, 400, "Error when mapping response", err)
			return
		}

		response = append(response, r)
	}

	cc.log.Info("Get All Company Response", "response", response)
	helper.SendSuccessJSON(c, 200, "Success get list company", response)
}

func (cc *CompanyController) GetCompany(c *gin.Context) {
	var idStr = c.Param("id")
	id, err := strconv.Atoi(idStr)

	rows, err := cc.repo.GetById(c, id)
	if err != nil {
		cc.log.Error("GET COMPANY ERROR", "error", err)
		helper.SendErrorJSON(c, 500, "Error when getting company", err)
		return
	}

	var response []model.GetAllCompanyResponse
	for rows.Next() {
		var r model.GetAllCompanyResponse
		err = rows.Scan(&r.Id, &r.Name, &r.Status, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			cc.log.Error("ERROR WHEN MAPPING RESPONSE", "error", err)
			helper.SendErrorJSON(c, 400, "Error when mapping response", err)
			return
		}

		response = append(response, r)
	}

	cc.log.Info("Get Company Response", "response", response)
	helper.SendSuccessJSON(c, 200, "Success get company", response)
}
