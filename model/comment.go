package model

import (
	"encoding/json"
	"errors"
	"html/template"

	"github.com/ego008/youdb"
	"github.com/missdeer/kani/util"
)

type Comment struct {
	ID       uint64 `json:"id"`
	AID      uint64 `json:"aid"`
	UID      uint64 `json:"uid"`
	Content  string `json:"content"`
	ClientIP string `json:"clientip"`
	AddTime  uint64 `json:"addtime"`
}

type CommentListItem struct {
	ID         uint64 `json:"id"`
	AID        uint64 `json:"aid"`
	UID        uint64 `json:"uid"`
	Name       string `json:"name"`
	Avatar     string `json:"avatar"`
	Content    string `json:"content"`
	ContentFmt template.HTML
	AddTime    uint64 `json:"addtime"`
	AddTimeFmt string `json:"addtimefmt"`
}

type CommentPageInfo struct {
	Items    []CommentListItem `json:"items"`
	HasPrev  bool              `json:"hasprev"`
	HasNext  bool              `json:"hasnext"`
	FirstKey uint64            `json:"firstkey"`
	LastKey  uint64            `json:"lastkey"`
}

func CommentGetByKey(db *youdb.DB, aid string, cid uint64) (Comment, error) {
	obj := Comment{}
	rs := db.Hget("article_comment:"+aid, youdb.I2b(cid))
	if rs.State != "ok" {
		return obj, errors.New(rs.State)
	}
	if err := json.Unmarshal(rs.Data[0], &obj); err != nil {
		return obj, err
	}
	return obj, nil
}

func CommentSetByKey(db *youdb.DB, aid string, cid uint64, obj Comment) error {
	jb, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return db.Hset("article_comment:"+aid, youdb.I2b(cid), jb)
}

func CommentDelByKey(db *youdb.DB, aid string, cid uint64) error {
	return db.Hdel("article_comment:"+aid, youdb.I2b(cid))
}

func CommentList(db *youdb.DB, cmd, tb, key string, limit, tz int) CommentPageInfo {
	var items []CommentListItem
	var citems []Comment
	userMap := map[uint64]UserMini{}
	var userKeys [][]byte
	var hasPrev, hasNext bool
	var firstKey, lastKey uint64

	keyStart := youdb.DS2b(key)
	if cmd == "hrscan" {
		rs := db.Hrscan(tb, keyStart, limit)
		if rs.State == "ok" {
			for i := len(rs.Data) - 2; i >= 0; i -= 2 {
				item := Comment{}
				json.Unmarshal(rs.Data[i+1], &item)
				citems = append(citems, item)
				userMap[item.UID] = UserMini{}
				userKeys = append(userKeys, youdb.I2b(item.UID))
			}
		}
	} else if cmd == "hscan" {
		rs := db.Hscan(tb, keyStart, limit)
		if rs.State == "ok" {
			for i := 0; i < (len(rs.Data) - 1); i += 2 {
				item := Comment{}
				json.Unmarshal(rs.Data[i+1], &item)
				citems = append(citems, item)
				userMap[item.UID] = UserMini{}
				userKeys = append(userKeys, youdb.I2b(item.UID))
			}
		}
	}

	if len(citems) > 0 {
		rs := db.Hmget("user", userKeys)
		if rs.State == "ok" {
			for i := 0; i < (len(rs.Data) - 1); i += 2 {
				item := UserMini{}
				json.Unmarshal(rs.Data[i+1], &item)
				userMap[item.ID] = item
			}
		}

		for _, citem := range citems {
			user := userMap[citem.UID]
			item := CommentListItem{
				ID:         citem.ID,
				AID:        citem.AID,
				UID:        citem.UID,
				Name:       user.Name,
				Avatar:     user.Avatar,
				AddTime:    citem.AddTime,
				AddTimeFmt: util.TimeFmt(citem.AddTime, "2006-01-02 15:04", tz),
				ContentFmt: template.HTML(util.ContentFmt(db, citem.Content)),
			}
			items = append(items, item)
			if firstKey == 0 {
				firstKey = item.ID
			}
			lastKey = item.ID
		}

		rs = db.Hrscan(tb, youdb.I2b(firstKey), 1)
		if rs.State == "ok" {
			hasPrev = true
		}
		rs = db.Hscan(tb, youdb.I2b(lastKey), 1)
		if rs.State == "ok" {
			hasNext = true
		}
	}

	return CommentPageInfo{
		Items:    items,
		HasPrev:  hasPrev,
		HasNext:  hasNext,
		FirstKey: firstKey,
		LastKey:  lastKey,
	}
}
