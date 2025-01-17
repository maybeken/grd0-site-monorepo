package data

import (
	"encoding/json"

	"fmt"

	"os"

	"grd0.net/api/schema"
)

func readFile[T any](filename string) (*T, error) {
	file, err := os.Open("data/" + filename)

	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}

	defer file.Close()

	var data T

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, err
	}

	return &data, nil
}

func ReadBlogs() ([]schema.Blog, error) {
	data, err := readFile[[]schema.Blog]("blogPost.json")

	return *data, err
}

func ReadGalleryDetail() ([]schema.GalleryDetail, error) {
	data, err := readFile[[]schema.GalleryDetail]("galleryDetail.json")

	return *data, err
}

func ReadGalleryCategory() (schema.GalleryCategory, error) {
	data, err := readFile[schema.GalleryCategory]("galleryCategory.json")

	return *data, err
}

func ReadMapLocation() ([]schema.MapLocation, error) {
	data, err := readFile[[]schema.MapLocation]("travelersMap.json")

	return *data, err
}

func ReadAsset() (schema.AssetFileList, error) {
	data, err := readFile[schema.AssetFileList]("files.json")

	return *data, err
}
