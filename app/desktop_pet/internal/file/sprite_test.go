package file

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTestFile(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
}

func decodeDataURI(t *testing.T, uri string) []byte {
	t.Helper()
	if !strings.HasPrefix(uri, "data:image/png;base64,") {
		t.Fatalf("unexpected prefix: %q", uri)
	}
	raw, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(uri, "data:image/png;base64,"))
	if err != nil {
		t.Fatalf("DecodeString: %v", err)
	}
	return raw
}

func TestNewSpriteTool_SetsDefaults(t *testing.T) {
	st := NewSpriteTool(SpriteTool{})
	if st.DefaultImageFolder != "default" {
		t.Fatalf("DefaultImageFolder = %q, want %q", st.DefaultImageFolder, "default")
	}
	if st.URISchemePrefix != "data:image/png;base64," {
		t.Fatalf("URISchemePrefix = %q, want %q", st.URISchemePrefix, "data:image/png;base64,")
	}
	if st.Suffix != ".png" {
		t.Fatalf("Suffix = %q, want %q", st.Suffix, ".png")
	}
}

func TestSpriteTool_LoadFrameFromDir_FallsBackToDefault(t *testing.T) {
	tmp := t.TempDir()
	want := []byte("not-a-real-png")

	st := NewSpriteTool(SpriteTool{})
	st.StaticAssetPath = tmp
	st.DefaultImageFolder = "default"
	st.SpecifiedImageFolder = "custom"

	writeTestFile(t, filepath.Join(tmp, "default", "icon.png"), want)

	gotURI, err := st.LoadFrameFromDir("icon")
	if err != nil {
		t.Fatalf("LoadFrameFromDir: %v", err)
	}
	got := decodeDataURI(t, gotURI)
	if string(got) != string(want) {
		t.Fatalf("decoded bytes = %q, want %q", string(got), string(want))
	}
}

func TestSpriteTool_LoadFramesFromDir_PrefersSpecifiedWhenAvailable(t *testing.T) {
	tmp := t.TempDir()

	st := NewSpriteTool(SpriteTool{})
	st.StaticAssetPath = tmp
	st.DefaultImageFolder = "default"
	st.SpecifiedImageFolder = "custom"

	writeTestFile(t, filepath.Join(tmp, "custom", "stand", "1.png"), []byte("custom-1"))
	writeTestFile(t, filepath.Join(tmp, "custom", "stand", "2.png"), []byte("custom-2"))
	writeTestFile(t, filepath.Join(tmp, "default", "stand", "1.png"), []byte("default-1"))

	frames, err := st.LoadFramesFromDir("stand")
	if err != nil {
		t.Fatalf("LoadFramesFromDir: %v", err)
	}
	if len(frames) != 2 {
		t.Fatalf("len(frames) = %d, want %d", len(frames), 2)
	}

	decoded := map[string]bool{}
	for _, f := range frames {
		decoded[string(decodeDataURI(t, f))] = true
	}
	if !decoded["custom-1"] || !decoded["custom-2"] || decoded["default-1"] {
		t.Fatalf("unexpected decoded frames: %#v", decoded)
	}
}

func TestSpriteTool_LoadFramesFromDir_FallsBackWhenSpecifiedEmpty(t *testing.T) {
	tmp := t.TempDir()

	st := NewSpriteTool(SpriteTool{})
	st.StaticAssetPath = tmp
	st.DefaultImageFolder = "default"
	st.SpecifiedImageFolder = "custom"

	writeTestFile(t, filepath.Join(tmp, "custom", "stand", "readme.txt"), []byte("ignore"))
	writeTestFile(t, filepath.Join(tmp, "default", "stand", "1.png"), []byte("default-1"))

	frames, err := st.LoadFramesFromDir("stand")
	if err != nil {
		t.Fatalf("LoadFramesFromDir: %v", err)
	}
	if len(frames) != 1 {
		t.Fatalf("len(frames) = %d, want %d", len(frames), 1)
	}
	if got := string(decodeDataURI(t, frames[0])); got != "default-1" {
		t.Fatalf("decoded bytes = %q, want %q", got, "default-1")
	}
}

func TestSpriteTool_LoadFramesFromDir_ErrorsWhenMissing(t *testing.T) {
	st := NewSpriteTool(SpriteTool{})
	st.StaticAssetPath = t.TempDir()
	st.DefaultImageFolder = "default"
	st.SpecifiedImageFolder = "custom"

	if _, err := st.LoadFramesFromDir("stand"); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

