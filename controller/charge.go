package controller

import (
	"net/http"
	"time"

	"github.com/missdeer/kani/model"
	"github.com/rs/xid"
)

func (h *BaseHandler) UserCharge(w http.ResponseWriter, r *http.Request) {
	currentUser, _ := h.CurrentUser(w, r)
	if currentUser.ID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	type pageData struct {
		PageData
		Uobj model.User
		Now  int64
	}

	tpl := h.CurrentTpl(r)
	evn := &pageData{}
	evn.SiteCf = h.App.Cf.Site
	evn.Title = "设置"
	evn.Keywords = ""
	evn.Description = ""
	evn.IsMobile = tpl == "mobile"
	evn.CurrentUser = currentUser

	evn.ShowSideAd = true
	evn.PageName = "user_setting"

	evn.Uobj = currentUser
	evn.Now = time.Now().UTC().Unix()
	h.SetCookie(w, "token", xid.New().String(), 1)
	h.Render(w, tpl, evn, "layout.html", "usersetting.html")
}
