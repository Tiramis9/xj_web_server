package handler

import (
	"context"
	"fmt"
	"xj_web_server/model"
	"xj_web_server/module/db/db"
	"xj_web_server/module/db/proto"
	"xj_web_server/util"
)

// Db : 用于实现DbServiceHandler接口的对象
type DbService struct{}

func (dbs *DbService) GetHost(ctx context.Context, req *proto.DbReqHost, res *proto.DbRespHost) error {

	fmt.Println("db.GetDB()", db.GetDB())

	hosts, err := model.GetHost(db.GetDB())
	fmt.Println("db.GetDB()1", db.GetDB())
	if err != nil {
		res.Code = 1
		res.Message = err.Error()
		return nil
	}

	var respHosts []*proto.DbHost

	for _, host := range hosts {
		respHosts = append(respHosts, &proto.DbHost{Id: 0, HostName: "", Ip: host.ServerAddr, Port: "0"})
	}

	res.Code = util.SuccessCode
	res.Message = "获取成功"
	res.Host = respHosts

	return nil
}
