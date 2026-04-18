// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// LogGameDao is the data access object for the table log_game.
type LogGameDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  LogGameColumns     // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// LogGameColumns defines and stores column names for the table log_game.
type LogGameColumns struct {
	Id        string // 自增id
	RoomId    string // 房间id
	Type      string // 房间类型：1段位房，2好友房
	Status    string // 房间状态
	UserId    string // 用户id
	Point     string // 用户分数
	Action    string // 行为
	Remain    string // 剩余的牌
	OutCards  string // 打出的牌
	Text      string // 完整信息
	CreatTime string // 创建时间
}

// logGameColumns holds the columns for the table log_game.
var logGameColumns = LogGameColumns{
	Id:        "id",
	RoomId:    "room_id",
	Type:      "type",
	Status:    "status",
	UserId:    "user_id",
	Point:     "point",
	Action:    "action",
	Remain:    "remain",
	OutCards:  "out_cards",
	Text:      "text",
	CreatTime: "creat_time",
}

// NewLogGameDao creates and returns a new DAO object for table data access.
func NewLogGameDao(handlers ...gdb.ModelHandler) *LogGameDao {
	return &LogGameDao{
		group:    "default",
		table:    "log_game",
		columns:  logGameColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *LogGameDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *LogGameDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *LogGameDao) Columns() LogGameColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *LogGameDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *LogGameDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *LogGameDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
