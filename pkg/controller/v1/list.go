package v1

//QueryRequestMeta query request meta
type QueryRequestMeta struct {
	PageNo   int `json:"pageNo"`
	PageSize int `json:"pageSize"`
}

//ListResponseMeta list response meta
type ListResponseMeta struct {
	PageNo     int   `json:"pageNo"`
	PageSize   int   `json:"pageSize"`
	TotalCount int64 `json:"totalCount"`
}
