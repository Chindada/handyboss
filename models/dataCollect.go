package models

import (
	"strconv"
	"time"
)

// Wise Wise
type Wise struct {
	ID              int64
	IdleTime        int64
	MacAddress      string
	Token           string
	PutTimeInterval int64
	IP              string
}

// DcPatchBody DcPatchBody
type DcPatchBody struct {
	TSt  int64
	TEnd int64
}

// DcPutBody DcPutBody
type DcPutBody struct {
	UID  int64
	MAC  int64
	TmF  int64
	Fltr int64
}

// DcGetLogBody DcGetLogBody
type DcGetLogBody struct {
	UID   int64
	MAC   int64
	TmF   int64
	SysTk int64
	Fltr  int64
	TSt   int64
	TEnd  int64
	Amt   int64
	Total int64
	TLst  int64
	TFst  int64
}

// WiseData WiseData
type WiseData struct {
	LogMsg []struct {
		PE     int64     `json:"PE"`
		UID    string    `json:"UID"`
		TIM    string    `json:"TIM"`
		SysTk  int64     `json:"SysTk"`
		Record [][]int64 `json:"Record"`
	} `json:"LogMsg"`
}

// ConvertToDi ConvertToDi
func (c *WiseData) ConvertToDi() (di []Di, err error) {
	for _, v := range c.LogMsg {
		timeStamp, err := strconv.Atoi(v.TIM)
		if err != nil {
			return di, err
		}
		tmp := Di{
			Di0:        v.Record[0][3],
			Di1:        v.Record[1][3],
			Di2:        v.Record[2][3],
			Di3:        v.Record[3][3],
			Di4:        v.Record[4][3],
			Di5:        v.Record[5][3],
			Di6:        v.Record[6][3],
			Di7:        v.Record[7][3],
			MacAddress: v.UID[10:22],
			Timestamp:  int64(timeStamp),
			SysTk:      v.SysTk,
		}
		di = append(di, tmp)
	}
	return di, err
}

// Di Di
type Di struct {
	ID         int64  `json:"id" orm:"column(id)"`
	MacAddress string `json:"macAddress"`
	Di0        int64  `json:"di0"`
	Di1        int64  `json:"di1"`
	Di2        int64  `json:"di2"`
	Di3        int64  `json:"di3"`
	Di4        int64  `json:"di4"`
	Di5        int64  `json:"di5"`
	Di6        int64  `json:"di6"`
	Di7        int64  `json:"di7"`
	Analyzed   int64  `json:"analyzed"`
	Timestamp  int64  `json:"timestamp"`
	SysTk      int64  `json:"sysTk"`
	// CreateTime time.Time `json:"createTime" orm:"auto_now;type(datetime)"`
}

// DcStatus DcStatus
type DcStatus struct {
	ID         int64     `json:"id" orm:"column(id)"`
	MacAddress string    `json:"macAddress"`
	Status     int64     `json:"stauts"`
	Timestamp  int64     `json:"timestamp"`
	CycleTime  float64   `json:"cycleTime"`
	CreateTime time.Time `json:"createTime" orm:"auto_now;type(datetime)"`
}
