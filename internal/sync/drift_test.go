package sync_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yourusername/vaultpull/internal/sync"
)

func TestNewDriftDetector_NotNil(t *testing.T) {
	store := synct.TempDir() + "/snap.json")
	detector := sync.NewDrifttrequire.NotNil(t, detector)
}

func TestDriftDetector_NoDrift_WhenSnapshotMatches(t *testing.T) {
	store := sync.NewSnapshotStore(t.TempDir() + "/snap.json")
	detector := sync.NewDriftDetector(store)

	current := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}

	err := store.Save("default", current)
	require.NoError(t, err)

	drifted, diff, err := detector.Detect("default", current)
	require.NoError(t, err)
	assert.False(t, drifted)
	assert.Empty(t, diff)
}

func TestDriftDetector_DetectsDrift_WhenValueChanged(t *testing.T) {
	store := sync.NewSnapshotStore(t.TempDir() + "/snap.json")
	detector := sync.NewDriftDetector(store)

	saved := map[string]string{"API_KEY": "old-key"}
	err := store.Save("prod", saved)
	require.NoError(t, err)

	current := map[string]string{"API_KEY": "new-key"}
	drifted, diff, err := detector.Detect("prod", current)
	require.NoError(t, err)
	assert.True(t, drifted)
	assert.Contains(t, diff, "API_KEY")
}

func TestDriftDetector_DetectsDrift_WhenKeyAdded(t *testing.T) {
	store := sync.NewSnapshotStore(t.TempDir() + "/snap.json")
	detector := sync.NewDriftDetector(store)

	saved := map[string]string{"FOO": "bar"}
	err := store.Save("dev", saved)
	require.NoError(t, err)

	current := map[string]string{"FOO": "bar", "NEW_KEY": "value"}
	drifted, diff, err := detector.Detect("dev", current)
	require.NoError(t, err)
	assert.True(t, drifted)
	assert.Contains(t, diff, "NEW_KEY")
}

func TestDriftDetector_DetectsDrift_WhenKeyRemoved(t *testing.T) {
	store := sync.NewSnapshotStore(t.TempDir() + "/snap.json")
	detector := sync.NewDriftDetector(store)

	saved := map[string]string{"FOO": "bar", "GONE": "value"}
	err := store.Save("staging", saved)
	require.NoError(t, err)

	current := map[string]string{"FOO": "bar"}
	drifted, diff, err := detector.Detect("staging", current)
	require.NoError(t, err)
	assert.True(t, drifted)
	assert.Contains(t, diff, "GONE")
}

func TestDriftDetector_NoSnapshot_ReturnsNoDrift(t *testing.T) {
	store := sync.NewSnapshotStore(t.TempDir() + "/snap.json")
	detector := sync.NewDriftDetector(store)

	current := map[string]string{"X": "1"}
	drifted, diff, err := detector.Detect("unknown-profile", current)
	require.NoError(t, err)
	assert.False(t, drifted)
	assert.Empty(t, diff)
}

func TestDriftDetector_LastChecked_UpdatedAfterDetect(t *testing.T) {
	store := sync.NewSnapshotStore(t.TempDir() + "/snap.json")
	detector := sync.NewDriftDetector(store)

	before := time.Now()
	_, _, err := detector.Detect("any", map[string]string{})
	require.NoError(t, err)

	last := detector.LastChecked()
	assert.True(t, last.Equal(before) || last.After(before))
}
