package Structs

type GroupData struct {
	Id              int
	Name            string
	BotsOnline      []DashboardBotData
	BotsOffline     []DashboardBotData
	BotsCount       int
	CompletedAction int
	BannedAction    int
	LockedAction    int
	AutoReplenish   int
}
