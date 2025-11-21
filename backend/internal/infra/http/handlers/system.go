package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"sort"

	"go.uber.org/zap"
)

type SystemHandler struct {
	logger *zap.Logger
}

func NewSystemHandler(logger *zap.Logger) *SystemHandler {
	return &SystemHandler{
		logger: logger,
	}
}

type FileEntry struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	IsDir bool   `json:"is_dir"`
}

// ListFiles godoc
// @Summary Lister les fichiers et dossiers
// @Description Liste le contenu d'un répertoire sur le serveur
// @Tags system
// @Produce json
// @Param path query string false "Chemin du répertoire (défaut: racine)"
// @Success 200 {array} handlers.FileEntry
// @Failure 400 {object} handlers.ErrorInfo
// @Failure 500 {object} handlers.ErrorInfo
// @Router /system/files [get]
func (h *SystemHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		// Default to home directory or root
		home, err := os.UserHomeDir()
		if err != nil {
			path = "/"
		} else {
			path = home
		}
	}

	// Security check: prevent traversing up from root (though root is allowed)
	// In a real app, we might want to restrict this to specific allowed directories.
	// For this personal tool, we allow full access but ensure path is clean.
	path = filepath.Clean(path)

	entries, err := os.ReadDir(path)
	if err != nil {
		h.logger.Error("failed to read directory", zap.String("path", path), zap.Error(err))
		WriteError(w, http.StatusBadRequest, "Failed to read directory: "+err.Error())
		return
	}

	var files []FileEntry

	// Add parent directory entry if not at root
	if path != "/" {
		parent := filepath.Dir(path)
		files = append(files, FileEntry{
			Name:  "..",
			Path:  parent,
			IsDir: true,
		})
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		// Only show directories for source selection, or maybe all files?
		// User wants to select a "dossier" (directory) for source.
		// Let's show everything but mark directories.

		fullPath := filepath.Join(path, entry.Name())
		files = append(files, FileEntry{
			Name:  entry.Name(),
			Path:  fullPath,
			IsDir: info.IsDir(),
		})
	}

	// Sort: Directories first, then files
	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir != files[j].IsDir {
			return files[i].IsDir
		}
		return files[i].Name < files[j].Name
	})

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"current_path": path,
		"entries":      files,
	})
}
