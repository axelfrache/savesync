package sftp

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"

	"github.com/axelfrache/savesync/internal/domain"
)

// Backend implements domain.Backend for SFTP storage
type Backend struct {
	sshClient  *ssh.Client
	sftpClient *sftp.Client
	basePath   string
}

// New creates a new SFTP backend
func New() *Backend {
	return &Backend{}
}

// Init initializes the backend with configuration
func (b *Backend) Init(cfg map[string]string) error {
	host, ok := cfg["host"]
	if !ok {
		return fmt.Errorf("host is required in config")
	}

	port := cfg["port"]
	if port == "" {
		port = "22"
	}

	user, ok := cfg["user"]
	if !ok {
		return fmt.Errorf("user is required in config")
	}

	basePath, ok := cfg["path"]
	if !ok {
		return fmt.Errorf("path is required in config")
	}
	b.basePath = basePath

	// Setup SSH client config
	var authMethods []ssh.AuthMethod

	// Try password auth
	if password, ok := cfg["password"]; ok {
		authMethods = append(authMethods, ssh.Password(password))
	}

	// Try key auth
	if keyPath, ok := cfg["key_path"]; ok {
		key, err := os.ReadFile(keyPath)
		if err != nil {
			return fmt.Errorf("failed to read SSH key: %w", err)
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return fmt.Errorf("failed to parse SSH key: %w", err)
		}

		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	if len(authMethods) == 0 {
		return fmt.Errorf("no authentication method provided (password or key_path required)")
	}

	sshConfig := &ssh.ClientConfig{
		User:            user,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: Add proper host key verification
	}

	// Connect to SSH server
	addr := fmt.Sprintf("%s:%s", host, port)
	sshClient, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SSH server: %w", err)
	}
	b.sshClient = sshClient

	// Create SFTP client
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		sshClient.Close()
		return fmt.Errorf("failed to create SFTP client: %w", err)
	}
	b.sftpClient = sftpClient

	// Create base directories
	if err := b.sftpClient.MkdirAll(filepath.Join(b.basePath, "chunks")); err != nil {
		return fmt.Errorf("failed to create chunks directory: %w", err)
	}

	if err := b.sftpClient.MkdirAll(filepath.Join(b.basePath, "manifests")); err != nil {
		return fmt.Errorf("failed to create manifests directory: %w", err)
	}

	return nil
}

// StoreChunk stores a chunk via SFTP
func (b *Backend) StoreChunk(ctx context.Context, hash string, data []byte) error {
	if len(hash) < 4 {
		return fmt.Errorf("invalid hash: too short")
	}

	// Create nested directory structure
	dir := filepath.Join(b.basePath, "chunks", hash[:2], hash[2:4])
	if err := b.sftpClient.MkdirAll(dir); err != nil {
		return fmt.Errorf("failed to create chunk directory: %w", err)
	}

	chunkPath := filepath.Join(dir, hash)

	// Check if chunk already exists (deduplication)
	if _, err := b.sftpClient.Stat(chunkPath); err == nil {
		return nil // Chunk already exists
	}

	file, err := b.sftpClient.Create(chunkPath)
	if err != nil {
		return fmt.Errorf("failed to create chunk file: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("failed to write chunk: %w", err)
	}

	return nil
}

// LoadChunk loads a chunk via SFTP
func (b *Backend) LoadChunk(ctx context.Context, hash string) ([]byte, error) {
	if len(hash) < 4 {
		return nil, fmt.Errorf("invalid hash: too short")
	}

	chunkPath := filepath.Join(b.basePath, "chunks", hash[:2], hash[2:4], hash)

	file, err := b.sftpClient.Open(chunkPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to open chunk: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read chunk: %w", err)
	}

	return data, nil
}

// DeleteChunk deletes a chunk via SFTP
func (b *Backend) DeleteChunk(ctx context.Context, hash string) error {
	if len(hash) < 4 {
		return fmt.Errorf("invalid hash: too short")
	}

	chunkPath := filepath.Join(b.basePath, "chunks", hash[:2], hash[2:4], hash)

	if err := b.sftpClient.Remove(chunkPath); err != nil {
		if os.IsNotExist(err) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("failed to delete chunk: %w", err)
	}

	return nil
}

// ChunkExists checks if a chunk exists via SFTP
func (b *Backend) ChunkExists(ctx context.Context, hash string) (bool, error) {
	if len(hash) < 4 {
		return false, fmt.Errorf("invalid hash: too short")
	}

	chunkPath := filepath.Join(b.basePath, "chunks", hash[:2], hash[2:4], hash)

	_, err := b.sftpClient.Stat(chunkPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check chunk: %w", err)
	}

	return true, nil
}

// StoreManifest stores a snapshot manifest via SFTP
func (b *Backend) StoreManifest(ctx context.Context, snapshotID string, manifest []byte) error {
	manifestPath := filepath.Join(b.basePath, "manifests", snapshotID+".json")

	file, err := b.sftpClient.Create(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to create manifest file: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(manifest); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}

	return nil
}

// LoadManifest loads a snapshot manifest via SFTP
func (b *Backend) LoadManifest(ctx context.Context, snapshotID string) ([]byte, error) {
	manifestPath := filepath.Join(b.basePath, "manifests", snapshotID+".json")

	file, err := b.sftpClient.Open(manifestPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to open manifest: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest: %w", err)
	}

	return data, nil
}

// DeleteManifest deletes a snapshot manifest via SFTP
func (b *Backend) DeleteManifest(ctx context.Context, snapshotID string) error {
	manifestPath := filepath.Join(b.basePath, "manifests", snapshotID+".json")

	if err := b.sftpClient.Remove(manifestPath); err != nil {
		if os.IsNotExist(err) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("failed to delete manifest: %w", err)
	}

	return nil
}

// Close closes the SFTP and SSH connections
func (b *Backend) Close() error {
	if b.sftpClient != nil {
		b.sftpClient.Close()
	}
	if b.sshClient != nil {
		b.sshClient.Close()
	}
	return nil
}
