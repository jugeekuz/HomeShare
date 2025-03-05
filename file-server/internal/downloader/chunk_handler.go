package downloader
import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"sync"
)

func (c *DefaultFileChunker) ChunkFile(filePath string) ([]ChunkMeta, error) {
	var chunks []ChunkMeta

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buffer := make([]byte, c.chunkSize)
	index := 0

	for {
		bytesRead, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if bytesRead == 0 {
			break
		}

		hash := md5.Sum(buffer[:bytesRead])
		hashString := hex.EncodeToString(hash[:])

		chunkFileName := fmt.Sprintf("%s.chunk.%d", filePath, index)

		chunkFile, err := os.Create(chunkFileName)
		if err != nil {
			return nil, err
		}
		_, err = chunkFile.Write(buffer[:bytesRead])
		if err != nil {
			return nil, err
		}

		chunks = append(chunks, ChunkMeta{FileName: chunkFileName, MD5Hash: hashString, Index: index})

		chunkFile.Close()

		index++
	}

	return chunks, nil
}

func (c *DefaultFileChunker) ChunklargeFile(filePath string) ([]ChunkMeta, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var chunks []ChunkMeta

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	numChunks := int(fileInfo.Size() / int64(c.chunkSize))
	if fileInfo.Size()%int64(c.chunkSize) != 0 {
		numChunks++
	}

	chunkChan := make(chan ChunkMeta, numChunks)
	errChan := make(chan error, numChunks)
	indexChan := make(chan int, numChunks)

	for i := 0; i < numChunks; i++ {
		indexChan <- i
	}
	close(indexChan)

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range indexChan {
				offset := int64(index) * int64(c.chunkSize)
				buffer := make([]byte, c.chunkSize)

				if _,err := file.Seek(offset, 0); err != nil {
					errChan <- err
					return
				}

				bytesRead, err := file.Read(buffer)
				if err != nil && err != io.EOF {
					errChan <- err
					return
				}

				if bytesRead == 0 {
					return
				};
				
				hash := md5.Sum(buffer[:bytesRead])
				hashString := hex.EncodeToString(hash[:])

				chunkFileName := fmt.Sprintf("%s.chunk.%d", filePath, index)

				chunkFile, err := os.Create(chunkFileName)
				if err != nil {
					errChan <- err
					return
				}
				_, err = chunkFile.Write(buffer[:bytesRead])
				if err != nil {
					errChan <- err
					return
				}

				chunk := ChunkMeta{
					FileName: chunkFileName,
					MD5Hash:  hashString,
					Index:    index,
				}
				mu.Lock()
				chunks = append(chunks, chunk)
				mu.Unlock()

				chunkFile.Close()

				chunkChan <- chunk
			}
		}()
	}

	go func() {
		wg.Wait()
		close(chunkChan)
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	return chunks, nil
}

