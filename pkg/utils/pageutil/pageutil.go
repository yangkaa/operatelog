package pageutil

import (
	"strconv"

	"github.com/emicklei/go-restful"
)

// ExtractPageParameter extracts page parameters from request.
func ExtractPageParameter(r *restful.Request) (page, pageSize int) {
	page, _ = strconv.Atoi(r.QueryParameter("page"))
	if page <= 0 {
		page = 1
	}
	pageNo, _ := strconv.Atoi(r.QueryParameter("pageNo"))
	if pageNo <= 0 {
		pageNo = 1
	}
	if pageNo > page {
		page = pageNo
	}
	pageSize, _ = strconv.Atoi(r.QueryParameter("pageSize"))
	if pageSize == 0 || pageSize < -1 {
		pageSize = 10
	}
	return page, pageSize
}

// PageParameter -
func PageParameter(ws *restful.WebService) *restful.Parameter {
	return ws.QueryParameter("page", "页码").DataType("integer")
}

// PageSizeParameter -
func PageSizeParameter(ws *restful.WebService) *restful.Parameter {
	return ws.QueryParameter("pageSize", "每页的大小. -1 表示不限制").DataType("integer")
}
