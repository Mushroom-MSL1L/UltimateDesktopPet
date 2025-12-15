package activities

import (
	"context"
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Mushroom-MSL1L/UltimateDesktopPet/app/desktop_pet/internal/database"
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

func TestActivityMeta_LoadAllActivities(t *testing.T) {
	meta := NewActivityMeta()
	meta.DB.InitDB(context.Background(), filepath.Join(t.TempDir(), "static_assets.db"), database.DBID(999))
	t.Cleanup(meta.DB.CloseDB)
	if err := meta.DB.GetDB().AutoMigrate(&Activity{}); err != nil {
		t.Fatalf("AutoMigrate: %v", err)
	}

	*meta.Activity = Activity{Name: "A", Path: "act_a"}
	if err := meta.Controller.Create(meta.DB.GetDB()); err != nil {
		t.Fatalf("Create(A): %v", err)
	}
	idA := meta.Activity.ID

	*meta.Activity = Activity{Name: "B", Path: "act_b"}
	if err := meta.Controller.Create(meta.DB.GetDB()); err != nil {
		t.Fatalf("Create(B): %v", err)
	}
	idB := meta.Activity.ID

	activities, err := meta.LoadAllActivities()
	if err != nil {
		t.Fatalf("LoadAllActivities: %v", err)
	}
	if len(activities) != 2 {
		t.Fatalf("len(activities) = %d, want %d", len(activities), 2)
	}

	ids := map[uint]bool{}
	for _, a := range activities {
		ids[a.ID] = true
	}
	if !ids[idA] || !ids[idB] {
		t.Fatalf("unexpected IDs: %#v (want %d and %d)", ids, idA, idB)
	}
}

func TestActivityMeta_LoadActivityFramesByID(t *testing.T) {
	meta := NewActivityMeta()
	meta.DB.InitDB(context.Background(), filepath.Join(t.TempDir(), "static_assets.db"), database.DBID(999))
	t.Cleanup(meta.DB.CloseDB)
	if err := meta.DB.GetDB().AutoMigrate(&Activity{}); err != nil {
		t.Fatalf("AutoMigrate: %v", err)
	}

	staticAssets := t.TempDir()
	meta.ST.StaticAssetPath = staticAssets
	meta.ST.DefaultImageFolder = "default"
	meta.ST.SpecifiedImageFolder = "custom"

	*meta.Activity = Activity{Name: "A", Path: "act_a"}
	if err := meta.Controller.Create(meta.DB.GetDB()); err != nil {
		t.Fatalf("Create(A): %v", err)
	}
	idA := meta.Activity.ID

	writeTestFile(t, filepath.Join(staticAssets, "default", "act_a", "1.png"), []byte("frame-1"))
	writeTestFile(t, filepath.Join(staticAssets, "default", "act_a", "2.png"), []byte("frame-2"))

	frames, err := meta.LoadActivityFramesByID(idA)
	if err != nil {
		t.Fatalf("LoadActivityFramesByID: %v", err)
	}
	if len(frames) != 2 {
		t.Fatalf("len(frames) = %d, want %d", len(frames), 2)
	}

	decoded := map[string]bool{}
	for _, f := range frames {
		decoded[string(decodeDataURI(t, f))] = true
	}
	if !decoded["frame-1"] || !decoded["frame-2"] {
		t.Fatalf("unexpected decoded frames: %#v", decoded)
	}
}
