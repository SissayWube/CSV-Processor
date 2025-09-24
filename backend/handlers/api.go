// Package handlers provides HTTP handlers for the CSV processing service.
package handlers

import (
	"csv_processor/services"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// clientOrigin defines the allowed origin for CORS requests.
const clientOrigin = "http://localhost:3000"

// serverBaseURL defines the base URL for the API server.
const serverBaseURL = "http://localhost:8080"

// Define the directory where processed files are stored.
const ProcessedFilesDir = "/home/sissay/Desktop/CSV-Processor/backend/processed_files"

// SetupRouter initializes and configures the Gin router with necessary middleware and routes.
func SetupRouter() *gin.Engine {
	api := gin.Default()

	// Configure CORS middleware to allow requests from the specified client origin.
	api.Use(cors.New(cors.Config{
		AllowOrigins:     []string{clientOrigin},
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Disposition"},
		AllowCredentials: true,
	}))

	// Define API routes.
	// Handles CSV file uploads.
	api.POST("/upload", UploadCSV)

	// Serves processed CSV files for download.
	api.GET("/download/:filename", DownloadCSV)

	return api
}

// UploadCSV handles the HTTP POST request for uploading a CSV file.
func UploadCSV(c *gin.Context) {
	file, err := c.FormFile("csv_file")
	if err != nil {
		c.JSON(400, gin.H{"error": "Could not get uploaded file."})
		return
	}

	uploadedFile, err := file.Open()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Could not open uploaded file: %v", err)})
		return
	}
	defer uploadedFile.Close()

	totalSales, err := services.ProcessCSV(uploadedFile)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Error processing file: %v", err)})
		return
	}

	// Ensure the processed files directory exists
	if err := os.MkdirAll(ProcessedFilesDir, 0755); err != nil {
		c.JSON(500, gin.H{"error": "Could not create processed files directory"})
		return
	}

	resultFileName := fmt.Sprintf("city_sales_%d.csv", time.Now().UnixNano())
	resultFilePath := filepath.Join(ProcessedFilesDir, resultFileName)
	outFile, err := os.Create(resultFilePath)
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not create result file"})
		return
	}
	defer outFile.Close()

	if err := services.WriteCitySalesCSV(*totalSales, outFile); err != nil {
		c.JSON(500, gin.H{"error": "Error generating CSV"})
		return
	}

	downloadURL := serverBaseURL + "/download/" + resultFileName
	c.JSON(200, gin.H{"download_url": downloadURL})
}

// DownloadCSV handles the HTTP GET request for downloading a processed CSV file.
func DownloadCSV(c *gin.Context) {
	filename := c.Param("filename")

	// Sanitize the filename to prevent directory traversal vulnerabilities.
	filePath, err := filepath.Abs(filepath.Join(ProcessedFilesDir, filename))
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}
	if !strings.HasPrefix(filePath, ProcessedFilesDir) {
		c.String(http.StatusBadRequest, "Invalid filename")
		return
	}

	// Serve the file for download.
	c.FileAttachment(filePath, filename)
}
