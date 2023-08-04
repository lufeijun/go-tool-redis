package datacube

import "github.com/lufeijun/go-tool-wechat/wechat/officialaccount/wcontext"

type reqDate struct {
	BeginDate string `json:"begin_date"`
	EndDate   string `json:"end_date"`
}

// DataCube 数据统计
type DataCube struct {
	*wcontext.Context
}

// NewCube 数据统计
func NewCube(context *wcontext.Context) *DataCube {
	dataCube := new(DataCube)
	dataCube.Context = context
	return dataCube
}
