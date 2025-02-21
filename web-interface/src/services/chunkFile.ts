import config from "../configs/config";
import { FileMeta } from "../types";


type Callback = () => void;

type ProgressCallback = (progress: number) => void;


interface UploadResponse {
    success:    boolean;
    error?:     string;
}

export const uploadChunk = async (chunkFormData: FormData, callback: Callback, retry: number = 0) : Promise<UploadResponse> => {

    const MAX_RETRIES = config.MAX_CHUNK_RETRIES;

    try {
        const response = await fetch(config.UPLOAD_URL, {
        method: "POST",
        body: chunkFormData,
        });

        if (!response.ok) throw new Error(`HTTP ${response.status}`);

        callback();
        return { success: true };
    } catch (error) {
        if (retry >= MAX_RETRIES) {
        return { success: false, error: `Max retries exceeded: ${(error as Error).message}` };
        }
        return uploadChunk(chunkFormData, callback, retry + 1);
    }
}

export const createChunk = (file: File, fileMeta: FileMeta, chunkIndex: number) : FormData => {
    
    if (!file) throw new Error('Invalid file');
    if (chunkIndex < 0) throw new Error('Invalid chunk number');

    const CHUNK_SIZE = config.MAX_CHUNK_SIZE_MB * 1024 * 1024;
    const totalChunks = Math.ceil(file.size / CHUNK_SIZE);

    const start = chunkIndex * CHUNK_SIZE;
    const end = Math.min(start + CHUNK_SIZE, file.size);
    const chunk = file.slice(start, end);
    const formData = new FormData();

    formData.append('fileId', fileMeta.fileId);
    formData.append('fileName', fileMeta.fileName);
    formData.append('fileExtension', fileMeta.fileExtension);
    formData.append('md5Hash', fileMeta.md5Hash);
    formData.append('chunkIndex', chunkIndex.toString());
    formData.append('totalChunks', totalChunks.toString());
    formData.append('chunk', chunk);

    return formData;
}

export const chunkAndUpload = async (onProgress: ProgressCallback, fileMeta: FileMeta, file: File) : Promise<UploadResponse> => {
    if (!file) return { success: false, error: "No file provided" };
    
    const CHUNK_SIZE = config.MAX_CHUNK_SIZE_MB * 1024 * 1024;
    const CONCURRENT_CHUNKS = config.MAX_CONCURRENT_CHUNKS;
  
    try {
      const totalChunks = Math.ceil(file.size / CHUNK_SIZE);
      let uploadedChunks = 0;
  
      const updateProgress = () => {
        uploadedChunks++;
        const percentage = (uploadedChunks / totalChunks) * 100;
        onProgress(percentage);
      };

      let uploadPromises: Promise<UploadResponse>[] = [];

      for (let chunkIndex = 0; chunkIndex < totalChunks; chunkIndex++) {
        const chunkFormData = createChunk(file, fileMeta, chunkIndex);

        uploadPromises.push(
          uploadChunk(chunkFormData, updateProgress)
        );
        
        // Upload N chunks concurrently
        if (uploadPromises.length >= CONCURRENT_CHUNKS || chunkIndex === totalChunks) {
          const results = await Promise.all(uploadPromises);
          if (results.some(res => !res.success)) throw new Error("Chunk upload failed");
          uploadPromises = []; // Reset for the next batch
        }
      }
  
      return { success: true };
    } catch (error) {
      return { success: false, error: (error as Error).message };
    }
  };

