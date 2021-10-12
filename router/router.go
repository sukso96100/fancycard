package router

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sukso96100/fancycard/render"
	"github.com/sukso96100/fancycard/tmpl"
)

func SetupRouter(e *gin.Engine) {
	e.GET("/url", RenderWithDataFromURL)
}

func RenderWithDataFromURL(c *gin.Context) {
	templatePath := c.DefaultQuery("template", "default")
	queryParams := c.Request.URL.Query()

	templateText, err := tmpl.LoadTemplate(templatePath)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Template not found",
		})
		return
	}
	compiledTemplate, err := tmpl.CompileTemplate(templateText, queryParams)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	externalTemplateURL := ""
	if strings.HasPrefix(templatePath, "http://") || strings.HasPrefix(templatePath, "https://") {
		externalTemplateURL = templatePath
	}

	imageBuff, err := render.RenderImage(compiledTemplate, externalTemplateURL, render.DefaultRenderOptions)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Data(200, "image/png", imageBuff)
}

func RenderWithDataFromMetaTags(c *gin.Context) {
	targetURL := c.DefaultQuery("url", "")
	if targetURL == "" {
		c.JSON(400, gin.H{
			"error": "Target URL is empty",
		})
	}

}
