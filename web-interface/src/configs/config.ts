const config = {
    MAX_CHUNK_SIZE_MB: 5,
    MAX_CHUNK_RETRIES: 5,
    MAX_CONCURRENT_CHUNKS: 6,
    MAX_FILE_SIZE_MB: (1024 * 5) * 1024 * 1024,
    BASE_URL: 'https://kuza.gr',
    UPLOAD_URL: 'https://kuza.gr/upload',
    SHARE_URL: 'https://kuza.gr/share',
    SHARING_POST_URL: 'https://kuza.gr/share-file',
    GET_SHARING_FILES_URL: 'https://kuza.gr/share-files',
    GET_DOWNLOAD_FILE_AVAILABLE_URL: 'https://kuza.gr/download-available',
    DOWNLOAD_URL: 'https://kuza.gr/download'
}

export default config