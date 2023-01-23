package Helpers

func ContainsInt(check int, numbers []int) bool {
	for _, item := range numbers {
		if item == check {
			return true
		}
	}
	return false
}
