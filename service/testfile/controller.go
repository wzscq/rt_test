package testfile

import (
	"log"
	"github.com/gin-gonic/gin"
	"rt_test_service/common"
	"rt_test_service/crv"
	"net/http"
	"strconv"
)

type GetContentReq struct {
	DeviceID string `json:"deviceID"`
	TimeStamp string `json:"timestamp"`
	From int64  		`json:"from"`
	To int64	`json:"to1"`
}

type TestFileController struct {
	OutPath string
}

func (tfc *TestFileController)GetContent(c *gin.Context){
	var header crv.CommonHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("TestFileController GetContent wrong request")
		return
	}	

	var rep GetContentReq
	if err := c.BindJSON(&rep); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("TestFileController GetContent with error")
		return
  	}	

	//获取文件内容
	timestamp,_:=strconv.ParseInt(rep.TimeStamp,10,64)
	tf:=GetTestFile(tfc.OutPath,rep.DeviceID,timestamp)
	if tf==nil {
		rsp:=common.CreateResponse(common.CreateError(common.ResultFileNotExist,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("TestFileController GetContent file not exist.")
		return
	}
	defer tf.CloseReadOnly()

	//获取文件内容
	//from,_:=strconv.ParseInt(rep.From,10,64)
	//to,_:=strconv.ParseInt(rep.To,10,64)
	content:=tf.GetContent(rep.From,rep.To)
	rsp:=common.CreateResponse(common.CreateError(common.ResultSuccess,nil),content)
	c.IndentedJSON(http.StatusOK, rsp)
}

func (tfc *TestFileController)GetPoints(c *gin.Context){
	var header crv.CommonHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("TestFileController GetPoints wrong request")
		return
	}	

	var rep GetContentReq
	if err := c.BindJSON(&rep); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("TestFileController GetPoints with error")
		return
  }	

	//获取文件内容
	timestamp,_:=strconv.ParseInt(rep.TimeStamp,10,64)
	tf:=GetTestFile(tfc.OutPath,rep.DeviceID,timestamp)
	if tf==nil {
		rsp:=common.CreateResponse(common.CreateError(common.ResultFileNotExist,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("TestFileController GetPoints file not exist.")
		return
	}
	defer tf.CloseReadOnly()

	//获取文件内容
	points:=tf.GetPoints()
	rsp:=common.CreateResponse(common.CreateError(common.ResultSuccess,nil),points)
	c.IndentedJSON(http.StatusOK, rsp)
}

//Bind bind the controller function to url
func (tfc *TestFileController) Bind(router *gin.Engine) {
	log.Println("Bind CacheFileController")
	router.POST("/testfile/GetContent", tfc.GetContent)
	router.POST("/testfile/GetPoints", tfc.GetPoints)
}