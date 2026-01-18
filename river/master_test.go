package river

import (
	"os"
	"testing"
	"time"

	"github.com/siddontang/go-mysql/mysql"
    "github.com/stretchr/testify/assert"
)

func TestMasterInfoSaveEmptyName(t *testing.T) {
	dataDir := "./var_test_master_info"
	os.RemoveAll(dataDir)
	defer os.RemoveAll(dataDir)

	m, err := loadMasterInfo(dataDir)
	assert.NoError(t, err)

    // 1. Save a valid position
	pos := mysql.Position{Name: "mysql-bin.000001", Pos: 1234}
	err = m.Save(pos)
	assert.NoError(t, err)
    
    // Force save (bypass time check)
    m.lastSaveTime = time.Unix(0, 0)

	// Verify it's saved
    assert.Equal(t, "mysql-bin.000001", m.Name)
    assert.Equal(t, uint32(1234), m.Pos)

	// 2. Save an empty name position (which caused the bug)
	emptyPos := mysql.Position{Name: "", Pos: 0}
	err = m.Save(emptyPos)
	assert.NoError(t, err)

	// 3. Verify it is NOT overwritten (fix works)
	assert.Equal(t, "mysql-bin.000001", m.Name)
	assert.Equal(t, uint32(1234), m.Pos)
}
