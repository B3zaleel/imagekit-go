package imagekit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Add tags to an array of files.
func (imgKit *ImageKit) AddTags(fileIds, tags []string) (updatedFileIds []string, err error) {
	reqBody := make(map[string][]string)
	reqBody["fileIds"] = fileIds
	reqBody["tags"] = tags
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/files/addTags", BASE_URL),
		bytes.NewBufferString(string(reqBodyBytes)),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resBodyStr, err := imgKit.DoRequest(req)
	if err != nil {
		return nil, err
	}
	updatedFileIds = []string{}
	err = json.Unmarshal([]byte(resBodyStr), &updatedFileIds)
	if err != nil {
		return nil, err
	}
	return updatedFileIds, nil
}

// Delete an array of files.
func (imgKit *ImageKit) DeleteFiles(fileIds []string) (deletedFileIds []string, err error) {
	reqBody := make(map[string][]string)
	reqBody["fileIds"] = fileIds
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/files/batch/deleteByFileIds", BASE_URL),
		bytes.NewBufferString(string(reqBodyBytes)),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resBodyStr, err := imgKit.DoRequest(req)
	if err != nil {
		return nil, err
	}
	deletedFileIds = []string{}
	err = json.Unmarshal([]byte(resBodyStr), &deletedFileIds)
	if err != nil {
		return nil, err
	}
	return deletedFileIds, nil
}

// Remove AI tags from an array of files.
func (imgKit *ImageKit) RemoveAITags(fileIds, aiTags []string) (updatedFileIds []string, err error) {
	reqBody := make(map[string][]string)
	reqBody["fileIds"] = fileIds
	reqBody["AITags"] = aiTags
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/files/removeAITags", BASE_URL),
		bytes.NewBufferString(string(reqBodyBytes)),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resBodyStr, err := imgKit.DoRequest(req)
	if err != nil {
		return nil, err
	}
	updatedFileIds = []string{}
	err = json.Unmarshal([]byte(resBodyStr), &updatedFileIds)
	if err != nil {
		return nil, err
	}
	return updatedFileIds, nil
}

// Remove tags from an array of files.
func (imgKit *ImageKit) RemoveTags(fileIds, tags []string) (updatedFileIds []string, err error) {
	reqBody := make(map[string][]string)
	reqBody["fileIds"] = fileIds
	reqBody["tags"] = tags
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/files/removeTags", BASE_URL),
		bytes.NewBufferString(string(reqBodyBytes)),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resBodyStr, err := imgKit.DoRequest(req)
	if err != nil {
		return nil, err
	}
	updatedFileIds = []string{}
	err = json.Unmarshal([]byte(resBodyStr), &updatedFileIds)
	if err != nil {
		return nil, err
	}
	return updatedFileIds, nil
}

// Get details of a bulk job.
func (imgKit *ImageKit) GetBulkJobStatus(
	jobId string) (jobDetails *JobDetails, err error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/bulkJobs/%s", BASE_URL, jobId),
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resBodyStr, err := imgKit.DoRequest(req)
	if err != nil {
		return nil, err
	}
	jobDetails = &JobDetails{}
	err = json.Unmarshal([]byte(resBodyStr), jobDetails)
	if err != nil {
		return nil, err
	}
	return jobDetails, nil
}
