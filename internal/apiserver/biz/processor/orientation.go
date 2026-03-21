package processor

import (
	"image"
	"os"

	"github.com/disintegration/imaging"
	exif "github.com/dsoprea/go-exif/v3"
)

// readOrientationFromFile reads the EXIF orientation tag from any image file format.
// Uses brute-force search so it works with JPEG, PNG (eXIf chunk), TIFF, etc.
// Returns 1 (normal) if no orientation tag is found or on error.
func readOrientationFromFile(filePath string) int {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return 1
	}
	return readOrientationFromBytes(data)
}

// readOrientationFromBytes reads the EXIF orientation tag from raw image bytes.
func readOrientationFromBytes(data []byte) int {
	rawExif, err := exif.SearchAndExtractExif(data)
	if err != nil || rawExif == nil {
		return 1
	}

	tags, _, err := exif.GetFlatExifDataUniversalSearch(rawExif, nil, true)
	if err != nil {
		return 1
	}

	for _, tag := range tags {
		if tag.TagId == 0x0112 { // Orientation tag
			if vals, ok := tag.Value.([]uint16); ok && len(vals) > 0 {
				orient := int(vals[0])
				if orient >= 1 && orient <= 8 {
					return orient
				}
			}
		}
	}

	return 1
}

// fixImageOrientation applies the EXIF orientation transform to an image.
// Returns the corrected image based on the orientation value (1-8).
func fixImageOrientation(img image.Image, orient int) image.Image {
	switch orient {
	case 1:
		// Normal, no transform needed
		return img
	case 2:
		// Flip horizontal
		return imaging.FlipH(img)
	case 3:
		// Rotate 180
		return imaging.Rotate180(img)
	case 4:
		// Flip vertical
		return imaging.FlipV(img)
	case 5:
		// Transpose (flip horizontal + rotate 270)
		return imaging.Transpose(img)
	case 6:
		// Rotate 270 (90 CW)
		return imaging.Rotate270(img)
	case 7:
		// Transverse (flip horizontal + rotate 90)
		return imaging.Transverse(img)
	case 8:
		// Rotate 90 (90 CCW)
		return imaging.Rotate90(img)
	default:
		return img
	}
}

// orientedDimensions returns the width and height after applying the EXIF orientation.
// Orientations 5-8 swap width and height.
func orientedDimensions(w, h, orient int) (int, int) {
	if orient >= 5 && orient <= 8 {
		return h, w
	}
	return w, h
}
