package common

type CommonRsp struct {
	ErrorCode int `json:"errorCode"`
	Message string `json:"message"`
	Error bool `json:"error"`
	Result interface{} `json:"result"`
	Params map[string]interface{} `json:"params"`
}

type CommonError struct {
	ErrorCode int `json:"errorCode"`
	Params map[string]interface{} `json:"params"`
}

const (
	ResultSuccess = 10000000
	ResultWrongRequest = 10000001
	ResultCreateDirError=10000027
	ResultBase64DecodeError=10000028
	ResultCreateFileError=10000029
	ResultGenerateReportError=10100030
	ResultNotSupportedReportType=10100031
	ResultSaveDataError = 10200005
	ResultQueryRequestError = 10200006
	ResultGetCachedTestFileInfoError = 10200007
	ResultNoLogFileError = 10200008
	ResultReadLogPathError = 10200009
	ResultOpenZipFileError = 10200010
)

var errMsg = map[int]CommonRsp{
	ResultSuccess:CommonRsp{
		ErrorCode:ResultSuccess,
		Message:"操作成功",
		Error:false,
	},
	ResultSaveDataError:CommonRsp{
		ErrorCode:ResultSaveDataError,
		Message:"保存数据到数据时发生错误，请与管理员联系处理",
		Error:true,
	},
	ResultQueryRequestError:CommonRsp{
		ErrorCode:ResultQueryRequestError,
		Message:"CRVClient发送查询参数请求失败，请与管理员联系处理",
		Error:true,
	},
	ResultWrongRequest:CommonRsp{
		ErrorCode:ResultWrongRequest,
		Message:"请求参数错误，请检查参数是否完整，参数格式是否正确",
		Error:true,
	},
	ResultNotSupportedReportType:CommonRsp{
		ErrorCode:ResultWrongRequest,
		Message:"创建报表时遇到不支持的报表类型，请确认报表类型是否正确",
		Error:true,
	},
	ResultCreateDirError:CommonRsp{
		ErrorCode:ResultCreateDirError,
		Message:"保存文件时创建文件夹失败，请与管理员联系处理",
		Error:true,
	},
	ResultBase64DecodeError:CommonRsp{
		ErrorCode:ResultBase64DecodeError,
		Message:"保存文件时文件内容Base64解码失败，请与管理员联系处理",
		Error:true,
	},
	ResultCreateFileError:CommonRsp{
		ErrorCode:ResultCreateFileError,
		Message:"创建文件失败，请与管理员联系处理",
		Error:true,
	},
	ResultGenerateReportError:CommonRsp{
		ErrorCode:ResultGenerateReportError,
		Message:"生成报表文件失败，请与管理员联系处理",
		Error:true,
	},
	ResultGetCachedTestFileInfoError:CommonRsp{
		ErrorCode:ResultGetCachedTestFileInfoError,
		Message:"获取测试文件信息失败，请与管理员联系处理",
		Error:true,
	},
	ResultNoLogFileError:CommonRsp{
		ErrorCode:ResultNoLogFileError,
		Message:"没有找到对应的日志文件，请确认测试文件是否存在",
		Error:true,
	},
	ResultReadLogPathError:CommonRsp{
		ErrorCode:ResultReadLogPathError,
		Message:"读取日志文件夹失败，请与管理员联系处理",
		Error:true,
	},
	ResultOpenZipFileError:CommonRsp{
		ErrorCode:ResultOpenZipFileError,
		Message:"打开zip文件失败，请与管理员联系处理",
		Error:true,
	},
}

func CreateResponse(err *CommonError,result interface{})(*CommonRsp){
	if err==nil {
		commonRsp:=errMsg[ResultSuccess]
		commonRsp.Result=result
		return &commonRsp
	}

	commonRsp:=errMsg[err.ErrorCode]
	commonRsp.Result=result
	commonRsp.Params=err.Params
	return &commonRsp
}

func CreateError(errorCode int,params map[string]interface{})(*CommonError){
	return &CommonError{
		ErrorCode:errorCode,
		Params:params,
	}
}