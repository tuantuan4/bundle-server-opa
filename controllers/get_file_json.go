package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"opa-test/models"
	"os"
)

func WriteJsonToFile(data interface{}, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	return nil
}

type ResponseDataOPA struct {
	Name     string `json:"name"`
	EndPoint string `json:"endPoint"`
	Role     string `json:"role"`
}

func GetFileJson(db *gorm.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var result []models.User
		if err := db.Find(&result).Error; err != nil {
			ctx.JSONP(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		var modifiedDatas []ResponseDataOPA
		for _, i := range result {
			res := ResponseDataOPA{
				Name:     i.Name,
				EndPoint: i.EndPoint,
				Role:     i.Role,
			}
			modifiedDatas = append(modifiedDatas, res)
		}
		ctx.JSON(http.StatusOK, modifiedDatas)
		WriteJsonToFile(modifiedDatas, "bundle/data/data.json")
	}
}
