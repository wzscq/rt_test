package logfile

import (
    "testing"
	"log"
	"os"
)

func TestExtractDDIBFile(t *testing.T) {
	logRoot:="../testfiles/testlog"
	logFile:="001_161921001529_BFEBFBFF000806EC.zip"
	extractRoot:="../testfiles/ddib"
	fileName,errorCode:=extractDDIBFile(logRoot,logFile,extractRoot)

	//判断文件是否存在
	_,err:=os.Stat(extractRoot+"/"+fileName)
	if err!=nil {
		log.Println("TestExtractDDIBFile error: file not exist,fileName=",fileName)
		t.Fail()
		return
	}
	log.Println("ddib fileName=",fileName," errorCode=",errorCode)
}

func TestMoveReportToExtractRootFolder(t *testing.T){
	extractRoot:="../testfiles/ddib"
	reportFile:="../testfiles/testlog/001_161921001529_BFEBFBFF000806EC.zip"
	destFileName:="testdestfile"
	fileName,err:=moveReportToExtractRootFolder(reportFile,extractRoot,destFileName)
	if err!=nil {
		log.Println("TestMoveReportToExtractRootFolder error: move file failed,reportFile=",reportFile," extractRoot=",extractRoot," destFileName=",destFileName)
		t.Fail()
		return
	}

	//判断文件是否存在
	_,err=os.Stat(extractRoot+"/"+fileName)
	if err!=nil {
		log.Println("TestExtractDDIBFile error: file not exist,fileName=",extractRoot+"/"+fileName)
		t.Fail()
		return
	}

	//移动回去
	os.Rename(extractRoot+"/"+fileName,reportFile)
}


