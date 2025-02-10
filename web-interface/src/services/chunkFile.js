import config from "../configs/config";

export const uploadChunk = async (chunkFormData, callback, retry = 0) => {
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
            return { success: false, error: `Max retries exceeded: ${error.message}` };
        }
        return uploadChunk(chunkFormData, callback, retry + 1);
    }
  };


export const createChunk = (file, fileId, chunkNumber) => {
    if (!file || !(file instanceof File)) throw new Error('Invalid file');
    if (typeof fileId !== 'string') throw new Error('Invalid file ID');
    if (chunkNumber <= 0) throw new Error('Invalid chunk number');

    const CHUNK_SIZE = config.MAX_CHUNK_SIZE_MB * 1024 * 1024;
    const totalChunks = Math.ceil(file.size / CHUNK_SIZE);

    const start = (chunkNumber - 1) * CHUNK_SIZE;
    const end = Math.min(start + CHUNK_SIZE, file.size);
    const chunk = file.slice(start, end);
    const formData = new FormData();

    formData.append('fileID', fileId);
    formData.append('chunkNumber', chunkNumber);
    formData.append('totalChunks', totalChunks);
    formData.append('chunk', chunk);

    return formData;
}

export const chunkAndUpload = async (onProgress, file) => {
    if (!file || !(file instanceof File)) return;
  
    const CHUNK_SIZE = config.MAX_CHUNK_SIZE_MB * 1024 * 1024;
    const CONCURRENT_REQUESTS = config.MAX_CONCURRENT_REQUESTS;
    const fileId = file.name;
  
    try {
      const totalChunks = Math.ceil(file.size / CHUNK_SIZE);
      let uploadedChunks = 0;
  
      const updateProgress = () => {
        uploadedChunks++;
        const percentage = (uploadedChunks / totalChunks) * 100;
        onProgress(percentage);
      };
  
      const uploadPromises = [];
      for (let chunkNumber = 1; chunkNumber <= totalChunks; chunkNumber++) {
        const chunkFormData = createChunk(file, fileId, chunkNumber);
        uploadPromises.push(
          uploadChunk(chunkFormData, updateProgress)
        );
  
        if (uploadPromises.length >= CONCURRENT_REQUESTS || chunkNumber === totalChunks) {
          const results = await Promise.all(uploadPromises);
          if (results.some(res => !res.success)) throw new Error("Chunk upload failed");
          uploadPromises.length = 0;
        }
      }
  
      return { success: true };
    } catch (error) {
      return { success: false, error: error.message };
    }
  };
  
  const generateUniqueFileId = (file) => {
    return `${file.name}-${Date.now()}-${Math.random().toString(36).slice(2)}`;
  };

