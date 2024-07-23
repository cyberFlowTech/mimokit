package statlog

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/cyberFlowTech/mimokit/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strconv"
	"time"
)

type StatLog struct {
	db sqlx.SqlConn
	//redis redis.UniversalClient
}

type Action string

const (
	// User Event
	UserLogin          Action = "UserLogin"
	UserRegister       Action = "UserRegister"
	UserCardVisited    Action = "UserCardVisited"
	UserCardScan       Action = "UserCardScan"
	UserCardH5Visited  Action = "UserCardH5Visited"
	UserCardH5Download Action = "UserCardH5Download"
	UserCardH5Active   Action = "UserCardH5Active"

	// RedPacket Event
	RedPacketCreate Action = "RedPacketCreate" //红包创建 ext: 红包分发份数 quantity:4
	RedPacketGet    Action = "RedPacketGet"    //获取红包

	// Club Event
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

	// Dynamic Event
	DynamicVisited   Action = "DynamicVisited"
	DynamicPublish   Action = "DynamicPublish"
	DynamicPraised   Action = "DynamicPraised"
	DynamicCommented Action = "DynamicCommented"
	DynamicShared    Action = "DynamicShared"
)

type Platform string

const (
	Platform_IOS     Platform = "10"
	Platform_Andriod Platform = "20"
	Platform_PC      Platform = "30"
	Platform_Unknown Platform = "99"
)

func GetPlatformTypeByAPi(api string) Platform {
	fromApi := utils.GetPlatformFromApi(api)
	platformType := Platform_Unknown
	switch fromApi {
	case "10":
		platformType = Platform_IOS
	case "20":
		platformType = Platform_Andriod
	case "30":
		platformType = Platform_PC
	case "99":
		platformType = Platform_Unknown

	}
	return platformType
}

func New(db sqlx.SqlConn) StatLog {
	return StatLog{db: db}
}

func getTableName() string {
	return "db_mime_log.t_stat_log"
}

const defaultJson = "{}"

type statLogEntity struct {
	Id           int
	PlatformType Platform
	UserId       int64
	ActionType   Action
	Ext          string
}

type LogEntity struct {
	PlatformType Platform
	UserId       int64
	ActionType   Action
	Ext          map[string]string
}

func (s StatLog) Log(d []LogEntity) error {
	for _, v := range d {
		switch v.ActionType {
		case UserLogin:
			_, _ = s.UserLoginEvent(v.PlatformType, v.UserId, v.Ext)
		case UserRegister:
			_, _ = s.UserRegisterEvent(v.PlatformType, v.UserId, v.Ext)
		case UserCardVisited:
			_, _ = s.UserCardVisitedEvent(v.PlatformType, v.UserId)
		case UserCardScan:
			_, _ = s.UserCardScanEvent(v.PlatformType, v.UserId)
		case UserCardH5Visited:
			_, _ = s.UserCardH5VisitedEvent(v.PlatformType, v.UserId)
		case UserCardH5Download:
			_, _ = s.UserCardH5DownloadEvent(v.PlatformType, v.UserId)
		case UserCardH5Active:
			_, _ = s.UserCardH5ActiveEvent(v.PlatformType, v.UserId)
		case RedPacketCreate:
			quantity, _ := strconv.Atoi(v.Ext["quantity"])
			_, _ = s.RedPacketCreateEvent(v.PlatformType, v.UserId, int64(quantity))
		case RedPacketGet:
			_, _ = s.RedPacketGetEvent(v.PlatformType, v.UserId)
		case ClubVisited:
			tbMid, _ := strconv.Atoi(v.Ext["tb_mid"])
			_, _ = s.ClubVisitedEvent(v.PlatformType, v.UserId, int64(tbMid))
		case ClubCreated:
			_, _ = s.ClubCreatedEvent(v.PlatformType, v.UserId)
		case ClubShared:
			_, _ = s.ClubSharedEvent(v.PlatformType, v.UserId)
		case ClubH5Visited:
			_, _ = s.ClubH5VisitedEvent(v.PlatformType, v.UserId)
		case ClubH5Download:
			_, _ = s.ClubH5DownloadEvent(v.PlatformType, v.UserId)
		case ClubH5Active:
			_, _ = s.ClubH5ActiveEvent(v.PlatformType, v.UserId)
		case ClubDissolution:
			_, _ = s.ClubDissolutionEvent(v.PlatformType, v.UserId)
		case ChannelDelete:
			_, _ = s.ChannelDeleteEvent(v.PlatformType, v.UserId)
		case ChannelCreate:
			_, _ = s.ChannelCreateEvent(v.PlatformType, v.UserId)
		case ClubQuit:
			_, _ = s.ClubQuitEvent(v.PlatformType, v.UserId)
		case ClubRoleCreate:
			_, _ = s.ClubRoleCreateEvent(v.PlatformType, v.UserId)
		case ClubBlockUser:
			_, _ = s.ClubBlockUserEvent(v.PlatformType, v.UserId)
		case ClubCompletedProfile:
			_, _ = s.ClubCompletedProfileEvent(v.PlatformType, v.UserId)
		case DynamicVisited:
			_, _ = s.DynamicVisitedEvent(v.PlatformType, v.UserId)
		case DynamicPublish:
			_, _ = s.DynamicPublishEvent(v.PlatformType, v.UserId)
		case DynamicPraised:
			_, _ = s.DynamicPraisedEvent(v.PlatformType, v.UserId)
		case DynamicCommented:
			_, _ = s.DynamicCommentedEvent(v.PlatformType, v.UserId)
		case DynamicShared:
			_, _ = s.DynamicSharedEvent(v.PlatformType, v.UserId)
		}
	}
	return nil
}

func (s StatLog) insertToDB(d statLogEntity) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (`platform`,`user_id`,`action`,`ext`,`ctime`) values (?, ?, ?, ?, ?)", getTableName())
	res, err := s.db.Exec(query, d.PlatformType, d.UserId, d.ActionType, d.Ext, time.Now().Unix())
	if err != nil {
		logx.Errorf("stat log insertToDB,d:%+v err:%v", d, err)
	}
	return res, err
}

func (s StatLog) UserLoginEvent(platform Platform, userId int64, extData map[string]string) (sql.Result, error) {
	ext := defaultJson
	if extData != nil {
		marshal, err := json.Marshal(extData)
		if err != nil {
			return nil, err
		}
		ext = string(marshal)
	}

	data := statLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   UserLogin,
		Ext:          ext,
	}
	return s.insertToDB(data)
}

func (s StatLog) UserRegisterEvent(platform Platform, userId int64, extData map[string]string) (sql.Result, error) {
	marshal, err := json.Marshal(extData)
	if err != nil {
		return nil, err
	}
	ext := string(marshal)
	data := statLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   UserRegister,
		Ext:          ext,
	}
	return s.insertToDB(data)
}

func (s StatLog) UserCardVisitedEvent(platform Platform, userId int64) (sql.Result, error) {
	data := statLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   UserCardVisited,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) UserCardScanEvent(platform Platform, userId int64) (sql.Result, error) {
	data := statLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   UserCardScan,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) UserCardH5VisitedEvent(platform Platform, userId int64) (sql.Result, error) {
	data := statLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   UserCardH5Visited,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) UserCardH5DownloadEvent(platform Platform, userId int64) (sql.Result, error) {
	data := statLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   UserCardH5Download,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) UserCardH5ActiveEvent(platform Platform, userId int64) (sql.Result, error) {
	data := statLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   UserCardH5Active,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

// RedPacketCreateEvent
// @param quantity 一个红包发的数量
func (s StatLog) RedPacketCreateEvent(platform Platform, userId, quantity int64) (sql.Result, error) {
	m := map[string]int64{}
	m["quantity"] = quantity
	marshal, _ := json.Marshal(m)
	data := statLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   RedPacketCreate,
		Ext:          string(marshal),
	}
	return s.insertToDB(data)
}

func (s StatLog) RedPacketGetEvent(platform Platform, userId int64) (sql.Result, error) {
	data := statLogEntity{
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
	data := statLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   ClubVisited,
		Ext:          string(marshal),
	}
	return s.insertToDB(data)
}

func (s StatLog) ClubCreatedEvent(platform Platform, userId int64) (sql.Result, error) {
	data := statLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   ClubCreated,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) ClubSharedEvent(platform Platform, userId int64) (sql.Result, error) {
	data := statLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   ClubShared,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) ClubH5VisitedEvent(platform Platform, userId int64) (sql.Result, error) {
	data := statLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   ClubH5Visited,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) ClubH5DownloadEvent(platform Platform, userId int64) (sql.Result, error) {
	data := statLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   ClubH5Download,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) ClubH5ActiveEvent(platform Platform, userId int64) (sql.Result, error) {
	data := statLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   ClubH5Active,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) ClubDissolutionEvent(platform Platform, userId int64) (sql.Result, error) {
	data := statLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   ClubDissolution,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) ChannelDeleteEvent(platform Platform, userId int64) (sql.Result, error) {
	data := statLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   ChannelDelete,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) ChannelCreateEvent(platform Platform, userId int64) (sql.Result, error) {
	data := statLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   ChannelCreate,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) ClubQuitEvent(platform Platform, userId int64) (sql.Result, error) {
	data := statLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   ClubQuit,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) ClubRoleCreateEvent(platform Platform, userId int64) (sql.Result, error) {
	data := statLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   ClubRoleCreate,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) ClubBlockUserEvent(platform Platform, userId int64) (sql.Result, error) {
	data := statLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   ClubBlockUser,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) ClubCompletedProfileEvent(platform Platform, userId int64) (sql.Result, error) {
	data := statLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   ClubCompletedProfile,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) DynamicVisitedEvent(platform Platform, userId int64) (sql.Result, error) {
	data := statLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   DynamicVisited,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) DynamicPublishEvent(platform Platform, userId int64) (sql.Result, error) {
	data := statLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   DynamicPublish,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) DynamicPraisedEvent(platform Platform, userId int64) (sql.Result, error) {
	data := statLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   DynamicPraised,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) DynamicCommentedEvent(platform Platform, userId int64) (sql.Result, error) {
	data := statLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   DynamicCommented,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}

func (s StatLog) DynamicSharedEvent(platform Platform, userId int64) (sql.Result, error) {
	data := statLogEntity{
		PlatformType: platform,
		UserId:       userId,
		ActionType:   DynamicShared,
		Ext:          defaultJson,
	}
	return s.insertToDB(data)
}
