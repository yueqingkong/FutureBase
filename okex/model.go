package okex

import "time"

const (
	CONTENT_TYPE = "Content-Type"
	ACCEPT       = "Accept"
	COOKIE       = "Cookie"
	LOCALE       = "locale="

	OK_ACCESS_KEY        = "OK-ACCESS-KEY"
	OK_ACCESS_SIGN       = "OK-ACCESS-SIGN"
	OK_ACCESS_TIMESTAMP  = "OK-ACCESS-TIMESTAMP"
	OK_ACCESS_PASSPHRASE = "OK-ACCESS-PASSPHRASE"

	APPLICATION_JSON      = "application/json"
	APPLICATION_JSON_UTF8 = "application/json; charset=UTF-8"

	ENGLISH = "en_US"
)

type Instrument []struct {
	BaseCurrency   string `json:"base_currency"`
	BaseIncrement  string `json:"base_increment"`
	BaseMinSize    string `json:"base_min_size"`
	InstrumentID   string `json:"instrument_id"`
	MinSize        string `json:"min_size"`
	ProductID      string `json:"product_id"`
	QuoteCurrency  string `json:"quote_currency"`
	QuoteIncrement string `json:"quote_increment"`
	SizeIncrement  string `json:"size_increment"`
	TickSize       string `json:"tick_size"`
}

type SpotCandles [][]interface{}

type Depth struct {
	Asks      [][]string `json:"asks"`
	Bids      [][]string `json:"bids"`
	Timestamp time.Time  `json:"timestamp"`
}

type SpotTicker struct {
	Ask            string `json:"ask"`
	BaseVolume24h  string `json:"base_volume_24h"`
	BestAsk        string `json:"best_ask"`
	BestBid        string `json:"best_bid"`
	Bid            string `json:"bid"`
	High24h        string `json:"high_24h"`
	InstrumentID   string `json:"instrument_id"`
	Last           string `json:"last"`
	Low24h         string `json:"low_24h"`
	Open24h        string `json:"open_24h"`
	ProductID      string `json:"product_id"`
	QuoteVolume24h string `json:"quote_volume_24h"`
	Timestamp      string `json:"timestamp"`
}

type Trade struct {
	Time      time.Time `json:"time"`
	Timestamp time.Time `json:"timestamp"`
	TradeID   string    `json:"trade_id"`
	Price     string    `json:"price"`
	Size      string    `json:"size"`
	Side      string    `json:"side"`
}

type Candles [][]interface{}

type FutureIndex struct {
	InstrumentID string    `json:"instrument_id"`
	Index        string    `json:"index"`
	Timestamp    time.Time `json:"timestamp"`
}

type FutureRate struct {
	InstrumentID string    `json:"instrument_id"`
	Rate         string    `json:"rate"`
	Timestamp    time.Time `json:"timestamp"`
}

type FutureEstimatedPrice struct {
	InstrumentID    string    `json:"instrument_id"`
	SettlementPrice string    `json:"settlement_price"`
	Timestamp       time.Time `json:"timestamp"`
}

type FutureOpenInterest struct {
	InstrumentID string    `json:"instrument_id"`
	Amount       string    `json:"amount"`
	Timestamp    time.Time `json:"timestamp"`
}

type FuturePriceLimit struct {
	InstrumentID string    `json:"instrument_id"`
	Highest      string    `json:"highest"`
	Lowest       string    `json:"lowest"`
	Timestamp    time.Time `json:"timestamp"`
}

type FutureMarkPrice struct {
	MarkPrice    string    `json:"mark_price"`
	InstrumentID string    `json:"instrument_id"`
	Timestamp    time.Time `json:"timestamp"`
}

type FutureLiquidation []struct {
	Loss         float32   `json:"loss"`
	Size         int       `json:"size"`
	Price        float32   `json:"price"`
	CreatedAt    time.Time `json:"created_at"`
	Type         int       `json:"type"`
	InstrumentID string    `json:"instrument_id"`
}

type FutureHold struct {
	InstrumentID string    `json:"instrument_id"`
	Timestamp    time.Time `json:"timestamp"`
	Amount       string    `json:"amount"`
}

type FutureFill struct {
	TradeID      string    `json:"trade_id"`
	InstrumentID string    `json:"instrument_id"`
	Price        string    `json:"price"`
	OrderQty     string    `json:"order_qty"`
	OrderID      string    `json:"order_id"`
	CreatedAt    time.Time `json:"created_at"`
	ExecType     string    `json:"exec_type"`
	Fee          string    `json:"fee"`
	Side         string    `json:"side"`
}

type FutureInstrument struct {
	InstrumentID        string `json:"instrument_id"`
	UnderlyingIndex     string `json:"underlying_index"`
	QuoteCurrency       string `json:"quote_currency"`
	TickSize            string `json:"tick_size"`
	ContractVal         string `json:"contract_val"`
	Listing             string `json:"listing"`
	Delivery            string `json:"delivery"`
	TradeIncrement      string `json:"trade_increment"`
	Alias               string `json:"alias"`
	IsInverse           string `json:"is_inverse"`
	ContractValCurrency string `json:"contract_val_currency"`
}

type FutureDepth struct {
	Bids      [][]int   `json:"bids"`
	Asks      [][]int   `json:"asks"`
	Timestamp time.Time `json:"timestamp"`
}

type FutureTickers struct {
	InstrumentID string    `json:"instrument_id"`
	Last         string    `json:"last"`
	BestBid      string    `json:"best_bid"`
	BestAsk      string    `json:"best_ask"`
	High24H      string    `json:"high_24h"`
	Low24H       string    `json:"low_24h"`
	Volume24H    string    `json:"volume_24h"`
	Timestamp    time.Time `json:"timestamp"`
}

type FutureTicker struct {
	InstrumentID string    `json:"instrument_id"`
	Last         string    `json:"last"`
	BestBid      string    `json:"best_bid"`
	BestAsk      string    `json:"best_ask"`
	High24H      string    `json:"high_24h"`
	Low24H       string    `json:"low_24h"`
	Volume24H    string    `json:"volume_24h"`
	Timestamp    time.Time `json:"timestamp"`
}

type FutureTrade struct {
	TradeID   string    `json:"trade_id"`
	Side      string    `json:"side"`
	Price     string    `json:"price"`
	Qty       string    `json:"qty"`
	Timestamp time.Time `json:"timestamp"`
}

type ServerTime struct {
	Iso   string `json:"iso"`
	Epoch string `json:"epoch"`
}

type ExchangeRate struct {
	InstrumentId string  `json:"instrument_id"`
	Rate         float64 `json:"rate,string"`
	Timestamp    string  `json:"timestamp"`
}

type BizWarmTips struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Result struct {
	Result bool `json:"result"`
}

type PageResult struct {
	From  int
	To    int
	Limit int
}

type FuturePosition struct {
	Result  bool `json:"result"`
	Holding [][]struct {
		LongQty              string    `json:"long_qty"`
		LongAvailQty         string    `json:"long_avail_qty"`
		LongMargin           string    `json:"long_margin"`
		LongLiquiPrice       string    `json:"long_liqui_price"`
		LongPnlRatio         string    `json:"long_pnl_ratio"`
		LongAvgCost          string    `json:"long_avg_cost"`
		LongSettlementPrice  string    `json:"long_settlement_price"`
		RealisedPnl          string    `json:"realised_pnl"`
		ShortQty             string    `json:"short_qty"`
		ShortAvailQty        string    `json:"short_avail_qty"`
		ShortMargin          string    `json:"short_margin"`
		ShortLiquiPrice      string    `json:"short_liqui_price"`
		ShortPnlRatio        string    `json:"short_pnl_ratio"`
		ShortAvgCost         string    `json:"short_avg_cost"`
		ShortSettlementPrice string    `json:"short_settlement_price"`
		InstrumentID         string    `json:"instrument_id"`
		LongLeverage         string    `json:"long_leverage"`
		ShortLeverage        string    `json:"short_leverage"`
		CreatedAt            time.Time `json:"created_at"`
		UpdatedAt            time.Time `json:"updated_at"`
		MarginMode           string    `json:"margin_mode"`
	} `json:"holding"`
}

type FutureInstrumentPosition struct {
	Result  bool `json:"result"`
	Holding []struct {
		LongQty              string    `json:"long_qty"`
		LongAvailQty         string    `json:"long_avail_qty"`
		LongAvgCost          string    `json:"long_avg_cost"`
		LongSettlementPrice  string    `json:"long_settlement_price"`
		RealisedPnl          string    `json:"realised_pnl"`
		ShortQty             string    `json:"short_qty"`
		ShortAvailQty        string    `json:"short_avail_qty"`
		ShortAvgCost         string    `json:"short_avg_cost"`
		ShortSettlementPrice string    `json:"short_settlement_price"`
		LiquidationPrice     string    `json:"liquidation_price"`
		InstrumentID         string    `json:"instrument_id"`
		Leverage             string    `json:"leverage"`
		CreatedAt            time.Time `json:"created_at"`
		UpdatedAt            time.Time `json:"updated_at"`
		MarginMode           string    `json:"margin_mode"`
		LongMargin           string    `json:"long_margin"`
		LongLiquiPrice       string    `json:"long_liqui_price"`
		LongPnlRatio         string    `json:"long_pnl_ratio"`
		ShortMargin          string    `json:"short_margin"`
		ShortLiquiPrice      string    `json:"short_liqui_price"`
		ShortPnlRatio        string    `json:"short_pnl_ratio"`
		LongLeverage         string    `json:"long_leverage"`
		ShortLeverage        string    `json:"short_leverage"`
	} `json:"holding"`
	MarginMode string `json:"margin_mode"`
}

type FutureAccount struct {
	Info struct {
		Btc struct {
			Contracts []struct {
				AvailableQty      string `json:"available_qty"`
				FixedBalance      string `json:"fixed_balance"`
				InstrumentID      string `json:"instrument_id"`
				MarginForUnfilled string `json:"margin_for_unfilled"`
				MarginFrozen      string `json:"margin_frozen"`
				RealizedPnl       string `json:"realized_pnl"`
				UnrealizedPnl     string `json:"unrealized_pnl"`
			} `json:"contracts"`
			Equity            string `json:"equity"`
			MarginMode        string `json:"margin_mode"`
			TotalAvailBalance string `json:"total_avail_balance"`
		} `json:"btc"`
		Eth struct {
			Contracts []struct {
				AvailableQty      string `json:"available_qty"`
				FixedBalance      string `json:"fixed_balance"`
				InstrumentID      string `json:"instrument_id"`
				MarginForUnfilled string `json:"margin_for_unfilled"`
				MarginFrozen      string `json:"margin_frozen"`
				RealizedPnl       string `json:"realized_pnl"`
				UnrealizedPnl     string `json:"unrealized_pnl"`
			} `json:"contracts"`
			Equity            string `json:"equity"`
			MarginMode        string `json:"margin_mode"`
			TotalAvailBalance string `json:"total_avail_balance"`
		} `json:"eth"`
		Eos struct {
			Contracts []struct {
				AvailableQty      string `json:"available_qty"`
				FixedBalance      string `json:"fixed_balance"`
				InstrumentID      string `json:"instrument_id"`
				MarginForUnfilled string `json:"margin_for_unfilled"`
				MarginFrozen      string `json:"margin_frozen"`
				RealizedPnl       string `json:"realized_pnl"`
				UnrealizedPnl     string `json:"unrealized_pnl"`
			} `json:"contracts"`
			Equity            string `json:"equity"`
			MarginMode        string `json:"margin_mode"`
			TotalAvailBalance string `json:"total_avail_balance"`
		} `json:"eos"`
		Ltc struct {
			Contracts []struct {
				AvailableQty      string `json:"available_qty"`
				FixedBalance      string `json:"fixed_balance"`
				InstrumentID      string `json:"instrument_id"`
				MarginForUnfilled string `json:"margin_for_unfilled"`
				MarginFrozen      string `json:"margin_frozen"`
				RealizedPnl       string `json:"realized_pnl"`
				UnrealizedPnl     string `json:"unrealized_pnl"`
			} `json:"contracts"`
			Equity            string `json:"equity"`
			MarginMode        string `json:"margin_mode"`
			TotalAvailBalance string `json:"total_avail_balance"`
		} `json:"ltc"`
	} `json:"info"`
}

type FutureSymbolAccount struct {
	TotalAvailBalance string `json:"total_avail_balance"`
	Contracts         []struct {
		AvailableQty      string `json:"available_qty"`
		FixedBalance      string `json:"fixed_balance"`
		InstrumentID      string `json:"instrument_id"`
		MarginForUnfilled string `json:"margin_for_unfilled"`
		MarginFrozen      string `json:"margin_frozen"`
		RealizedPnl       string `json:"realized_pnl"`
		UnrealizedPnl     string `json:"unrealized_pnl"`
	} `json:"contracts"`
	Equity     string `json:"equity"`
	MarginMode string `json:"margin_mode"`
}

type FutureTimes struct {
	BTCUSD190222 struct {
		LongLeverage  int `json:"long_leverage"`
		ShortLeverage int `json:"short_leverage"`
	} `json:"BTC-USD-190222"`
	MarginMode   string `json:"margin_mode"`
	BTCUSD190329 struct {
		LongLeverage  int `json:"long_leverage"`
		ShortLeverage int `json:"short_leverage"`
	} `json:"BTC-USD-190329"`
	BTCUSD190215 struct {
		LongLeverage  int `json:"long_leverage"`
		ShortLeverage int `json:"short_leverage"`
	} `json:"BTC-USD-190215"`
}

type FutureLedger []struct {
	LedgerID  string    `json:"ledger_id"`
	Timestamp time.Time `json:"timestamp"`
	Amount    string    `json:"amount"`
	Balance   string    `json:"balance"`
	Currency  string    `json:"currency"`
	Type      string    `json:"type"`
	Details   struct {
		OrderID      int64  `json:"order_id"`
		InstrumentID string `json:"instrument_id"`
	} `json:"details"`
}

type FutureOrder struct {
	Result       bool   `json:"result"`
	ErrorMessage string `json:"error_message"`
	ErrorCode    string `json:"error_code"`
	OrderID      string `json:"order_id"`
}

type FutureCancel struct {
	Result       bool   `json:"result"`
	OrderID      string `json:"order_id"`
	InstrumentID string `json:"instrument_id"`
}

type FutureList struct {
	Result    bool `json:"result"`
	OrderInfo []struct {
		InstrumentID string    `json:"instrument_id"`
		Size         string    `json:"size"`
		Timestamp    time.Time `json:"timestamp"`
		FilledQty    string    `json:"filled_qty"`
		Fee          string    `json:"fee"`
		OrderID      string    `json:"order_id"`
		Price        string    `json:"price"`
		PriceAvg     string    `json:"price_avg"`
		Status       string    `json:"status"`
		Type         string    `json:"type"`
		ContractVal  string    `json:"contract_val"`
		Leverage     string    `json:"leverage"`
	} `json:"order_info"`
}

type FutureOrderInfo struct {
	InstrumentID string    `json:"instrument_id"`
	Size         string    `json:"size"`
	Timestamp    time.Time `json:"timestamp"`
	FilledQty    string    `json:"filled_qty"`
	Fee          string    `json:"fee"`
	OrderID      string    `json:"order_id"`
	Price        string    `json:"price"`
	PriceAvg     string    `json:"price_avg"`
	Status       string    `json:"status"`
	Type         string    `json:"type"`
	ContractVal  string    `json:"contract_val"`
	Leverage     string    `json:"leverage"`
}

type Trades struct {
	TradeID   string    `json:"trade_id"`
	Side      string    `json:"side"`
	Price     string    `json:"price"`
	Qty       string    `json:"qty"`
	Timestamp time.Time `json:"timestamp"`
}

type FutureCandles [][]interface{}

type Account struct {
	Frozen    string `json:"frozen"`
	Hold      string `json:"hold"`
	ID        string `json:"id"`
	Currency  string `json:"currency"`
	Balance   string `json:"balance"`
	Available string `json:"available"`
	Holds     string `json:"holds"`
}

type Ledger struct {
	Timestamp time.Time `json:"timestamp"`
	LedgerID  string    `json:"ledger_id"`
	CreatedAt time.Time `json:"created_at"`
	Currency  string    `json:"currency"`
	Amount    string    `json:"amount"`
	Balance   string    `json:"balance"`
	Type      string    `json:"type"`
	Details   []struct {
		InstrumentID string `json:"instrument_id"`
		OrderID      string `json:"order_id"`
		ProductID    string `json:"product_id"`
	} `json:"details,omitempty"`
}

type Order struct {
	ClientOid    string `json:"client_oid"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	OrderID      string `json:"order_id"`
	Result       bool   `json:"result"`
}

type OrderItem struct {
	ClientOid      string    `json:"client_oid"`
	CreatedAt      time.Time `json:"created_at"`
	FilledNotional string    `json:"filled_notional"`
	FilledSize     string    `json:"filled_size"`
	Funds          string    `json:"funds"`
	InstrumentID   string    `json:"instrument_id"`
	Notional       string    `json:"notional"`
	OrderID        string    `json:"order_id"`
	OrderType      string    `json:"order_type"`
	Price          string    `json:"price"`
	ProductID      string    `json:"product_id"`
	Side           string    `json:"side"`
	Size           string    `json:"size"`
	State          string    `json:"state"`
	Status         string    `json:"status"`
	Timestamp      time.Time `json:"timestamp"`
	Type           string    `json:"type"`
}

type CancelOrder struct {
	ClientOid    string `json:"client_oid"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	OrderID      string `json:"order_id"`
	Result       bool   `json:"result"`
}

type OrderPending struct {
	ClientOid      string    `json:"client_oid"`
	CreatedAt      time.Time `json:"created_at"`
	FilledNotional string    `json:"filled_notional"`
	FilledSize     string    `json:"filled_size"`
	Funds          string    `json:"funds"`
	InstrumentID   string    `json:"instrument_id"`
	Notional       string    `json:"notional"`
	OrderID        string    `json:"order_id"`
	OrderType      string    `json:"order_type"`
	Price          string    `json:"price"`
	ProductID      string    `json:"product_id"`
	Side           string    `json:"side"`
	Size           string    `json:"size"`
	State          string    `json:"state"`
	Status         string    `json:"status"`
	Timestamp      time.Time `json:"timestamp"`
	Type           string    `json:"type"`
}

type OrderInfo struct {
	ClientOid      string    `json:"client_oid"`
	CreatedAt      time.Time `json:"created_at"`
	FilledNotional string    `json:"filled_notional"`
	FilledSize     string    `json:"filled_size"`
	Funds          string    `json:"funds"`
	InstrumentID   string    `json:"instrument_id"`
	Notional       string    `json:"notional"`
	OrderID        string    `json:"order_id"`
	OrderType      string    `json:"order_type"`
	Price          string    `json:"price"`
	ProductID      string    `json:"product_id"`
	Side           string    `json:"side"`
	Size           string    `json:"size"`
	State          string    `json:"state"`
	Status         string    `json:"status"`
	Timestamp      time.Time `json:"timestamp"`
	Type           string    `json:"type"`
}

type Deal []interface{}

type SwapTicker struct {
	InstrumentID string    `json:"instrument_id"`
	Last         string    `json:"last"`
	BestBid      string    `json:"best_bid"`
	BestAsk      string    `json:"best_ask"`
	High24H      string    `json:"high_24h"`
	Low24H       string    `json:"low_24h"`
	Volume24H    string    `json:"volume_24h"`
	Timestamp    time.Time `json:"timestamp"`
}

type SwapOrder struct {
	ClientOid    string `json:"client_oid"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	OrderID      string `json:"order_id"`
	Result       string `json:"result"`
}
