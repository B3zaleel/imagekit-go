package imagekit

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Uploads a file to ImageKit.io.
func (imgKit *ImageKit) Upload(
	file,
	fileName string,
	options *FileOptions) (result *FileDetails, err error) {
	boundary := fmt.Sprintf("%s%s", strings.Repeat("-", 15), getRandomHex(32))
	body, err := getBody(file, fileName, boundary, options)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(
		http.MethodPost,
		UPLOAD_URL,
		bytes.NewBufferString(body),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	contentType := fmt.Sprintf("multipart/form-data; boundary=%s", boundary)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Content-Length", strconv.Itoa(len(body)))
	if err != nil {
		return nil, err
	}
	bodyStr, err := imgKit.DoRequest(req)
	if err != nil {
		return nil, err
	}
	result = &FileDetails{}
	err = json.Unmarshal([]byte(bodyStr), result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Checks the given file and fileName and transforms the file if necessary.
func getFile(
	file,
	fileName string) (newFile string, err error) {
	fileParamLen := len(bytes.Trim([]byte(file), " "))
	fileNameParamLen := len(bytes.Trim([]byte(fileName), " "))
	if !(fileParamLen > 0 && fileNameParamLen > 0) {
		return "", errors.New("file and fileName must not be empty")
	}
	_, err = url.ParseRequestURI(file)
	if err == nil {
		// valid url
		return file, nil
	}
	st, err := os.Stat(file)
	if err == nil && st.Mode().IsRegular() {
		// valid file path
		fileContents, err := os.ReadFile(file)
		if err != nil {
			return "", err
		}
		return base64.RawStdEncoding.EncodeToString(fileContents), nil
	}
	if _, err = base64.RawStdEncoding.DecodeString(file); err == nil {
		// valid base64 encoded file
		return file, nil
	}
	return "", errors.New("file is invalid")
}

// Creates the body for the upload request.
func getBody(
	file,
	fileName,
	boundary string,
	options *FileOptions) (body string, err error) {
	validFile, err := getFile(file, fileName)
	if err != nil {
		return "", err
	}
	formBodyBuilder := new(strings.Builder)
	dataFields := make(map[string]string)
	if options != nil {
		dataFields, err = options.ToDict()
		if err != nil {
			return "", err
		}
	}
	dataFields["file"] = validFile
	dataFields["fileName"] = fileName
	i := 0
	for key, val := range dataFields {
		formBodyBuilder.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		formBodyBuilder.WriteString(fmt.Sprintf(
			"Content-Disposition: form-data; name=\"%s\"\r\n",
			key,
		))
		formBodyBuilder.WriteString(fmt.Sprintf("Content-Type: text/plain"))
		formBodyBuilder.WriteString("\r\n\r\n")
		formBodyBuilder.WriteString(val)
		formBodyBuilder.WriteString("\r\n")
		if i == len(dataFields)-1 {
			formBodyBuilder.WriteString(fmt.Sprintf("--%s--", boundary))
		}
		i++
	}
	return formBodyBuilder.String(), nil
}

// Generates a random hexadecimal string with a given length.
func getRandomHex(length int) string {
	randInts := make([]byte, length)
	randSrc := rand.NewSource(time.Now().Unix())
	for i := 0; i < length; i++ {
		n := rand.New(randSrc).Intn(16)
		if n > 9 {
			randInts[i] = byte('A') + byte(n-10)
		} else {
			randInts[i] = byte('0') + byte(n)
		}
	}
	return string(randInts)
}
