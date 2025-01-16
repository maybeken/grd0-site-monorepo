package data

import (
	"encoding/json"

	"fmt"

	"os"

	"grd0.net/api/schema"
)

func ReadBlogs() ([]schema.Blog, error) {
	file, err := os.Open("data/blogPost.json")

	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}

	defer file.Close()

	var data []schema.Blog

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, err
	}

	return data, nil
}

func ReadGalleryDetail() ([]schema.GalleryDetail, error) {
	file, err := os.Open("data/galleryDetail.json")

	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}

	defer file.Close()

	var data []schema.GalleryDetail

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, err
	}

	return data, nil
}

func ReadGalleryCategory() ([]schema.GalleryCategory, error) {
	file, err := os.Open("data/galleryCategory.json")

	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}

	defer file.Close()

	var data []schema.GalleryCategory

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, err
	}

	return data, nil
}

func ReadMapLocation() ([]schema.MapLocation, error) {
	file, err := os.Open("data/travelersMap.json")

	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}

	defer file.Close()

	var data []schema.MapLocation

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, err
	}

	return data, nil
}
