package usecase

// ValidOrderNumber - check number is valid or not based on Luhn algorithm
func ValidOrderNumber(number uint64) bool {
	return (number%10+checksum(number/10))%10 == 0
}

// checksum -.
func checksum(number uint64) uint64 {
	var luhn uint64

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 { // even
			cur *= 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number /= 10
	}

	return luhn % 10
}
