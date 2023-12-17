package report

const (
	CMD_KPIReport="KPIReport"
	CMD_CustomReport="CustomReport"
	CMD_JOIN="="
	CMD_RET="R"
)

type Command struct {
	CommandType string `json:"CommandType"`
	Params interface{} `json:"Params"`
}

type Response struct {
	CommandType string `json:"CommandType"`
	Params []string `json:"Params"`
}

func GetKPICommand(files []string,template string)(*Command){
	params:=[]interface{}{
		files,
		"",
		"",
		template,
	}
	command:=Command{CommandType:CMD_KPIReport,Params:params}
	return &command
}