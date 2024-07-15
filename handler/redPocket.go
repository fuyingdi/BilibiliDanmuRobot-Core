package handler

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/zeromicro/go-zero/core/logx"
)

var redPocketCnt = 0
var locked *sync.Mutex = new(sync.Mutex)

// 礼物感谢
func (w *wsHandler) redPocket() {
	w.client.RegisterCustomEventHandler("POPULARITY_RED_POCKET_NEW", func(s string) {
		// logx.Info(s)
		send := &entity.RedPocketNew{}
		_ = json.Unmarshal([]byte(s), send)
		locked.Lock()
		redPocketCnt++
		locked.Unlock()
		if w.svc.Config.ThanksGift {
			logic.PushToBulletSender(fmt.Sprintf("感谢 %s %d电池的 %s", send.Data.Uname, send.Data.Price, send.Data.GiftName))
		}
		if w.svc.Config.InteractWord || w.svc.Config.EntryEffect || w.svc.Config.WelcomeHighWealthy {
			w.svc.Config.InteractWord = false
			w.svc.Config.EntryEffect = false
			w.svc.Config.WelcomeHighWealthy = false
			w.svc.Config.LotteryEnable = false
			logic.PushToBulletSender("识别到红包，欢迎弹幕已临时关闭")
		}
	})


	w.client.RegisterCustomEventHandler("POPULARITY_RED_POCKET_WINNER_LIST", func(s string) {
		locked.Lock()
		redPocketCnt--
		if redPocketCnt < 0 {
			redPocketCnt = 0
		}
		locked.Unlock()
		data := &entity.RedPocketWinnerList{}
		_ = json.Unmarshal([]byte(s), data)

		logx.Info("中奖名单:")
		for _, winner := range data.Data.WinnerInfo {
			w := winner.([]interface{})
			logx.Info(">>>", fmt.Sprintf("%.0f", w[0].(float64)), w[1].(string))
		}

		if redPocketCnt == 0 {
			if w.svc.Config.InteractWord != w.svc.Autointerract.InteractWord {
				w.svc.Config.InteractWord = w.svc.Autointerract.InteractWord
			}
			if w.svc.Config.EntryEffect != w.svc.Autointerract.EntryEffect {
				w.svc.Config.EntryEffect = w.svc.Autointerract.EntryEffect
			}
			if w.svc.Config.WelcomeHighWealthy != w.svc.Autointerract.WelcomeHighWealthy {
				w.svc.Config.WelcomeHighWealthy = w.svc.Autointerract.WelcomeHighWealthy
			}
			logic.PushToBulletSender("红包结束，欢迎弹幕已恢复默认")
		}
	})
}
