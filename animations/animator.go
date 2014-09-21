package animations

import (
	"errors"
)

type Animator interface {
	BeginAnimation()
	EndAnimation()
}

var (
	ErrInvalidAnimationName = errors.New("invalid animation name")
)

func Animate(animationName string) error {
	var animator Animator
	switch {
	case animationName == "runner":
		animator = new(Runner);
	}

	if animator == nil {
		return ErrInvalidAnimationName
	}

	// Stop existing animations

	// Start the animation on its own go routine


	return nil
}
