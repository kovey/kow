package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanPath(t *testing.T) {
	assert.Equal(t, "/", CleanPath(""))
	assert.Equal(t, "/user", CleanPath("user"))
	assert.Equal(t, "/user/name", CleanPath("./user/index/../name"))
	assert.Equal(t, "/name", CleanPath("../user/../index/../name"))
	assert.Equal(t, "/user/name", CleanPath("./user/name"))
	assert.Equal(t, "/user/name", CleanPath("/user/./name"))
}
