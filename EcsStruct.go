package rushtool

import "github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

type AccessAliYun struct {
	AccessKeyId     string `yaml:"accessKeyId"`
	AccessKeySecret string `yaml:"accessKeySecret"`
	Region          string `yaml:"region"`
	Client          ecs.Client
}
type AccessAliYunStruct struct {
	AccessKeyId     string `yaml:"accessKeyId"`
	AccessKeySecret string `yaml:"accessKeySecret"`
	Region          string `yaml:"region"`
}
type AnsibleInventoryAll struct {
	Hosts    []string            `yaml:"hosts,omitempty"`
	Vars     AnsibleVars         `yaml:"vars,omitempty"`
	Children map[string]Children `yaml:"children"`
}
type Children struct {
	Hosts map[string]string `yaml:"hosts"`
	Vars  AnsibleVars       `yaml:"vars"`
}
type AnsibleInventory struct {
	All AnsibleInventoryAll `yaml:"all"`
}
type AnsibleVars struct {
	AnsibleSSHHost           string `yaml:"ansible_ssh_host,omitempty"`
	AnsibleSSHPass           string `yaml:"ansible_ssh_pass,omitempty"`
	AnsibleSSHUser           string `yaml:"ansible_ssh_user"`
	AnsibleSSHPrivateKeyFile string `yaml:"ansible_ssh_private_key_file"`
}
type AttrSessionSecureCRT struct {
	Tpl1         string `yaml:"tpl1"`
	Tpl2         string `yaml:"tpl2"`
	ConfigPath   string `yaml:"configPath"`
	PersonalPath string `yaml:"personalPath"`
}
type ProcessAttr struct {
	Protocol string
	Port string
	Name string
}
