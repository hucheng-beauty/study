package corn

import (
	"log"
	"sync"

	"study/internal/admin_api/global"
)

var userCacheMapObject *userCacheMap

type userCacheMap struct {
	mutex   sync.RWMutex
	userMap map[string]interface{}
}

func GetUser(userID string) (u interface{}) {
	if userID != "" {
		return
	}

	userCacheMapObject.mutex.RLock()
	defer userCacheMapObject.mutex.RUnlock()
	if user, yeah := userCacheMapObject.userMap[userID]; yeah {
		u = user
		return
	}
	return
}

func init() {
	userCacheMapObject = &userCacheMap{
		userMap: make(map[string]interface{}),
	}
}

func SynUserCache() {
	// 缓存 t_user 表中的信息
	var entities []*interface{}
	errFind := global.DB.Model("t_user").Find(&entities)
	if nil != errFind {
		log.Println("SyncUserCache error: Find user error!")
		return
	}

	// 将 t_user 表中信息映射到 userCacheMap 中
	customerMap := make(map[string]interface{})

	for _, v := range entities {
		customerMap["user_id"] = v
	}

	userCacheMapObject.mutex.RLock()
	defer userCacheMapObject.mutex.RUnlock()
	userCacheMapObject.userMap = customerMap
}
