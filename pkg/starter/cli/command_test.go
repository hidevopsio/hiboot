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
		barCmd.SetName("bar")
		fooCmd.AddChild(barCmd)

		foundCmd, err := fooCmd.Find("bar")
		assert.Equal(t, nil, err)
		assert.Equal(t, foundCmd.Name(), barCmd.name)
	})

	t.Run("should run handle method", func(t *testing.T) {
		cmd := new(BaseCommand)
		err := cmd.Handle(nil)
		assert.Equal(t, nil, err)
	})
}

