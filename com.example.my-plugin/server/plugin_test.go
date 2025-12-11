package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/stretchr/testify/assert"
)

func TestServeHTTP(t *testing.T) {
	assert := assert.New(t)
	plugin := Plugin{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/v1/hello", nil)
	r.Header.Set("Mattermost-User-ID", "test-user-id")

	plugin.ServeHTTP(nil, w, r)

	result := w.Result()
	assert.NotNil(result)
	defer result.Body.Close()
	bodyBytes, err := io.ReadAll(result.Body)
	assert.Nil(err)
	bodyString := string(bodyBytes)

	assert.Equal("Hello, world!", bodyString)
	assert.Equal("Hello, world!", bodyString)
}

func TestFileWillBeUploaded(t *testing.T) {
	plugin := Plugin{}
	
	// Test case 1: Allowed file (PNG)
	info := &model.FileInfo{
		Name:     "test.png",
		MimeType: "image/png",
	}
	result, err := plugin.FileWillBeUploaded(nil, info, nil, nil)
	assert.Empty(t, err)
	assert.NotNil(t, result)

	// Test case 2: Allowed file (JPEG extension)
	info = &model.FileInfo{
		Name:     "test.jpg",
		MimeType: "image/unknown", // Specific mimetype check fail, but extension pass
	}
	result, err = plugin.FileWillBeUploaded(nil, info, nil, nil)
	assert.Empty(t, err)
	assert.NotNil(t, result)

	// Test case 3: Disallowed file (PDF)
	info = &model.FileInfo{
		Name:     "test.pdf",
		MimeType: "application/pdf",
	}
	result, err = plugin.FileWillBeUploaded(nil, info, nil, nil)
	assert.Equal(t, "只允許上傳 PNG, JPG, JPEG, SVG 格式的圖片檔案。", err)
	assert.Nil(t, result)

	// Test case 4: Disallowed file (Exe)
	info = &model.FileInfo{
		Name:     "test.exe",
		MimeType: "application/x-msdownload",
	}
	result, err = plugin.FileWillBeUploaded(nil, info, nil, nil)
	assert.Equal(t, "只允許上傳 PNG, JPG, JPEG, SVG 格式的圖片檔案。", err)
	assert.Nil(t, result)
}
