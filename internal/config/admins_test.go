package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadAdminVKUserIDs(t *testing.T) {
	t.Run("missing file", func(t *testing.T) {
		ids, err := LoadAdminVKUserIDs(filepath.Join(t.TempDir(), "admins.txt"))
		require.NoError(t, err)
		assert.Empty(t, ids)
	})

	t.Run("empty file", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "admins.txt")
		require.NoError(t, os.WriteFile(path, nil, 0o644))

		ids, err := LoadAdminVKUserIDs(path)
		require.NoError(t, err)
		assert.Empty(t, ids)
	})

	t.Run("only comments", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "admins.txt")
		require.NoError(t, os.WriteFile(path, []byte("# нет админов\n\n"), 0o644))

		ids, err := LoadAdminVKUserIDs(path)
		require.NoError(t, err)
		assert.Empty(t, ids)
	})

	t.Run("success", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "admins.txt")
		body := "# главные админы\n\n  100  \n100\n200\n"
		require.NoError(t, os.WriteFile(path, []byte(body), 0o644))

		ids, err := LoadAdminVKUserIDs(path)
		require.NoError(t, err)
		assert.Equal(t, map[int64]struct{}{100: {}, 200: {}}, ids)
	})

	t.Run("failure (text)", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "admins.txt")
		require.NoError(t, os.WriteFile(path, []byte("abc\n"), 0o644))

		ids, err := LoadAdminVKUserIDs(path)
		assert.Error(t, err)
		assert.Nil(t, ids)
	})

	t.Run("failure (negative)", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "admins.txt")
		require.NoError(t, os.WriteFile(path, []byte("-5\n"), 0o644))

		ids, err := LoadAdminVKUserIDs(path)
		assert.Error(t, err)
		assert.Nil(t, ids)
	})

	t.Run("failure (bad id)", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "admins.txt")
		require.NoError(t, os.WriteFile(path, []byte("0\n"), 0o644))

		ids, err := LoadAdminVKUserIDs(path)
		assert.Error(t, err)
		assert.Nil(t, ids)
	})
}
