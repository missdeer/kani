package model

import (
	"encoding/json"
	"errors"

	"github.com/ego008/youdb"
)

type User struct {
	ID                uint64 `json:"id"`
	Name              string `json:"name"`
	Gender            string `json:"gender"`
	Flag              int    `json:"flag"`
	Avatar            string `json:"avatar"`
	Password          string `json:"password"`
	Email             string `json:"email"`
	URL               string `json:"url"`
	Telephone         string `json:"telephone"`
	Coin              uint64 `json:"coin"`
	Articles          uint64 `json:"articles"`
	Replies           uint64 `json:"replies"`
	RegTime           uint64 `json:"regtime"`
	LastPostTime      uint64 `json:"lastposttime"`
	LastReplyTime     uint64 `json:"lastreplytime"`
	LastLoginTime     uint64 `json:"lastlogintime"`
	About             string `json:"about"`
	Notice            string `json:"notice"`
	NoticeNum         int    `json:"noticenum"`
	EmailVerified     bool   `json:"emailverified"`
	TelephoneVerified bool   `json:"telephoneverified"`
	Hidden            bool   `json:"hidden"`
	Session           string `json:"session"`
}

type UserMini struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type UserPageInfo struct {
	Items    []User `json:"items"`
	HasPrev  bool   `json:"hasprev"`
	HasNext  bool   `json:"hasnext"`
	FirstKey uint64 `json:"firstkey"`
	LastKey  uint64 `json:"lastkey"`
}

func UserGetByID(db *youdb.DB, uid uint64) (User, error) {
	obj := User{}
	rs := db.Hget("user", youdb.I2b(uid))
	if rs.State == "ok" {
		if err := json.Unmarshal(rs.Data[0], &obj); err != nil {
			return obj, err
		}
		return obj, nil
	}
	return obj, errors.New(rs.State)
}

func UserUpdate(db *youdb.DB, obj User) error {
	jb, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return db.Hset("user", youdb.I2b(obj.ID), jb)
}

func UserGetByName(db *youdb.DB, name string) (User, error) {
	obj := User{}
	rs := db.Hget("user_name2uid", []byte(name))
	if rs.State != "ok" {
		return obj, errors.New(rs.State)
	}
	rs2 := db.Hget("user", rs.Data[0])
	if rs2.State != "ok" {
		return obj, errors.New(rs2.State)
	}
	if err := json.Unmarshal(rs2.Data[0], &obj); err != nil {
		return obj, err
	}
	return obj, nil
}

func UserGetIDByName(db *youdb.DB, name string) string {
	rs := db.Hget("user_name2uid", []byte(name))
	if rs.State == "ok" {
		return youdb.B2ds(rs.Data[0])
	}
	return ""
}

func UserListByFlag(db *youdb.DB, cmd, tb, key string, limit int) UserPageInfo {
	var items []User
	var keys [][]byte
	var hasPrev, hasNext bool
	var firstKey, lastKey uint64

	keyStart := youdb.DS2b(key)
	if cmd == "hrscan" {
		rs := db.Hrscan(tb, keyStart, limit)
		if rs.State == "ok" {
			for i := 0; i < (len(rs.Data) - 1); i += 2 {
				keys = append(keys, rs.Data[i])
			}
		}
	} else if cmd == "hscan" {
		rs := db.Hscan(tb, keyStart, limit)
		if rs.State == "ok" {
			for i := len(rs.Data) - 2; i >= 0; i -= 2 {
				keys = append(keys, rs.Data[i])
			}
		}
	}

	if len(keys) > 0 {
		rs := db.Hmget("user", keys)
		if rs.State == "ok" {
			for i := 0; i < (len(rs.Data) - 1); i += 2 {
				item := User{}
				json.Unmarshal(rs.Data[i+1], &item)
				items = append(items, item)
				if firstKey == 0 {
					firstKey = item.ID
				}
				lastKey = item.ID
			}

			rs = db.Hscan(tb, youdb.I2b(firstKey), 1)
			if rs.State == "ok" {
				hasPrev = true
			}
			rs = db.Hrscan(tb, youdb.I2b(lastKey), 1)
			if rs.State == "ok" {
				hasNext = true
			}
		}
	}

	return UserPageInfo{
		Items:    items,
		HasPrev:  hasPrev,
		HasNext:  hasNext,
		FirstKey: firstKey,
		LastKey:  lastKey,
	}
}
