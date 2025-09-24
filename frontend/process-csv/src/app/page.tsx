"use client";
import React, { useState, useRef } from "react";

const CSVUploader: React.FC = () => {
  const [file, setFile] = useState<File | null>(null);
  const [isUploading, setIsUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const [downloadUrl, setDownloadUrl] = useState<string | null>(null);
  const inputRef = useRef<HTMLInputElement | null>(null);

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

    try {
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
                if (inputRef.current) {
                  inputRef.current.value = "";
                }
                resolve();
              } else {
                reject(new Error("No download URL in response."));
              }
            } catch {
              reject(new Error("Invalid server response."));
            }
          } else {
            reject(new Error(`Upload failed: ${xhr.status} - ${xhr.statusText}`));
          }
        });

        xhr.addEventListener("error", () => {
          reject(new Error("Network error during upload."));
        });

        xhr.open("POST", process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/upload");
        xhr.send(formData);
      });
    } catch (err: unknown) {
      if (err instanceof Error) {
        setError(err.message || "Error uploading file.");
        console.error("Upload error:", err);
      } else {
        setError("Error uploading file.");
        console.error("Upload error:", err);
      }
    } finally {
      setIsUploading(false);
    }
  };

  return (
    <div className="bg-gray-900 border border-gray-700 rounded-xl p-8 shadow-lg max-w-md w-full flex flex-col gap-6 text-gray-100">
      <h2 className="text-2xl font-bold mb-4 text-cyan-400 text-center">
        CSV File Processor
      </h2>
      <form
        onSubmit={handleUpload}
        className="flex flex-col gap-4"
      >
        <label
          htmlFor="csvFileInput"
          className={`block bg-gray-800 border-2 border-dashed border-cyan-400 rounded-lg p-6 text-center cursor-pointer transition-colors ease-in-out duration-300 hover:bg-gray-700 text-lg ${file ? "text-white" : "text-gray-400"
            }`}
        >
          {file
            ? `Selected: ${file.name}`
            : "Click to upload the CSV file"}
          <input
            id="csvFileInput"
            type="file"
            accept=".csv,text/csv"
            onChange={handleFileChange}
            className="hidden"
            ref={inputRef}
          />
        </label>

        {error && (
          <p className="text-red-400 text-sm text-center">
            {error}
          </p>
        )}

        {isUploading && (
          <div className="w-full text-center flex flex-col items-center gap-2">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-cyan-400"></div>
            <p className="mb-2">
              Uploading... {Math.round(uploadProgress)}%
            </p>
            <div
              className="w-full bg-gray-800 rounded-full h-2.5 overflow-hidden"
            >
              <div
                className="bg-cyan-400 h-2.5 rounded-full transition-all duration-300 ease-in-out"
                style={{ width: `${uploadProgress}%` }}
              ></div>
            </div>
          </div>
        )}

        <button
          type="submit"
          disabled={!file || isUploading}
          className="bg-cyan-400 text-gray-900 font-bold py-3 px-6 rounded-lg text-lg transition-colors ease-in-out mt-2 disabled:bg-gray-600 disabled:text-gray-400 disabled:cursor-not-allowed hover:bg-cyan-300"
        >
          {isUploading ? "Uploading..." : "Upload & Process"}
        </button>
      </form>

      {downloadUrl && (
        <a
          href={downloadUrl}
          target="_blank"
          rel="noopener noreferrer" // Added for security
          className="mt-4 inline-block bg-cyan-400 text-gray-900 py-3 px-6 rounded-lg font-bold no-underline text-center transition-colors duration-300 hover:bg-cyan-300"
        >
          Download Processed CSV
        </a>
      )}
    </div>
  );
};

export default function Home() {
  return (
    <div className="bg-black min-h-screen flex items-center justify-center p-8 font-sans">
      <main className="flex flex-col items-center gap-8 w-full">
        <CSVUploader />
      </main>
    </div>
  );
}