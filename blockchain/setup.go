package blockchain

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	//"github.com/hyperledger/fabric-sdk-go/api/apitxn/resmgmtclient"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"fmt"
	//"github.com/hyperledger/fabric-sdk-go/api/apitxn/chmgmtclient"
	"time"
)

//定义结构体
type FabricSetup struct {
	ConfigFile    string                           //sdk配置文件所在路径
	ChannelID     string                           //应用通道名称
	ChannelConfig string                           //应用通道交易配置文件所在路径
	OrgAdmin      string                           // 组织管理员名称
	OrgName       string                           //组织名称
	Initialized   bool                             //是否初始化
	Admin         resmgmtclient.ResourceMgmtClient //fabric环境中资源管理者
	SDK           *fabsdk.FabricSDK                //SDK实例
}

//1. 创建SDK实例并使用SDK实例创建应用通道，将Peer节点加入到创建的应用通道中
func (f *FabricSetup) Initialize() error {
	//判断是否已经初始化
	if f.Initialized {
		return fmt.Errorf("SDK已被实例化")
	}

	//创建SDK对象
	sdk, err := fabsdk.New(config.FromFile(f.ConfigFile))

	if err != nil {
		return fmt.Errorf("SDK实例化失败:%v", err)
	}
	f.SDK = sdk

	//创建一个具有管理权限的应用通道客户端管理对象
	//prepare contexts
	ordererClientContext := sdk.Context(fabsdk.WithUser(f.OrgAdmin), fabsdk.WithOrg(f.OrgName))
	chmgmtclient, err := resmgmt.New(ordererClientContext)
	if err != nil {
		return fmt.Errorf("创建应用通道管理客户端管理对象失败,%v", err)
	}

	//获取当前的会话用户对象
	sessionClientContext := sdk.Context(fabsdk.WithUser(f.OrgAdmin), fabsdk.WithOrg(f.OrgName))
	session, err := resmgmt.New(sessionClientContext)
	if err != nil {
		return fmt.Errorf("获取当前会话用户对象失败%v", err)
	}

	orgAdminUser := session
	//指定创建应用通道所需要的所有参数
	/*
	$ peer channel create -o orderer.example.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/channel.tx --tls --cafile \
	/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
	 */
	chReq := chmgmtclient.SaveChannelRequest{ChannelID: f.ChannelID, ChannelConfig: f.ChannelConfig, SigningIdentity: orgAdminUser}

	//创建应用通道

	err = chmClient.SaveChannel(chReq)
	if err != nil {
		return fmt.Errorf("创建应用通道失败:%v", err)
	}

	time.Sleep(time.Second * 5)

	//创建一个管理资源的客户端对象
	f.Admin, err = f.SDK.NewClient(fabsdk.WithUser(f.OrgAdmin)).ResourceMgmt()
	if err != nil {
		return fmt.Errorf("创建资源管理对象失败:%v", err)
	}
	//将peer 节点加入到应用通道中
	err = f.Admin.JoinChannel(f.ChannelID)
	if err != nil {
		return fmt.Errorf("peer加入节点失败:%v", err)
	}

	f.Initialized = true
	fmt.Println("SDK实例化成功")

	return nil

}