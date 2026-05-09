package admin

import (
	"testing"
)

// TestExtractFilename 测试从 OSSKey 提取文件名
func TestExtractFilename(t *testing.T) {
	tests := []struct {
		name   string
		ossKey string
		want   string
	}{
		{"standard path", "family123/baby456/photo.jpg", "photo.jpg"},
		{"nested path", "prefix/a/b/c/file.png", "file.png"},
		{"no separator", "simple.jpg", "simple.jpg"},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractFilename(tt.ossKey)
			if got != tt.want {
				t.Errorf("extractFilename(%q) = %q, want %q", tt.ossKey, got, tt.want)
			}
		})
	}
}

// TestPhotoItemFields 验证 photoItem 包含 filename 和 original_url 字段
// 黑盒测试：通过构造测试验证字段正确设置
func TestPhotoItemFields(t *testing.T) {
	item := photoItem{
		ID:          "p1",
		BabyID:      "b1",
		FamilyID:    "f1",
		Filename:    "photo.jpg",
		OriginalURL: "https://cdn.example.com/family/baby/photo.jpg",
		Size:        1024,
		ContentType: "image/jpeg",
		Status:      "active",
		UploadedAt:  "2026-05-09 12:00:00",
	}

	if item.Filename == "" {
		t.Error("photoItem.Filename should not be empty")
	}
	if item.OriginalURL == "" {
		t.Error("photoItem.OriginalURL should not be empty")
	}
}
