package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildMetadata(t *testing.T) {
	assert.Equal(t, Version, "0.9")
	assert.Equal(t, BuildTime, "Sun, 11 Dec 2022 18:42:01 +0000")
}
