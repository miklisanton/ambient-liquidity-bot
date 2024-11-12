package services

import (
    "testing"
)


var ambientService *AmbientService

func TestMain(m *testing.M) {
    ambientService = NewAmbientService()
    m.Run()
}

func TestGetUserPools(t *testing.T) {
    pools, err := ambientService.GetUserPools("0x2e1389F741dD0651ea8D235F0460d6C5d15cDFFe")
    if err != nil {
        t.Errorf("failed to get user pools: %v", err)
    }
    for _, pool := range pools {
        t.Logf("%+v", pool)
        t.Log(len(pool.PositionId))
    }
}

