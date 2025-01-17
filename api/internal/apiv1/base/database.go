package base

import (
	"strconv"

	"github.com/gotomicro/ego-component/egorm"
	"github.com/kl7sn/toolkit/xgo"
	"github.com/spf13/cast"

	"github.com/clickvisual/clickvisual/api/internal/invoker"
	"github.com/clickvisual/clickvisual/api/internal/service"
	"github.com/clickvisual/clickvisual/api/internal/service/event"
	"github.com/clickvisual/clickvisual/api/internal/service/permission"
	"github.com/clickvisual/clickvisual/api/internal/service/permission/pmsplugin"
	"github.com/clickvisual/clickvisual/api/pkg/component/core"
	"github.com/clickvisual/clickvisual/api/pkg/model/db"
	"github.com/clickvisual/clickvisual/api/pkg/model/view"
)

func DatabaseCreate(c *core.Context) {
	iid := cast.ToInt(c.Param("iid"))
	if iid == 0 {
		c.JSONE(core.CodeErr, "invalid parameter", nil)
		return
	}
	var req view.ReqDatabaseCreate
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "invalid parameter: "+err.Error(), nil)
		return
	}

	if err := permission.Manager.CheckNormalPermission(view.ReqPermission{
		UserId:      c.Uid(),
		ObjectType:  pmsplugin.PrefixInstance,
		ObjectIdx:   strconv.Itoa(iid),
		SubResource: pmsplugin.InstanceBase,
		Acts:        []string{pmsplugin.ActEdit},
	}); err != nil {
		c.JSONE(1, err.Error(), nil)
		return
	}

	obj := db.Database{
		Iid:            iid,
		Name:           req.Name,
		Cluster:        req.Cluster,
		Uid:            c.Uid(),
		IsCreateByMogo: 1,
	}
	op, err := service.InstanceManager.Load(iid)
	if err != nil {
		c.JSONE(core.CodeErr, err.Error(), nil)
		return
	}
	tx := invoker.Db.Begin()
	if err = db.DatabaseCreate(tx, &obj); err != nil {
		c.JSONE(1, "create failed 01: "+err.Error(), nil)
		return
	}
	if req.Cluster != "" {
		xgo.Go(func() {
			_ = op.DatabaseCreate(req.Name, req.Cluster)
		})
		if err = tx.Commit().Error; err != nil {
			tx.Rollback()
			c.JSONE(1, "alarm create failed 03: "+err.Error(), nil)
			return
		}
		c.JSONOK("cluster database creation takes more time")
	} else {
		err = op.DatabaseCreate(req.Name, req.Cluster)
		if err != nil {
			tx.Rollback()
			c.JSONE(core.CodeErr, "create failed: "+err.Error(), nil)
			return
		}
		if err = tx.Commit().Error; err != nil {
			tx.Rollback()
			c.JSONE(1, "alarm create failed 03: "+err.Error(), nil)
			return
		}
		c.JSONOK()
	}
	event.Event.AlarmCMDB(c.User(), db.OpnDatabasesCreate, map[string]interface{}{"database": obj})
	return
}

func DatabaseExistList(c *core.Context) {
	iid := cast.ToInt(c.Param("iid"))
	if iid == 0 {
		c.JSONE(1, "param error: missing iid", nil)
		return
	}
	op, err := service.InstanceManager.Load(iid)
	if err != nil {
		c.JSONE(core.CodeErr, err.Error(), nil)
		return
	}
	res, err := op.Databases()
	if err != nil {
		c.JSONE(core.CodeErr, err.Error(), nil)
		return
	}
	c.JSONE(core.CodeOK, "succ", res)
	return
}

func DatabaseList(c *core.Context) {
	iid := cast.ToInt(c.Param("iid"))
	conds := egorm.Conds{}
	if iid != 0 {
		conds["iid"] = iid
	}
	dl, err := db.DatabaseList(invoker.Db, conds)
	if err != nil {
		c.JSONE(core.CodeErr, err.Error(), nil)
		return
	}
	res := make([]view.RespDatabaseItem, 0)
	for _, row := range dl {
		if !service.InstanceManager.ReadPermissionDatabase(c.Uid(), row.Iid, row.ID) {
			continue
		}
		tmp := view.RespDatabaseItem{
			Id:   row.ID,
			Iid:  row.Iid,
			Name: row.Name,
			Uid:  row.Uid,
		}
		if row.Instance != nil {
			tmp.DatasourceType = row.Instance.Datasource
			tmp.InstanceName = row.Instance.Name
			tmp.Mode = row.Instance.Mode
			tmp.Clusters = row.Instance.Clusters
		}
		res = append(res, tmp)
	}
	c.JSONE(core.CodeOK, "succ", res)
	return
}

func DatabaseDelete(c *core.Context) {
	id := cast.ToInt(c.Param("id"))
	if id == 0 {
		c.JSONE(1, "invalid parameter", nil)
		return
	}
	database, err := db.DatabaseInfo(invoker.Db, id)
	if err != nil {
		c.JSONE(1, "failed to delete database: "+err.Error(), nil)
		return
	}
	if err = permission.Manager.CheckNormalPermission(view.ReqPermission{
		UserId:      c.Uid(),
		ObjectType:  pmsplugin.PrefixInstance,
		ObjectIdx:   strconv.Itoa(database.Iid),
		SubResource: pmsplugin.InstanceBase,
		Acts:        []string{pmsplugin.ActDelete},
		DomainType:  pmsplugin.PrefixDatabase,
		DomainId:    strconv.Itoa(database.ID),
	}); err != nil {
		c.JSONE(1, err.Error(), nil)
		return
	}
	conds := egorm.Conds{}
	conds["did"] = id
	tables, err := db.TableList(invoker.Db, conds)
	if len(tables) > 0 {
		c.JSONE(1, "you should delete all tables before delete database", nil)
		return
	}
	if database.IsCreateByMogo == 1 {
		op, errLoad := service.InstanceManager.Load(database.Iid)
		if errLoad != nil {
			c.JSONE(core.CodeErr, errLoad.Error(), nil)
			return
		}
		err = op.DropDatabase(database.Name, database.Cluster)
		if err != nil {
			c.JSONE(core.CodeErr, err.Error(), nil)
			return
		}
	}
	err = db.DatabaseDelete(invoker.Db, id)
	if err != nil {
		c.JSONE(1, "failed to delete database, corresponding record does not exist in database: "+err.Error(), nil)
		return
	}
	event.Event.AlarmCMDB(c.User(), db.OpnDatabasesDelete, map[string]interface{}{"database": database})
	c.JSONOK()
}
