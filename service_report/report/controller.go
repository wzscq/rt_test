package report

import (
	"rt_report/common"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type generateReportRequest struct {
	ReportType string `json:"reportType"`
	FileName []string `json:"fileName"`
	Template string `json:"template`
}

type ReportController struct {
	DingliClient *DingliClient
}

func (cl *ReportController)generateReport(c *gin.Context){
	var rep generateReportRequest
	if err := c.BindJSON(&rep); err != nil {
		log.Println(err)
		params:=map[string]interface{}{
			"error":err,
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("ReportController end generateReport with error")
		return
  }
	
	var err error
	if rep.ReportType == CMD_KPIReport {
		err=cl.DingliClient.GetKPIReport(rep.FileName,rep.Template)
	} else if rep.ReportType == CMD_CustomReport {
		err=cl.DingliClient.GetCustomerReport(rep.FileName,rep.Template)
	} else {
		params:=map[string]interface{}{
			"report type":rep.ReportType,
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultNotSupportedReportType,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("ReportController end generateReport with error: not supported report type ",rep.ReportType)
		return
	}

	if err!=nil{
		params:=map[string]interface{}{
			"error":err,
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultGenerateReportError,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("ReportController end generateReport with error")
		return
	}

	rsp:=common.CreateResponse(nil,nil)
	c.IndentedJSON(http.StatusOK, rsp)
	log.Println("ReportController end generateReport success")
}

//Bind bind the controller function to url
func (cl *ReportController) Bind(router *gin.Engine) {
	log.Println("Bind SmartLockController")
	router.POST("report/generate", cl.generateReport)
}