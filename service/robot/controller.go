package robot

import (
	"log"
	"github.com/gin-gonic/gin"
	"rt_test_service/common"
	"rt_test_service/crv"
	"net/http"
	"encoding/base64"
  "io/ioutil"
)

type RobotController struct {
	RobotClient *RobotClient
	CRVClient *crv.CRVClient
}

func (rtc *RobotController)getRobotList(c *gin.Context){
	var header crv.CommonHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("RobotController getRobotList wrong request")
		return
	}	

	//获取token
	oauthRsp, err:=rtc.RobotClient.Oauth()
	if err!=nil {
		log.Println("RobotController getRobotList Oauth error",err)
		params:=map[string]interface{}{
			"error":err,
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultGetRobotPlatformTokenError,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return	
	}

	if oauthRsp.Result==nil {
		log.Println("RobotController getRobotList Oauth error",oauthRsp)
		params:=map[string]interface{}{
			"message":oauthRsp.Message,
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultGetRobotPlatformTokenError,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return	
	}

	//获取机器人列表
	getRobotListRsp, err:=rtc.RobotClient.GetRobotList(oauthRsp.Result.Token)
	if err!=nil {
		log.Println("RobotController getRobotList GetRobotList error",err)
		params:=map[string]interface{}{
			"error":err,
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultGetRobotPlatformAPIError,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return	
	}

	if getRobotListRsp.Result==nil {
		log.Println("RobotController getRobotList GetRobotList error",getRobotListRsp)
		params:=map[string]interface{}{
			"message":getRobotListRsp.Message,
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultGetRobotPlatformAPIError,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return	
	}

	//机器人信息写入数据库
	rsp:=UpdateRobotList(rtc.CRVClient,getRobotListRsp.Result,header.Token)
	if rsp!=nil {
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	rsp=common.CreateResponse(common.CreateError(common.ResultSuccess,nil),nil)
	c.IndentedJSON(http.StatusOK, rsp)
	log.Println("RobotController getRobotList success")
}

func (rtc *RobotController)getCurrentRobotStatus(c *gin.Context){
	var header crv.CommonHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("RobotController getCurrentRobotStatus wrong request")
		return
	}	

	var rep crv.CommonReq
	if err := c.BindJSON(&rep); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("RobotController getCurrentRobotStatus with error")
		return
  }	

	if rep.SelectedRowKeys ==nil || len(*rep.SelectedRowKeys)==0 {
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("RobotController getCurrentRobotStatus：request SelectedRowKeys is empty")
		return
	}

	//获取token
	oauthRsp, err:=rtc.RobotClient.Oauth()
	if err!=nil {
		log.Println("RobotController getCurrentRobotStatus：request Oauth error",err)
		params:=map[string]interface{}{
			"error":err,
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultGetRobotPlatformTokenError,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return	
	}

	if oauthRsp.Result==nil {
		log.Println("RobotController getCurrentRobotStatus：request Oauth error",oauthRsp)
		params:=map[string]interface{}{
			"message":oauthRsp.Message,
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultGetRobotPlatformTokenError,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return	
	}

	//循环获取机器人状态
	for _,robotID:=range *rep.SelectedRowKeys {
		getCurrentRobotStatusRsp, err:=rtc.RobotClient.GetCurrentRobotStatus(oauthRsp.Result.Token,robotID)
		if err!=nil {
			log.Println("RobotController getCurrentRobotStatus：request GetCurrentRobotStatus error",err)
			params:=map[string]interface{}{
				"error":err,
			}
			rsp:=common.CreateResponse(common.CreateError(common.ResultGetRobotPlatformAPIError,params),nil)
			c.IndentedJSON(http.StatusOK, rsp)
			return	
		}
	
		if getCurrentRobotStatusRsp.Result==nil {
			log.Println("RobotController getCurrentRobotStatus：request GetCurrentRobotStatus error",getCurrentRobotStatusRsp)
			params:=map[string]interface{}{
				"message":getCurrentRobotStatusRsp.Message,
			}
			rsp:=common.CreateResponse(common.CreateError(common.ResultGetRobotPlatformAPIError,params),nil)
			c.IndentedJSON(http.StatusOK, rsp)
			return	
		}

		rsp:=UpdateRobotStatus(rtc.CRVClient,getCurrentRobotStatusRsp.Result,header.Token)
		if rsp!=nil {
			c.IndentedJSON(http.StatusOK, rsp)
			return
		}
	}

	rsp:=common.CreateResponse(common.CreateError(common.ResultSuccess,nil),nil)
	c.IndentedJSON(http.StatusOK, rsp)
	log.Println("RobotController getRobotList success")
}

func (rtc *RobotController)mapUpload(c *gin.Context){
	var header ServerHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("RobotController mapUpload wrong request")
		return
	}	

	var rep UploadMapReq
	if err := c.ShouldBind(&rep); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("RobotController mapUpload with error")
		return
  }	

	file, err := c.FormFile("file")
	if err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("RobotController mapUpload with error")
		return
	}

	f, _ := file.Open()

	// 读取文件的内容到一个字节切片中
	content, _ := ioutil.ReadAll(f)
	// 将字节切片转换为Base64编码
	encoded := base64.StdEncoding.EncodeToString(content)

	//上传服务器
	rsp:=SaveRobotMap(rtc.CRVClient,&rep,encoded,header.Token)
	if rsp != nil {
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	rsp=common.CreateResponse(common.CreateError(common.ResultSuccess,nil),nil)
	c.IndentedJSON(http.StatusOK, rsp)
	log.Println("RobotController getRobotList success")
}

//Bind bind the controller function to url
func (rtc *RobotController) Bind(router *gin.Engine) {
	log.Println("Bind RobotController")
	router.POST("/getRobotList", rtc.getRobotList)
	router.POST("/getCurrentRobotStatus", rtc.getCurrentRobotStatus)
	router.POST("/mapUpload", rtc.mapUpload)
}