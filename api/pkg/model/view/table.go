package view

type ReqTableCreateExist struct {
	DatabaseName  string `form:"databaseName" json:"databaseName" binding:"required"`
	TableName     string `form:"tableName" json:"tableName" binding:"required"`
	TimeField     string `form:"timeField" json:"timeField"`
	TimeFieldType int    `form:"timeFieldType" json:"timeFieldType"`
	Desc          string `form:"desc" json:"desc"`
}

type ReqTableCreateExistBatch struct {
	TableList []ReqTableCreateExist `form:"tableList" json:"tableList"`
}

type ReqTableCreate struct {
	TableName string `form:"tableName" binding:"required"`
	Typ       int    `form:"typ" binding:"required"`
	Days      int    `form:"days" binding:"required"`
	Brokers   string `form:"brokers" binding:"required"`
	Topics    string `form:"topics" binding:"required"`
	Consumers int    `form:"consumers" binding:"required"`
	Desc      string `form:"desc"`
}

type ReqTableId struct {
	Instance   string `form:"instance" binding:"required"`
	Database   string `form:"database" binding:"required"`
	Table      string `form:"table" binding:"required"`
	Datasource string `form:"datasource" binding:"required"`
}

type RespTableSimple struct {
	Id         int    `json:"id"`
	TableName  string `json:"tableName"`
	CreateType int    `json:"createType"`
	Desc       string `json:"desc"`
}

type RespTableDetail struct {
	Did        int    `json:"did"`     // 数据库 id
	Name       string `json:"name"`    // table
	Typ        int    `json:"typ"`     // table 类型 1 app 2 ego 3 ingress
	Days       int    `json:"days"`    // 数据过期时间
	Brokers    string `json:"brokers"` // kafka broker
	Topic      string `json:"topic"`   // kafka topic
	Uid        int    `json:"uid"`     // 操作人
	Desc       string `json:"desc"`    //
	SQLContent struct {
		Keys []string          `json:"keys"`
		Data map[string]string `json:"data"`
	} `json:"sqlContent"`
	Database   RespDatabaseItem `json:"database"`
	CreateType int              `json:"createType"`
	TimeField  string           `json:"timeField"`
	Ctime      int64            `json:"ctime"`
	Utime      int64            `json:"utime"`
}

type RespColumn struct {
	Name     string `json:"name"`
	TypeDesc string `json:"typeDesc"`
	Type     int    `json:"type"`
}

type RespDatabaseSelfBuilt struct {
	Name   string                 `json:"name"`
	Tables []*RespTablesSelfBuilt `json:"tables"`
}

type RespTablesSelfBuilt struct {
	Name string `json:"name"`
}
