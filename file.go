package imagekit

import (
	"fmt"
	"bytes"
	"net/http"
	"encoding/json"
)

// List and search files.
func (imgKit *ImageKit) GetFiles(
	params *FilesFetchParams) (fileDetails *[]FileDetails, err error) {
	query := ""
	if params != nil {
		query, err = params.BuildURLQuery()
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/files%s", BASE_URL, query),
		bytes.NewBufferString(""),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resBodyStr, err := imgKit.DoRequest(req)
	if err != nil {
		return nil, err
	}
	fileDetails = &[]FileDetails{}
	err = json.Unmarshal([]byte(resBodyStr), fileDetails)
	if err != nil {
		return nil, err
	}
	return fileDetails, nil
}

// Get details of a file.
func (imgKit *ImageKit) GetFileDetails(fileId string) (fileDetail *FileDetails, err error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/files/%s/details", BASE_URL, fileId),
		bytes.NewBufferString(""),
	)
	req.Header.Set("Content-Type", "application/json")
	resBodyStr, err := imgKit.DoRequest(req)
	if err != nil {
		return nil, err
	}
	fileDetail = &FileDetails{}
	err = json.Unmarshal([]byte(resBodyStr), fileDetail)
	if err != nil {
		return nil, err
	}
	return fileDetail, nil
}

// Update details of a file.
func (imgKit *ImageKit) UpdateFileDetails(
	fileId string,
	options *FileOptions) (fileDetail *FileDetails, err error) {
	reqBody := ""
	if options != nil {
		reqBody, err = options.ToJSON()
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("%s/files/%s/details", BASE_URL, fileId),
		bytes.NewBufferString(reqBody),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resBodyStr, err := imgKit.DoRequest(req)
	if err != nil {
		return nil, err
	}
	fileDetail = &FileDetails{}
	err = json.Unmarshal([]byte(resBodyStr), fileDetail)
	if err != nil {
		return nil, err
	}
	return fileDetail, nil
}

// Delete a file.
func (imgKit *ImageKit) DeleteFile(fileId string) (err error) {
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/files/%s", BASE_URL, fileId),
		bytes.NewBufferString(""),
	)
	if err != nil {
		return err
	}
	_, err = imgKit.DoRequest(req)
	return err
}

// Copy a file.
func (imgKit *ImageKit) CopyFile(srcFilePath, destFolderPath string) (err error) {
	reqBody := make(map[string]string)
	reqBody["sourceFilePath"] = srcFilePath
	reqBody["destinationPath"] = destFolderPath
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/files/copy", BASE_URL),
		bytes.NewBufferString(string(reqBodyBytes)),
	)
	req.Header.Set("Content-Type", "application/json")
	_, err = imgKit.DoRequest(req)
	if err != nil {
		return err
	}
	return nil
}

// Move a file.
func (imgKit *ImageKit) MoveFile(srcFilePath, destFolderPath string) (err error) {
	reqBody := make(map[string]string)
	reqBody["sourceFilePath"] = srcFilePath
	reqBody["destinationPath"] = destFolderPath
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/files/move", BASE_URL),
		bytes.NewBufferString(string(reqBodyBytes)),
	)
	req.Header.Set("Content-Type", "application/json")
	_, err = imgKit.DoRequest(req)
	if err != nil {
		return err
	}
	return nil
}

// Rename a file.
func (imgKit *ImageKit) RenameFile(srcFilePath, newFileName string, purgeCache ...bool) (purgeRequestId string, err error) {
	reqBody := make(map[string]interface{})
	reqBody["filePath"] = srcFilePath
	reqBody["newFileName"] = newFileName
	reqBody["purgeCache"] = false
	if len(purgeCache) > 0 {
		reqBody["purgeCache"] = purgeCache[0]
	}
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/files/rename", BASE_URL),
		bytes.NewBufferString(string(reqBodyBytes)),
	)
	req.Header.Set("Content-Type", "application/json")
	resBodyStr, err := imgKit.DoRequest(req)
	if err != nil {
		return "", err
	}
	renameResponseFields := &map[string]string{}
	err = json.Unmarshal([]byte(resBodyStr), renameResponseFields)
	if err != nil {
		return "", err
	}
	return (*renameResponseFields)["purgeRequestId"], nil
}
