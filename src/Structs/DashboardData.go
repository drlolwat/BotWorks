package Structs

type DashboardData struct {
	RunningBots    int
	StoppedBots    int
	UniqueScripts  int
	CurrentTime    string
	OnlineBotData  []DashboardBotData
	OfflineBotData []DashboardBotData
	GroupData      []GroupData
	LeftMessage    string
	RightMessage   string
}
