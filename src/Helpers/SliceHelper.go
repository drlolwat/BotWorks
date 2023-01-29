package Helpers

import "BotWorks/src/Structs"

func ContainsInt(check int, numbers []int) bool {
	for _, item := range numbers {
		if item == check {
			return true
		}
	}
	return false
}

func ContainsBotData(check Structs.BotData, haystack []Structs.BotData) bool {
	for _, item := range haystack {
		if item == check {
			return true
		}
	}
	return false
}

func ContainsGroupData(check Structs.GroupData, haystack []Structs.GroupData) bool {
	for _, item := range haystack {
		if item.Id == check.Id {
			return true
		}
	}
	return false
}
