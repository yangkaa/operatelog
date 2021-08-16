package controller

import (
	"fmt"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"goodrain.com/operatelog/cmd/operatelog/config"
	v1 "goodrain.com/operatelog/pkg/controller/v1"
	"goodrain.com/operatelog/pkg/usecase"
	"goodrain.com/operatelog/pkg/utils/pageutil"
	restfulutil "goodrain.com/operatelog/pkg/utils/restful"
	"net/http"
)

// MkyAuditLogController -
type MkyAuditLogController struct {
	MkyAuditLogUcase usecase.MkyAuditLogUsecase `inject:""`
}

// NewMkyAuditLogController creates a new MkyAuditLogController.
func NewMkyAuditLogController() *MkyAuditLogController {
	return &MkyAuditLogController{}
}

func (m *MkyAuditLogController) DisablePrefix() bool {
	return false
}

// WebService -
func (m *MkyAuditLogController) WebService(ws *restful.WebService) {
	tags := []string{"MkyAuditLog"}

	ws.Route(ws.POST("/mkyAuditLogs").To(m.createMkyAuditLog).
		Doc("创建日志").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(v1.MkyAuditLog{}).
		Returns(200, "ok", nil))
	ws.Route(ws.GET("/mkyAuditLogs").To(m.listMkyAuditLogs).
		Doc("日志列表").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("page", "页码").DataType("integer")).
		Param(ws.QueryParameter("pageSize", "每页的大小. -1 表示不限制").DataType("integer")).
		Param(ws.QueryParameter("startTime", "起始时间").DataType("string")).
		Param(ws.QueryParameter("endTime", "截止时间").DataType("string")).
		Param(ws.QueryParameter("query", "模糊搜索").DataType("string")).
		Returns(200, "ok", v1.ListMkyAuditLogResp{}))
	ws.Route(ws.POST("/batchLogs").To(m.batchLogs).
		Doc("批量创建日志").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(v1.BatchLogsReq{}).
		Returns(200, "ok", nil))
}

func (m *MkyAuditLogController) createMkyAuditLog(r *restful.Request, w *restful.Response) {
	token := r.HeaderParameter("Authorization")
	if token == "" || token != config.C.AccessToken {
		restfulutil.Error(w, http.StatusForbidden, fmt.Errorf("无权访问本api，请授权后在进行访问"))
		return
	}
	var req v1.MkyAuditLog
	if err := r.ReadEntity(&req); err != nil {
		restfulutil.Error(w, http.StatusBadRequest, fmt.Errorf("请求错误，请检查请求参数"))
		return
	}
	err := m.MkyAuditLogUcase.Create(&req)
	if err != nil {
		restfulutil.Error(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeaderAndEntity(http.StatusOK, &v1.Response{
		Message: "请求成功",
		Code:    http.StatusOK,
	})
}

func (m *MkyAuditLogController) listMkyAuditLogs(r *restful.Request, w *restful.Response) {
	query := r.QueryParameter("query")
	startTime := r.QueryParameter("startTime")
	endTime := r.QueryParameter("endTime")
	page, pageSeize := pageutil.ExtractPageParameter(r)
	resp, err := m.MkyAuditLogUcase.List(startTime, endTime, query, page, pageSeize)
	if err != nil {
		restfulutil.Error(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteAsJson(resp)
}

func (m *MkyAuditLogController) batchLogs(r *restful.Request, w *restful.Response) {
	token := r.HeaderParameter("Authorization")
	if token == "" || token != config.C.AccessToken {
		restfulutil.Error(w, http.StatusForbidden, fmt.Errorf("无权访问本api，请授权后在进行访问"))
		return
	}
	var req v1.BatchLogsReq
	if err := r.ReadEntity(&req); err != nil {
		restfulutil.Error(w, http.StatusBadRequest, fmt.Errorf("请求错误，请检查请求参数"))
		return
	}
	err := m.MkyAuditLogUcase.BatchCreateLogs(&req)
	if err != nil {
		restfulutil.Error(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeaderAndEntity(http.StatusOK, &v1.Response{
		Message: "请求成功",
		Code:    http.StatusOK,
	})
}
