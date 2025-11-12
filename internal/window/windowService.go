package window

import (
	"context"
	"sync"

	pp "UltimateDesktopPet/pkg/print"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// PetHitbox defines the sprite bounds reported from the frontend.
type PetHitbox struct {
	Left             float64 `json:"left"`
	Top              float64 `json:"top"`
	Width            float64 `json:"width"`
	Height           float64 `json:"height"`
	DevicePixelRatio float64 `json:"devicePixelRatio"`
}

// Bounds represents the current window rectangle in screen coordinates.
type Bounds struct {
	Left   float64
	Top    float64
	Width  float64
	Height float64
}

// WindowService centralizes window sizing logic and exposes it directly to Wails bindings.
type WindowService struct {
	ctx    context.Context
	mu     sync.RWMutex
	hitbox PetHitbox
}

func NewWindowService() *WindowService {
	return &WindowService{}
}

// SetContext should be called from the main OnStartup hook so runtime helpers
// can access the valid lifecycle context later on.
func (w *WindowService) SetContext(ctx context.Context) {
	w.mu.Lock()
	w.ctx = ctx
	w.mu.Unlock()
}

func (w *WindowService) UpdatePetHitbox(hitbox PetHitbox) error {
	if hitbox.DevicePixelRatio <= 0 {
		hitbox.DevicePixelRatio = 1
	}

	w.mu.Lock()
	w.hitbox = hitbox
	w.mu.Unlock()
	return nil
}

func (w *WindowService) AdjustWindowfromLeftBottom(width, height float64) {
	ctx := w.runtimeCtx()
	if ctx == nil {
		pp.Warn(pp.System, "window context not initialised yet; skipping resize request")
		return
	}

	current := w.currentBounds(ctx)
	newBound := Bounds{
		Top:    current.Top,
		Left:   current.Left,
		Width:  width,
		Height: height,
	}

	applyWindowBounds(ctx, newBound)
}

func (w *WindowService) currentBounds(ctx context.Context) Bounds {
	w.mu.RLock()
	hitbox := w.hitbox
	w.mu.RUnlock()

	windowX, windowY := runtime.WindowGetPosition(ctx)

	return Bounds{
		Left:   float64(windowX),
		Top:    float64(windowY),
		Width:  hitbox.Width,
		Height: hitbox.Height,
	}
}

func applyWindowBounds(ctx context.Context, bounds Bounds) Bounds {
	runtime.WindowSetPosition(ctx, int(bounds.Left), int(bounds.Top))
	runtime.WindowSetSize(ctx, int(bounds.Width), int(bounds.Height))
	return bounds
}

func (w *WindowService) runtimeCtx() context.Context {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.ctx
}
