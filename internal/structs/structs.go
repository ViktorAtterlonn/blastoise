package structs

type RequestResult struct {
	StatusCode int
	Duration   int
}

type Ctx struct {
	Url        string
	Rps        int
	Duration   int
	Method     string
	Body       string
	ResultChan chan []*RequestResult
	AbortChan  chan bool
}
