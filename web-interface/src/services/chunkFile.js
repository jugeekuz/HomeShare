import config from "../configs/config";

export const uploadChunk = async (chunk, callback, retry=0) => {
    const MAX_RETRIES = config.MAX_CHUNK_RETRIES;
    const uploadUrl = config.UPLOAD_URL

    try {
        const response = await fetch(uploadUrl, {
            method: 'POST',
            body: chunk,
        });
        if (response.ok) {
            callback();
            return { success: true };
        } 
        if (retry >= MAX_RETRIES) {
            return { success: false, error: 'Max retries exceeded' };
        }
        return await uploadChunk(chunk, callback, retry + 1);

    } catch(error) {
        if (retry >= MAX_RETRIES) {
            return { success: false, error: 'Network error: Max retries exceeded' };
        }
        return await uploadChunk(chunk, callback, retry + 1);
    }

}


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

export const chunkAndUpload = async (setProgress, file) => {
    if (!file || !(file instanceof File)) return;

    const CHUNK_SIZE = config.MAX_CHUNK_SIZE_MB * 1024 * 1024;
    const CONCURRENT_REQUESTS = config.MAX_CONCURRENT_REQUESTS;

    try {
        if (!file) return;

        const totalChunks = Math.ceil(file.size / CHUNK_SIZE);

        const callback = (chunkIndex) => {
            const percentage = ((chunkIndex+1)/totalChunks)*100
            if (percentage === 100) setTimeout(() => setProgress(0), 800);
            setProgress(percentage);
        }

        const fileId = `${file.name}`;

        const totalPromises = Math.ceil(totalChunks/CONCURRENT_REQUESTS);

        for (let i = 0; i < totalPromises; i++) {
            const promiseStart = i*CONCURRENT_REQUESTS;
            const promiseEnd = Math.min(totalChunks, (i+1)*CONCURRENT_REQUESTS);
            
            const promises = Array.from({length: (promiseEnd-promiseStart)}, (_,i) => {
                const chunkIndex = promiseStart + i + 1
                const chunk = createChunk(file, fileId, chunkIndex)
                return uploadChunk(chunk, () => callback(chunkIndex))
            })
            
            const results = await Promise.all(promises)
            if (results.some((result) => result.success === false)) {
                throw new Error('Some chunks failed to upload');
            }
        }
        return {success: true}

    } catch (error) {
        return { success: false, error: 'Max retries exceeded' };
    }
}

