package main

import (
	"fmt"
	imageprocessing "goroutines_pipeline/image_processing"
	"image"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Job struct {
	InputPath string
	Image     image.Image
	OutPath   string
}

func loadImage(paths []string, isGoroutines bool) <-chan Job {
	out := make(chan Job)
	if isGoroutines {
		go func() {
			// For each input path create a job and add it to
			// the out channel
			for _, p := range paths {
				// *********************************************************************************************
				// My part
				// Checking if the input file exists ie if the image folder is empty, we will throw an error
				// *********************************************************************************************
				if _, err := os.Stat(p); os.IsNotExist(err) {
					log.Printf("Input file does not exist in the path: %s", p)
					continue
				}
				job := Job{InputPath: p, OutPath: strings.Replace(p, "images/", "images/output/", 1)}
				job.Image = imageprocessing.ReadImage(p)
				// *********************************************************************************************
				// My part
				// Checking if we were able to read the image. If not, throw an error.
				// *********************************************************************************************
				if job.Image == nil {
					log.Printf("Failed to read image file: %s", p)
					continue
				}
				out <- job
			}
			close(out)
		}()
	} else {
		go func() { // Use a goroutine to avoid deadlock
			for _, p := range paths {
				// *********************************************************************************************
				// My part
				// Checking if the input file exists ie if the image folder is empty, we will throw an error
				// *********************************************************************************************
				if _, err := os.Stat(p); os.IsNotExist(err) {
					log.Printf("Input file does not exist in the path: %s", p)
					continue
				}
				job := Job{InputPath: p, OutPath: strings.Replace(p, "images/", "images/output/", 1)}
				job.Image = imageprocessing.ReadImage(p)
				// *********************************************************************************************
				// My part
				// Checking if we were able to read the image. If not, throw an error.
				// *********************************************************************************************
				if job.Image == nil {
					log.Printf("Failed to read image file: %s", p)
					continue
				}
				out <- job
			}
			close(out)
		}()
	}
	return out
}

func resize(input <-chan Job, isGoroutines bool) <-chan Job {
	out := make(chan Job)
	if isGoroutines {
		go func() {
			// For each input job, create a new job after resize and add it to
			// the out channel
			for job := range input { // Read from the channel
				job.Image = imageprocessing.Resize(job.Image)
				out <- job
			}
			close(out)
		}()
	} else {
		go func() { // Use a goroutine to avoid deadlock
			for job := range input { // Read from the channel
				job.Image = imageprocessing.Resize(job.Image)
				out <- job
			}
			close(out)
		}()
	}
	return out
}

func convertToGrayscale(input <-chan Job, isGoroutines bool) <-chan Job {
	out := make(chan Job)
	if isGoroutines {
		go func() {
			for job := range input { // Read from the channel
				job.Image = imageprocessing.Grayscale(job.Image)
				out <- job
			}
			close(out)
		}()
	} else {
		go func() { // Use a goroutine to avoid deadlock
			for job := range input { // Read from the channel
				job.Image = imageprocessing.Grayscale(job.Image)
				out <- job
			}
			close(out)
		}()
	}
	return out
}

// trying different color conversion and it works
// replace channel3 := convertToGrayscale(channel2, isGoroutines) with channel3 := convertToCustomColorConversion(channel2, isGoroutines)
func convertToCustomColorConversion(input <-chan Job, isGoroutines bool) <-chan Job {
	out := make(chan Job)
	if isGoroutines {
		go func() {
			for job := range input { // Read from the channel
				job.Image = imageprocessing.CustomColorConversion(job.Image)
				out <- job
			}
			close(out)
		}()
	} else {
		go func() { // Use a goroutine to avoid deadlock
			for job := range input { // Read from the channel
				job.Image = imageprocessing.CustomColorConversion(job.Image)
				out <- job
			}
			close(out)
		}()
	}
	return out
}

func saveImage(input <-chan Job, isGoroutines bool) <-chan bool {
	out := make(chan bool)
	if isGoroutines {
		go func() {
			for job := range input { // Read from the channel
				// *********************************************************************************************
				// My part
				// Check if the output directory exists
				// *********************************************************************************************
				outputDirectory := strings.Replace(job.OutPath, "/"+filepath.Base(job.OutPath), "", 1)
				if _, err := os.Stat(outputDirectory); os.IsNotExist(err) {
					log.Printf("Output directory does not exist: %s", outputDirectory)
					out <- false
					continue
				}
				imageprocessing.WriteImage(job.OutPath, job.Image)
				out <- true
			}
			close(out)
		}()
	} else {
		go func() { // Use a goroutine to avoid deadlock
			for job := range input { // Read from the channel
				// *********************************************************************************************
				// My part
				// Check if the output directory exists
				// *********************************************************************************************
				outputDirectory := strings.Replace(job.OutPath, "/"+filepath.Base(job.OutPath), "", 1)
				if _, err := os.Stat(outputDirectory); os.IsNotExist(err) {
					log.Printf("Output directory does not exist: %s", outputDirectory)
					out <- false
					continue
				}
				imageprocessing.WriteImage(job.OutPath, job.Image)
				out <- true
			}
			close(out)
		}()
	}
	return out
}

func runPipeline(isGoroutines bool) {
	// imagePaths := []string{"images/image1.jpeg",
	//     "images/image2.jpeg",
	//     "images/image3.jpeg",
	//     "images/image4.jpeg",
	// }

	// *********************************************************************************************
	// My part
	// Instead of reading only four images, lets read every image in the path
	// *********************************************************************************************

	imagePaths, err := filepath.Glob("images/*.jpeg")
	if err != nil {
		log.Fatalf("Failed to read images: %v", err)
	}

	channel1 := loadImage(imagePaths, isGoroutines)
	channel2 := resize(channel1, isGoroutines)
	channel3 := convertToGrayscale(channel2, isGoroutines)
	// Or use custom color conversion
	// channel3 := convertToCustomColorConversion(channel2, isGoroutines)
	writeResults := saveImage(channel3, isGoroutines)

	for success := range writeResults {
		if success {
			fmt.Println("Success!")
		} else {
			fmt.Println("Failed!")
		}
	}
}

func main() {
	start := time.Now()
	runPipeline(true) // With goroutines
	elapsed := time.Since(start)
	fmt.Printf("Pipeline with goroutines, time elapsed is: %s\n", elapsed)

	start = time.Now()
	runPipeline(false) // Without goroutines
	elapsed = time.Since(start)
	fmt.Printf("Pipeline without goroutines, time elapsed is: %s\n", elapsed)
}
