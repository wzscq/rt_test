package testlog

import (
	"rt_report/common"
	"rt_report/crv"
	"rt_report/logfile"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

type TestLogController struct {
	TestLogConf *common.TestLogConf
	CRVClient *crv.CRVClient
}

//Bind bind the controller function to url
func (cl *TestLogController) Bind(router *gin.Engine) {
	log.Println("Bind SmartLockController")
	router.POST("testlog/getDownloadOperation", cl.getDownloadOperation)
	router.GET("testlog/download/:fileName", cl.download)
}

func (cl *TestLogController)download(c *gin.Context){
	fileName:=c.Param("fileName")
	if len(fileName)==0 {
		log.Println("TestLogController download error: fileName is empty")
		c.IndentedJSON(http.StatusBadRequest, nil)
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+fileName)
	fullFileName:=cl.TestLogConf.Root+"/"+fileName
	//判单文件是否存在
	if _,err:=os.Stat(fullFileName);err!=nil {
		log.Println("TestLogController download error: file not exist,fileName=",fileName)
		c.IndentedJSON(http.StatusNotFound, nil)
		return
	}

	c.File(fullFileName)
}

func (cl *TestLogController)getDownloadOperation(c *gin.Context){
	var header crv.CommonHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("TestLogController download with wrong request")
		return
	}	

	var rep crv.CommonReq
	if err := c.BindJSON(&rep); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("end TestLogController download with error")
		return
  	}	

	if rep.SelectedRowKeys==nil || len(*rep.SelectedRowKeys)==0 {
		log.Println("end TestLogController download with error:SelectedRowKeys is empty")
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	cachedTestFileID:=(*rep.SelectedRowKeys)[0]
	cachedTestFile,errorCode:=logfile.GetCachedTestFile(cl.CRVClient,cachedTestFileID,header.Token)
	if errorCode!=common.ResultSuccess {
		log.Println("end TestLogController download with error:errorCode=",errorCode)
		rsp:=common.CreateResponse(common.CreateError(errorCode,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	log.Println(cachedTestFile)

	logFileName,errorCode:=logfile.GetLogFileNameByCachedTestFile(cl.TestLogConf.Root,cachedTestFile)
	if errorCode!=common.ResultSuccess {
		log.Println("end TestLogController download with error:errorCode=",errorCode)
		rsp:=common.CreateResponse(common.CreateError(errorCode,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	log.Println(logFileName)
	
	//获取测试文件下载操作返回结果
	res:=crv.GetDownloadFileResult(cl.TestLogConf.DownloadUrl,logFileName)	
	rsp:=common.CreateResponse(nil,res)
	c.IndentedJSON(http.StatusOK, rsp)
	log.Println("end TestLogController download with success")
}