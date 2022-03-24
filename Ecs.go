package rushtool

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/flosch/pongo2/v5"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

func (srcAc *AccessAliYun) Init() {
	client, _ := ecs.NewClientWithAccessKey(srcAc.Region, srcAc.AccessKeyId, srcAc.AccessKeySecret)
	srcAc.Client = *client
}

func (srcAc *AccessAliYun) DescSecurityGroup() {
	response, err := srcAc.Client.DescribeSecurityGroups(ecs.CreateDescribeSecurityGroupsRequest())
	CheckErr(err)
	for _, i := range response.SecurityGroups.SecurityGroup {
		fmt.Println(i.SecurityGroupId, i.Description)
		request := ecs.CreateDescribeSecurityGroupAttributeRequest()
		request.SecurityGroupId = i.SecurityGroupId
		res, err := srcAc.Client.DescribeSecurityGroupAttribute(request)
		CheckErr(err)
		for _, rule := range res.Permissions.Permission {
			fmt.Println(rule)
		}
	}

}
func (srcAc *AccessAliYun) DescLoadBalancer() {
	client, err := slb.NewClientWithAccessKey(srcAc.Region, srcAc.AccessKeyId, srcAc.AccessKeySecret)
	CheckErr(err)
	req := slb.CreateDescribeLoadBalancerAttributeRequest()
	res, err := client.DescribeLoadBalancerAttribute(req)
	for _, slb := range res.LoadBalancerName {
		fmt.Println(slb)
	}
}
func (srcAc *AccessAliYun) ModifySecurityGroup(secureId string, rules []ProcessAttr) {
	for _, p := range rules {
		requestsAli := ecs.ModifySecurityGroupEgressRuleRequest{}
		requestsAli.RegionId = srcAc.Region
		requestsAli.SourceCidrIp = "0.0.0.0/0"
		requestsAli.SecurityGroupId = secureId
		requestsAli.PortRange = fmt.Sprintf("%s/%s", p.Port, p.Port)
		requestsAli.IpProtocol = p.Protocol
		requestsAli.Description = p.Name
		requestsAli.Policy = "accept"
		res, err := srcAc.Client.ModifySecurityGroupEgressRule(&requestsAli)
		CheckErr(err)
		fmt.Println(res)
	}
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
func GetHosts(name string, ip string) error {
	contentHosts := ReadFile("/etc/hosts")
	lineLists := strings.Split(string(contentHosts), "\n")
	hostMap := make(map[string]string)
	for _, i := range lineLists {
		if strings.Index(i, "#") == 0 {
			continue
		} else {
			if i == "\n" || i == " \n" || i == " " {
				continue
			} else {
				r, _ := regexp.Compile(`(?P<ip>[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)\s+(?P<host>\S+)`)
				if result := r.FindStringSubmatch(i); result != nil {
					//pp.Print(result)

					if _, ok := hostMap[result[2]]; !ok {
						hostMap[result[2]] = result[1]
					}
				}
				//fmt.Println(i)
			}

		}
	}
	addHostFormat := "%s %s\n"
	addHostStr := ""
	if _, ok := hostMap[name]; !ok {
		addHostStr = addHostStr + fmt.Sprintf(addHostFormat, ip, name)
		//fmt.Println(addHostStr)
	}
	WriteFile("/etc/hosts", string(contentHosts)+addHostStr)
	return nil
}
func (srcAc *AccessAliYun) LstM(result map[string]ecs.InstancesInDescribeInstances, keyMap map[string]string) {
	for _, v := range result {
		if len(v.Instance) > 0 {
			//fmt.Println(k, len(v.Instance))
			for _, i := range v.Instance {
				if i.OSType == "linux" {
					if len(keyMap) > 0 {
						fmt.Println(
							i.Tags.Tag,
							//i.Tags.Tag[0].TagValue,
							//i.Tags.Tag[0].TagKey,
							i.VpcAttributes.PrivateIpAddress.IpAddress[0],
							i.InstanceName,
							i.SecurityGroupIds.SecurityGroupId,
							keyMap[i.KeyPairName])
						GetHosts(i.InstanceName, i.VpcAttributes.PrivateIpAddress.IpAddress[0])
					} else {
						fmt.Println(
							i.Tags.Tag,
							//i.Tags.Tag[0].TagValue,
							//i.Tags.Tag[0].TagKey,
							i.VpcAttributes.PrivateIpAddress.IpAddress[0],
							i.InstanceName,
							i.SecurityGroupIds.SecurityGroupId,
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
	err := ioutil.WriteFile("inventory.ini", b, 0777)
	CheckErr(err)
}
func PrintInventoryYaml(f AnsibleInventory) []byte {
	b, err := yaml.Marshal(&f)
	if err != nil {
		fmt.Println(err)
	}
	return b
}

func (srcAc *AccessAliYun) GenerateSecureCRTSession(result map[string]ecs.InstancesInDescribeInstances, filePath AttrSessionSecureCRT, keyMap map[string]string) error {
	for _, v := range result {
		for _, i := range v.Instance {
			if i.OSType == "linux" {
				configTpl, err := pongo2.FromBytes(AliReadFile(filePath.Tpl1))
				CheckErr(err)
				sessionConfig, err := configTpl.Execute(pongo2.Context{"name": map[string]interface{}{
					"host":     i.VpcAttributes.PrivateIpAddress.IpAddress[0],
					"pem":      keyMap[i.KeyPairName],
					"username": "root",
				}})
				CheckErr(err)
				err = ioutil.WriteFile(fmt.Sprintf("%s/%s%s", filePath.ConfigPath, i.InstanceName, ".ini"), []byte(sessionConfig), 0644)
				CheckErr(err)
				personalTpl, err := pongo2.FromBytes(AliReadFile(filePath.Tpl2))
				CheckErr(err)
				personalConfig, err := personalTpl.Execute(pongo2.Context{"name": map[string]interface{}{
					"host":     i.VpcAttributes.PrivateIpAddress.IpAddress[0],
					"pem":      keyMap[i.KeyPairName],
					"username": "root",
				}})
				err = ioutil.WriteFile(fmt.Sprintf("%s/%s%s", filePath.PersonalPath, i.InstanceName, ".ini"), []byte(personalConfig), 0644)
				CheckErr(err)
			}
		}
	}
	return nil
}
func AliReadFile(filePath string) []byte {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Read error")
	}
	return content
}
