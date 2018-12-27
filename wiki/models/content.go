package models

import (
	"time"
)

type Content struct {
	Platform            string    `xorm:"VARCHAR(255)"`
	Channelname         string    `xorm:"VARCHAR(255)"`
	Channel             string    `xorm:"VARCHAR(255)"`
	Label               string    `xorm:"VARCHAR(255)"`
	Describe            string    `xorm:"VARCHAR(2048)"`
	Downloadurl         string    `xorm:"unique VARCHAR(1024)"`
	Category            string    `xorm:"VARCHAR(255)"`
	Size                int       `xorm:"INT(11)"`
	Addtime             time.Time `xorm:"DATETIME"`
	Id                  int       `xorm:"not null pk autoincr INT(11)"`
	Poster              string    `xorm:"VARCHAR(255)"`
	Extra               string    `xorm:"VARCHAR(255)"`
	IsMainland          int       `xorm:"not null comment('大陆是否可以下载') TINYINT(1)"`
	IssuedStatus        string    `xorm:"not null default 'unissued' comment('下发状态') index VARCHAR(255)"`
	LastIssuedTime      time.Time `xorm:"DATETIME"`
	MineStatus          string    `xorm:"not null default 'unminer' comment('挖矿状态') index VARCHAR(255)"`
	IssuedMachineNumber int       `xorm:"not null default 00000000000 comment('下发矿机个数') INT(11)"`
	MineMachineNumber   int       `xorm:"not null default 00000000000 comment('参与挖矿的矿机个数，这个值只是统计上报挖矿的主机个数，不去重') INT(11)"`
}
