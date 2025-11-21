package window

import (
	pp "UltimateDesktopPet/pkg/print"
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// PetHitbox defines the sprite bounds reported from the frontend.
type PetHitbox struct {
	Left             int `json:"left"`
	Top              int `json:"top"`
	Width            int `json:"width"`
	Height           int `json:"height"`
	DevicePixelRatio int `json:"devicePixelRatio"`
}

// Bounds represents the current window rectangle in screen coordinates.
type Bounds struct {
	Left   int
	Top    int
	Width  int
	Height int
}

// WindowService centralizes window sizing logic and exposes it directly to Wails bindings.
type WindowService struct {
	ctx    context.Context
	devicePixelRatio float64
}

func NewWindowService() *WindowService {
	return &WindowService{}
}

// SetContext should be called from the main OnStartup hook so runtime helpers
// can access the valid lifecycle context later on.
func (w *WindowService) SetContext(ctx context.Context) {
	w.ctx = ctx
}

func (w *WindowService) SetDevicePixelRatio(devicePixelRatio float64) {
	w.devicePixelRatio = devicePixelRatio
}

func (w *WindowService) AdjustWindowFromLeftBottom(width, height int) {
	current := w.currentBounds()
	newBound := Bounds{
		Top:    current.Top,
		Left:   current.Left,
		Width:  width,
		Height: height,
	}

	applyWindowBounds(w.ctx, newBound)
}

func (w *WindowService) AdjustWindowFromBottom(width, height int) {

	current := w.currentBounds()

	newBound := Bounds{
		Top:    current.Top,
		Left:   int(float64(current.Left) + float64(current.Width-width)/2*w.devicePixelRatio),
		Width:  width,
		Height: height,
	}

	pp.Assert(pp.System, "current bounds: %+v", current)
	pp.Assert(pp.System, "new bounds:     %+v", newBound)
	applyWindowBounds(w.ctx, newBound)
}

func (w *WindowService) currentBounds() Bounds {
	windowWidth, windowHeight := runtime.WindowGetSize(w.ctx)
	windowX, windowY := runtime.WindowGetPosition(w.ctx)

	return Bounds{
		Left:   windowX,
		Width:  windowWidth,
		Height: windowHeight,
		Top:    windowY,
	}
}

func applyWindowBounds(ctx context.Context, bounds Bounds) Bounds {
	runtime.WindowSetSize(ctx, int(bounds.Width), int(bounds.Height))
	runtime.WindowSetPosition(ctx, int(bounds.Left), int(bounds.Top))
	return bounds
}
