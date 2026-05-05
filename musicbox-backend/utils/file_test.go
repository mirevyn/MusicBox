package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNormalizeUploadPath(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "relative", in: "uploads/avatars/a.jpg", want: "uploads/avatars/a.jpg"},
		{name: "docker absolute", in: "/app/uploads/music/a.mp3", want: "uploads/music/a.mp3"},
		{name: "legacy absolute", in: "/opt/musicbox/uploads/covers/a.jpg", want: "uploads/covers/a.jpg"},
		{name: "public absolute", in: "/uploads/lyrics/a.lrc", want: "uploads/lyrics/a.lrc"},
		{name: "windows", in: `C:\musicbox\uploads\avatars\a.jpg`, want: "uploads/avatars/a.jpg"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizeUploadPath(tt.in); got != tt.want {
				t.Fatalf("NormalizeUploadPath(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestDeleteFileOnlyDeletesUploadsFiles(t *testing.T) {
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd returned error: %v", err)
	}
	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Chdir returned error: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalWd); err != nil {
			t.Fatalf("restore Chdir returned error: %v", err)
		}
	})

	targetDir := filepath.Join(tempDir, "uploads", "avatars")
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	targetPath := filepath.Join(targetDir, "a.jpg")
	if err := os.WriteFile(targetPath, []byte("avatar"), 0644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	if err := DeleteFile("/app/uploads/avatars/a.jpg"); err != nil {
		t.Fatalf("DeleteFile returned error: %v", err)
	}
	if _, err := os.Stat(targetPath); !os.IsNotExist(err) {
		t.Fatalf("target file still exists or stat failed unexpectedly: %v", err)
	}

	outsidePath := filepath.Join(tempDir, "outside.txt")
	if err := os.WriteFile(outsidePath, []byte("outside"), 0644); err != nil {
		t.Fatalf("WriteFile outside returned error: %v", err)
	}
	if err := DeleteFile("uploads/../outside.txt"); err == nil {
		t.Fatal("DeleteFile returned nil error for path outside uploads")
	}
	if _, err := os.Stat(outsidePath); err != nil {
		t.Fatalf("outside file was changed: %v", err)
	}
}
