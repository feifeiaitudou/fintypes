package comm

import "github.com/shawnwyckoff/commpkg/apputil/errorz"

type (
	HighlyRelevant struct {
		Items []PairExt
	}

	ModeratelyRelevant struct {
		Highs []HighlyRelevant
	}

	LowlyRelevant struct {
		Mids []ModeratelyRelevant
	}

	Irrelevant struct {
		Lows []LowlyRelevant
	}

	Portfolio struct {
		RelevantSet *Irrelevant `json:"RelevantSet,omitempty"`
		DirectSet   []PairExt   `json:"DirectSet,omitempty"`
	}
)

func (p Portfolio) Verify() error {
	if len(p.RelevantSet.Lows) == 0 && len(p.DirectSet) == 0 {
		return errorz.Errorf("empty Portfolio")
	}

	if len(p.RelevantSet.Lows) > 0 {
		rMap := map[PairExt]bool{}
		total := 0
		for _, LrSet := range p.RelevantSet.Lows {
			for _, MidSet := range LrSet.Mids {
				for _, HighSet := range MidSet.Highs {
					for _, target := range HighSet.Items {
						rMap[target] = true
						total++
					}
				}
			}
		}
		if total != len(rMap) {
			return errorz.Errorf("these is duplicated targets")
		}
	}

	return nil
}

func (p Portfolio) Targets() []PairExt {
	if len(p.DirectSet) > 0 {
		return p.DirectSet
	}

	if len(p.RelevantSet.Lows) > 0 {
		rMap := map[PairExt]bool{}
		total := 0
		for _, LrSet := range p.RelevantSet.Lows {
			for _, MidSet := range LrSet.Mids {
				for _, HighSet := range MidSet.Highs {
					for _, target := range HighSet.Items {
						rMap[target] = true
						total++
					}
				}
			}
		}

		var r []PairExt
		for k := range rMap {
			r = append(r, k)
		}
		return r
	}
	return nil
}
