package Structs

type GroupData struct {
	Id          int
	Name        string
	BotsOnline  []DashboardBotData
	BotsOffline []DashboardBotData
}
