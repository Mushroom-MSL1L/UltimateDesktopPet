package database

import (
	"path/filepath"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type testModel struct {
	ID   uint
	Name string
}

func openTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dbFile := filepath.Join(t.TempDir(), "test.db")
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm.Open: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db.DB: %v", err)
	}
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})
	if err := db.AutoMigrate(&testModel{}); err != nil {
		t.Fatalf("AutoMigrate: %v", err)
	}
	return db
}

func TestBaseController_Create_ErrorsWithoutModel(t *testing.T) {
	db := openTestDB(t)

	var c BaseController[testModel]
	if err := c.Create(db); err == nil {
		t.Fatalf("expected error, got nil")
	}

	var m *testModel
	c = BaseController[testModel]{Model: &m}
	if err := c.Create(db); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestBaseController_ReadAll_ValidatesInputs(t *testing.T) {
	db := openTestDB(t)

	var nilController *BaseController[testModel]
	if _, err := nilController.ReadAll(db); err == nil {
		t.Fatalf("expected error, got nil")
	}

	c := &BaseController[testModel]{}
	if _, err := c.ReadAll(db); err == nil {
		t.Fatalf("expected error, got nil")
	}

	model := &testModel{}
	c = &BaseController[testModel]{Model: &model}
	if _, err := c.ReadAll(nil); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestBaseController_CRUD(t *testing.T) {
	db := openTestDB(t)

	model := &testModel{Name: "first"}
	c := &BaseController[testModel]{Model: &model}

	if err := c.Create(db); err != nil {
		t.Fatalf("Create(first): %v", err)
	}
	id1 := model.ID
	if id1 == 0 {
		t.Fatalf("expected ID set after Create")
	}

	*model = testModel{Name: "second"}
	if err := c.Create(db); err != nil {
		t.Fatalf("Create(second): %v", err)
	}
	id2 := model.ID
	if id2 == 0 || id2 == id1 {
		t.Fatalf("expected a new ID after Create; id1=%d id2=%d", id1, id2)
	}

	got1, err := c.Read(db, id1)
	if err != nil {
		t.Fatalf("Read(id1): %v", err)
	}
	if got1.ID != id1 || got1.Name != "first" {
		t.Fatalf("Read(id1) = %#v, want ID=%d Name=%q", got1, id1, "first")
	}

	all, err := c.ReadAll(db)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(*all) != 2 {
		t.Fatalf("len(ReadAll) = %d, want %d", len(*all), 2)
	}

	latest, err := c.ReadFirst(db)
	if err != nil {
		t.Fatalf("ReadFirst: %v", err)
	}
	if latest.ID != id2 || latest.Name != "second" {
		t.Fatalf("ReadFirst = %#v, want ID=%d Name=%q", latest, id2, "second")
	}
}
