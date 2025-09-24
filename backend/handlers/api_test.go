package handlers_test

import (
	"bytes"
	"csv_processor/handlers"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	router := handlers.SetupRouter()
	return router
}

func createTestFile(t *testing.T, content string) (string, func()) {
	t.Helper()
	tmpFile, err := os.CreateTemp("", "test.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	return tmpFile.Name(), func() {
		err := os.Remove(tmpFile.Name())
		if err != nil {
			log.Printf("Failed to remove temp file: %v", err)
		}
	}
}

func performRequest(r http.Handler, method, path string, body io.Reader, contentType string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestUploadCSV(t *testing.T) {
	router := setupRouter()

	t.Run("Valid CSV Upload", func(t *testing.T) {
		csvContent := "New York,2022-01-01,1200.50\nLos Angeles,2022-01-01,800.00"
		filePath, cleanup := createTestFile(t, csvContent)
		defer cleanup()

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		file, err := os.Open(filePath)
		if err != nil {
			t.Fatalf("Failed to open test file: %v", err)
		}
		defer file.Close()
		fileField, err := writer.CreateFormFile("csv_file", filepath.Base(filePath))
		if err != nil {
			t.Fatalf("Failed to create form file: %v", err)
		}
		_, err = io.Copy(fileField, file)
		if err != nil {
			t.Fatalf("Failed to copy file content: %v", err)
		}
		writer.Close()

		w := performRequest(router, http.MethodPost, "/upload", body, writer.FormDataContentType())

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "download_url")

		// Parse the download_url from the response
		var resp struct {
			DownloadURL string `json:"download_url"`
		}
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		if err != nil {
			t.Fatalf("Failed to parse JSON response: %v", err)
		}
		assert.NotEmpty(t, resp.DownloadURL)

		// Extract the filename from the download URL
		parts := strings.Split(resp.DownloadURL, "/")
		filename := parts[len(parts)-1]

		// Now test downloading the file
		downloadPath := fmt.Sprintf("/download/%s", filename)
		w2 := performRequest(router, http.MethodGet, downloadPath, nil, "")

		assert.Equal(t, http.StatusOK, w2.Code)
		assert.Contains(t, w2.Body.String(), "New York")
		assert.Contains(t, w2.Body.String(), "Los Angeles")
	})

	t.Run("Invalid File Type Upload", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("csv_file", "not a csv")
		writer.Close()

		w := performRequest(router, http.MethodPost, "/upload", body, writer.FormDataContentType())

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Could not get uploaded file")
	})

	t.Run("Missing File Upload", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.Close()

		w := performRequest(router, http.MethodPost, "/upload", body, writer.FormDataContentType())

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Could not get uploaded file")
	})
}

func TestDownloadCSV_FileNotFound(t *testing.T) {
	router := setupRouter()

	w := performRequest(router, http.MethodGet, "/download/nonexistent.csv", nil, "")
	assert.Equal(t, http.StatusNotFound, w.Code)
}
