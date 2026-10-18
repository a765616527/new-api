package service

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestReportCurrentSystemInstanceIncludesActiveConnections(t *testing.T) {
	previousDB, previousNodeName := model.DB, common.NodeName
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.SystemInstance{}))
	model.DB = db
	common.NodeName = "test-node"
	common.AddActiveConnections(2)
	t.Cleanup(func() {
		common.AddActiveConnections(-2)
		model.DB = previousDB
		common.NodeName = previousNodeName
	})

	require.NoError(t, ReportCurrentSystemInstance())

	instances, err := model.ListSystemInstances()
	require.NoError(t, err)
	require.Len(t, instances, 1)
	var info SystemInstanceInfo
	require.NoError(t, common.UnmarshalJsonStr(instances[0].Info, &info))
	assert.Equal(t, "test-node", info.Node.Name)
	assert.Equal(t, int64(2), info.HTTP.ActiveConnections)
}
