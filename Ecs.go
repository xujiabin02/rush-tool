package rushtool

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
	"strings"
)

type AccessAliYun struct {
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

func ListInstance(ac AccessAliYun, instance ecs.Region, inputChan chan map[string]ecs.InstancesInDescribeInstances) {
	client, _ := ecs.NewClientWithAccessKey(instance.RegionId, ac.AccessKeyId, ac.AccessKeySecret)
	request := ecs.CreateDescribeInstancesRequest()
	response, err := client.DescribeInstances(request)
	CheckErr(err)
	count := response.TotalCount
	if count > 0 {
		if count > 10 {
			//fmt.Println(instance.RegionId, count)
			client, _ := ecs.NewClientWithAccessKey(instance.RegionId, ac.AccessKeyId, ac.AccessKeySecret)
			request := ecs.CreateDescribeInstancesRequest()
			request.PageSize = requests.Integer(strconv.Itoa(50))
			response, err := client.DescribeInstances(request)
			CheckErr(err)
			inputChan <- map[string]ecs.InstancesInDescribeInstances{instance.RegionId: response.Instances}
			//ListI(response.Instances)

		} else {
			//fmt.Println(instance.RegionId, count)
			inputChan <- map[string]ecs.InstancesInDescribeInstances{instance.RegionId: response.Instances}
			//ListI(response.Instances)
		}
	} else {
		inputChan <- map[string]ecs.InstancesInDescribeInstances{}
	}

}
func (srcAc *AccessAliYun) ListEcs() map[string]ecs.InstancesInDescribeInstances {
	ecsChan := make(chan map[string]ecs.InstancesInDescribeInstances)
	client, _ := ecs.NewClientWithAccessKey(srcAc.Region, srcAc.AccessKeyId, srcAc.AccessKeySecret)
	regions := ecs.CreateDescribeRegionsRequest()
	responseRegions, err := client.DescribeRegions(regions)
	CheckErr(err)
	var regionLst = ecs.Regions{}
	for _, i := range responseRegions.Regions.Region {
		if strings.HasPrefix(i.RegionId, "cn-") {
			regionLst.Region = append(regionLst.Region, i)
		}
	}
	for _, i := range regionLst.Region {
		go ListInstance(*srcAc, i, ecsChan)

	}
	result := map[string]ecs.InstancesInDescribeInstances{}
	for n := 0; n < len(regionLst.Region); n++ {
		value := <-ecsChan
		for k, v := range value {
			result[k] = v
		}
	}
	return result
}
func (srcAc *AccessAliYun) LstM(result map[string]ecs.InstancesInDescribeInstances, keyMap map[string]string) {
	for _, v := range result {
		if len(v.Instance) > 0 {
			//fmt.Println(k, len(v.Instance))
			for _, i := range v.Instance {
				if i.OSType == "linux" {
					if len(keyMap) > 0 {
						fmt.Println(
							i.Tags.Tag[0].TagValue,
							i.Tags.Tag[0].TagKey,
							i.VpcAttributes.PrivateIpAddress.IpAddress[0],
							i.InstanceName,
							keyMap[i.KeyPairName])
					} else {
						fmt.Println(
							i.Tags.Tag[0].TagValue,
							i.Tags.Tag[0].TagKey,
							i.VpcAttributes.PrivateIpAddress.IpAddress[0],
							i.InstanceName,
							i.KeyPairName)
					}

				}
			}
		}
	}
}
func (srcAc *AccessAliYun) LstHosts(result map[string]ecs.InstancesInDescribeInstances) {
	for _, v := range result {
		if len(v.Instance) > 0 {
			for _, i := range v.Instance {
				fmt.Println(
					i.VpcAttributes.PrivateIpAddress.IpAddress[0],
					i.InstanceName,
				)
			}

		}
	}
}

func (srcAc *AccessAliYun) GenerateAnsibleInventory(result map[string]ecs.InstancesInDescribeInstances, keyMap map[string]string) {
	c := make(map[string]Children, 0)
	a := AnsibleInventory{All: AnsibleInventoryAll{Children: c}}
	for _, v := range result {
		if len(v.Instance) > 0 {
			for _, i := range v.Instance {
				if i.OSType == "linux" {

					tmpChildren := Children{Hosts: make(map[string]string, 0)}
					//tmpChildren.Hosts[i.VpcAttributes.PrivateIpAddress.IpAddress[0]] = nil
					tmpChildren.Hosts[i.InstanceName] = ""
					tmpChildren.Vars.AnsibleSSHUser = "root"
					tmpChildren.Vars.AnsibleSSHPrivateKeyFile = keyMap[i.KeyPairName]
					groupName := fmt.Sprintf("%s%s", i.Tags.Tag[0].TagKey, i.Tags.Tag[0].TagValue)
					fmt.Println(tmpChildren)
					fmt.Println(groupName)
					if _, ok := a.All.Children[groupName]; ok {
						//a.All.Children[groupName].Hosts[i.VpcAttributes.PrivateIpAddress.IpAddress[0]] = nil
						a.All.Children[groupName].Hosts[i.InstanceName] = ""
					} else {
						a.All.Children[groupName] = tmpChildren
					}
				}
			}
		}
	}
	b := PrintInventoryYaml(a)
	err := ioutil.WriteFile("/w/h", b, 0777)
	CheckErr(err)
}
func PrintInventoryYaml(f AnsibleInventory) []byte {
	b, err := yaml.Marshal(&f)
	if err != nil {
		fmt.Println(err)
	}
	return b
}
