package main

import "time"

// MarketInfo is the data cg returns when listing all coins
type MarketInfo struct {
	ID                           string    `json:"id"`
	Symbol                       string    `json:"symbol"`
	Name                         string    `json:"name"`
	Image                        string    `json:"image"`
	CurrentPrice                 float64   `json:"current_price"`
	MarketCap                    int64     `json:"market_cap"`
	MarketCapRank                int       `json:"market_cap_rank"`
	FullyDilutedValuation        float64   `json:"fully_diluted_valuation"`
	TotalVolume                  float64   `json:"total_volume"`
	High24H                      float64   `json:"high_24h"`
	Low24H                       float64   `json:"low_24h"`
	PriceChange24H               float64   `json:"price_change_24h"`
	PriceChangePercentage24H     float64   `json:"price_change_percentage_24h"`
	MarketCapChange24H           float64   `json:"market_cap_change_24h"`
	MarketCapChangePercentage24H float64   `json:"market_cap_change_percentage_24h"`
	CirculatingSupply            float64   `json:"circulating_supply"`
	TotalSupply                  float64   `json:"total_supply"`
	MaxSupply                    float64   `json:"max_supply"`
	Ath                          float64   `json:"ath"`
	AthChangePercentage          float64   `json:"ath_change_percentage"`
	AthDate                      time.Time `json:"ath_date"`
	Atl                          float64   `json:"atl"`
	AtlChangePercentage          float64   `json:"atl_change_percentage"`
	AtlDate                      time.Time `json:"atl_date"`
	LastUpdated                  time.Time `json:"last_updated"`
}

// CoinData is a stripped down representation of cg data
type CoinData struct {
	ID              string      `json:"id"`
	Symbol          string      `json:"symbol"`
	Name            string      `json:"name"`
	AssetPlatformID interface{} `json:"asset_platform_id"`
	Platforms       struct {
		string `json:""`
	} `json:"platforms"`
	BlockTimeInMinutes int           `json:"block_time_in_minutes"`
	HashingAlgorithm   string        `json:"hashing_algorithm"`
	Categories         []string      `json:"categories"`
	PublicNotice       interface{}   `json:"public_notice"`
	AdditionalNotices  []interface{} `json:"additional_notices"`
	Localization       struct {
		En   string `json:"en"`
		De   string `json:"de"`
		Es   string `json:"es"`
		Fr   string `json:"fr"`
		It   string `json:"it"`
		Pl   string `json:"pl"`
		Ro   string `json:"ro"`
		Hu   string `json:"hu"`
		Nl   string `json:"nl"`
		Pt   string `json:"pt"`
		Sv   string `json:"sv"`
		Vi   string `json:"vi"`
		Tr   string `json:"tr"`
		Ru   string `json:"ru"`
		Ja   string `json:"ja"`
		Zh   string `json:"zh"`
		ZhTw string `json:"zh-tw"`
		Ko   string `json:"ko"`
		Ar   string `json:"ar"`
		Th   string `json:"th"`
		ID   string `json:"id"`
	} `json:"localization"`
	Description struct {
		En   string `json:"en"`
		De   string `json:"de"`
		Es   string `json:"es"`
		Fr   string `json:"fr"`
		It   string `json:"it"`
		Pl   string `json:"pl"`
		Ro   string `json:"ro"`
		Hu   string `json:"hu"`
		Nl   string `json:"nl"`
		Pt   string `json:"pt"`
		Sv   string `json:"sv"`
		Vi   string `json:"vi"`
		Tr   string `json:"tr"`
		Ru   string `json:"ru"`
		Ja   string `json:"ja"`
		Zh   string `json:"zh"`
		ZhTw string `json:"zh-tw"`
		Ko   string `json:"ko"`
		Ar   string `json:"ar"`
		Th   string `json:"th"`
		ID   string `json:"id"`
	} `json:"description"`
	Links struct {
		Homepage                    []string    `json:"homepage"`
		BlockchainSite              []string    `json:"blockchain_site"`
		OfficialForumURL            []string    `json:"official_forum_url"`
		ChatURL                     []string    `json:"chat_url"`
		AnnouncementURL             []string    `json:"announcement_url"`
		TwitterScreenName           string      `json:"twitter_screen_name"`
		FacebookUsername            string      `json:"facebook_username"`
		BitcointalkThreadIdentifier interface{} `json:"bitcointalk_thread_identifier"`
		TelegramChannelIdentifier   string      `json:"telegram_channel_identifier"`
		SubredditURL                string      `json:"subreddit_url"`
		ReposURL                    struct {
			Github    []string      `json:"github"`
			Bitbucket []interface{} `json:"bitbucket"`
		} `json:"repos_url"`
	} `json:"links"`
	Image struct {
		Thumb string `json:"thumb"`
		Small string `json:"small"`
		Large string `json:"large"`
	} `json:"image"`
	CountryOrigin                string  `json:"country_origin"`
	GenesisDate                  string  `json:"genesis_date"`
	SentimentVotesUpPercentage   float64 `json:"sentiment_votes_up_percentage"`
	SentimentVotesDownPercentage float64 `json:"sentiment_votes_down_percentage"`
	MarketCapRank                int     `json:"market_cap_rank"`
	CoingeckoRank                int     `json:"coingecko_rank"`
	CoingeckoScore               float64 `json:"coingecko_score"`
	DeveloperScore               float64 `json:"developer_score"`
	CommunityScore               float64 `json:"community_score"`
	LiquidityScore               float64 `json:"liquidity_score"`
	PublicInterestScore          float64 `json:"public_interest_score"`
	MarketData                   struct {
		CurrentPrice struct {
			Usd float64 `json:"usd"`
		} `json:"current_price"`
		TotalValueLocked interface{} `json:"total_value_locked"`
		McapToTvlRatio   interface{} `json:"mcap_to_tvl_ratio"`
		FdvToTvlRatio    interface{} `json:"fdv_to_tvl_ratio"`
		Roi              interface{} `json:"roi"`
		Ath              struct {
			Usd float64 `json:"usd"`
		} `json:"ath"`
		AthChangePercentage struct {
			Usd float64 `json:"usd"`
		} `json:"ath_change_percentage"`
		AthDate struct {
			Usd time.Time `json:"usd"`
		} `json:"ath_date"`
		Atl struct {
			Usd float64 `json:"usd"`
		} `json:"atl"`
		AtlChangePercentage struct {
			Usd float64 `json:"usd"`
		} `json:"atl_change_percentage"`
		AtlDate struct {
			Usd time.Time `json:"usd"`
		} `json:"atl_date"`
		MarketCap struct {
			Usd float64 `json:"usd"`
		} `json:"market_cap"`
		MarketCapRank         int `json:"market_cap_rank"`
		FullyDilutedValuation struct {
			Usd float64 `json:"usd"`
		} `json:"fully_diluted_valuation"`
		TotalVolume struct {
			Usd float64 `json:"usd"`
		} `json:"total_volume"`
		High24H struct {
			Usd float64 `json:"usd"`
		} `json:"high_24h"`
		Low24H struct {
			Usd float64 `json:"usd"`
		} `json:"low_24h"`
		PriceChange24H               float64 `json:"price_change_24h"`
		PriceChangePercentage24H     float64 `json:"price_change_percentage_24h"`
		PriceChangePercentage7D      float64 `json:"price_change_percentage_7d"`
		PriceChangePercentage14D     float64 `json:"price_change_percentage_14d"`
		PriceChangePercentage30D     float64 `json:"price_change_percentage_30d"`
		PriceChangePercentage60D     float64 `json:"price_change_percentage_60d"`
		PriceChangePercentage200D    float64 `json:"price_change_percentage_200d"`
		PriceChangePercentage1Y      float64 `json:"price_change_percentage_1y"`
		MarketCapChange24H           float64 `json:"market_cap_change_24h"`
		MarketCapChangePercentage24H float64 `json:"market_cap_change_percentage_24h"`
		PriceChange24HInCurrency     struct {
			Usd float64 `json:"usd"`
		} `json:"price_change_24h_in_currency"`
		PriceChangePercentage1HInCurrency struct {
			Usd float64 `json:"usd"`
		} `json:"price_change_percentage_1h_in_currency"`
		PriceChangePercentage24HInCurrency struct {
			Usd float64 `json:"usd"`
		} `json:"price_change_percentage_24h_in_currency"`
		PriceChangePercentage7DInCurrency struct {
			Usd float64 `json:"usd"`
		} `json:"price_change_percentage_7d_in_currency"`
		PriceChangePercentage14DInCurrency struct {
			Usd float64 `json:"usd"`
		} `json:"price_change_percentage_14d_in_currency"`
		PriceChangePercentage30DInCurrency struct {
			Usd float64 `json:"usd"`
		} `json:"price_change_percentage_30d_in_currency"`
		PriceChangePercentage60DInCurrency struct {
			Usd float64 `json:"usd"`
		} `json:"price_change_percentage_60d_in_currency"`
		PriceChangePercentage200DInCurrency struct {
			Usd float64 `json:"usd"`
		} `json:"price_change_percentage_200d_in_currency"`
		PriceChangePercentage1YInCurrency struct {
			Usd float64 `json:"usd"`
		} `json:"price_change_percentage_1y_in_currency"`
		MarketCapChange24HInCurrency struct {
			Usd float64 `json:"usd"`
		} `json:"market_cap_change_24h_in_currency"`
		MarketCapChangePercentage24HInCurrency struct {
			Usd float64 `json:"usd"`
		} `json:"market_cap_change_percentage_24h_in_currency"`
		TotalSupply       float64   `json:"total_supply"`
		MaxSupply         float64   `json:"max_supply"`
		CirculatingSupply float64   `json:"circulating_supply"`
		LastUpdated       time.Time `json:"last_updated"`
	} `json:"market_data"`
	CommunityData struct {
		FacebookLikes            interface{} `json:"facebook_likes"`
		TwitterFollowers         int         `json:"twitter_followers"`
		RedditAveragePosts48H    float64     `json:"reddit_average_posts_48h"`
		RedditAverageComments48H float64     `json:"reddit_average_comments_48h"`
		RedditSubscribers        int         `json:"reddit_subscribers"`
		RedditAccountsActive48H  int         `json:"reddit_accounts_active_48h"`
		TelegramChannelUserCount interface{} `json:"telegram_channel_user_count"`
	} `json:"community_data"`
	DeveloperData struct {
		Forks                        int `json:"forks"`
		Stars                        int `json:"stars"`
		Subscribers                  int `json:"subscribers"`
		TotalIssues                  int `json:"total_issues"`
		ClosedIssues                 int `json:"closed_issues"`
		PullRequestsMerged           int `json:"pull_requests_merged"`
		PullRequestContributors      int `json:"pull_request_contributors"`
		CodeAdditionsDeletions4Weeks struct {
			Additions int `json:"additions"`
			Deletions int `json:"deletions"`
		} `json:"code_additions_deletions_4_weeks"`
		CommitCount4Weeks              int   `json:"commit_count_4_weeks"`
		Last4WeeksCommitActivitySeries []int `json:"last_4_weeks_commit_activity_series"`
	} `json:"developer_data"`
	PublicInterestStats struct {
		AlexaRank   int         `json:"alexa_rank"`
		BingMatches interface{} `json:"bing_matches"`
	} `json:"public_interest_stats"`
	StatusUpdates []interface{} `json:"status_updates"`
	LastUpdated   time.Time     `json:"last_updated"`
}
