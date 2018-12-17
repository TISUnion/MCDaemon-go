package command

import (
	"sync"
)

var Group *PlayerGroup

type PlayerGroup struct {
	players map[string][]string
	lock    sync.Mutex
}

func init() {
	Group = &PlayerGroup{}
	Group.players = make(map[string][]string)
}

//向组内添加玩家
func (pg *PlayerGroup) AddPlayer(groupname string, player string) {
	pg.lock.Lock()
	defer pg.lock.Unlock()
	pg.players[groupname] = append(pg.players[groupname], player)
}

//删除组内指定玩家
func (pg *PlayerGroup) DelPlayer(groupname string, player string) {
	pg.lock.Lock()
	defer pg.lock.Unlock()
	//如果在组内，存在才删除
	if pg._hasPlayer(groupname, player, false) {
		var index int
		for k, val := range pg.players[groupname] {
			if val == player {
				index = k
				break
			}
		}
		//去除查询到的玩家
		pg.players[groupname] = append(pg.players[groupname][0:index], pg.players[groupname][index+1:]...)
	}
}

//查询组内是否有指定玩家
func (pg *PlayerGroup) HasPlayer(groupname string, player string) bool {
	return pg._hasPlayer(groupname, player, true)
}

func (pg *PlayerGroup) GetPlayer() map[string][]string {
	pg.lock.Lock()
	defer pg.lock.Unlock()
	return pg.players
}

func (pg *PlayerGroup) _hasPlayer(groupname string, player string, isLock bool) bool {
	if isLock {
		pg.lock.Lock()
		defer pg.lock.Unlock()
	}
	if len(pg.players) > 0 && pg.players[groupname] != nil {
		for _, val := range pg.players[groupname] {
			if val == player {
				return true
			}
		}
	}
	return false
}
