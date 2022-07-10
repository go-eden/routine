package g

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestG(t *testing.T) {
	gp1 := G()
	assert.True(t, gp1 != nil)

	t.Run("G in another goroutine", func(t *testing.T) {
		gp2 := G()
		assert.True(t, gp2 != nil)
		assert.True(t, gp1 != gp2)
	})

	gType := reflect.TypeOf(G0())
	sf, ss := gType.FieldByName("labels")
	assert.True(t, ss && sf.Offset > 0)
}
