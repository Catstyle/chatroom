package utils_test

import (
	"testing"

	"github.com/catstyle/chatroom/utils"
	"github.com/stretchr/testify/assert"
)

func TestFilterSensitiveText(t *testing.T) {
	trietree := utils.NewTrieTree()
	assert.NotNil(t, trietree)

	words := []string{
		"hell",
		"asshole",
		"assholes",
	}
	trietree.AddWord(words...)

	assert.Equal(t, "****boy", trietree.Filter("hellboy", '*'))
	assert.Equal(t, "*******", trietree.Filter("asshole", '*'))
	assert.Equal(t, "********", trietree.Filter("assholes", '*'))
	assert.Equal(t, "****boy, *******", trietree.Filter("hellboy, asshole", '*'))
	assert.Equal(t, "assboy", trietree.Filter("assboy", '*'))
}
