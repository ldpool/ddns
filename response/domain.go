package response

type DomainAll struct {
	Response DomainResponse `json:"Response"`
}

type DomainResponse struct {
	DomainCountInfo DomainCount  `json:"DomainCountInfo"`
	DomainList      []DomainInfo `json:"DomainList"`
	RequestId       string       `json:"RequestId"`
}

type DomainCount struct {
	DomainTotal int `json:"DomainTotal"`
}

type DomainInfo struct {
	DomainId    int    `json:"DomainId"`
	Name        string `json:"Name"`
	Status      string `json:"Status"`
	GradeTitle  string `json:"GradeTitle"`
	Grade       string `json:"Grade"`
	RecordCount int    `json:"RecordCount"`
	CreatedOn   string `json:"CreatedOn"`
	UpdatedOn   string `json:"UpdatedOn"`
}

func (d *DomainInfo) MapStatus() *DomainInfo {
	statusMap := map[string]string{
		"ENABLE": "启用",
		"PAUSE":  "暂停",
		"SPAM":   "封禁",
		"LOCK":   "锁定",
	}
	cnStatus, found := statusMap[d.Status]
	if found {
		d.Status = cnStatus
	}

	return d
}

// 记录
type RecordAll struct {
	Response RecordResponse `json:"Response"`
}

type RecordResponse struct {
	RecordCountInfo RecordCount  `json:"RecordCountInfo"`
	RecordList      []RecordInfo `json:"RecordList"`
	RequestId       string       `json:"RequestId"`
}

type RecordCount struct {
	TotalCount int `json:"TotalCount"`
}

type RecordInfo struct {
	RecordId  int    `json:"RecordId"`
	Name      string `json:"Name"`
	Status    string `json:"Status"`
	Value     string `json:"Value"`
	TTL       int    `json:"TTL"`
	MX        int    `json:"MX"`
	Type      string `json:"Type"`
	Line      string `json:"Line"`
	UpdatedOn string `json:"UpdatedOn"`
}

func (d *RecordInfo) MapStatus() *RecordInfo {
	statusMap := map[string]string{
		"ENABLE":  "启用",
		"DISABLE": "暂停",
		"SPAM":    "封禁",
		"LOCK":    "锁定",
	}
	cnStatus, found := statusMap[d.Status]
	if found {
		d.Status = cnStatus
	}

	return d
}

type RecordEdit struct {
	Response RecordInfoE `json:"Response"`
}

type RecordInfoE struct {
	RecordEditInfo RecordEditInfo `json:"RecordInfo"`
}

type RecordEditInfo struct {
	Id         int    `json:"Id"`
	Domain     string `json:"domain"`
	DomainId   int    `json:"DomainId"`
	SubDomain  string `json:"SubDomain"`
	RecordType string `json:"RecordType"`
	RecordLine string `json:"RecordLine"`
	TTL        int    `json:"TTL"`
	MX         int    `json:"MX"`
	Type       string `json:"Type"`
	Value      string `json:"Value"`
	Enabled    int    `json:"Enabled"`
	UpdatedOn  string `json:"UpdatedOn"`
}

type GradeResponse struct {
	Response struct {
		TypeList  []string `json:"TypeList"`
		RequestId string   `json:"RequestId"`
	} `json:"Response"`
}

type Line struct {
	Name   string `json:"Name"`
	LineId string `json:"LineId"`
}

type LineListResponse struct {
	Response struct {
		LineList []Line `json:"LineList"`
	} `json:"Response"`
}
