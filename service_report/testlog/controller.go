package testlog

import (
	"rt_report/common"
	"rt_report/crv"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type TestLogController struct {
	TestLogConf *common.TestLogConf
}

//Bind bind the controller function to url
func (cl *TestLogController) Bind(router *gin.Engine) {
	log.Println("Bind SmartLockController")
	router.POST("testlog/download", cl.download)
}

func (cl *TestLogController)download(c *gin.Context){
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

	


}