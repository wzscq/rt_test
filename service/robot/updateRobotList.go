package robot

import (
	"log"
	"rt_test_service/common"
	"rt_test_service/crv"
)

type ProjectRobot struct {
	ID string `json:"id"`
	Version string `json:"version"`
	Name string `json:"name"`
	ProjectID string `json:"rt_project_id"`
}	

var MOdElID_ROBOT="rt_project_robot"

var	queryRobotFields=[]map[string]interface{}{
	{"field": "id"},
	{"field": "version"},
}

func GetRobotList(crvClient *crv.CRVClient,token string)([]ProjectRobot){
	commonRep:=crv.CommonReq{
		ModelID:MOdElID_ROBOT,
		Fields:&queryRobotFields,
	}

	rsp,commonErr:=crvClient.Query(&commonRep,token)
	if commonErr!=common.ResultSuccess {
		return nil
	}

	if rsp.Error == true {
		log.Println("GetRobotList error:",rsp.ErrorCode,rsp.Message)
		return nil
	}

	resLst,ok:=rsp.Result["list"].([]interface{})
	if !ok {
		log.Println("GetRobotList error: no list in rsp.")
		return nil
	}

	robots:=[]ProjectRobot{}

	for _,res:=range resLst {
		resMap,ok:=res.(map[string]interface{})
		if !ok {
			log.Println("GetRobotList error: can not convert item to map.")
			return nil
		}

		robotItem:=ProjectRobot{
			ID:resMap["id"].(string),
			Version:resMap["version"].(string),
		}

		robots=append(robots,robotItem)
	}

	return robots
}

func UpdateRobotList(crvClient *crv.CRVClient,robotList []RobotInfo,token string)(int){
	//获取已经存在的机器人列表
	currentRobotList:=GetRobotList(crvClient,token)

	saveRobotList:=[]map[string]interface{}{}
	//比较机器人列表，找出需要新增的机器人和需要删除的机器人
	//找出需要更新和删除的机器人
	for _,currentRobot:=range currentRobotList {
		//在新的机器人列表中找到当前机器人
		var robotInfo *RobotInfo 
		robotInfo=nil
		for _,robot:=range robotList {
			if currentRobot.ID==robot.RobotId {
				robotInfo=&robot
				break
			}
		}
		if robotInfo != nil {
			saveRobot:=map[string]interface{}{
				"id":currentRobot.ID,
				"version":currentRobot.Version,
				"name":robotInfo.RobotName,
				"_save_type":"update",
			}
			saveRobotList=append(saveRobotList,saveRobot)
		} else {
			//需要删除
			saveRobot:=map[string]interface{}{
				"id":currentRobot.ID,
				"version":currentRobot.Version,
				"_save_type":"delete",
			}
			saveRobotList=append(saveRobotList,saveRobot)
		}
	}
	//找出新增的机器人
	for _,robot:=range robotList {
		var currentRobot *ProjectRobot 
		currentRobot=nil
		for _,currentRobotItem:=range currentRobotList {
			if currentRobotItem.ID==robot.RobotId {
				currentRobot=&currentRobotItem
				break
			}
		}
		if currentRobot == nil {
			//需要新增
			saveRobot:=map[string]interface{}{
				"id":robot.RobotId,
				"name":robot.RobotName,
				"_save_type":"create",
			}
			saveRobotList=append(saveRobotList,saveRobot)
		}
	}

	//保存机器人列表
	commonRep:=crv.CommonReq{
		ModelID:MOdElID_ROBOT,
		List:&saveRobotList,
	}

	rsp,commonErr:=crvClient.Save(&commonRep,token)
	if commonErr!=common.ResultSuccess {
		return commonErr
	}

	if rsp.Error == true {
		log.Println("UpdateRobotList error:",rsp.ErrorCode,rsp.Message)
		return rsp.ErrorCode
	}

	return common.ResultSuccess
}

