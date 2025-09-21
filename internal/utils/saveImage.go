package utils

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
)

func SaveImage(ctx *gin.Context, file *multipart.FileHeader, prefix string) string {
	const maxSize = 2 * 1024 * 1024
	if file.Size > maxSize {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"message": "File too large (max 2MB)",
		})
		return ""
	}

	ext := filepath.Ext(file.Filename)
	re := regexp.MustCompile(`(?i)\.(png|jpg|jpeg|webp)$`)
	if !re.MatchString(ext) {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"message": "Invalid file type (only PNG, JPG, JPEG, WEBP allowed)",
		})
		return ""
	}

	filename := fmt.Sprintf("%s_%d%s", prefix, time.Now().UnixNano(), ext)
	location := filepath.Join("public", filename)

	if err := ctx.SaveUploadedFile(file, location); err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"message": "Failed to save file",
		})
		return ""
	}

	return filename
}
