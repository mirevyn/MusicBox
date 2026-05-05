package service

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

func TestPublicUploadPath(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "relative upload path",
			in:   "uploads/avatars/a.jpg",
			want: "uploads/avatars/a.jpg",
		},
		{
			name: "relative upload path with dot prefix",
			in:   "./uploads/covers/a.jpg",
			want: "uploads/covers/a.jpg",
		},
		{
			name: "docker app upload path",
			in:   "/app/uploads/avatars/a.jpg",
			want: "uploads/avatars/a.jpg",
		},
		{
			name: "legacy opt upload path",
			in:   "/opt/musicbox/uploads/music/a.mp3",
			want: "uploads/music/a.mp3",
		},
		{
			name: "public upload path with leading slash",
			in:   "/uploads/lyrics/a.lrc",
			want: "uploads/lyrics/a.lrc",
		},
		{
			name: "windows upload path",
			in:   `C:\musicbox\uploads\avatars\a.jpg`,
			want: "uploads/avatars/a.jpg",
		},
		{
			name: "windows upload path with forward slashes",
			in:   "C:/musicbox/uploads/covers/a.jpg",
			want: "uploads/covers/a.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := publicUploadPath(tt.in); got != tt.want {
				t.Fatalf("publicUploadPath(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestSaveUploadedFileReturnsPublicPath(t *testing.T) {
	dstDir := filepath.Join(t.TempDir(), "uploads", "avatars")
	header := testMultipartFileHeader(t, "avatar.png", testPNGBytes())

	got, err := SaveUploadedFile(nil, header, dstDir)
	if err != nil {
		t.Fatalf("SaveUploadedFile returned error: %v", err)
	}
	if !strings.HasPrefix(got, "uploads/avatars/") {
		t.Fatalf("SaveUploadedFile returned %q, want uploads/avatars/ prefix", got)
	}
	if !strings.HasSuffix(got, ".png") {
		t.Fatalf("SaveUploadedFile returned %q, want .png suffix", got)
	}

	savedPath := filepath.Join(dstDir, path.Base(got))
	data, err := os.ReadFile(savedPath)
	if err != nil {
		t.Fatalf("saved file was not created at %s: %v", savedPath, err)
	}
	if !bytes.Equal(data, testPNGBytes()) {
		t.Fatalf("saved file content = %v, want %v", data, testPNGBytes())
	}
}

func TestSaveUploadedFileRejectsInvalidImageContent(t *testing.T) {
	dstDir := filepath.Join(t.TempDir(), "uploads", "avatars")
	header := testMultipartFileHeader(t, "avatar.png", []byte("not an image"))

	if _, err := SaveUploadedFile(nil, header, dstDir); err == nil {
		t.Fatal("SaveUploadedFile returned nil error for invalid image content")
	} else if !errors.Is(err, ErrInvalidUpload) {
		t.Fatalf("SaveUploadedFile error = %v, want ErrInvalidUpload", err)
	}
}

func TestSaveUploadedFileRejectsUnsupportedExtension(t *testing.T) {
	dstDir := filepath.Join(t.TempDir(), "uploads", "avatars")
	header := testMultipartFileHeader(t, "avatar.exe", testPNGBytes())

	if _, err := SaveUploadedFile(nil, header, dstDir); err == nil {
		t.Fatal("SaveUploadedFile returned nil error for unsupported extension")
	} else if !errors.Is(err, ErrInvalidUpload) {
		t.Fatalf("SaveUploadedFile error = %v, want ErrInvalidUpload", err)
	}
}

func testPNGBytes() []byte {
	return []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n', 0, 0, 0, 0}
}

func testMultipartFileHeader(t *testing.T, filename string, contents []byte) *multipart.FileHeader {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("CreateFormFile returned error: %v", err)
	}
	if _, err := part.Write(contents); err != nil {
		t.Fatalf("writing multipart content returned error: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("closing multipart writer returned error: %v", err)
	}

	req := httptest.NewRequest("POST", "/upload", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err := req.ParseMultipartForm(int64(body.Len() + 1024)); err != nil {
		t.Fatalf("ParseMultipartForm returned error: %v", err)
	}

	files := req.MultipartForm.File["file"]
	if len(files) != 1 {
		t.Fatalf("multipart form has %d files, want 1", len(files))
	}
	return files[0]
}
