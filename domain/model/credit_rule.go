package model

const FreeCreditAmount = 200
const StarterCreditAmount = 2000
const CreatorCreditAmount = 8000
const StudioCreditAmount = 20000

var PlanCredits = map[string]int{
	"pri_01kfxyf1fet153za7jmgj4dd0p": StarterCreditAmount,
	"pri_01kfxyg0qcrqk571v48vjzzy08": CreatorCreditAmount,
	"pri_01kfxybm6ee1ccp1fa4hd8gf7d": StudioCreditAmount,
}

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
