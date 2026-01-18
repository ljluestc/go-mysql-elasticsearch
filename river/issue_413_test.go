package river

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/siddontang/go-mysql/mysql"
	"github.com/stretchr/testify/require"
)

func TestIssue413(t *testing.T) {
	// Create a minimal config
	cfg := new(Config)
	cfg.DataDir = t.TempDir()
	cfg.MyAddr = "127.0.0.1:3306"
	cfg.MyUser = "root"
	cfg.MyPassword = ""
	cfg.ESAddr = "127.0.0.1:9200"

	// Create a dummy master.info file
	masterInfoPath := cfg.DataDir + "/master.info"
	initialContent := `bin_name = "mysql-bin.013708"
bin_pos = 1234
`
	err := os.WriteFile(masterInfoPath, []byte(initialContent), 0644)
	require.NoError(t, err)

	// Initialize River
	r := &River{
		c:      cfg,
		syncCh: make(chan interface{}, 1024),
		ctx:    nil,
	}
	r.ctx, r.cancel = context.WithCancel(context.Background())
	r.wg.Add(1)

	// Load master info manually
	r.master, err = loadMasterInfo(cfg.DataDir)
	require.NoError(t, err)

	// Verify initial position
	pos := r.master.Position()
	require.Equal(t, "mysql-bin.013708", pos.Name)
	require.Equal(t, uint32(1234), pos.Pos)

	// Run syncLoop in a goroutine
	go r.syncLoop()

	// Send an empty posSaver
	// This simulates what is suspected: an empty position being passed to syncLoop
	r.syncCh <- posSaver{
		pos:   mysql.Position{Name: "", Pos: 0},
		force: true,
	}

	// Give it some time to process
	time.Sleep(100 * time.Millisecond)

	// Stop syncLoop
	r.cancel()
	r.wg.Wait()

	// Check master info
	newPos := r.master.Position()
	
	// Expectation for the fix: it should NOT be empty.
	require.NotEmpty(t, newPos.Name, "Position name should not be empty after syncLoop")
	require.Equal(t, "mysql-bin.013708", newPos.Name)
}
