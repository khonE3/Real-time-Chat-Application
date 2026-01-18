package handler

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/khonE3/chat-backend/internal/model"
	"github.com/khonE3/chat-backend/internal/repository"
)

// Maximum file size (10MB)
const MaxFileSize = 10 * 1024 * 1024

// Allowed file types
var AllowedMimeTypes = map[string]bool{
	"image/jpeg":      true,
	"image/png":       true,
	"image/gif":       true,
	"image/webp":      true,
	"image/svg+xml":   true,
	"video/mp4":       true,
	"video/webm":      true,
	"application/pdf": true,
	"text/plain":      true,
}

// FileHandler handles file upload operations
type FileHandler struct {
	fileRepo   *repository.GormFileRepository
	uploadPath string
	baseURL    string
}

// NewFileHandler creates a new file handler
func NewFileHandler(fileRepo *repository.GormFileRepository, uploadPath, baseURL string) *FileHandler {
	// Create upload directory if not exists
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		log.Printf("Warning: Could not create upload directory: %v", err)
	}

	return &FileHandler{
		fileRepo:   fileRepo,
		uploadPath: uploadPath,
		baseURL:    baseURL,
	}
}

// Upload handles file upload
func (h *FileHandler) Upload(c *fiber.Ctx) error {
	// Get user ID from query
	userIDStr := c.Query("userId")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userId is required",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid userId",
		})
	}

	// Get the file from form
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	// Check file size
	if file.Size > MaxFileSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("File too large. Maximum size is %d MB", MaxFileSize/1024/1024),
		})
	}

	// Check mime type
	mimeType := file.Header.Get("Content-Type")
	if !AllowedMimeTypes[mimeType] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File type not allowed",
		})
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().UnixNano(), ext)

	// Save file
	filePath := filepath.Join(h.uploadPath, filename)
	if err := c.SaveFile(file, filePath); err != nil {
		log.Printf("Failed to save file: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save file",
		})
	}

	// Create file record
	fileRecord := &model.File{
		UserID:       userID,
		Filename:     filename,
		OriginalName: file.Filename,
		MimeType:     mimeType,
		Size:         file.Size,
		URL:          fmt.Sprintf("%s/uploads/%s", h.baseURL, filename),
	}

	if err := h.fileRepo.Create(c.Context(), fileRecord); err != nil {
		// Clean up file if database insert fails
		os.Remove(filePath)
		log.Printf("Failed to create file record: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create file record",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.FileUploadResponse{
		ID:           fileRecord.ID.String(),
		Filename:     fileRecord.Filename,
		OriginalName: fileRecord.OriginalName,
		MimeType:     fileRecord.MimeType,
		Size:         fileRecord.Size,
		URL:          fileRecord.URL,
	})
}

// UploadMultiple handles multiple file uploads
func (h *FileHandler) UploadMultiple(c *fiber.Ctx) error {
	// Get user ID from query
	userIDStr := c.Query("userId")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userId is required",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid userId",
		})
	}

	// Get files from form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No files uploaded",
		})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No files uploaded",
		})
	}

	var results []model.FileUploadResponse

	for _, file := range files {
		// Check file size
		if file.Size > MaxFileSize {
			continue
		}

		// Check mime type
		mimeType := file.Header.Get("Content-Type")
		if !AllowedMimeTypes[mimeType] {
			continue
		}

		// Generate unique filename
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().UnixNano(), ext)

		// Save file
		filePath := filepath.Join(h.uploadPath, filename)
		if err := c.SaveFile(file, filePath); err != nil {
			log.Printf("Failed to save file: %v", err)
			continue
		}

		// Create file record
		fileRecord := &model.File{
			UserID:       userID,
			Filename:     filename,
			OriginalName: file.Filename,
			MimeType:     mimeType,
			Size:         file.Size,
			URL:          fmt.Sprintf("%s/uploads/%s", h.baseURL, filename),
		}

		if err := h.fileRepo.Create(c.Context(), fileRecord); err != nil {
			os.Remove(filePath)
			log.Printf("Failed to create file record: %v", err)
			continue
		}

		results = append(results, model.FileUploadResponse{
			ID:           fileRecord.ID.String(),
			Filename:     fileRecord.Filename,
			OriginalName: fileRecord.OriginalName,
			MimeType:     fileRecord.MimeType,
			Size:         fileRecord.Size,
			URL:          fileRecord.URL,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"files": results,
		"count": len(results),
	})
}

// GetFile retrieves file info
func (h *FileHandler) GetFile(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file ID",
		})
	}

	file, err := h.fileRepo.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "File not found",
		})
	}

	return c.JSON(file)
}

// ServeFile serves the actual file
func (h *FileHandler) ServeFile(c *fiber.Ctx) error {
	filename := c.Params("filename")

	// Prevent directory traversal
	if strings.Contains(filename, "..") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid filename",
		})
	}

	filePath := filepath.Join(h.uploadPath, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "File not found",
		})
	}

	// Get file from database for mime type
	file, err := h.fileRepo.GetByFilename(c.Context(), filename)
	if err == nil {
		c.Set("Content-Type", file.MimeType)
	}

	return c.SendFile(filePath)
}

// Delete deletes a file
func (h *FileHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file ID",
		})
	}

	// Get file info first
	file, err := h.fileRepo.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "File not found",
		})
	}

	// Delete from filesystem
	filePath := filepath.Join(h.uploadPath, file.Filename)
	if err := os.Remove(filePath); err != nil {
		log.Printf("Failed to delete file from filesystem: %v", err)
	}

	// Delete from database
	if err := h.fileRepo.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete file",
		})
	}

	return c.JSON(fiber.Map{
		"message": "File deleted successfully",
	})
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
