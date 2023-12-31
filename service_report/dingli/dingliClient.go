package dingli

import (
	"rt_report/common"
	"encoding/json"
	"log"
	"net"
	"sync"
	"time"
	"errors"
)

var connect_mutex sync.Mutex

type DingliClient struct {
	DingliServer *common.DingliServerConf
}

func (cl *DingliClient)CreateReport(reportType,templete string,files []string)(string,error){
	if reportType == CMD_KPIReport {
		return cl.GetKPIReport(files,templete)
	} else if reportType == CMD_CustomReport {
		return cl.GetCustomerReport(files,templete)
	}
	return "",errors.New("not supported report type: "+reportType)
}

func (cl *DingliClient)GetKPIReport(files []string,template string)(string,error){
	//这里需要加锁，防止多个线程同时调用
	connect_mutex.Lock()
    defer connect_mutex.Unlock()
	//链接服务器
	// 建立客户端连接
	timeoutDuration,_ := time.ParseDuration(cl.DingliServer.Timeout) 
	server:=cl.DingliServer.Server+":"+cl.DingliServer.Port
	log.Println("connect to dingli server:", server)
    conn, err := net.DialTimeout("tcp",server ,timeoutDuration)
    if err != nil {
        //connect_mutex.Unlock()
        log.Println("connect to dingli server error:", err)
		return "",err 
    }
    defer conn.Close()

	cmd:=GetKPICommand(files,template)

	//结构转换为json
	jsonCmd,err:=json.Marshal(cmd)
	if err!=nil{
		log.Println("json marshal error:", err)
		return "",err
	}

	//发送命令
	log.Println("send command:", string(jsonCmd))
	_, err = conn.Write(jsonCmd)
	if err != nil {
		log.Println("send command error:", err)
		return "",err
	}

	//接收返回值
	buffer := make([]byte, 1024)
	length, err := conn.Read(buffer)
	if err != nil {
		log.Println("receive response error:", err) 
		return "",err
	}

	var rsp Response
	err = json.Unmarshal(buffer[:length], &rsp)
	if err != nil {
		log.Println("json unmarshal error:", err)
		return "",err
	}

	log.Println("receive response:", rsp)

	return rsp.Params[2],nil
}

func (cl *DingliClient)GetCustomerReport(files []string,template string)(string,error){
	//这里需要加锁，防止多个线程同时调用
	connect_mutex.Lock()
  defer connect_mutex.Unlock()
	//链接服务器
	// 建立客户端连接
	timeoutDuration,_ := time.ParseDuration(cl.DingliServer.Timeout) 
	server:=cl.DingliServer.Server+":"+cl.DingliServer.Port
	log.Println("connect to dingli server:", server)
	conn, err := net.DialTimeout("tcp",server ,timeoutDuration)
	if err != nil {
		//connect_mutex.Unlock()
		log.Println("connect to dingli server error:", err)
		return "",err 
	}
	defer conn.Close()

	cmd:=GetCustomerCommand(files,template)

	//结构转换为json
	jsonCmd,err:=json.Marshal(cmd)
	if err!=nil{
		log.Println("json marshal error:", err)
		return "",err
	}

	//发送命令
	log.Println("send command:", string(jsonCmd))
	_, err = conn.Write(jsonCmd)
	if err != nil {
		log.Println("send command error:", err)
		return "",err
	}

	//接收返回值
	buffer := make([]byte, 1024)
	length, err := conn.Read(buffer)
	if err != nil {
		log.Println("receive response error:", err) 
		return "",err
	}

	var rsp Response
	err = json.Unmarshal(buffer[:length], &rsp)
	if err != nil {
		log.Println("json unmarshal error:", err)
		return "",err
	}

	log.Println("receive response:", rsp)

	return rsp.Params[2],nil
}



