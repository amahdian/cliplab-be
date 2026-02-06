package model

const FreeCreditAmount = 200

type CreditRule struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Amount   int    `json:"amount"`
	Key      string `json:"key"`
}

const (
	CreditKeyEngagementRate      = "engagement_rate"
	CreditKeyEngagementBreakdown = "engagement_breakdown"
	CreditKeyReelAnalyze         = "reel_analyze"
)

var CreditRules = []CreditRule{
	{
		Name:     "Engagement Rate",
		Category: "Engagement Rate",
		Amount:   0,
		Key:      CreditKeyEngagementRate,
	},
	{
		Name:     "Engagement Rate Breakdown",
		Category: "Engagement Rate",
		Amount:   1,
		Key:      CreditKeyEngagementBreakdown,
	},
	{
		Name:     "Reel Analyze",
		Category: "Reel Analyze",
		Amount:   10,
		Key:      CreditKeyReelAnalyze,
	},
}

func GetCreditRule(key string) *CreditRule {
	for _, rule := range CreditRules {
		if rule.Key == key {
			return &rule
		}
	}
	return nil
}
