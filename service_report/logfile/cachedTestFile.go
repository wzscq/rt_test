package logfile

import (
	"rt_report/common"
	"rt_report/crv"
	"log"
)

type cachedTestFile struct {
	Timestamp string
	DeviceID string
	ID string
}

var	testFileFields=[]map[string]interface{}{
	{"field": "id"},
	{"field": "timestamp"},
	{"field": "device_id"},
}

var MODELID_CACHED_TESTFILE="rt_cache_test_file"

func GetCachedTestFile(crvClient *crv.CRVClient,cachedTestFileID,token string)(*cachedTestFile,int){
	//调用crv接口获取测试文件信息
	filter:=map[string]interface{}{
		"id":cachedTestFileID,
	}

	commonRep:=crv.CommonReq{
		ModelID:MODELID_CACHED_TESTFILE,
		Fields:&testFileFields,
		Filter:&filter,
	}

	rsp,errorCode:=crvClient.Query(&commonRep,token)
	if errorCode!=common.ResultSuccess {
		return nil,errorCode
	}

	//解析返回结果
	if rsp.Result==nil {
		log.Println("getCachedTestFile error: the result in rsp is nil.")
		return nil,common.ResultGetCachedTestFileInfoError
	}

	resLst,ok:=rsp.Result.(map[string]interface{})["list"].([]interface{})
	if !ok {
		log.Println("getCachedTestFile error: no list in result.")
		return nil,common.ResultGetCachedTestFileInfoError
	}

	if len(resLst)==0 {
		log.Println("getCachedTestFile error: the list of result is empty.")
		return nil,common.ResultGetCachedTestFileInfoError
	}

	res,ok:=resLst[0].(map[string]interface{})
	if !ok {
		log.Println("getCachedTestFile error: the first element of list is not map.")
		return nil,common.ResultGetCachedTestFileInfoError
	}

	timestamp,ok:=res["timestamp"].(string)
	if !ok {
		log.Println("getCachedTestFile error: the timestamp of result is not string.")
		return nil,common.ResultGetCachedTestFileInfoError
	}

	deviceID,ok:=res["device_id"].(string)
	if !ok {
		log.Println("getCachedTestFile error: the device_id of result is not string.")
		return nil,common.ResultGetCachedTestFileInfoError
	}

	id,ok:=res["id"].(string)
	if !ok {
		log.Println("getCachedTestFile error: the id of result is not string.")
		return nil,common.ResultGetCachedTestFileInfoError
	}

	//解析返回结果
	cachedTestFile:=&cachedTestFile{
		Timestamp:timestamp,
		DeviceID:deviceID,
		ID:id,
	}

	return cachedTestFile,common.ResultSuccess
}