package game

import (
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

// Global pool for DrawImageOptions to reduce allocations
var drawOptionsPool = sync.Pool{
	New: func() interface{} {
		return &ebiten.DrawImageOptions{}
	},
}

// GetDrawOptions retrieves a DrawImageOptions from the pool
// The returned object is reset to a clean state
func GetDrawOptions() *ebiten.DrawImageOptions {
	op := drawOptionsPool.Get().(*ebiten.DrawImageOptions)
	// Reset to clean state
	op.GeoM.Reset()
	op.ColorM.Reset()
	return op
}

// PutDrawOptions returns a DrawImageOptions to the pool for reuse
func PutDrawOptions(op *ebiten.DrawImageOptions) {
	if op != nil {
		drawOptionsPool.Put(op)
	}
}
