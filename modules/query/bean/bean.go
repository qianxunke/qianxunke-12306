package bean

type Train struct {
	SecretStr string
	TrainCode string
	Num       string
	From      string
	To        string
	FindFrom  string
	FindTo    string
	StartTime string
	EndTime   string
	CostTime  string
	CanBuy    string
	StartDate string
}

type QueryItem struct {
	Data       Data
	Httpstatus int
	Messages   string
	status     bool
}

type Data struct {
	Flag   string
	Map    map[string]string
	Result []string
}

type QueryTimeResult struct {
	OK        string
	WaitTime  int64
	WaitCount int64
	OrderId   string
}
