package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	sdelete "gorm.io/plugin/soft_delete"
)

const (
	TableMogoAlarm                = "mogo_alarm"
	TableNameUser                 = "mogo_user"
	TableMogoEvent                = "mogo_event"
	TableNameK8SConfigMap         = "mogo_k8s_cm"
	TableNameCluster              = "mogo_cluster"
	TableNameView                 = "mogo_base_view"
	TableNameTable                = "mogo_base_table"
	TableNameIndex                = "mogo_base_index"
	TableMogoAlarmFilter          = "mogo_alarm_filter"
	TableNameDatabase             = "mogo_base_database"
	TableNameInstance             = "mogo_base_instance"
	TableNameConfiguration        = "mogo_configuration"
	TableMogoAlarmHistory         = "mogo_alarm_history"
	TableMogoAlarmChannel         = "mogo_alarm_channel"
	TableMogoAlarmCondition       = "mogo_alarm_condition"
	TableNameConfigurationHistory = "mogo_configuration_history"
	TableNameConfigurationPublish = "mogo_configuration_publish"

	TableNamePmsRole         = "mogo_pms_role"
	TableNamePmsRoleDetail   = "mogo_pms_role_detail"
	TableNamePmsRoleRef      = "mogo_pms_role_ref"
	TableNamePmsRoleRefGrant = "mogo_pms_role_ref_grant"
	TableNamePmsCasbinRule   = "mogo_pms_casbin_rule"
	TableNamePmsDefaultRole  = "mogo_pms_default_role"
	TableNamePmsCustomRole   = "mogo_pms_custom_role"
)

type BaseModel struct {
	ID    int               `gorm:"not null;primary_key;AUTO_INCREMENT" json:"id"`
	Ctime int64             `gorm:"bigint;autoCreateTime;comment:创建时间" json:"ctime"`
	Utime int64             `gorm:"bigint;autoUpdateTime;comment:更新时间" json:"utime"`
	Dtime sdelete.DeletedAt `gorm:"bigint;comment:删除时间" json:"dtime"`
}

type ReqPage struct {
	Current  int `json:"current" form:"current"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

func (r *ReqPage) Valid() error {
	if r.Current == 0 {
		r.Current = 1
	}
	if r.PageSize == 0 {
		r.PageSize = 10
	}
	if r.Current < 0 {
		return fmt.Errorf("current MUST be larger than 0")
	}
	if r.PageSize < 0 {
		return fmt.Errorf("invalid pageSize")
	}
	return nil
}

type String2String map[string]string

func (t String2String) Value() (driver.Value, error) {
	b, err := json.Marshal(t)
	return string(b), err
}

func (t *String2String) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), t)
}

type Strings []string

func (t Strings) Value() (driver.Value, error) {
	b, err := json.Marshal(t)
	return string(b), err
}

func (t *Strings) Scan(input interface{}) error {
	in := input.([]byte)
	if len(in) == 0 {
		in = []byte("[]")
	}
	return json.Unmarshal(in, t)
}

type Ints []int

func (t Ints) Value() (driver.Value, error) {
	b, err := json.Marshal(t)
	return string(b), err
}

func (t *Ints) Scan(input interface{}) error {
	if len(input.([]byte)) == 0 {
		return json.Unmarshal([]byte("[]"), t)
	}
	if err := json.Unmarshal(input.([]byte), t); err != nil {
		return json.Unmarshal([]byte("[]"), t)
	}
	return json.Unmarshal(input.([]byte), t)
}
