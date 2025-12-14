package items

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

func TestItemsMeta_LoadAllItemsAndFrames(t *testing.T) {
	meta := NewItemMeta()
	meta.DB.InitDB(context.Background(), filepath.Join(t.TempDir(), "static_assets.db"), database.DBID(999))
	t.Cleanup(meta.DB.CloseDB)
	if err := meta.DB.GetDB().AutoMigrate(&Item{}); err != nil {
		t.Fatalf("AutoMigrate: %v", err)
	}

	staticAssets := t.TempDir()
	meta.ST.StaticAssetPath = staticAssets
	meta.ST.DefaultImageFolder = "default"
	meta.ST.SpecifiedImageFolder = "custom"

	*meta.Item = Item{Name: "I1", Path: "item_1"}
	if err := meta.Controller.Create(meta.DB.GetDB()); err != nil {
		t.Fatalf("Create(I1): %v", err)
	}
	id1 := meta.Item.ID

	all, err := meta.LoadAllItems()
	if err != nil {
		t.Fatalf("LoadAllItems: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("len(all) = %d, want %d", len(all), 1)
	}

	got, err := meta.LoadItemByID(id1)
	if err != nil {
		t.Fatalf("LoadItemByID: %v", err)
	}
	if got.ID != id1 || got.Name != "I1" || got.Path != "item_1" {
		t.Fatalf("LoadItemByID = %#v, want ID=%d Name=%q Path=%q", got, id1, "I1", "item_1")
	}

	writeTestFile(t, filepath.Join(staticAssets, "default", "item_1.png"), []byte("frame-1"))
	frame, err := meta.LoadItemFrameByID(id1)
	if err != nil {
		t.Fatalf("LoadItemFrameByID: %v", err)
	}
	if gotBytes := string(decodeDataURI(t, frame)); gotBytes != "frame-1" {
		t.Fatalf("decoded bytes = %q, want %q", gotBytes, "frame-1")
	}
}
