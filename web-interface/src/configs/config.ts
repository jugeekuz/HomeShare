const DOMAIN_NAME = import.meta.env.VITE_DOMAIN || "mydomain.com";
if (!DOMAIN_NAME) {
    throw new Error('❌ Missing VITE_DOMAIN – please set it in your .env file');
}
const config = {
    MAX_CHUNK_SIZE_MB: 5,
    MAX_CHUNK_RETRIES: 5,
    MAX_CONCURRENT_CHUNKS: 6,
    MAX_FILE_SIZE_MB: (1024 * 5) * 1024 * 1024,
    OTP_LENGTH: 6,
    BASE_URL: `https://api.${DOMAIN_NAME}`,
    UPLOAD_URL: `https://api.${DOMAIN_NAME}/upload`,
    AUTH_SHARE_URL: `https://api.${DOMAIN_NAME}/auth-share`,
    SHARE_URL: `https://api.${DOMAIN_NAME}/share`,
    SHARING_POST_URL: `https://api.${DOMAIN_NAME}/share-file`,
    GET_SHARING_FILES_URL: `https://api.${DOMAIN_NAME}/share-files`,
    GET_DOWNLOAD_FILE_AVAILABLE_URL: `https://api.${DOMAIN_NAME}/download-available`,
    DOWNLOAD_URL: `https://api.${DOMAIN_NAME}/download`
}

export default config