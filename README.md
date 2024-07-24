# Image Processing Pipeline

This project demonstrates an image processing pipeline implemented in Go, with the ability to run with and without goroutines for parallel processing. The pipeline includes the following stages:
- Load Images
- Resize Images
- Convert Images to Grayscale (or Custom Color Conversion)
- Save Images

## Structure

- `main.go`: Contains the main pipeline logic.
- `image_processing/image_processing.go`: Contains the image processing functions.
- `main_test.go`: Contains the unit tests for the pipeline.

## Setup

### Prerequisites

- Go (version 1.15 or later)
- A directory named `images` with some JPEG images for processing.


### Project Structure

.
├── README.md
├── main.go
├── main_test.go
└── image_processing
└── image_processing.go


### Creating the Project

1. **Clone the Repository:**

```bash
git clone <repository-url>
cd <repository-directory>


- Ensure you have a directory named images in the root of the project, and place some JPEG images in it.
- Build the Program: