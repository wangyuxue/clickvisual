package db

import (
	"errors"

	"github.com/gotomicro/ego-component/egorm"
	"github.com/gotomicro/ego/core/elog"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/clickvisual/clickvisual/api/internal/invoker"
)

type K8SConfigMap struct {
	BaseModel

	ClusterId int    `gorm:"column:cluster_id;type:int(11);index:uix_cluster_id_name_namespace,unique" json:"clusterId"` // 集群ID
	Name      string `gorm:"column:name;type:varchar(128);index:uix_cluster_id_name_namespace,unique" json:"name"`
	Namespace string `gorm:"column:namespace;type:varchar(128);index:uix_cluster_id_name_namespace,unique" json:"namespace"`
}

func (m *K8SConfigMap) TableName() string {
	return TableNameK8SConfigMap
}

// K8SConfigMapCreate CRUD
func K8SConfigMapCreate(db *gorm.DB, data *K8SConfigMap) (err error) {
	if err = db.Create(data).Error; err != nil {
		invoker.Logger.Error("create cluster error", zap.Error(err))
		return
	}
	return nil
}

// K8SConfigMapUpdate ...
func K8SConfigMapUpdate(db *gorm.DB, paramId int, ups map[string]interface{}) (err error) {
	var sql = "`id`=?"
	var binds = []interface{}{paramId}
	if err = db.Table(TableNameK8SConfigMap).Where(sql, binds...).Updates(ups).Error; err != nil {
		invoker.Logger.Error("update cluster error", zap.Error(err))
		return
	}
	return
}

// K8SConfigMapInfoX get single item by condition
func K8SConfigMapInfoX(conds map[string]interface{}) (resp K8SConfigMap, err error) {
	sql, binds := egorm.BuildQuery(conds)
	if err = invoker.Db.Table(TableNameK8SConfigMap).Where(sql, binds...).First(&resp).Error; err != nil && err != gorm.ErrRecordNotFound {
		invoker.Logger.Error("K8SConfigMapInfoX infoX error", zap.Error(err))
		return
	}
	return
}

// K8SConfigMapListX get single item by condition
func K8SConfigMapListX(conds map[string]interface{}) (resp []K8SConfigMap, err error) {
	sql, binds := egorm.BuildQuery(conds)
	if err = invoker.Db.Table(TableNameK8SConfigMap).Where(sql, binds...).Find(&resp).Error; err != nil && err != gorm.ErrRecordNotFound {
		invoker.Logger.Error("K8SConfigMapListX infoX error", zap.Error(err))
		return
	}
	return
}

func K8SConfigMapLoadOrSave(db *gorm.DB, data *K8SConfigMap) (resp *K8SConfigMap, err error) {
	conds := egorm.Conds{}
	conds["cluster_id"] = data.ClusterId
	conds["name"] = data.Name
	conds["namespace"] = data.Namespace
	respLoad, errLoad := K8SConfigMapInfoX(conds)
	if errLoad != nil {
		if errors.Is(errLoad, gorm.ErrRecordNotFound) {
			// Save
			errSave := K8SConfigMapCreate(db, data)
			if errSave != nil {
				return nil, errSave
			}
			return data, nil
		}
		return nil, errLoad
	}
	return &respLoad, nil
}

func K8SConfigMapInfo(paramId int) (resp K8SConfigMap, err error) {
	var sql = "`id`= ?"
	var binds = []interface{}{paramId}
	if err = invoker.Db.Table(TableNameK8SConfigMap).Where(sql, binds...).First(&resp).Error; err != nil && err != gorm.ErrRecordNotFound {
		invoker.Logger.Error("cluster info error", zap.Error(err))
		return
	}
	return
}

// K8SConfigMapDelete soft delete item by id
func K8SConfigMapDelete(db *gorm.DB, id int) (err error) {
	if err = db.Model(K8SConfigMap{}).Delete(&K8SConfigMap{}, id).Error; err != nil {
		invoker.Logger.Error("cluster delete error", zap.Error(err))
		return
	}
	return
}

// K8SConfigMapList return item list by condition
func K8SConfigMapList(conds egorm.Conds) (resp []*K8SConfigMap, err error) {
	sql, binds := egorm.BuildQuery(conds)
	// Fetch record with Rancher Info....
	if err = invoker.Db.Table(TableNameK8SConfigMap).Where(sql, binds...).Find(&resp).Error; err != nil && err != gorm.ErrRecordNotFound {
		invoker.Logger.Error("list clusters error", elog.String("err", err.Error()))
		return
	}
	return
}

// K8SConfigMapListPage return item list by pagination
func K8SConfigMapListPage(conds egorm.Conds, reqList *ReqPage) (total int64, respList []*K8SConfigMap) {
	respList = make([]*K8SConfigMap, 0)
	if reqList.PageSize == 0 {
		reqList.PageSize = 10
	}
	if reqList.Current == 0 {
		reqList.Current = 1
	}
	sql, binds := egorm.BuildQuery(conds)
	db := invoker.Db.Table(TableNameK8SConfigMap).Where(sql, binds...)
	db.Count(&total)
	db.Offset((reqList.Current - 1) * reqList.PageSize).Limit(reqList.PageSize).Find(&respList)
	return
}
