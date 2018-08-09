package cli

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestCommand(t *testing.T) {

	t.Run("should add child command and found the child", func(t *testing.T) {
		fooCmd := new(fooCommand)
		barCmd := new(barCommand)
		fooCmd.SetName("foo")
		fooCmd.Add(barCmd)

		foundCmd, err := fooCmd.Find("bar")
		assert.Equal(t, nil, err)
		assert.Equal(t, foundCmd.GetName(), barCmd.GetName())
	})

	t.Run("should run handle method", func(t *testing.T) {
		cmd := new(BaseCommand)
		err := cmd.Run(nil)
		assert.Equal(t, nil, err)
	})
}

