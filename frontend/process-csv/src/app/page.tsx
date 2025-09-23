"use client";
import React, { useState } from "react";

const CSVUploader: React.FC = () => {
  const [file, setFile] = useState<File | null>(null);
  const [isUploading, setIsUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const [downloadUrl, setDownloadUrl] = useState<string | null>(null);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      setFile(e.target.files[0]);
      setError(null);
      setDownloadUrl(null);
    }
  };

  const handleUpload = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!file) {
      setError("Please select a file to upload.");
      return;
    }

    setIsUploading(true);
    setUploadProgress(0);
    setError(null);
    setDownloadUrl(null);

    const formData = new FormData();
    formData.append("csv_file", file);

    const xhr = new XMLHttpRequest();

    await new Promise<void>((resolve, reject) => {
      xhr.upload.addEventListener("progress", (event) => {
        if (event.lengthComputable) {
          const percentComplete = (event.loaded / event.total) * 100;
          setUploadProgress(percentComplete);
        }
      });

      xhr.addEventListener("load", () => {
        if (xhr.status >= 200 && xhr.status < 300) {
          try {
            const response = JSON.parse(xhr.responseText);
            if (response.download_url) {
              setDownloadUrl(response.download_url);
              setFile(null);
              resolve();
            } else {
              reject(new Error("No download URL in response."));
            }
          } catch (err) {
            reject(new Error("Invalid server response."));
          }
        } else {
          reject(new Error(`Upload failed: ${xhr.status} - ${xhr.statusText}`));
        }
      });

      xhr.addEventListener("error", () => {
        reject(new Error("Network error during upload."));
      });

      xhr.open("POST", "http://localhost:8080/upload");
      xhr.send(formData);
    }).catch((err: any) => {
      setError(err.message || "Error uploading file.");
      console.error("Upload error:", err);
    });

    setIsUploading(false);
  };

  return (
    <div
      style={{
        backgroundColor: "#1a1a1a",
        border: "1px solid #333",
        borderRadius: "8px",
        padding: "2rem",
        boxShadow: "0 4px 8px rgba(0, 0, 0, 0.4)",
        maxWidth: "500px",
        width: "100%",
        display: "flex",
        flexDirection: "column",
        gap: "1.5rem",
        color: "#f0f0f0",
      }}
    >
      <h2 style={{ fontSize: "1.8rem", marginBottom: "1rem", color: "#61dafb" }}>
        CSV File Processor
      </h2>
      <form onSubmit={handleUpload} style={{ display: "flex", flexDirection: "column", gap: "1rem" }}>
        <label
          htmlFor="csvFileInput"
          style={{
            display: "block",
            backgroundColor: "#282c34",
            border: "2px dashed #61dafb",
            borderRadius: "5px",
            padding: "1.5rem",
            textAlign: "center",
            cursor: "pointer",
            transition: "background-color 0.3s ease",
            fontSize: "1.1rem",
            color: file ? "#ffffff" : "#cccccc",
          }}
          onMouseOver={(e) => (e.currentTarget.style.backgroundColor = "#3a3f47")}
          onMouseOut={(e) => (e.currentTarget.style.backgroundColor = "#282c34")}
        >
          {file ? `Selected: ${file.name}` : "Click or Drag & Drop CSV File Here"}
          <input
            id="csvFileInput"
            type="file"
            accept=".csv,text/csv"
            onChange={handleFileChange}
            style={{ display: "none" }}
          />
        </label>

        {error && (
          <p style={{ color: "#ff6b6b", fontSize: "0.9rem", textAlign: "center" }}>
            {error}
          </p>
        )}

        {isUploading && (
          <div style={{ width: "100%", textAlign: "center" }}>
            <p style={{ marginBottom: "0.5rem" }}>
              Uploading... {Math.round(uploadProgress)}%
            </p>
            <div
              style={{
                width: "100%",
                backgroundColor: "#282c34",
                borderRadius: "5px",
                overflow: "hidden",
              }}
            >
              <div
                style={{
                  width: `${uploadProgress}%`,
                  height: "10px",
                  backgroundColor: "#61dafb",
                  borderRadius: "5px",
                  transition: "width 0.3s ease-in-out",
                }}
              ></div>
            </div>
          </div>
        )}

        <button
          type="submit"
          disabled={!file || isUploading}
          style={{
            backgroundColor: isUploading ? "#555" : "#61dafb",
            color: isUploading ? "#ccc" : "#1a1a1a",
            border: "none",
            borderRadius: "5px",
            padding: "0.8rem 1.5rem",
            fontSize: "1.1rem",
            fontWeight: "bold",
            cursor: isUploading ? "not-allowed" : "pointer",
            transition: "background-color 0.3s ease",
            marginTop: "0.5rem",
          }}
        >
          {isUploading ? "Uploading..." : "Upload & Process"}
        </button>
      </form>

      {downloadUrl && (
        <a
          href={downloadUrl}
          target="_blank"
          rel="noopener noreferrer"
          style={{
            marginTop: "1rem",
            display: "inline-block",
            backgroundColor: "#61dafb",
            color: "#1a1a1a",
            padding: "0.8rem 1.5rem",
            borderRadius: "5px",
            fontWeight: "bold",
            textDecoration: "none",
            textAlign: "center",
            transition: "background-color 0.3s",
          }}
        >
          Download Processed CSV
        </a>
      )}
    </div>
  );
};

export default function Home() {
  return (
    <div
      style={{
        backgroundColor: "#000000",
        minHeight: "100vh",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        padding: "2rem",
        fontFamily: "sans-serif",
      }}
    >
      <main style={{
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        gap: "2rem",
        width: "100%"
      }}>
        <CSVUploader />
      </main>
    </div>
  );
}