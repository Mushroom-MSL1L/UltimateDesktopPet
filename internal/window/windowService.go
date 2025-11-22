package window

import (
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
	ctx              context.Context
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

	applyWindowBounds(w.ctx, w.ensureWithinScreen(newBound))
}

func (w *WindowService) AdjustWindowFromBottom(width, height int) {
	current := w.currentBounds()

	newBound := Bounds{
		Top:    current.Top,
		Left:   int(float64(current.Left) + float64(current.Width-width)/2*w.devicePixelRatio),
		Width:  width,
		Height: height,
	}

	applyWindowBounds(w.ctx, w.ensureWithinScreen(newBound))
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

func (w *WindowService) ensureWithinScreen(bounds Bounds) Bounds {
	screenWidth, screenHeight, ok := w.screenSize()
	if !ok {
		return bounds
	}

	if bounds.Width > screenWidth {
		bounds.Width = screenWidth
	}
	if bounds.Height > screenHeight {
		bounds.Height = screenHeight
	}

	if bounds.Left < 0 {
		bounds.Left = 0
	}
	if bounds.Top < 0 {
		bounds.Top = 0
	}

	if right := float64(bounds.Left)/w.devicePixelRatio + float64(bounds.Width); right > float64(screenWidth) {
		bounds.Left = int(float64(screenWidth-bounds.Width) * (w.devicePixelRatio))
	}
	if bottom := float64(bounds.Top)/w.devicePixelRatio + float64(bounds.Height); bottom > float64(screenHeight) {
		bounds.Top = int(float64(screenHeight-bounds.Height) * (w.devicePixelRatio))
	}

	return bounds
}

func (w *WindowService) screenSize() (int, int, bool) {
	screens, err := runtime.ScreenGetAll(w.ctx)
	if err != nil || len(screens) == 0 {
		return 0, 0, false
	}

	screen := selectScreen(screens)
	width := screen.Size.Width
	height := screen.Size.Height

	if width == 0 || height == 0 {
		width = screen.Width
		height = screen.Height
	}

	if width == 0 || height == 0 {
		return 0, 0, false
	}

	return width, height, true
}

func selectScreen(screens []runtime.Screen) runtime.Screen {
	for _, screen := range screens {
		if screen.IsCurrent {
			return screen
		}
	}

	for _, screen := range screens {
		if screen.IsPrimary {
			return screen
		}
	}

	return screens[0]
}

func applyWindowBounds(ctx context.Context, bounds Bounds) Bounds {
	runtime.WindowSetSize(ctx, int(bounds.Width), int(bounds.Height))
	runtime.WindowSetPosition(ctx, int(bounds.Left), int(bounds.Top))
	return bounds
}
