/**
 * This package was copied directly from https://github.com/jubnzv/go-taskwarrior/blob/master/task.go and modified to suit my needs
 **/
package taskwarrior

type Task struct {
	Id          int     `json:"id"`
	Description string  `json:"description"`
	Project     string  `json:"project"`
	Status      string  `json:"status"`
	Uuid        string  `json:"uuid"`
	Urgency     float32 `json:"urgency"`
	Priority    string  `json:"priority"`
	Due         string  `json:"due"`
	End         string  `json:"end"`
	Entry       string  `json:"entry"`
	Modified    string  `json:"modified"`
}
