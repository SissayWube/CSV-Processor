// Package handlers provides HTTP handlers for the CSV processing service.
package handlers

import (
	"csv_processor/services"
	"fmt"
	"log"
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
const processedFilesDir = "/home/sissay/Desktop/CSV-Processor/backend/processed_files"

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
	// Retrieve the uploaded file from the form data.
	file, err := c.FormFile("csv_file")
	if err != nil {
		log.Printf("Error getting form file: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not get uploaded file."})
		return
	}

	// Validate that the uploaded file is indeed a CSV.
	if file.Header.Get("Content-Type") != "text/csv" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Please upload a CSV file."})
		return
	}

	// Open the uploaded file for reading.
	uploadedFile, err := file.Open()
	if err != nil {
		log.Printf("Could not open uploaded file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not process uploaded file."})
		return
	}
	defer uploadedFile.Close() // Ensure the file is closed after processing.

	// Process the CSV file content using the business logic in the services package.
	totalSales, err := services.ProcessCSV(uploadedFile)
	if err != nil {
		log.Printf("Error processing CSV file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing CSV file."})
		return
	}

	// Generate a unique filename for the processed result to avoid naming conflicts.
	resultFileName := fmt.Sprintf("city_sales_%d.csv", time.Now().UnixNano())

	// Ensure the directory exists.
	if err := os.MkdirAll(processedFilesDir, 0755); err != nil {
		log.Printf("Could not create directory for processed files: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not set up storage for processed files."})
		return
	}
	resultFilePath := filepath.Join(processedFilesDir, resultFileName)

	// Create the output file to store the processed CSV data.
	outFile, err := os.Create(resultFilePath)
	if err != nil {
		log.Printf("Could not create result file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create result file."})
		return
	}
	defer outFile.Close() // Ensure the output file is closed.

	// Write the processed sales data to the output CSV file.
	if err := services.WriteCitySalesCSV(*totalSales, outFile); err != nil {
		log.Printf("Error writing processed CSV: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating result CSV."})
		return
	}

	// Construct the download URL for the processed file.
	downloadURL := fmt.Sprintf("%s/download/%s", serverBaseURL, resultFileName)
	c.JSON(http.StatusOK, gin.H{"download_url": downloadURL})
}

// DownloadCSV handles the HTTP GET request for downloading a processed CSV file.
func DownloadCSV(c *gin.Context) {
	filename := c.Param("filename")

	// Sanitize the filename to prevent directory traversal vulnerabilities.

	if strings.Contains(filename, "..") || strings.ContainsRune(filename, os.PathSeparator) {
		c.String(http.StatusBadRequest, "Invalid filename")
		return
	}

	// Construct the full path to the file.
	filePath := filepath.Join(processedFilesDir, filename)

	// Serve the file for download.
	c.FileAttachment(filePath, filename)
}
