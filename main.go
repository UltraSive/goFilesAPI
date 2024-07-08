package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FileMetadata struct {
	Type     string    `json:"type"`
	Name     string    `json:"name"`
	Size     int64     `json:"size"`
	Modified time.Time `json:"modified"`
}

func main() {
	r := gin.Default()

	files := r.Group("/files")
	{
		files.GET("/contents", getServerFileContents)
		files.GET("/list-directory", getServerListDirectory)
		files.PUT("/rename", putServerRenameFiles)
		files.POST("/copy", postServerCopyFile)
		files.POST("/write", postServerWriteFile)
		files.POST("/delete", postServerDeleteFiles)
		files.POST("/compress", postServerCompressFiles)
		files.POST("/decompress", postServerDecompressFiles)
		files.POST("/chmod", postServerChmodFile)
	}

	r.Run() // listen and serve on 0.0.0.0:8080
}

// Handler for getting file contents
func getServerFileContents(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Path query parameter is required"})
		return
	}

	content, err := os.ReadFile(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"content": string(content)})
}

// Handler for listing directory contents
func getServerListDirectory(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Path query parameter is required"})
		return
	}

	files, err := os.ReadDir(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var fileMetadata []FileMetadata
	for _, file := range files {
		//filePath := filepath.Join(path, file.Name())

		// Get file info
		fileInfo, err := file.Info()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Determine file type
		fileType := "file"
		if fileInfo.IsDir() {
			fileType = "directory"
		}

		// Format size
		size := fileInfo.Size()

		// Create FileMetadata object
		meta := FileMetadata{
			Type:     fileType,
			Name:     file.Name(),
			Size:     size,
			Modified: fileInfo.ModTime(),
		}

		fileMetadata = append(fileMetadata, meta)
	}

	c.JSON(http.StatusOK, gin.H{"files": fileMetadata})
}

// Handler for renaming files
func putServerRenameFiles(c *gin.Context) {
	oldPath := c.Query("old_path")
	newPath := c.Query("new_path")
	if oldPath == "" || newPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Old path and new path query parameters are required"})
		return
	}

	err := os.Rename(oldPath, newPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File renamed successfully"})
}

// Handler for copying files
func postServerCopyFile(c *gin.Context) {
	src := c.Query("src")
	dest := c.Query("dest")
	if src == "" || dest == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Source and destination query parameters are required"})
		return
	}

	input, err := os.ReadFile(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = os.WriteFile(dest, input, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File copied successfully"})
}

// Handler for writing files
func postServerWriteFile(c *gin.Context) {
	path := c.Query("path")
	content := c.PostForm("content")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Path query parameter is required"})
		return
	}

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File written successfully"})
}

// Handler for deleting files
func postServerDeleteFiles(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Path query parameter is required"})
		return
	}

	err := os.Remove(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}

// Handler for compressing files
func postServerCompressFiles(c *gin.Context) {
	// Implement your compression logic here
}

// Handler for decompressing files
func postServerDecompressFiles(c *gin.Context) {
	// Implement your decompression logic here
}

// Handler for changing file permissions
func postServerChmodFile(c *gin.Context) {
	path := c.Query("path")
	mode := c.Query("mode")
	if path == "" || mode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Path and mode query parameters are required"})
		return
	}

	parsedMode, err := strconv.ParseUint(mode, 8, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mode"})
		return
	}

	err = os.Chmod(path, os.FileMode(parsedMode))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File permissions changed successfully"})
}
