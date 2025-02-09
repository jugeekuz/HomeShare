// FileUploader.js
import React, { useState } from 'react';

const CHUNK_SIZE = 1024 * 1024; // 1 MB

function FileUploader() {
  const [file, setFile] = useState(null);
  const [uploadProgress, setUploadProgress] = useState(0);

  const handleFileChange = (e) => {
    setFile(e.target.files[0]);
  };

  const uploadFile = async () => {
    if (!file) return;

    const totalChunks = Math.ceil(file.size / CHUNK_SIZE);
    // Create a simple unique ID using the file name and current time.
    const fileID = `${file.name}-${Date.now()}`;
    let uploadedChunks = 0;

    // Loop through the file and send each chunk
    for (let chunkNumber = 1; chunkNumber <= totalChunks; chunkNumber++) {
      const start = (chunkNumber - 1) * CHUNK_SIZE;
      const end = Math.min(start + CHUNK_SIZE, file.size);
      const chunk = file.slice(start, end);

      const formData = new FormData();
      formData.append('fileID', fileID);
      formData.append('chunkNumber', chunkNumber);
      formData.append('totalChunks', totalChunks);
      formData.append('chunk', chunk);

      try {
        const response = await fetch('http://localhost:8080/upload', {
          method: 'POST',
          body: formData,
        });
        if (response.ok) {
          uploadedChunks++;
          setUploadProgress(Math.round((uploadedChunks / totalChunks) * 100));
        } else {
          console.error('Chunk upload failed');
          break;
        }
      } catch (error) {
        console.error('Error uploading chunk', error);
        break;
      }
    }
  };

  return (
    <div>
      <h1>Chunked File Upload</h1>
      <input type="file" onChange={handleFileChange} />
      <button onClick={uploadFile} disabled={!file}>
        Upload
      </button>
      <div>Progress: {uploadProgress}%</div>
    </div>
  );
}

export default FileUploader;
