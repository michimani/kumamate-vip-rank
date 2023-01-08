package kvr

type VIPRank struct {
	Min   uint64
	Max   uint64
	Title string
}

type VIPRankList []VIPRank

func (vl *VIPRankList) GetTitle(power uint64) string {
	for _, v := range *vl {
		if (v.Max == 0 && v.Min <= power) || (v.Max > 0 && v.Min <= power && v.Max >= power) {
			return v.Title
		}
	}

	return ""
}
