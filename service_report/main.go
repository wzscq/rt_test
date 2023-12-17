package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"rt_report/common"
	"rt_report/report"
	"log"
	"os"
)

func main() {
	//设置log打印文件名和行号
  	log.SetFlags(log.Lshortfile | log.LstdFlags)

	confFile:="conf/conf.json"
	if len(os.Args)>1 {
		confFile=os.Args[1]
		log.Println(confFile)
	}

	conf:=common.InitConfig(confFile)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:true,
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
  	}))

	dingliClient:=report.DingliClient{DingliServer:&conf.DingliServer}
	reportController:=report.ReportController{DingliClient:&dingliClient}
	reportController.Bind(router)

	router.Run(conf.Service.Port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}