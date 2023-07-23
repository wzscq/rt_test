package robot

import (
	"log"
	"github.com/gin-gonic/gin"
	"rt_test_service/common"
	"rt_test_service/crv"
	"net/http"
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
			"response":oauthRsp,
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
			"response":getRobotListRsp,
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultGetRobotPlatformAPIError,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return	
	}

	//机器人信息写入数据库
	UpdateRobotList(rtc.CRVClient,getRobotListRsp.Result,header.Token)

	rsp:=common.CreateResponse(common.CreateError(common.ResultSuccess,nil),nil)
	c.IndentedJSON(http.StatusOK, rsp)
	log.Println("RobotController getRobotList success")
}

//Bind bind the controller function to url
func (rtc *RobotController) Bind(router *gin.Engine) {
	log.Println("Bind RobotController")
	router.POST("/getRobotList", rtc.getRobotList)
}