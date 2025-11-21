package backupservice

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

const (
	// DefaultChunkSize is the default chunk size (4MB)
	DefaultChunkSize = 4 * 1024 * 1024
)

// Chunker handles file chunking
type Chunker struct {
	chunkSize int
}

// NewChunker creates a new chunker
func NewChunker(chunkSize int) *Chunker {
	if chunkSize <= 0 {
		chunkSize = DefaultChunkSize
	}
	return &Chunker{chunkSize: chunkSize}
}

// ChunkFile chunks a file and returns chunk hashes and data
func (c *Chunker) ChunkFile(filePath string) ([]ChunkInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var chunks []ChunkInfo
	buffer := make([]byte, c.chunkSize)

	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}
		if n == 0 {
			break
		}

		// Calculate hash for this chunk
		chunkData := buffer[:n]
		hash := sha256.Sum256(chunkData)
		hashStr := hex.EncodeToString(hash[:])

		chunks = append(chunks, ChunkInfo{
			Hash: hashStr,
			Size: int64(n),
			Data: append([]byte(nil), chunkData...), // Copy data
		})
	}

	return chunks, nil
}

// ChunkInfo contains information about a chunk
type ChunkInfo struct {
	Hash string
	Size int64
	Data []byte
}

// HashFile calculates the SHA256 hash of an entire file
func HashFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to hash file: %w", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
