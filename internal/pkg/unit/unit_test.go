package unit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseDigitalStorageWithBytes(t *testing.T) {
	t.Parallel()

	size, err := ParseDigitalStorage("100B")
	require.NoError(t, err)
	require.Equal(t, 100, size)

	size, err = ParseDigitalStorage("100b")
	require.NoError(t, err)
	require.Equal(t, 100, size)

	size, err = ParseDigitalStorage("100")
	require.NoError(t, err)
	require.Equal(t, 100, size)
}

func TestParseDigitalStorageWithKiloBytes(t *testing.T) {
	t.Parallel()

	size, err := ParseDigitalStorage("1024KB")
	require.NoError(t, err)
	require.Equal(t, 1024*1024, size)

	size, err = ParseDigitalStorage("1024Kb")
	require.NoError(t, err)
	require.Equal(t, 1024*1024, size)

	size, err = ParseDigitalStorage("1024kb")
	require.NoError(t, err)
	require.Equal(t, 1024*1024, size)
}

func TestParseDigitalStorageWithMegaBytes(t *testing.T) {
	t.Parallel()

	size, err := ParseDigitalStorage("512MB")
	require.NoError(t, err)
	require.Equal(t, 512*1024*1024, size)

	size, err = ParseDigitalStorage("512Mb")
	require.NoError(t, err)
	require.Equal(t, 512*1024*1024, size)

	size, err = ParseDigitalStorage("512mb")
	require.NoError(t, err)
	require.Equal(t, 512*1024*1024, size)
}

func TestParseDigitalStorageWithGigaBytes(t *testing.T) {
	t.Parallel()

	size, err := ParseDigitalStorage("10GB")
	require.NoError(t, err)
	require.Equal(t, 10*1024*1024*1024, size)

	size, err = ParseDigitalStorage("10Gb")
	require.NoError(t, err)
	require.Equal(t, 10*1024*1024*1024, size)

	size, err = ParseDigitalStorage("10gb")
	require.NoError(t, err)
	require.Equal(t, 10*1024*1024*1024, size)
}

func TestParseIncorrectDigitalStorage(t *testing.T) {
	t.Parallel()

	_, err := ParseDigitalStorage("-10")
	require.Error(t, err)

	_, err = ParseDigitalStorage("-10TB")
	require.Error(t, err)

	_, err = ParseDigitalStorage("abc")
	require.Error(t, err)

	_, err = ParseDigitalStorage("10TB")
	require.Equal(t, err.Error(), "incorrect size")
}
