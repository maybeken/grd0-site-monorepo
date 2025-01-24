package data

import (
	"encoding/json"

	"fmt"
	"os"
	"sync"

	"grd0.net/api/schema"
)

var (
	cache = make(map[string]interface{})
	mu    sync.RWMutex
)

func readFile[T any](filename string) (*T, error) {
	mu.RLock()
	if cachedData, found := cache[filename]; found {
		mu.RUnlock()
		return cachedData.(*T), nil
	}
	mu.RUnlock()

	file, err := os.Open("data/files/" + filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	fmt.Println("Cache not found, read file from disk:", filename)
	defer file.Close()

	var data T
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, err
	}

	// Cache the data
	mu.Lock()
	cache[filename] = &data
	mu.Unlock()

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

func ReadMusic() ([]schema.Music, error) {
	data, err := readFile[[]schema.Music]("music.json")

	return *data, err
}
