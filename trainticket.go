package main

type ticketsInfo struct {
	Httpstatus int `json:"httpstatus"`
	Data       struct {
		Result []string `json:"result"`
		Flag   string   `json:"flag"`
		Map    struct {
			QSW string `json:"QSW"`
			XIW string `json:"XIW"`
		} `json:"map"`
	} `json:"data"`
	Messages string `json:"messages"`
	Status   bool   `json:"status"`
}
