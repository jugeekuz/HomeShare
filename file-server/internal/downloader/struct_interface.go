package downloader

type ChunkMeta struct {
	FileName string `json:"file_name"`
	MD5Hash  string `json:"md5_hash"`
	Index    int    `json:"index"`
}

type Config struct {
	ChunkSize int
	ServerURL string
}

type DefaultFileChunker struct {
	chunkSize int
}

type DefaultUploader struct {
	serverURL string
}

type DefaultMetadataManager struct{}

type FileChunker interface {
	ChunkFile(filePath string) ([]ChunkMeta, error)
	ChunklargeFile(filePath string) ([]ChunkMeta, error)
}

type Uploader interface {
	UploadChunk(chunk ChunkMeta) error
}

type MetadataManager interface {
	LoadMetadata(filePath string) (map[string]ChunkMeta, error)
	SaveMetadata(filePath string, metadata map[string]ChunkMeta) error
}