package rabbit

//mq queue
const QUEUE_MASTER_EXP = "/yi/queue_master_exp"

type MasterExpMq struct {
	//大师id(uid)
	MasterID int64 `json:"master_id"`
	// 评价结果 1：好 2 ：中 3：差
	ExpResult int `json:"exp_result"`
}
