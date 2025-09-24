# CSV Processor

A full-stack application for uploading, processing, and downloading CSV files. The backend is written in Go (Gin), and the frontend is a Next.js (React) app.

---

## How to Run the App

### **Backend**

1. **Install dependencies:**
   ```sh
   cd backend
   go mod tidy
   ```

2. **Run the backend server:**
   ```sh
   go run main.go
   ```
   The server will start on `http://localhost:8080`.

### **Frontend**

1. **Install dependencies:**
   ```sh
   cd frontend/process-csv
   npm install
   ```

2. **Run the frontend app:**
   ```sh
   npm run dev
   ```
   The app will be available at `http://localhost:3000`.

---

## How to Test

### **Backend Unit Tests**

Run all backend tests with coverage:
```sh
cd backend
go test -v ./... -cover
```

To generate an HTML coverage report:
```sh
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### **Frontend**

The frontend is a Next.js app. You can add tests using your preferred React testing library (e.g., Jest, React Testing Library).

---

## Algorithm Explanation & Memory Efficiency

### **Processing Algorithm**

- The backend reads the uploaded CSV file **row by row** using a streaming parser (`encoding/csv`).
- For each row, it extracts the city and sales value.
- It maintains a `map[string]float64` to sum sales per city.
- After processing, it writes the result to a new CSV file.

### **Memory Efficiency Strategy**

- **Streaming:** The CSV is processed as a stream (`io.Reader`), so the entire file is **not loaded into memory** at once.
- **Aggregation:** Only the city sales map is kept in memory, which is efficient unless there are millions of unique cities.
- **Output:** The result is written directly to disk, not kept in memory.

### **Estimated Big O Complexity**

- **Time Complexity:** O(n), where n is the number of rows in the CSV file.
- **Space Complexity:** O(c), where c is the number of unique cities (since only the aggregation map is stored).

---

## Unit Test Coverage (Backend)

- **Services:**  
  - `ProcessCSV` and `WriteCitySalesCSV` are covered with various cases (valid, malformed, empty, etc.).
  - Coverage: **~93%** (see `go test -cover` output).

- **Handlers:**  
  - `UploadCSV` and `DownloadCSV` are tested for valid uploads, missing/invalid files, and download (including not found).
  - Tests check both the upload and download workflow, including parsing the returned download URL.

---

## Directory Structure

```
backend/
  main.go
  handlers/
    api.go
    api_test.go
  services/
    processor.go
    processor_test.go
  processed_files/
frontend/
  process-csv/
    src/app/
      page.tsx
```

---

## Notes

- The backend saves processed files in `backend/processed_files/`.
- The download link is returned in the upload response and used by the frontend for downloading the result.
- CORS is enabled for local development.

---