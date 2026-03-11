package collection

// iif is a simple inline if statement that simplifies some of the code in this package.
func iif[T any](condition bool, trueValue T, falseValue T) T {

	if condition {
		return trueValue
	}

	return falseValue
}
