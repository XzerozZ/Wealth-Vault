package utils

func MaskBankAccount(accNum string) string {
	if len(accNum) <= 4 {
		return accNum
	}

	return "xxx-x-" + accNum[len(accNum)-4:]
}
