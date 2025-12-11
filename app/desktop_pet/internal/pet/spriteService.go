package pet

import (
	pp "github.com/Mushroom-MSL1L/UltimateDesktopPet/pkg/print"
)

func (p *PetMeta) PetFrames() ([]string, error) {
	return p.PetFramesBy(defaultPetAnimationType)
}

func (p *PetMeta) PetFramesBy(animationType string) ([]string, error) {
	frames, err := p.ST.LoadFramesFromDir(animationType)
	if err != nil {
		pp.Fatal(pp.Pet, "PetFramesBy: failed to load frames for animation type %s: %v", animationType, err)
		return nil, err
	}
	return frames, nil
}

func (p *PetMeta) PetFramesDrag() ([]string, error) {
	return p.PetFramesBy("drag")
}

func (p *PetMeta) PetFramesDrop() ([]string, error) {
	return p.PetFramesBy("drop")
}

func (p *PetMeta) PetFramesMoveLeft() ([]string, error) {
	return p.PetFramesBy("move_left")
}

func (p *PetMeta) PetFramesMoveRight() ([]string, error) {
	return p.PetFramesBy("move_right")
}

func (p *PetMeta) PetFramesMoveFar() ([]string, error) {
	return p.PetFramesBy("move_far")
}

func (p *PetMeta) PetFramesStand() ([]string, error) {
	return p.PetFramesBy("stand")
}
