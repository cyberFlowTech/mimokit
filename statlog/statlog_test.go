package statlog

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"testing"
)

func TestStatLog_UserLoginEvent(t *testing.T) {
	type fields struct {
		db sqlx.SqlConn
	}
	type args struct {
		platform Platform
		userId   int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		//want    sql.Result
		wantErr bool
	}{
		{
			name: "xx",
			fields: fields{
				db: sqlx.NewMysql("root:root@tcp(192.168.31.57:3306)/db_mime_community?charset=utf8mb4&parseTime=true"),
			},
			args: args{
				platform: Platform_IOS,
				userId:   111,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := StatLog{
				db: tt.fields.db,
			}
			_, err := s.UserLoginEvent(tt.args.platform, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserLoginEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("UserLoginEvent() got = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestStatLog_RedPackageCreateEvent(t *testing.T) {
	type fields struct {
		db sqlx.SqlConn
	}
	type args struct {
		platform Platform
		userId   int64
		quantity int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		//want    sql.Result
		wantErr bool
	}{
		{
			name: "创建红包",
			fields: fields{
				db: sqlx.NewMysql("root:root@tcp(192.168.31.57:3306)/db_mime_community?charset=utf8mb4&parseTime=true"),
			},
			args: args{
				platform: Platform_Andriod,
				userId:   111,
				quantity: 20,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//s := StatLog{
			//	db:    tt.fields.db,
			//}
			data := LogEntity{
				PlatformType: tt.args.platform,
				UserId:       tt.args.userId,
				ActionType:   RedPacketCreate,
				Ext: map[string]string{
					"quantity": "20", //创建红包的数量
				},
			}
			New(tt.fields.db).Log([]LogEntity{data})

			//_, err := s.RedPacketCreateEvent(tt.args.platform, tt.args.userId, tt.args.quantity)
			//if (err != nil) != tt.wantErr {
			//	t.Errorf("RedPacketCreateEvent() error = %v, wantErr %v", err, tt.wantErr)
			//	return
			//}

			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("RedPacketCreateEvent() got = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestStatLog_RedPackageCreateExample(t *testing.T) {
	type fields struct {
		db sqlx.SqlConn
	}
	type args struct {
		platform Platform
		userId   int64
		quantity int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		//want    sql.Result
		wantErr bool
	}{
		{
			name: "创建红包",
			fields: fields{
				db: sqlx.NewMysql("root:root@tcp(192.168.31.57:3306)/db_mime_community?charset=utf8mb4&parseTime=true"),
			},
			args: args{
				platform: Platform_Andriod,
				userId:   111,
				quantity: 20,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := StatLog{
				db: tt.fields.db,
			}
			_, err := s.RedPacketCreateEvent(tt.args.platform, tt.args.userId, tt.args.quantity)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedPacketCreateEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("RedPacketCreateEvent() got = %v, want %v", got, tt.want)
			//}
		})
	}
}
