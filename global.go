package irapi

import "context"

type GlobalStatsResponse struct {
	DriverCounts struct {
		Total    uint64
		LapCount StringifiedUint64
	}
}

func (i *IRacing) GlobalStats(ctx context.Context) {

	// path := "/membersite/member/GetDriverCounts"

}
