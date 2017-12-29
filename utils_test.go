package growl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsFileExist(t *testing.T) {
	assert.Equal(t, true, IsFileExist("utils.go"))
	assert.Equal(t, false, IsFileExist("utils_dummy.go"))
}
