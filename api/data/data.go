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
