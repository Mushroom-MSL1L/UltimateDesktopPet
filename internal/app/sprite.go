package app

import (
	pp "UltimateDesktopPet/pkg/print"
)

func (a *App) PetFrames() ([]string, error) {
	frames, err := a.petMeta.LoadDefaultFrames()
	if err != nil {
		pp.Fatal(pp.App, "PetFrames: failed to load default frames: %v", err)
		return nil, err
	}
	return frames, nil
}

func (a *App) PetFramesBy(animationType string) ([]string, error) {
	frames, err := a.petMeta.LoadFramesByType(animationType)
	if err != nil {
		pp.Fatal(pp.App, "PetFramesBy: failed to load frames for animation type %s: %v", animationType, err)
		return nil, err
	}
	return frames, nil
}

func (a *App) PetFramesDrag() ([]string, error) {
	return a.PetFramesBy("drag")
}

func (a *App) PetFramesDrop() ([]string, error) {
	return a.PetFramesBy("drop")
}

func (a *App) PetFramesMoveLeft() ([]string, error) {
	return a.PetFramesBy("move_left")
}

func (a *App) PetFramesMoveRight() ([]string, error) {
	return a.PetFramesBy("move_right")
}

func (a *App) PetFramesMoveFar() ([]string, error) {
	return a.PetFramesBy("move_far")
}

func (a *App) PetFramesStand() ([]string, error) {
	return a.PetFramesBy("stand")
}
