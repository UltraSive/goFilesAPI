package main

import (
	"io/ioutil"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

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
		c.JSON(400, gin.H{"error": "Path query parameter is required"})
		return
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"content": string(content)})
}

// Handler for listing directory contents
func getServerListDirectory(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(400, gin.H{"error": "Path query parameter is required"})
		return
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	fileNames := make([]string, len(files))
	for i, file := range files {
		fileNames[i] = file.Name()
	}

	c.JSON(200, gin.H{"files": fileNames})
}

// Handler for renaming files
func putServerRenameFiles(c *gin.Context) {
	oldPath := c.Query("old_path")
	newPath := c.Query("new_path")
	if oldPath == "" || newPath == "" {
		c.JSON(400, gin.H{"error": "Old path and new path query parameters are required"})
		return
	}

	err := os.Rename(oldPath, newPath)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "File renamed successfully"})
}

// Handler for copying files
func postServerCopyFile(c *gin.Context) {
	src := c.Query("src")
	dest := c.Query("dest")
	if src == "" || dest == "" {
		c.JSON(400, gin.H{"error": "Source and destination query parameters are required"})
		return
	}

	input, err := ioutil.ReadFile(src)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = ioutil.WriteFile(dest, input, 0644)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "File copied successfully"})
}

// Handler for writing files
func postServerWriteFile(c *gin.Context) {
	path := c.Query("path")
	content := c.PostForm("content")
	if path == "" {
		c.JSON(400, gin.H{"error": "Path query parameter is required"})
		return
	}

	err := ioutil.WriteFile(path, []byte(content), 0644)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "File written successfully"})
}

// Handler for deleting files
func postServerDeleteFiles(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(400, gin.H{"error": "Path query parameter is required"})
		return
	}

	err := os.Remove(path)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "File deleted successfully"})
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
		c.JSON(400, gin.H{"error": "Path and mode query parameters are required"})
		return
	}

	parsedMode, err := strconv.ParseUint(mode, 8, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid mode"})
		return
	}

	err = os.Chmod(path, os.FileMode(parsedMode))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "File permissions changed successfully"})
}
