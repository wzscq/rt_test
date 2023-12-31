package logfile

import (
	"rt_report/common"
	"archive/zip"
	"strings"
	"log"
	"os"
	"io"
) 

func MoveReportToExtractRootFolder(reportFile,destFolder,destFileName string)(string,error){
	//获取reportFile的扩展名
	fileExt:=strings.ToLower(reportFile[strings.LastIndex(reportFile,"."):])
	destFileName=destFileName+fileExt
	fullDestFileName:=destFolder+"/"+destFileName

	//将reportFile移动到destFolder下
	err:=os.Rename(reportFile,fullDestFileName)
	if err!=nil {
		log.Println("moveReportToExtractRootFolder error: move file failed,reportFile=",reportFile," destFolder=",destFolder," destFileName=",destFileName)
		return "",err
	}

	return destFileName,nil
}

func ExtractDDIBFile(logPath,logFileName,extractPath string)(string,int){
	fullFileName:=logPath+"/"+logFileName
	//读取压缩文件
	r,err:=zip.OpenReader(fullFileName)
	if err!=nil {
		log.Println("extractDDIBFile error: open zip file failed,fileName=",fullFileName)
		return "",common.ResultOpenZipFileError
	}
	defer r.Close()

	//遍历压缩文件中的文件
	ddbibFileName:=""
	for _,f:=range r.File {
		//判断文件扩展名是否为ddib
		filePath := f.Name
		if !strings.HasSuffix(filePath,".ddib") {
			continue
		}
		log.Println("extractDDIBFile find ddib file:",filePath)
		ddbibFileName=filePath
		//将文件解压到指定目录
		dstFile, err := os.OpenFile(extractPath+"/"+ddbibFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			log.Println("extractDDIBFile error: open zip file failed,fileName=",ddbibFileName)
			return "",common.ResultOpenZipFileError
		}
		file, err := f.Open()
		if err != nil {
			log.Println("extractDDIBFile error: open zip file failed,fileName=",ddbibFileName)
			return "",common.ResultOpenZipFileError
		}
    	// 写入到解压到的目标文件
		if _, err := io.Copy(dstFile, file); err != nil {
			log.Println("extractDDIBFile error: open zip file failed,fileName=",ddbibFileName)
			return "",common.ResultOpenZipFileError
		}
		dstFile.Close()
		file.Close()
		break;
	}

	return ddbibFileName,common.ResultSuccess
}

func GetLogFileNameByCachedTestFile(logPath string,cachedTestFile *cachedTestFile)(string,int){
	//在logPath下查找文件名称中包含timestamp和deviceID的文件
	//如果找不到文件，返回空字符串
	//如果找到一个文件，返回文件名
	//如果找到多个文件，返回第一个文件名
	dir,err:=os.Open(logPath)
	if err!=nil {
		log.Println("getLogFileNameByCachedTestFile error: open logPath failed,logPath=",logPath)
		return "",common.ResultReadLogPathError
	}
	defer dir.Close()

	fileNames,err:=dir.Readdirnames(-1)
	if err!=nil {
		log.Println("getLogFileNameByCachedTestFile error: read logPath failed,logPath=",logPath)
		return "",common.ResultReadLogPathError
	}

	var logFileName string
	for _,fileName:=range fileNames {
		if !strings.HasSuffix(fileName,".zip") {
			continue
		}

		if strings.Contains(fileName,cachedTestFile.Timestamp) && strings.Contains(fileName,cachedTestFile.DeviceID) {
			logFileName=fileName
			break
		}
	}

	if logFileName=="" {
		log.Println("getLogFileNameByCachedTestFile error: can not find log file with timestamp:",cachedTestFile.Timestamp," and deviceID:",cachedTestFile.DeviceID)
		return "",common.ResultNoLogFileError
	}

	return logFileName,common.ResultSuccess
}