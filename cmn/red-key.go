package cmn

import "fmt"

type RedisKey struct {
}

var RedisKeyX = new(RedisKey)

func (r *RedisKey) PushRegID(uid int64) string {
	pat := "/yi/mod-push-reg-id/%d"
	return fmt.Sprintf(pat, uid)
}
