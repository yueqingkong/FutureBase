package trade

import (
	"github.com/yueqingkong/FutureBase/base"
	"github.com/yueqingkong/FutureBase/okex"
	"testing"
)

func TestBaseTrade_Start(t *testing.T) {
	baseTrade := BaseTrade{Mode: func() string {
		return "test"
	}, DB: func() (s string, s2 string, s3 string) {
		return "margin_test", "root", "qwer1234@QW"
	}, InitKeys: func() []string {
		return []string{"", "", ""}
	}, Account: func() map[string]string {
		account := make(map[string]string, 0)
		account["btc"] = "0.01"
		return account
	}, Port: func() string {
		return "9003"
	}, SymbolPair: func() base.SYMBOL {
		return base.ETHUSD
	}, ContractPeriod: func() base.CONTRACT_PERIOD {
		return base.QUARTER
	}, Plat: func() base.PlatBase {
		return okex.NewOkexFuture()
	}}

	baseTrade.Start(NewSimple())
}
