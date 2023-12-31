package report

import (
	"rt_report/common"
	"rt_report/crv"
	"rt_report/logfile"
	"rt_report/dingli"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

type generateReportRequest struct {
	ReportType string `json:"reportType"`
	FileName []string `json:"fileName"`
	Template string `json:"template`
}

type ReportController struct {
	DingliClient *dingli.DingliClient
	CRVClient *crv.CRVClient
	TestLogConf *common.TestLogConf
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
	if rep.ReportType == dingli.CMD_KPIReport {
		_,err=cl.DingliClient.GetKPIReport(rep.FileName,rep.Template)
	} else if rep.ReportType == dingli.CMD_CustomReport {
		_,err=cl.DingliClient.GetCustomerReport(rep.FileName,rep.Template)
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

func (cl *ReportController)getDownloadOperation(c *gin.Context){
	var header crv.CommonHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("TestLogController getDownloadOperation with wrong request")
		return
	}	

	var rep crv.CommonReq
	if err := c.BindJSON(&rep); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("end TestLogController getDownloadOperation with error")
		return
  	}	

	if rep.SelectedRowKeys==nil || len(*rep.SelectedRowKeys)==0 {
		log.Println("end TestLogController getDownloadOperation with error:SelectedRowKeys is empty")
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	if rep.List==nil || len(*rep.List)<=0 {
		log.Println("end TestLogController getDownloadOperation with error:list is empty")
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	//获取测试报表类型
	reportType,ok:=(*rep.List)[0]["report_type"].(string)
	if !ok {
		log.Println("end TestLogController getDownloadOperation with error:report_type is empty")
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	//获取测试报表模板
	template,ok:=(*rep.List)[0]["report_template"].(string)
	if !ok {
		log.Println("end TestLogController getDownloadOperation with error:report_template is empty")
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	cachedTestFileID:=(*rep.SelectedRowKeys)[0]

	log.Println("ReportController getDownloadOperation reportType=",reportType," template=",template," cachedTestFileID=",cachedTestFileID)

	cachedTestFile,errorCode:=logfile.GetCachedTestFile(cl.CRVClient,cachedTestFileID,header.Token)
	if errorCode!=common.ResultSuccess {
		log.Println("end TestLogController getDownloadOperation GetCachedTestFile with error:errorCode=",errorCode)
		rsp:=common.CreateResponse(common.CreateError(errorCode,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	log.Println(cachedTestFile)

	logFileName,errorCode:=logfile.GetLogFileNameByCachedTestFile(cl.TestLogConf.Root,cachedTestFile)
	if errorCode!=common.ResultSuccess {
		log.Println("end TestLogController getDownloadOperation GetLogFileNameByCachedTestFile with error:errorCode=",errorCode)
		rsp:=common.CreateResponse(common.CreateError(errorCode,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	log.Println("logFileName=",logFileName)

	fileName,errorCode:=logfile.ExtractDDIBFile(cl.TestLogConf.Root,logFileName,cl.TestLogConf.ExtractRoot)
	if errorCode!=common.ResultSuccess {
		log.Println("end TestLogController getDownloadOperation extractDDIBFile with error:errorCode=",errorCode)
		rsp:=common.CreateResponse(common.CreateError(errorCode,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}
	log.Println("ddib fileName=",fileName)

	reportFileName,err:=cl.DingliClient.CreateReport(reportType,template,[]string{cl.TestLogConf.ExtractRoot+"/"+fileName})	
	if err!=nil{
		params:=map[string]interface{}{
			"error":err.Error(),
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultGenerateReportError,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("ReportController getDownloadOperation CreateReport error:",err.Error())
		return
	}

	//获取报告文件名，将报告迁移到extractRoot目录下
	log.Println(reportFileName)
	destFileName:=cachedTestFile.DeviceID+"_"+cachedTestFile.Timestamp+"_"+reportType+"_"+template
	reportFileName,err=logfile.MoveReportToExtractRootFolder(reportFileName,cl.TestLogConf.ExtractRoot,destFileName)
	if err!=nil{
		params:=map[string]interface{}{
			"error":err.Error(),
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultGenerateReportError,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("ReportController getDownloadOperation moveReportToExtractRootFolder error:",err.Error())
		return
	}

	//构造下载操作对象
	res:=crv.GetDownloadFileResult(cl.TestLogConf.ReportDownloadUrl,reportFileName)	
	rsp:=common.CreateResponse(nil,res)
	c.IndentedJSON(http.StatusOK, rsp)
	log.Println("ReportController end getDownloadOperation success")
}

func (cl *ReportController)download(c *gin.Context){
	fileName:=c.Param("fileName")
	if len(fileName)==0 {
		log.Println("ReportController download error: fileName is empty")
		c.IndentedJSON(http.StatusBadRequest, nil)
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+fileName)
	fullFileName:=cl.TestLogConf.ExtractRoot+"/"+fileName
	//判单文件是否存在
	if _,err:=os.Stat(fullFileName);err!=nil {
		log.Println("ReportController download error: file not exist,fileName=",fileName)
		c.IndentedJSON(http.StatusNotFound, nil)
		return
	}

	c.File(fullFileName)
}

//Bind bind the controller function to url
func (cl *ReportController) Bind(router *gin.Engine) {
	log.Println("Bind SmartLockController")
	router.POST("report/generate", cl.generateReport)
	router.POST("report/getDownloadOperation", cl.getDownloadOperation)
	router.GET("report/download/:fileName", cl.download)
}