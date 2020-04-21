package huobi

const (
	ACCESS_KEY         = "AccessKeyId"
	SIGNATURE          = "Signature"
	SIGNNATURE_METHOD  = "SignatureMethod"
	SIGNNATURE_VERSION = "SignatureVersion"
	TIMESTAMP          = "Timestamp"
)

type SpotSymbol struct {
	Amount_precision int     `json:"amount-precision"`
	Base_currency    string  `json:"base-currency"`
	Max_order_amt    int     `json:"max-order-amt"`
	Min_order_amt    float32 `json:"min-order-amt"`
	Min_order_value  float32 `json:"min-order-value"`
	Price_precision  int     `json:"price-precision"`
	Quote_currency   string  `json:"quote-currency"`
	State            string  `json:"state"`
	Symbol           string  `json:"symbol"`
	Symbol_partition string  `json:"symbol-partition"`
	Value_precision  int     `json:"value-precision"`
}

type Chain struct {
	Chain                  string `json:"chain"`
	DepositStatus          string `json:"depositStatus"`
	MaxWithdrawAmt         string `json:"maxWithdrawAmt"`
	MinDepositAmt          string `json:"minDepositAmt"`
	MinWithdrawAmt         string `json:"minWithdrawAmt"`
	NumOfConfirmations     int    `json:"numOfConfirmations"`
	NumOfFastConfirmations int    `json:"numOfFastConfirmations"`
	TransactFeeWithdraw    string `json:"transactFeeWithdraw"`
	WithdrawFeeType        string `json:"withdrawFeeType"`
	WithdrawPrecision      int    `json:"withdrawPrecision"`
	WithdrawQuotaPerDay    string `json:"withdrawQuotaPerDay"`
	WithdrawQuotaPerYear   string `json:"withdrawQuotaPerYear"`
	WithdrawQuotaTotal     string `json:"withdrawQuotaTotal"`
	WithdrawStatus         string `json:"withdrawStatus"`
}

type Currency struct {
	Chains     []Chain `json:"chains"`
	Currency   string  `json:"currency"`
	InstStatus string  `json:"instStatus"`
}

type KLine struct {
	Amount float32 `json:"amount"`
	Close  float32 `json:"close"`
	Count  int     `json:"count"`
	High   float32 `json:"high"`
	ID     int     `json:"id"`
	Low    float32 `json:"low"`
	Open   float32 `json:"open"`
	Vol    float32 `json:"vol"`
}

type SpotMergin struct {
	Amount  float32   `json:"amount"`
	Ask     []float32 `json:"ask"`
	Bid     []float32 `json:"bid"`
	Close   float32   `json:"close"`
	Count   int       `json:"count"`
	High    float32   `json:"high"`
	ID      int       `json:"id"`
	Low     float32   `json:"low"`
	Open    float32   `json:"open"`
	Version int       `json:"version"`
	Vol     float32   `json:"vol"`
}

type SpotTick struct {
	Amount float32 `json:"amount"`
	Close  float32 `json:"close"`
	Count  float32 `json:"count"`
	High   float32 `json:"high"`
	Low    float32 `json:"low"`
	Open   float32 `json:"open"`
	Symbol string  `json:"symbol"`
	Vol    float32 `json:"vol"`
}

type SpotDepth struct {
	Asks    [][]float64 `json:"asks"`
	Bids    [][]float64 `json:"bids"`
	ID      int         `json:"id"`
	Ts      int         `json:"ts"`
	Version int         `json:"version"`
}

type SpotTrade struct {
	Record []struct {
		Amount    float32 `json:"amount"`
		Direction string  `json:"direction"`
		Price     float32 `json:"price"`
		Trade_id  int64   `json:"trade-id"`
		Ts        int64   `json:"ts"`
	} `json:"data"`
	ID int64 `json:"id"`
	Ts int64 `json:"ts"`
}

type SpotDetail struct {
	Amount  float64 `json:"amount"`
	Close   float64 `json:"close"`
	Count   int     `json:"count"`
	High    float64 `json:"high"`
	ID      int     `json:"id"`
	Low     float64 `json:"low"`
	Open    float64 `json:"open"`
	Version int     `json:"version"`
	Vol     float64 `json:"vol"`
}

type SpotAccount struct {
	ID      int    `json:"id"`
	State   string `json:"state"`
	Subtype string `json:"subtype"`
	Type    string `json:"type"`
}

type SpotBalance struct {
	ID   int `json:"id"`
	List []struct {
		Balance  string `json:"balance"`
		Currency string `json:"currency"`
		Type     string `json:"type"`
	} `json:"list"`

	State string `json:"state"`
	Type  string `json:"type"`
}

type SpotAggregateBalance struct {
	Balance  string `json:"balance"`
	Currency string `json:"currency"`
	Type     string `json:"type"`
}
