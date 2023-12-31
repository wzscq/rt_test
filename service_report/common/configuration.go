package common

import (
	"log"
	"os"
	"encoding/json"
)

type crvConf struct {
	Server string `json:"server"`
  AppID string `json:"appID"`
	Token string `json:"token"`
}

type DingliServerConf struct {
	Server string `json:"server"`
	Port string `json:"port"`
	Timeout string `json:"timeout"`
}

type serviceConf struct {
	Port string `json:"port"`
}

type TestLogConf struct {
	Root string `json:"root"`
	DownloadUrl string `json:"downloadUrl"`
	ExtractRoot string `json:"extractRoot"`
	ReportDownloadUrl string `json:"reportDownloadUrl"`
}

type Config struct {
	Service serviceConf `json:"service"`
	DingliServer DingliServerConf `json:"dingliServer"`
	TestLog TestLogConf `json:"testLog"`
	CRV crvConf `json:"crv"`
}

var gConfig Config

func InitConfig(confFile string)(*Config){
	log.Println("init configuation start ...")
	//获取用户账号
	//获取用户角色信息
	//根据角色过滤出功能列表
	fileName := confFile
	filePtr, err := os.Open(fileName)
	if err != nil {
        log.Fatal("Open file failed [Err:%s]", err.Error())
    }
    defer filePtr.Close()

	// 创建json解码器
    decoder := json.NewDecoder(filePtr)
    err = decoder.Decode(&gConfig)
	if err != nil {
		log.Println("json file decode failed [Err:%s]", err.Error())
	}
	log.Println("init configuation end")
	return &gConfig
}

func GetConfig()(*Config){
	return &gConfig
}