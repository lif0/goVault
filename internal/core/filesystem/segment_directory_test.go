package filesystem

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func cleanUp(path string) {
	if parentDir := filepath.Dir(path); parentDir != "." {
		path = parentDir
	}

	os.RemoveAll(path)
}

func TestSegmentsDirectoryDirectoryNotExists(t *testing.T) {
	t.Parallel()

	segmentsCount := 0
	expectedSegmentsCount := 3

	segmentDir, err := NewSegmentsDirectory("data/wal")
	assert.Nil(t, err)

	generateTestData(segmentDir.directory+"/"+t.Name()+"1.txt", 2)
	generateTestData(segmentDir.directory+"/"+t.Name()+"2.txt", 2)
	generateTestData(segmentDir.directory+"/"+t.Name()+"3.txt", 2)

	err = segmentDir.ForEach(func(data []byte) error {
		assert.True(t, len(data) != 0)
		segmentsCount++
		return nil
	})

	require.NoError(t, err)
	assert.Equal(t, expectedSegmentsCount, segmentsCount)

	cleanUp(segmentDir.directory)
}

func TestSegmentsDirectoryForEach(t *testing.T) {
	t.Parallel()

	segmentsCount := 0
	expectedSegmentsCount := 3

	segmentDir, err := NewSegmentsDirectory(t.TempDir())
	assert.Nil(t, err)

	generateTestData(segmentDir.directory+"/"+t.Name()+"1.txt", 2)
	generateTestData(segmentDir.directory+"/"+t.Name()+"2.txt", 2)
	generateTestData(segmentDir.directory+"/"+t.Name()+"3.txt", 2)

	err = segmentDir.ForEach(func(data []byte) error {
		assert.True(t, len(data) != 0)
		segmentsCount++
		return nil
	})

	require.NoError(t, err)
	assert.Equal(t, expectedSegmentsCount, segmentsCount)

	cleanUp(segmentDir.directory)
}

func TestSegmentsDirectoryForEachWithBreak(t *testing.T) {
	t.Parallel()

	segmentDir, err := NewSegmentsDirectory("test_data")
	assert.Nil(t, err)

	generateTestData(segmentDir.directory+"/"+t.Name()+"1.txt", 2)
	generateTestData(segmentDir.directory+"/"+t.Name()+"2.txt", 2)
	generateTestData(segmentDir.directory+"/"+t.Name()+"3.txt", 2)

	err = segmentDir.ForEach(func([]byte) error {
		return errors.New("error")
	})

	assert.Error(t, err, "error")

	cleanUp(segmentDir.directory)
}

func generateTestData(pathWithFile string, countLine int) {
	file, _ := os.Create(pathWithFile)
	for i := 1; i <= countLine; i++ {
		line := fmt.Sprintf("some data %d\n", i)
		_, _ = file.WriteString(line)
	}
	_ = file.Sync()
	_ = file.Close()
}
