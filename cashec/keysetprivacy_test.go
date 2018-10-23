package cashec

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"errors"
)

func TestGenerateKey(t *testing.T) {
	keySet := KeySet{}
	keySet.GenerateKey([]byte{1, 2, 3})
	sign, _ := keySet.Sign([]byte("hello"))
	valid, _ := keySet.Verify([]byte("hello"), sign)
	if valid != true {
		assert.Error(t, errors.New("Sign fail"))
	}
}
