package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"amasd/core"
	"amasd/models"
	"strconv"
)

type Spider struct {
	core.BaseController
}

func (s *Spider) Index(c *gin.Context) {
	if core.IsAjax(c) {
		projectId, _ := strconv.Atoi(c.DefaultPostForm("project_id", "1"))
		version := c.DefaultPostForm("version", "")
		page, _ := strconv.Atoi(c.DefaultPostForm("pagination[page]", "1"))
		pageLength, _ := strconv.Atoi(c.DefaultPostForm("pagination[perpage]", "10"))
		spider := new(models.Spider)
		spiders, totalCount := spider.FindPageSpiders(projectId, version, page, pageLength, "")
		c.JSON(http.StatusOK, gin.H{
			"data": spiders,
			"meta": gin.H{
				"page":    page,
				"total":   totalCount,
				"pages":   core.CalculationPages(totalCount, pageLength),
				"perpage": pageLength,
			},
		})
	} else {
		c.HTML(http.StatusOK, "spider/index", gin.H{
			"projects": new(models.Project).Find(),
		})
	}

}
