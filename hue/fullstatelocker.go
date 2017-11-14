package hue

import (
	"sync"
)

type fullStateLocker struct {
	fs *fullState
	m  sync.Mutex
}

func NewFullStateLocker(persistedHueStatePath string, nii NetworkInterfaceInfo) *fullStateLocker {

	return &fullStateLocker{
		fs: getFullState(persistedHueStatePath, nii),
	}
}

func (fsl *fullStateLocker) Lock() *fullState {

	fsl.m.Lock()
	return fsl.fs
}

func (fsl *fullStateLocker) Unlock() {

	fsl.m.Unlock()
}
