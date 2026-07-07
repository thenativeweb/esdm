package version_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thenativeweb/esdm/cmd/cmdutils"
	"github.com/thenativeweb/esdm/cmd/esdm/commands/version"
)

func TestCommand(t *testing.T) {
	t.Run("prints the version and git revision", func(t *testing.T) {
		cmdutils.Version = "v1.2.3"
		cmdutils.GitVersion = "abc123"

		var stdout bytes.Buffer
		version.Command.SetOut(&stdout)
		version.Command.SetArgs([]string{})

		err := version.Command.Execute()
		require.NoError(t, err)

		output := stdout.String()
		assert.Contains(t, output, "esdm: v1.2.3")
		assert.Contains(t, output, "Revision: abc123")
	})
}
