// Package processor provides image processing capabilities.
package processor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	blurhash "github.com/buckket/go-blurhash"
	"github.com/disintegration/imaging"
	"k8s.io/apimachinery/pkg/api/resource"
)

// Config holds processor configuration.
type Config struct {
	// BlurHashEnabled enables BlurHash generation.
	BlurHashEnabled bool

	// Variants defines the image variants to generate.
	Variants []VariantConfig

	// MaxFileSize is the maximum allowed file size (e.g. "50Mi").
	MaxFileSize resource.Quantity

	// AllowedTypes are the allowed MIME types.
	AllowedTypes []string
}

// VariantConfig defines a specific image variant configuration.
type VariantConfig struct {
	Name      string // "small", "medium", "original"
	MaxWidth  int
	MaxHeight int
}

// DefaultConfig returns the default processor configuration.
func DefaultConfig() Config {
	return Config{
		BlurHashEnabled: true,
		Variants: []VariantConfig{
			{Name: "small", MaxWidth: 400, MaxHeight: 400},
			{Name: "medium", MaxWidth: 1200, MaxHeight: 1200},
		},
		MaxFileSize: resource.MustParse("50Mi"),
		AllowedTypes: []string{
			"image/jpeg",
			"image/png",
			"image/gif",
			"image/webp",
		},
	}
}

// Processor handles image processing operations.
type Processor struct {
	config Config
}

// NewProcessor creates a new image processor.
func NewProcessor(config Config) *Processor {
	return &Processor{config: config}
}

// NewDefaultProcessor creates a processor with default configuration.
func NewDefaultProcessor() *Processor {
	return &Processor{config: DefaultConfig()}
}

// ImageInfo holds information about an uploaded image.
type ImageInfo struct {
	Format   string
	Width    int
	Height   int
	Size     int64
	MimeType string
}

// ProcessResult holds the result of image processing.
type ProcessResult struct {
	OriginalPath string
	Blurhash     string
	Variants     []VariantResult
}

// VariantResult holds information about a generated variant.
type VariantResult struct {
	Name     string
	Path     string
	Width    int
	Height   int
	FileSize int64
}

// Validate reads the image and validates it against the configuration.
func (p *Processor) Validate(ctx context.Context, reader io.Reader, mimeType string, size int64) (*ImageInfo, error) {
	// Check file size
	if maxSize := p.config.MaxFileSize.Value(); maxSize > 0 && size > maxSize {
		return nil, fmt.Errorf("file size %d exceeds maximum %d", size, maxSize)
	}

	// Check MIME type
	if !p.isAllowedType(mimeType) {
		return nil, fmt.Errorf("MIME type %s is not allowed", mimeType)
	}

	// Read image to decode dimensions
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, reader); err != nil {
		return nil, fmt.Errorf("failed to read image: %w", err)
	}

	img, format, err := image.Decode(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	bounds := img.Bounds()

	return &ImageInfo{
		Format:   format,
		Width:    bounds.Dx(),
		Height:   bounds.Dy(),
		Size:     size,
		MimeType: mimeType,
	}, nil
}

// GenerateBlurHash generates a BlurHash for the given image.
func (p *Processor) GenerateBlurHash(imgPath string) (string, error) {
	if !p.config.BlurHashEnabled {
		return "", nil
	}

	file, err := os.Open(imgPath)
	if err != nil {
		return "", fmt.Errorf("failed to open image: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	// BlurHash components: X components (horizontal detail), Y components (vertical detail)
	// 4x3 is a good balance between size and quality
	hash, err := blurhash.Encode(4, 3, img)
	if err != nil {
		return "", fmt.Errorf("failed to generate blurhash: %w", err)
	}

	return hash, nil
}

// GenerateVariants creates resized variants of the source image.
func (p *Processor) GenerateVariants(ctx context.Context, sourcePath string, outputDir string) ([]VariantResult, error) {
	sourceImg, err := imaging.Open(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open source image: %w", err)
	}

	bounds := sourceImg.Bounds()
	originalW := bounds.Dx()
	originalH := bounds.Dy()

	var results []VariantResult

	for _, variantCfg := range p.config.Variants {
		// Calculate target dimensions
		targetW, targetH := p.calculateDimensions(originalW, originalH, variantCfg.MaxWidth, variantCfg.MaxHeight)

		// Skip if image is already smaller than target
		if originalW <= targetW && originalH <= targetH {
			continue
		}

		// Create output path
		ext := filepath.Ext(sourcePath)
		variantPath := filepath.Join(outputDir, variantCfg.Name, filepath.Base(sourcePath))

		// Ensure output directory exists
		if err := os.MkdirAll(filepath.Dir(variantPath), 0755); err != nil {
			return nil, fmt.Errorf("failed to create output directory: %w", err)
		}

		// Resize image
		resized := imaging.Resize(sourceImg, targetW, targetH, imaging.Lanczos)

		// Save based on original format
		if err := p.saveImage(resized, variantPath, ext); err != nil {
			return nil, fmt.Errorf("failed to save variant %s: %w", variantCfg.Name, err)
		}

		// Get file size
		fileInfo, _ := os.Stat(variantPath)

		results = append(results, VariantResult{
			Name:     variantCfg.Name,
			Path:     variantPath,
			Width:    targetW,
			Height:   targetH,
			FileSize: fileInfo.Size(),
		})
	}

	return results, nil
}

// Process performs all processing steps on an uploaded image.
func (p *Processor) Process(ctx context.Context, sourcePath string, outputDir string) (*ProcessResult, error) {
	// Generate BlurHash
	blurhash, err := p.GenerateBlurHash(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to generate blurhash: %w", err)
	}

	// Generate variants
	variants, err := p.GenerateVariants(ctx, sourcePath, outputDir)
	if err != nil {
		return nil, fmt.Errorf("failed to generate variants: %w", err)
	}

	return &ProcessResult{
		OriginalPath: sourcePath,
		Blurhash:     blurhash,
		Variants:     variants,
	}, nil
}

// SaveOriginal saves the uploaded image to the specified path.
func (p *Processor) SaveOriginal(reader io.Reader, destPath string) error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	// Create file
	file, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy content
	if _, err := io.Copy(file, reader); err != nil {
		return err
	}

	return nil
}

// calculateDimensions calculates the target dimensions maintaining aspect ratio.
func (p *Processor) calculateDimensions(originalW, originalH, maxWidth, maxHeight int) (int, int) {
	if originalW <= maxWidth && originalH <= maxHeight {
		return originalW, originalH
	}

	ratioW := float64(maxWidth) / float64(originalW)
	ratioH := float64(maxHeight) / float64(originalH)
	ratio := min(ratioW, ratioH)

	return int(float64(originalW) * ratio), int(float64(originalH) * ratio)
}

// saveImage saves an image to disk in the specified format.
func (p *Processor) saveImage(img *image.NRGBA, path, ext string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		return jpeg.Encode(file, img, &jpeg.Options{Quality: 85})
	case ".png":
		return png.Encode(file, img)
	default:
		return jpeg.Encode(file, img, &jpeg.Options{Quality: 85})
	}
}

// isAllowedType checks if the MIME type is allowed.
func (p *Processor) isAllowedType(mimeType string) bool {
	for _, allowed := range p.config.AllowedTypes {
		if mimeType == allowed {
			return true
		}
	}
	return false
}

// min returns the minimum of two floats.
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// ValidateFileExtension checks if the file extension is valid.
func ValidateFileExtension(filename string) error {
	ext := strings.ToLower(filepath.Ext(filename))
	validExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	if !validExts[ext] {
		return errors.New("invalid file extension")
	}
	return nil
}

// GetMIMEType returns the MIME type for a given file extension.
func GetMIMEType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	mimeTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
	}

	if mt, ok := mimeTypes[ext]; ok {
		return mt
	}
	return "application/octet-stream"
}
