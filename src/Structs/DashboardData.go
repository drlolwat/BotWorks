package Structs

type DashboardData struct {
	RunningBots   int
	StoppedBots   int
	UniqueScripts int
	CurrentTime   string
	OnlineBotData []OnlineBotData
	GroupData     []GroupData
}
