package imagekit

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	BASE_URL = "https://api.imagekit.io/v1"
	UPLOAD_URL = "https://upload.imagekit.io/api/v1/files/upload"
	MIN_LIMIT_VALUE = 1
	MIN_SKIP_VALUE = 0
	MAX_LIMIT_VALUE = 1000
)
var VALID_TYPES = []string{"all", "file", "folder"}
var VALID_FILE_TYPES = []string{"all", "image", "non-image"}
var VALID_SORT_FIELDS = []string{
	"ASC_NAME",
	"DESC_NAME",
	"ASC_CREATED",
	"DESC_CREATED",
	"ASC_UPDATED",
	"DESC_UPDATED",
	"ASC_HEIGHT",
	"DESC_HEIGHT",
	"ASC_WIDTH",
	"DESC_WIDTH",
	"ASC_SIZE",
	"DESC_SIZE",
}

// Represents a struct with routines for managing assets on imagekit.io.
type ImageKit struct {
	PublicKey, PrivateKey, UrlEndpoint string
}

// Represents a typed 32-bit integer value.
type Int32 int32

// Represents a typed boolean value.
type Bool bool

// Represents a typed string value.
type String string

// Represents options for creating or modifying files.
type FileOptions struct {
	UseUniqueFileName *Bool
	Tags *[]String
	Folder *String
	IsPrivateFile *Bool
	CustomCoordinates *String
	ResponseFields *[]String
	Extensions *[]interface{}
	WebhookUrl *String
	OverwriteFile *Bool
	OverwriteAITags *Bool
	OverwriteTags *Bool
	OverwriteCustomMetadata *Bool
	CustomMetadata *interface{}
}

// Represents details about a file.
type FileDetails struct {
	FileId *String `json:"fileId" binding:"-"`
	Type *String `json:"type" binding:"-"`
	Name *String `json:"name" binding:"-"`
	FilePath *String `json:"filePath" binding:"-"`
	Tags *[]String `json:"tags" binding:"-"`
	AITags *[]interface{} `json:"AITags" binding:"-"`
	IsPrivateFile *Bool `json:"isPrivateFile" binding:"-"`
	CustomCoordinates *String `json:"customCoordinates" binding:"-"`
	Url *String `json:"url" binding:"-"`
	Thumbnail *String `json:"thumbnail" binding:"-"`
	FileType *String `json:"fileType" binding:"-"`
	Mime *String `json:"mime" binding:"-"`
	Height Int32 `json:"height" binding:"-"`
	Width Int32 `json:"width" binding:"-"`
	Size Int32 `json:"size" binding:"-"`
	HasAlpha *Bool `json:"hasAlpha" binding:"-"`
	CustomMetadata *interface{} `json:"customMetadata" binding:"-"`
	EmbeddedMetadata *interface{} `json:"embeddedMetadata" binding:"-"`
	CreatedAt *time.Time `json:"createdAt" binding:"-" time_format:"YYYY-MM-DDTHH:mm:ss.sssZ"`
	UpdatedAt *time.Time `json:"updatedAt" binding:"-" time_format:"YYYY-MM-DDTHH:mm:ss.sssZ"`
	ExtensionStatus map[string]string `json:"extensionStatus" binding:"-"`
}

// Represents query parameters for fetching files from imagekit.io.
type FilesFetchParams struct {
	Type, Sort, Path, SearchQuery *String
	Tags *[]String
	FileType *String
	Limit, Skip *Int32
}

// Runs an http request.
func (imgKit *ImageKit) DoRequest(req *http.Request) (body string, err error) {
	client := &http.Client{}
	fetchRequest := true
	res := new(http.Response)
	req.SetBasicAuth(imgKit.PrivateKey, "")
	for ; fetchRequest; {
		res, err = client.Do(req)
		if err != nil {
			return "", err
		}
		if res.StatusCode != 429 {
			fetchRequest = false
		}
		if res.StatusCode == 429 {
			waitTime, err := strconv.Atoi(res.Header.Get("X-RateLimit-Reset"))
			if err != nil {
				return "", err
			}
			time.Sleep(time.Millisecond * (time.Duration(waitTime) / time.Millisecond))
		}
	}
	buf := new(bytes.Buffer)
	if buf == nil || res.Body == nil {
		return "", errors.New("failed to create bytes buffer")
	}
	buf.ReadFrom(res.Body)
	body = buf.String()
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return "", errors.New(body)
	}
	return body, nil
}

// Joins a typed string array to a string.
func joinStringArray(strArr []String, sep string, escape bool) string {
	n := len(strArr)
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteString(sep)
		}
		if escape {
			sb.WriteString(strings.ReplaceAll(string(strArr[i]), "\"", "\\\""))
		} else {
			sb.WriteString(string(strArr[i]))
		}
	}
	return sb.String()
}

// Checks if a string is in an array of strings.
func (str String) StringInArray(arr []string) bool {
	for i := 0; i < len(arr); i++ {
		if arr[i] == string(str) {
			return true
		}
	}
	return false
}

// Converts options to a map of optionName-optionValue pairs.
func (options FileOptions) ToDict() (fields map[string]string, err error) {
	fields = make(map[string]string)
	if options.UseUniqueFileName != nil {
		fields["useUniqueFileName"] = fmt.Sprintf("%t", *options.UseUniqueFileName)
	}
	if options.Tags != nil {
		fields["tags"] = joinStringArray(*options.Tags, ",", false)
	}
	if options.Folder != nil {
		fields["folder"] = string(*options.Folder)
	}
	if options.IsPrivateFile != nil {
		fields["isPrivateFile"] = fmt.Sprintf("%t", *options.IsPrivateFile)
	}
	if options.CustomCoordinates != nil {
		fields["customCoordinates"] = string(*options.CustomCoordinates)
	}
	if options.ResponseFields != nil {
		fields["responseFields"] = joinStringArray(
			*options.ResponseFields,
			",",
			false,
		)
	}
	if options.Extensions != nil {
		customExtensionsJSON, err := json.Marshal(*options.Extensions)
		if err != nil {
			return nil, err
		}
		fields["extensions"] = string(customExtensionsJSON)
	}
	if options.WebhookUrl != nil {
		fields["webhookUrl"] = string(*options.WebhookUrl)
	}
	if options.OverwriteFile != nil {
		fields["overwriteFile"] = fmt.Sprintf("%t", *options.OverwriteFile)
	}
	if options.OverwriteAITags != nil {
		fields["overwriteAITags"] = fmt.Sprintf("%t", *options.OverwriteAITags)
	}
	if options.OverwriteTags != nil {
		fields["overwriteTags"] = fmt.Sprintf("%t", *options.OverwriteTags)
	}
	if options.OverwriteCustomMetadata != nil {
		fields["overwriteCustomMetadata"] = strconv.FormatBool(
			bool(*options.OverwriteCustomMetadata),
		)
	}
	if options.CustomMetadata != nil {
		customMetadataJSON, err := json.Marshal(*options.CustomMetadata)
		if err != nil {
			return nil, err
		}
		fields["customMetadata"] = string(customMetadataJSON)
	}
	return fields, nil
}

// Converts options to its JSON representation.
func (options FileOptions) ToJSON() (jsonStr string, err error) {
	var jsonBuilder strings.Builder
	jsonBuilder.WriteRune('{')
	delim := '0'
	jsonPut := func (str string) {
		if delim != '0' { jsonBuilder.WriteRune(delim) }
		jsonBuilder.WriteString(str)
		delim = ','
	}
	if options.UseUniqueFileName != nil {
		jsonPut(fmt.Sprintf(
			`"useUniqueFileName": %t`,
			bool(*options.UseUniqueFileName),
		))
	}
	if options.Tags != nil {
		jsonPut(fmt.Sprintf(
			`"tags": ["%s"]`,
			joinStringArray(*options.Tags, `","`, true),
		))
	}
	if options.Folder != nil {
		jsonPut(fmt.Sprintf(
			`"folder": "%s"`,
			strings.ReplaceAll(string(*options.Folder), `"`, `\"`),
		))
		delim = ','
	}
	if options.IsPrivateFile != nil {
		jsonPut(fmt.Sprintf(
			`"isPrivateFile": %t`,
			bool(*options.IsPrivateFile),
		))
	}
	if options.CustomCoordinates != nil {
		jsonPut(fmt.Sprintf(
			`"customCoordinates": "%s"`,
			string(*options.CustomCoordinates),
		))
	}
	if options.ResponseFields != nil {
		jsonPut(fmt.Sprintf(
			`"responseFields": "%s"`,
			joinStringArray(*options.ResponseFields, ",", true),
		))
	}
	if options.Extensions != nil {
		customExtensionsJSON, err := json.Marshal(*options.Extensions)
		if err != nil {
			return "", err
		}
		jsonPut(fmt.Sprintf(
			`"extensions": %s`,
			string(customExtensionsJSON),
		))
	}
	if options.WebhookUrl != nil {
		jsonPut(fmt.Sprintf(
			`"webhookUrl": "%s"`,
			string(*options.WebhookUrl),
		))
	}
	if options.OverwriteFile != nil {
		jsonPut(fmt.Sprintf(
			`"overwriteFile": %t`,
			bool(*options.OverwriteFile),
		))
	}
	if options.OverwriteAITags != nil {
		jsonPut(fmt.Sprintf(
			`"overwriteAITags": %t`,
			bool(*options.OverwriteAITags),
		))
	}
	if options.OverwriteTags != nil {
		jsonPut(fmt.Sprintf(
			`"overwriteTags": %t`,
			bool(*options.OverwriteTags),
		))
	}
	if options.OverwriteCustomMetadata != nil {
		jsonPut(fmt.Sprintf(
			`"overwriteCustomMetadata": %t`,
			bool(*options.OverwriteCustomMetadata),
		))
	}
	if options.CustomMetadata != nil {
		customMetadataJSON, err := json.Marshal(*options.CustomMetadata)
		if err != nil {
			return "", err
		}
		jsonPut(fmt.Sprintf(
			`"customMetadata": %s`,
			string(customMetadataJSON),
		))
	}
	jsonBuilder.WriteRune('}')
	return jsonBuilder.String(), nil
}

// Creates a URL query from the FilesFetchParams.
func (params FilesFetchParams) BuildURLQuery() (query string, err error) {
	var queryBuilder strings.Builder
	delim := '0'
	queryPut := func (str string) {
		if delim != '0' { queryBuilder.WriteRune(delim) }
		queryBuilder.WriteString(url.QueryEscape(str))
		delim = '&'
	}
	queryBuilder.WriteRune('?')
	if params.Type != nil {
		if !(*params.Type).StringInArray(VALID_TYPES) {
			return "", errors.New("invalid type value")
		}
		queryPut(fmt.Sprintf("type=%s", *params.Type))
	}
	if params.Sort != nil {
		if !(*params.Sort).StringInArray(VALID_SORT_FIELDS) {
			return "", errors.New("invalid sort value")
		}
		queryPut(fmt.Sprintf("sort=%s", *params.Sort))
	}
	if params.Path != nil {
		queryPut(fmt.Sprintf("path=%s", *params.Path))
	}
	if params.SearchQuery != nil {
		queryPut(fmt.Sprintf("searchQuery=%s", *params.SearchQuery))
	}
	if params.Tags != nil {
		queryPut(fmt.Sprintf("tags=%s", joinStringArray(*params.Tags, ",", false)))
	}
	if params.FileType != nil {
		if !(*params.FileType).StringInArray(VALID_FILE_TYPES) {
			return "", nil
		}
		queryPut(fmt.Sprintf("fileType=%s", *params.FileType))
	}
	if params.Limit != nil {
		if *params.Limit < MIN_LIMIT_VALUE || *params.Limit > MAX_LIMIT_VALUE {
			return "", errors.New("limit is out of bounds")
		}
		queryPut(fmt.Sprintf("limit=%d", *params.Limit))
	}
	if params.Skip != nil {
		if *params.Skip < MIN_SKIP_VALUE {
			return "", errors.New("skip is out of bounds")
		}
		queryPut(fmt.Sprintf("skip=%d", *params.Skip))
	}
	return queryBuilder.String(), nil
}
