package irapi

import (
	"context"
	"net/http"
)

// UserProfile represents the profile for the **current** user.
//
// iRacing provides much more data, in a different format, for the current user.
type UserProfile struct {
	ID       uint64 `json:"custID"`
	Licenses []struct {
		Group         int    `json:"licGroup"`
		Level         int    `json:"licLevel"`
		SRSub         int    `json:"srSub,string"`
		SRPrime       int    `json:"srPrime,string"`
		TTRating      Rating `json:"ttRating"`
		MPRRaces      int    `json:"mprNumRaces"`
		MPRTimeTrials int    `json:"mprNumTTs"`
		IRating       Rating `json:"iRating"`
		CatID         int    `json:"catId"`
		DisplayName   string `json:"licLevelDisplayName"`
		Color         string `json:"licColor"`
		GroupName     string `json:"licGroupDisplayName"`
	} `json:"licenses"`

	DisplayName            string `json:"DisplayName"`
	HasReadPrivacyPolicy   bool   `json:"hasReadPP"`
	HasReadTermsConditions bool   `json:"hasReadTC"`
}

// GetProfile gets the user's profile
func (c *IRacing) GetProfile(ctx context.Context) (*UserProfile, error) {

	profile := &UserProfile{}

	err := c.json(ctx, http.MethodGet, "/membersite/member/GetMember", nil, profile)

	return profile, err
}
