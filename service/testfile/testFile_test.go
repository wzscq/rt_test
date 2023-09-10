package testfile

import (
    "testing"
	"time"
	"strconv"
	"log"
	"encoding/json"
	"rt_test_service/crv"
)

var reportData=map[string]interface{}{
	"data": []map[string]interface{}{
			{
				"radio": map[string]interface{}{
					"measures_common": map[string]interface{}{
						"Current Network Type": "",
						"IMSI": "460011895631209",
						"IMEI": "864451045593183",
						"imsi": "460011895631209",
						"imei": "864451045593183",
						"duplex_mode": "",
						"net_Tye": "NR_SA",
						"plmn": "46001",
						"operator": "\u4e2d\u56fd\u8054\u901a",
						"operator_country": "cn",
						"operator_short": "CHN-UNICOM",
						"phone_name": "Redmi K30 5G Speed",
					},
					"measures_nr": map[string]interface{}{
						"NR PHY Throughput UL(CA Total)": "231528",
						"NR PHY Throughput DL(CA Total)": "3723656",
						"NR C-RNTI": "",
						"NR PUSCH Initial BLER": "7.69",
						"NR PUSCH BLER": "7.14",
						"NR PDSCH Initial BLER": "2.25",
						"NR PDSCH BLER": "2.2",
						"NR UL Avg MCS": "24.21",
						"NR DL Avg MCS": "10.76",
						"NR High Modulation UL/s": "3",
						"NR High Modulation DL/s": "4",
						"NR PUSCH TxPower": "15",
						"NR PUCCH TxPower": "19",
						"NR PRACH TxPower": "",
						"NR PDCCH UL GrantCount": "45",
						"NR PDCCH DL GrantCount": "99",
						"SS-SINR": "11.37",
						"SS-RSRP": "-84.59",
						"NR Network State": "NR Connected",
						"UE Category": "",
						"NR Slot Config UL Total": "",
						"NR Slot Config DL Total": "",
						"NR CSI RS Period": "",
						"NR SSB Period": "",
						"SSB Config Beam Num": "",
						"NR PCI": "53",
						"SSB ARFCN": "627264",
						"NR SSB Frequency(MHz)": "3408960",
						"NR WB CQI": "11.31",
						"NCI": "",
						"NR Bandwidth(MHz)": "",
						"NR Band": "",
						"NR TAC": "",
					},
				},
				"event": []map[string]interface{}{
					{
						"Logfile": "testRecoder_0818-212920",
						"pointIndex": 204,
						"EventIndex": 11,
						"EventTime": "2023-08-18 21:29:29.121",
						"EventCode": "0x8D",
						"name": "FTP Download TCPSlow",
					},
				},
				"msg": []map[string]interface{}{
					{
						"Logfile": "testRecoder_0818-212920",
						"pointIndex": 430,
						"MsgTime": "2023-08-18 21:29:33.894",
						"MsgCode": "0x50000200",
						"name": "NR->MeasurementReport",
					},
				},
				"time": 1692365376,
				"case_progress": map[string]interface{}{
					"session_id": 161920005569,
					"full_path": "d:\\test_recoder\\testReport_001_161920005569\\progress_460011895631209.csv",
					"imsi": "460011895631209",
					"DevId": "28e79a383c89153af87780ec745e2481",
					"Status": "FTPDownload",
					"Times": "1/1",
					"Progress": "8/30s",
					"FailTimes": 0,
					"OtherInfo": "Inst:2,729.30 kbps  Avg:2,799.36 kbps SuccRatio:100%(1/1)",
				},
			},
	},
	"robot_info": map[string]interface{}{
			"id": 1,
			"robot_id": "2bee174b7d7c36e9b98bd8772e66af5e",
			"map_id": "9ec5a9c61a374be4b0570077d7d6dd05",
			"pixel_x": 616.645741420779,
			"pixel_y": 285.72210123115497,
			"pixel_theta": -110.04722431915482,
			"record": "2023-08-18 21:27:05",
	},
	"pcTime": 1692365376,
}

func _TestCreateFile(t *testing.T) {
	outPath:="../localcache/"
	deviceID:="device1"
	timeStamp:=time.Now().Unix()
	tf:=GetTestFile(outPath,deviceID,timeStamp)
	if tf==nil {
		t.Error("GetTestFile failed")
		return
	}

	for i:=0;i<10;i++ {
		lineContent:="line"+strconv.Itoa(i)
		tf.WriteLine(lineContent)
	}
	
	tf.Close()
}

func _TestReadFile(t *testing.T){
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	outPath:="../localcache/"
	deviceID:="device1"
	var timeStamp int64 
	timeStamp=1693636403
	tf:=GetTestFile(outPath,deviceID,timeStamp)
	if tf==nil {
		t.Error("GetTestFile failed")
		return
	}
	defer tf.Close()

	idxContent:=tf.GetIdxContent()
	log.Println(idxContent)

	for i:=0;i<10;i++ {
		lines:=tf.GetContent(i,i+1)
		if lines==nil || len(lines)==0 {
			t.Error("GetContent failed i:="+strconv.Itoa(i))
			return
		}

		if lines[0]!="line"+strconv.Itoa(i) {
			t.Error("GetContent failed i:="+strconv.Itoa(i))
			return
		}
	}
}

func TestTestFilePool(t *testing.T){
	
	crvClient:=&crv.CRVClient{
		Server:"http://127.0.0.1:8200",
		Token:"rt_test_service",
		AppID:"rt_test",
	}
	
	//reportData转换为JSON字符串
	reportDataJson,_:=json.Marshal(reportData)
	//创建TestFilePool
	tfp:=InitTestFilePool("../localcache/","3s",crvClient)
	//写入数据
	tfp.WriteDeviceTestLine("device2",string(reportDataJson))
	//等待5秒
	for i:=0;i<5;i++ {
		time.Sleep(2*time.Second)
		tfp.WriteDeviceTestLine("device2",string(reportDataJson))
	}
	time.Sleep(8*time.Second)
}

