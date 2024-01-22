package pfsense_client

import "fmt"

type HostOverride struct {
	Host   string   `json:"host"`
	Domain string   `json:"domain"`
	IP     []string `json:"ip"`
	Tag    string   `json:"descr"`
}

type HostOverrides []HostOverride

func (hos HostOverrides) Len() int {
	return len(hos)
}

func (hos HostOverrides) Swap(i, j int) {
	hos[i], hos[j] = hos[j], hos[i]
}

func (hos HostOverrides) Less(i, j int) bool {
	if hos[i].Domain == hos[j].Domain {
		return hos[i].Host < hos[j].Host
	}
	return hos[i].Domain < hos[j].Domain
}

func (ho *HostOverride) Name() string {
	return fmt.Sprintf("%s.%s", ho.Host, ho.Domain)
}
