package util

/**
 * 价格相等 合理偏差 0.001
 * dest 目标价格
 */
func PriceEqual(dest float32, price float32) bool {
	var b = false
	if Abs(dest-price)/dest <= 0.002 {
		b = true
	}
	return b
}
