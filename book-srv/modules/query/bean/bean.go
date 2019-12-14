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
	Swtdz     string //32
	Ydz       string //31
	Edz       string //30
	Gjrw      string //21
	Rw        string //23
	Dw        string //33
	Yw        string //28
	Rz        string //24
	Yz        string //29
	Wz        string //26
	Qt        string //22
	Bz        string //1

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
