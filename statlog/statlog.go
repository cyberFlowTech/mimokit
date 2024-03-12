package statlog

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
)

type StatLog struct {
	db    sqlx.SqlConn
	redis redis.UniversalClient
}

type Action string

const (
	UserLogin            Action = "UserLogin"
	UserRegister         Action = "UserRegister"
	UserCardVisited      Action = "UserCardVisited"
	UserCardScan         Action = "UserCardScan"
	UserCardH5Visited    Action = "UserCardH5Visited"
	UserCardH5Download   Action = "UserCardH5Download"
	UserCardH5Active     Action = "UserCardH5Active"
	RedPacketCreate      Action = "RedPacketCreate"
	RedPacketGet         Action = "RedPacketGet"
	ClubVisited          Action = "ClubVisited"
	ClubCreated          Action = "ClubCreated"
	ClubShared           Action = "ClubShared"
	ClubH5Visited        Action = "ClubH5Visited"
	ClubH5Download       Action = "ClubH5Download"
	ClubH5Active         Action = "ClubH5Active"
	ClubDissolution      Action = "ClubDissolution"
	ChannelDelete        Action = "ChannelDelete"
	ChannelCreate        Action = "ChannelCreate"
	ClubQuit             Action = "ClubQuit"
	ClubRoleCreate       Action = "ClubRoleCreate"
	ClubBlockUser        Action = "ClubBlockUser"
	ClubCompletedProfile Action = "ClubCompletedProfile"
	DynamicVisited       Action = "DynamicVisited"
	DynamicPublish       Action = "DynamicPublish"
	DynamicPraised       Action = "DynamicPraised"
	DynamicCommented     Action = "DynamicCommented"
	DynamicShared        Action = "DynamicShared"
)

type Platform string

const (
	Platform_IOS     Platform = "10"
	Platform_Andriod Platform = "20"
	Platform_Unknown Platform = "99"
)

func New(db sqlx.SqlConn, redis redis.UniversalClient) StatLog {
	return StatLog{db: db, redis: redis}
}

func getTableName() string {
	return "db_mime_log.t_stat_log"
}

const defaultJson = "{}"

type StatLogEntity struct {
	Id           int
	PlatformType Platform
	UserId       int64
	ActionType   Action
	Ext          string
}

func (s StatLog) insertToDB(d StatLogEntity) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (`platform`,`user_id`,`action`,`ext`,`ctime`) values (?, ?, ?, ?, ?)", getTableName())
	res, err := s.db.Exec(query, d.PlatformType, d.UserId, d.ActionType, d.Ext, time.Now().Unix())
	if err != nil {
		logx.Errorf("stat log insertToDB,d:%+v err:%v", d, err)
	}
	return res, err
}

func (s StatLog) UserLoginEvent(platform Platform, userId int64) (sql.Result, error) {
	data := StatLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   UserLogin,
		Ext:          defaultJson,
	}
	res, err := s.insertToDB(data)
	s.redis.SAdd(context.Background(), "se:zapry:statlog:loginUsers:"+time.Now().Format("20060102"), userId)
	return res, err
}

func (s StatLog) UserRegisterEvent(platform Platform, userId int64) (sql.Result, error) {
	data := StatLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   UserRegister,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) UserCardVisitedEvent(platform Platform, userId int64) (sql.Result, error) {
	data := StatLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   UserCardVisited,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) UserCardScanEvent(platform Platform, userId int64) (sql.Result, error) {
	data := StatLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   UserCardScan,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) UserCardH5VisitedEvent(platform Platform, userId int64) (sql.Result, error) {
	data := StatLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   UserCardH5Visited,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) UserCardH5DownloadEvent(platform Platform, userId int64) (sql.Result, error) {
	data := StatLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   UserCardH5Download,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) UserCardH5ActiveEvent(platform Platform, userId int64) (sql.Result, error) {
	data := StatLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   UserCardH5Active,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) RedPacketCreateEvent(platform Platform, userId, quantity int64) (sql.Result, error) {
	m := map[string]int64{}
	m["quantity"] = quantity
	marshal, _ := json.Marshal(m)
	data := StatLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   RedPacketCreate,
		Ext:          string(marshal),
	}
	return s.insertToDB(data)
}

func (s StatLog) RedPacketGetEvent(platform Platform, userId int64) (sql.Result, error) {
	data := StatLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   RedPacketGet,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) ClubVisitedEvent(platform Platform, userId, tbMid int64) (sql.Result, error) {
	m := map[string]int64{}
	m["tb_mid"] = tbMid
	marshal, _ := json.Marshal(m)
	data := StatLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   ClubVisited,
		Ext:          string(marshal),
	}
	return s.insertToDB(data)
}

func (s StatLog) ClubCreatedEvent(platform Platform, userId int64) (sql.Result, error) {
	data := StatLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   ClubCreated,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) ClubSharedEvent(platform Platform, userId int64) (sql.Result, error) {
	data := StatLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   ClubShared,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) ClubH5VisitedEvent(platform Platform, userId int64) (sql.Result, error) {
	data := StatLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   ClubH5Visited,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) ClubH5DownloadEvent(platform Platform, userId int64) (sql.Result, error) {
	data := StatLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   ClubH5Download,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) ClubH5ActiveEvent(platform Platform, userId int64) (sql.Result, error) {
	data := StatLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   ClubH5Active,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) ClubDissolutionEvent(platform Platform, userId int64) (sql.Result, error) {
	data := StatLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   ClubDissolution,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) ChannelDeleteEvent(platform Platform, userId int64) (sql.Result, error) {
	data := StatLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   ChannelDelete,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) ChannelCreateEvent(platform Platform, userId int64) (sql.Result, error) {
	data := StatLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   ChannelCreate,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) ClubQuitEvent(platform Platform, userId int64) (sql.Result, error) {
	data := StatLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   ClubQuit,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) ClubRoleCreateEvent(platform Platform, userId int64) (sql.Result, error) {
	data := StatLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   ClubRoleCreate,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) ClubBlockUserEvent(platform Platform, userId int64) (sql.Result, error) {
	data := StatLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   ClubBlockUser,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) ClubCompletedProfileEvent(platform Platform, userId int64) (sql.Result, error) {
	data := StatLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   ClubCompletedProfile,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) DynamicVisitedEvent(platform Platform, userId int64) (sql.Result, error) {
	data := StatLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   DynamicVisited,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) DynamicPublishEvent(platform Platform, userId int64) (sql.Result, error) {
	data := StatLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   DynamicPublish,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) DynamicPraisedEvent(platform Platform, userId int64) (sql.Result, error) {
	data := StatLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   DynamicPraised,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) DynamicCommentedEvent(platform Platform, userId int64) (sql.Result, error) {
	data := StatLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   DynamicCommented,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) DynamicSharedEvent(platform Platform, userId int64) (sql.Result, error) {
	data := StatLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   DynamicShared,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}
