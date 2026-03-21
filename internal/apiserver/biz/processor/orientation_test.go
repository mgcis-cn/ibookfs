package processor

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// buildJPEGWithOrientation creates a JPEG file with a specific EXIF orientation tag.
func buildJPEGWithOrientation(w, h int, orient uint16) []byte {
	// Create a simple image
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	// Draw a red top-left corner to visually verify orientation
	for y := 0; y < h/4; y++ {
		for x := 0; x < w/4; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}

	// Build minimal EXIF with orientation
	exifData := buildMinimalEXIF(orient)

	var buf bytes.Buffer
	// SOI
	buf.Write([]byte{0xFF, 0xD8})
	// APP1 marker
	buf.Write([]byte{0xFF, 0xE1})
	// APP1 size (2 bytes for size + exif data)
	size := uint16(len(exifData) + 2)
	binary.Write(&buf, binary.BigEndian, size)
	buf.Write(exifData)

	// Encode the rest as JPEG (we need to strip SOI from jpeg.Encode output)
	var jpegBuf bytes.Buffer
	jpeg.Encode(&jpegBuf, img, &jpeg.Options{Quality: 90})
	jpegBytes := jpegBuf.Bytes()
	// Skip SOI (first 2 bytes) from encoded JPEG
	buf.Write(jpegBytes[2:])

	return buf.Bytes()
}

// buildMinimalEXIF creates minimal EXIF data with just the orientation tag.
func buildMinimalEXIF(orient uint16) []byte {
	var buf bytes.Buffer
	// Exif header: "Exif\0\0"
	buf.Write([]byte{0x45, 0x78, 0x69, 0x66, 0x00, 0x00})
	// TIFF header: little-endian
	buf.Write([]byte{0x49, 0x49})                           // II = little-endian
	binary.Write(&buf, binary.LittleEndian, uint16(0x002A)) // magic
	binary.Write(&buf, binary.LittleEndian, uint32(0x08))   // offset to IFD0

	// IFD0: 1 tag
	binary.Write(&buf, binary.LittleEndian, uint16(1)) // number of tags

	// Orientation tag (0x0112), SHORT type (3), count=1
	binary.Write(&buf, binary.LittleEndian, uint16(0x0112)) // tag
	binary.Write(&buf, binary.LittleEndian, uint16(3))      // type = SHORT
	binary.Write(&buf, binary.LittleEndian, uint32(1))      // count
	binary.Write(&buf, binary.LittleEndian, uint16(orient)) // value
	binary.Write(&buf, binary.LittleEndian, uint16(0))      // padding

	// Next IFD offset = 0 (no more IFDs)
	binary.Write(&buf, binary.LittleEndian, uint32(0))

	return buf.Bytes()
}

// buildPNGWithEXIFOrientation creates a PNG file with an eXIf chunk containing orientation.
func buildPNGWithEXIFOrientation(w, h int, orient uint16) []byte {
	// Create a simple image
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h/4; y++ {
		for x := 0; x < w/4; x++ {
			img.Set(x, y, color.RGBA{0, 0, 255, 255})
		}
	}

	// Encode as PNG
	var pngBuf bytes.Buffer
	png.Encode(&pngBuf, img)
	pngBytes := pngBuf.Bytes()

	// Build EXIF data (TIFF format, without "Exif\0\0" header for PNG eXIf chunk)
	exifPayload := buildTIFFWithOrientation(orient)

	// Insert eXIf chunk before IEND
	// PNG structure: signature (8 bytes) + chunks... + IEND
	// Find IEND position
	iendPos := bytes.Index(pngBytes, []byte("IEND"))
	if iendPos < 4 {
		return pngBytes // fallback
	}
	iendStart := iendPos - 4 // 4 bytes for length field before "IEND"

	var result bytes.Buffer
	result.Write(pngBytes[:iendStart])

	// Write eXIf chunk
	// Chunk: length(4) + type(4) + data + CRC(4)
	chunkType := []byte("eXIf")
	binary.Write(&result, binary.BigEndian, uint32(len(exifPayload)))
	result.Write(chunkType)
	result.Write(exifPayload)
	// CRC (simplified - just write zeros, PNG decoders may warn but go-exif brute-force search doesn't care)
	crc := pngCRC(chunkType, exifPayload)
	binary.Write(&result, binary.BigEndian, crc)

	// Write IEND
	result.Write(pngBytes[iendStart:])

	return result.Bytes()
}

// buildTIFFWithOrientation creates a minimal TIFF structure with orientation tag.
func buildTIFFWithOrientation(orient uint16) []byte {
	var buf bytes.Buffer
	// TIFF header: little-endian
	buf.Write([]byte{0x49, 0x49})                           // II = little-endian
	binary.Write(&buf, binary.LittleEndian, uint16(0x002A)) // magic
	binary.Write(&buf, binary.LittleEndian, uint32(0x08))   // offset to IFD0

	// IFD0: 1 tag
	binary.Write(&buf, binary.LittleEndian, uint16(1))

	// Orientation tag
	binary.Write(&buf, binary.LittleEndian, uint16(0x0112))
	binary.Write(&buf, binary.LittleEndian, uint16(3))      // SHORT
	binary.Write(&buf, binary.LittleEndian, uint32(1))      // count
	binary.Write(&buf, binary.LittleEndian, uint16(orient)) // value
	binary.Write(&buf, binary.LittleEndian, uint16(0))      // padding

	// No next IFD
	binary.Write(&buf, binary.LittleEndian, uint32(0))

	return buf.Bytes()
}

// pngCRC computes the CRC for a PNG chunk.
func pngCRC(chunkType, data []byte) uint32 {
	// CRC-32 table
	var table [256]uint32
	for i := 0; i < 256; i++ {
		c := uint32(i)
		for j := 0; j < 8; j++ {
			if c&1 != 0 {
				c = 0xEDB88320 ^ (c >> 1)
			} else {
				c >>= 1
			}
		}
		table[i] = c
	}

	crc := uint32(0xFFFFFFFF)
	for _, b := range chunkType {
		crc = table[(crc^uint32(b))&0xFF] ^ (crc >> 8)
	}
	for _, b := range data {
		crc = table[(crc^uint32(b))&0xFF] ^ (crc >> 8)
	}
	return crc ^ 0xFFFFFFFF
}

func TestReadOrientationFromBytes_JPEG(t *testing.T) {
	tests := []struct {
		name   string
		orient uint16
		want   int
	}{
		{"normal", 1, 1},
		{"rotate90", 8, 8},
		{"rotate180", 3, 3},
		{"rotate270", 6, 6},
		{"flipH", 2, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := buildJPEGWithOrientation(200, 100, tt.orient)
			got := readOrientationFromBytes(data)
			if got != tt.want {
				t.Errorf("readOrientationFromBytes() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestReadOrientationFromBytes_PNG(t *testing.T) {
	tests := []struct {
		name   string
		orient uint16
		want   int
	}{
		{"normal", 1, 1},
		{"rotate90", 8, 8},
		{"rotate180", 3, 3},
		{"rotate270", 6, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := buildPNGWithEXIFOrientation(200, 100, tt.orient)
			got := readOrientationFromBytes(data)
			if got != tt.want {
				t.Errorf("readOrientationFromBytes() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestReadOrientationFromBytes_NoEXIF(t *testing.T) {
	// Plain PNG without EXIF
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	var buf bytes.Buffer
	png.Encode(&buf, img)

	got := readOrientationFromBytes(buf.Bytes())
	if got != 1 {
		t.Errorf("readOrientationFromBytes(plain PNG) = %d, want 1", got)
	}
}

func TestOrientedDimensions(t *testing.T) {
	// Orientations 1-4: no swap
	w, h := orientedDimensions(200, 100, 1)
	if w != 200 || h != 100 {
		t.Errorf("orient=1: got %dx%d, want 200x100", w, h)
	}

	// Orientations 5-8: swap
	w, h = orientedDimensions(200, 100, 6)
	if w != 100 || h != 200 {
		t.Errorf("orient=6: got %dx%d, want 100x200", w, h)
	}

	w, h = orientedDimensions(200, 100, 8)
	if w != 100 || h != 200 {
		t.Errorf("orient=8: got %dx%d, want 100x200", w, h)
	}
}

func TestReadOrientationFromFile(t *testing.T) {
	// Write a JPEG with orientation 6 to a temp file
	data := buildJPEGWithOrientation(200, 100, 6)
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.jpg")
	os.WriteFile(path, data, 0644)

	got := readOrientationFromFile(path)
	if got != 6 {
		t.Errorf("readOrientationFromFile() = %d, want 6", got)
	}

	// Write a PNG with orientation 8
	pngData := buildPNGWithEXIFOrientation(200, 100, 8)
	pngPath := filepath.Join(tmpDir, "test.png")
	os.WriteFile(pngPath, pngData, 0644)

	got = readOrientationFromFile(pngPath)
	if got != 8 {
		t.Errorf("readOrientationFromFile(png) = %d, want 8", got)
	}
}
