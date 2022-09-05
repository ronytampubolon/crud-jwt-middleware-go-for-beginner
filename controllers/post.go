package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rony-tampubolon/sample-rest-api/models"
)

type CreatePostInput struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}
type UpdatePostInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func CreatePost(c *gin.Context) {
	var input CreatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}
	post := models.Post{Title: input.Title, Content: input.Content}
	models.DB.Create(&post)

	c.JSON(http.StatusOK, gin.H{"data": post})
}

func GetAllPost(c *gin.Context) {
	var posts []models.Post
	models.DB.Find(&posts)
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "list all posts",
		"data":    posts})

}

func GetById(c *gin.Context) {
	var post models.Post
	if err := models.DB.Where("id = ?", c.Param("id")).First(&post).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Find by ID",
		"data":    post,
	})
}

func UpdatePost(c *gin.Context) {
	var post models.Post
	if err := models.DB.Where("id = ?", c.Param("id")).First(&post).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	var input UpdatePostInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	// update Post
	updatePost := models.Post{Title: input.Title, Content: input.Content}
	models.DB.Model(&post).Updates(&updatePost)
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Updated successfully",
		"data":    post,
	})
}

// Delete Post
func DeletePost(c *gin.Context) {
	var deletedPost models.Post
	if err := models.DB.Where("id = ?", c.Param("id")).First(&deletedPost).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	models.DB.Delete(&deletedPost)
	c.JSON(http.StatusNoContent, gin.H{})
}
