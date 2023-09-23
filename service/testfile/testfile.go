package testfile

import (
	"os"
	"log"
	"bufio"
	"strconv"
	"io"
	"encoding/json"
)

type TestFile struct {
	OutPath string
	DeviceID string
	TimeStamp int64
	ContentFile *os.File
	LineCount int64
	//下面两个字段用于控制文件超时关闭
	lastLineCount int64
}

func (tf *TestFile) Close() {
	indexFileName:=tf.OutPath+"/"+tf.DeviceID+"_"+strconv.FormatInt(tf.TimeStamp,10)+".idx"
	idxFile,err:=os.OpenFile(indexFileName,os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Printf("Open file failed [Err:%s]\n", err.Error())
	} else {
		idxFile.WriteString(strconv.FormatInt(tf.LineCount,10)+"\n")
		idxFile.Close()
	}
	tf.ContentFile.Close()
}

func (tf *TestFile) CloseReadOnly() {
	tf.ContentFile.Close()
}

func (tf *TestFile) WriteLine(lineContent string) {
	tf.ContentFile.WriteString(lineContent+"\n")
	tf.LineCount++
}

func (tf *TestFile)GetIdxContent()(string){
	indexFileName:=tf.OutPath+"/"+tf.DeviceID+"_"+strconv.FormatInt(tf.TimeStamp,10)+".idx"
	idxFile,err:=os.Open(indexFileName)
	if err != nil {
		log.Printf("Open file failed [Err:%s]\n", err.Error())
	}
	defer idxFile.Close()
		
	bytes:=make([]byte,1024)
	n,err:=idxFile.Read(bytes)
	if err != nil {
		log.Printf("Read file failed [Err:%s]\n", err.Error())
		return ""
	}

	if n==0 {
		return ""
	}

	return string(bytes[:n])
}

func (tf *TestFile)GetContent(from,to int64)([]string){
	lines:=[]string{}
	//文件复位
	_, err := tf.ContentFile.Seek(0, io.SeekStart)
	if err != nil {
		log.Fatal(err)
		return lines
	}

	scanner := bufio.NewScanner(tf.ContentFile)
	var n int64
	n = -1
	for scanner.Scan() {
		n++
		lineStr:=string(scanner.Bytes())
		log.Printf("from:%d,to:%d,n:%d,line:%s",from,to,n,lineStr)
        if n < from {
            continue
        }
		lines=append(lines,lineStr)
		if n >= to {
            break
        }
		log.Println(lines)
    }
	log.Println(lines)
	return lines
}

func (tf *TestFile)GetPoints()([]map[string]interface{}){
	//文件复位
	_, err := tf.ContentFile.Seek(0, io.SeekStart)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	lines:=[]map[string]interface{}{}
	scanner := bufio.NewScanner(tf.ContentFile)
	for scanner.Scan() {
		var result map[string]interface{}
		err := json.Unmarshal(scanner.Bytes(), &result)
		if err != nil {
			log.Println(err)
			return nil
		}

		robotInfo:=result["robot_info"].(map[string]interface{})
		lines=append(lines,robotInfo)
  }
	return lines
}

func GetTestFile(outPath string,deviceID string,timeStamp int64) *TestFile {
	contentFileName:=outPath+"/"+deviceID+"_"+strconv.FormatInt(timeStamp,10)+".content"
	contentFile,err:=os.OpenFile(contentFileName,os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Printf("Open file failed [Err:%s]\n", err.Error())
		return nil
	}

	return &TestFile{
		DeviceID:deviceID,
		TimeStamp:timeStamp,
		ContentFile:contentFile,
		OutPath:outPath,		
	}
}