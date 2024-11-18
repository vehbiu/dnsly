package models

type BlockedRequest struct {
	Domain string `json:"domain"`
	IP     string `json:"ip"`
	Source string `json:"source"`
}
