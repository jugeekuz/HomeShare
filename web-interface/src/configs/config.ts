const config = {
    MAX_CHUNK_SIZE_MB: 5,
    MAX_CHUNK_RETRIES: 5,
    MAX_CONCURRENT_CHUNKS: 6,
    MAX_FILE_SIZE_MB: (1024 * 5) * 1024 * 1024,
    OTP_LENGTH: 6,
    BASE_URL: 'https://api.homeshare.pro',
    UPLOAD_URL: 'https://api.homeshare.pro/upload',
    SHARE_URL: 'https://api.homeshare.pro/share',
    SHARING_POST_URL: 'https://api.homeshare.pro/share-file',
    GET_SHARING_FILES_URL: 'https://api.homeshare.pro/share-files',
    GET_DOWNLOAD_FILE_AVAILABLE_URL: 'https://api.homeshare.pro/download-available',
    DOWNLOAD_URL: 'https://api.homeshare.pro/download'
}

export default config