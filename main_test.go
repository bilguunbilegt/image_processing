package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func createTestImage(path string) error {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			img.Set(x, y, color.RGBA{uint8(x), uint8(y), 255, 255})
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return jpeg.Encode(file, img, nil)
}

func TestLoadImage(t *testing.T) {
	// Setup: create test images
	imagePaths := []string{"images/test_image1.jpeg", "images/test_image2.jpeg"}
	for _, path := range imagePaths {
		err := createTestImage(path)
		if err != nil {
			t.Fatalf("Failed to create test image: %v", err)
		}
	}

	// Add an invalid image path
	imagePaths = append(imagePaths, "images/invalid_image.jpeg")

	jobs := loadImage(imagePaths, true) // Test with goroutines
	for job := range jobs {
		if strings.Contains(job.InputPath, "invalid") {
			t.Errorf("Expected to skip invalid image: %s", job.InputPath)
		}
		if job.Image == nil {
			t.Errorf("Failed to load image: %s", job.InputPath)
		}
		expectedOutPath := strings.Replace(job.InputPath, "images/", "images/output/", 1)
		if job.OutPath != expectedOutPath {
			t.Errorf("Expected OutPath: %s, but got: %s", expectedOutPath, job.OutPath)
		}
	}

	jobs = loadImage(imagePaths, false) // Test without goroutines
	for job := range jobs {
		if strings.Contains(job.InputPath, "invalid") {
			t.Errorf("Expected to skip invalid image: %s", job.InputPath)
		}
		if job.Image == nil {
			t.Errorf("Failed to load image: %s", job.InputPath)
		}
		expectedOutPath := strings.Replace(job.InputPath, "images/", "images/output/", 1)
		if job.OutPath != expectedOutPath {
			t.Errorf("Expected OutPath: %s, but got: %s", expectedOutPath, job.OutPath)
		}
	}

	// Clean up
	for _, path := range imagePaths {
		os.Remove(path)
	}
}

func TestSaveImage(t *testing.T) {
	// Setup: create a dummy image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	outPath := "images/output/test_image.jpeg"
	job := Job{
		Image:   img,
		OutPath: outPath,
	}

	// Ensure the output directory exists
	outDir := filepath.Dir(outPath)
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		err := os.MkdirAll(outDir, os.ModePerm)
		if err != nil {
			t.Fatalf("Failed to create output directory: %v", err)
		}
	}

	// Create a channel and pass the job to it
	input := make(chan Job, 1)
	input <- job
	close(input)

	// Test saveImage with goroutines
	results := saveImage(input, true)
	for success := range results {
		if !success {
			t.Errorf("Failed to save image: %s", job.OutPath)
		} else {
			// Check if the file exists
			if _, err := os.Stat(job.OutPath); os.IsNotExist(err) {
				t.Errorf("Output file does not exist: %s", job.OutPath)
			}
		}
	}

	// Test saveImage without goroutines
	input = make(chan Job, 1)
	input <- job
	close(input)

	results = saveImage(input, false)
	for success := range results {
		if !success {
			t.Errorf("Failed to save image: %s", job.OutPath)
		} else {
			// Check if the file exists
			if _, err := os.Stat(job.OutPath); os.IsNotExist(err) {
				t.Errorf("Output file does not exist: %s", job.OutPath)
			}
		}
	}

	// Clean up
	os.Remove(job.OutPath)
}

func TestSaveImageInvalidOutputDir(t *testing.T) {
	// Setup: create a dummy image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	outPath := "invalid_output_dir/test_image.jpeg"
	job := Job{
		Image:   img,
		OutPath: outPath,
	}

	// Create a channel and pass the job to it
	input := make(chan Job, 1)
	input <- job
	close(input)

	// Test saveImage with goroutines
	results := saveImage(input, true)
	for success := range results {
		if success {
			t.Errorf("Expected to fail saving image to invalid directory: %s", job.OutPath)
		}
	}

	// Test saveImage without goroutines
	input = make(chan Job, 1)
	input <- job
	close(input)

	results = saveImage(input, false)
	for success := range results {
		if success {
			t.Errorf("Expected to fail saving image to invalid directory: %s", job.OutPath)
		}
	}
}
