package imagekit

import (
	"fmt"
	"bytes"
	"net/http"
	"encoding/json"
)

// Create a folder.
func (imgKit *ImageKit) CreateFolder(folderName, parentFolderPath string) (err error) {
	reqBody := make(map[string]string)
	reqBody["folderName"] = folderName
	reqBody["parentFolderPath"] = parentFolderPath
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/folder", BASE_URL),
		bytes.NewBufferString(string(reqBodyBytes)),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	_, err = imgKit.DoRequest(req)
	return err
}

// Delete a folder.
func (imgKit *ImageKit) DeleteFolder(folderPath string) (err error) {
	reqBody := make(map[string]string)
	reqBody["folderPath"] = folderPath
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/folder", BASE_URL),
		bytes.NewBufferString(string(reqBodyBytes)),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	_, err = imgKit.DoRequest(req)
	return err
}
