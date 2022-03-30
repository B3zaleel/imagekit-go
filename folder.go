package imagekit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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

// Copy a folder.
func (imgKit *ImageKit) CopyFolder(sourceFolderPath, destinationPath string) (jobId string, err error) {
	reqBody := make(map[string]string)
	reqBody["sourceFolderPath"] = sourceFolderPath
	reqBody["destinationPath"] = destinationPath
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/bulkJobs/copyFolder", BASE_URL),
		bytes.NewBufferString(string(reqBodyBytes)),
	)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resBodyStr, err := imgKit.DoRequest(req)
	if err != nil {
		return "", err
	}
	responseFields := &map[string]string{}
	err = json.Unmarshal([]byte(resBodyStr), responseFields)
	if err != nil {
		return "", err
	}
	return (*responseFields)["jobId"], nil
}

// Move a folder.
func (imgKit *ImageKit) MoveFolder(sourceFolderPath, destinationPath string) (jobId string, err error) {
	reqBody := make(map[string]string)
	reqBody["sourceFolderPath"] = sourceFolderPath
	reqBody["destinationPath"] = destinationPath
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/bulkJobs/moveFolder", BASE_URL),
		bytes.NewBufferString(string(reqBodyBytes)),
	)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resBodyStr, err := imgKit.DoRequest(req)
	if err != nil {
		return "", err
	}
	responseFields := &map[string]string{}
	err = json.Unmarshal([]byte(resBodyStr), responseFields)
	if err != nil {
		return "", err
	}
	return (*responseFields)["jobId"], nil
}
