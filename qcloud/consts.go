package qcloud

const SecretId  = "ReplaceMeUseYourSecretId"
const SecretKey  = "ReplaceMeUseYourSecretKey"
const DomainBase  = "api.qcloud.com"
const HostPath  = "/v2/index.php"
const Method = "GET"

const (
	SUCCESS int = 0
)

type Action string
const (
	//----------bm------------------------
	DescribeDevice Action = "DescribeDevice"
	//-------------------------------------

	//-------- lb ---------------------------------------------------------------
	CreateBmLoadBalancer Action = "CreateBmLoadBalancer"
	DescribeBmLoadBalancers Action = "DescribeBmLoadBalancers"
	DeleteBmLoadBalancers Action  = "DeleteBmLoadBalancers"
	ModifyBmLoadBalancerAttributes Action = "ModifyBmLoadBalancerAttributes"
	CreateBmListeners  Action = "CreateBmListeners"
	DescribeBmListeners Action = "DescribeBmListeners"
	BindBmL4ListenerRs Action = "BindBmL4ListenerRs"
	DeleteBmListeners Action = "DeleteBmListeners"
	DescribeBmLoadBalancersTaskResult Action = "DescribeBmLoadBalancersTaskResult"
	DescribeBmBindInfo Action = "DescribeBmBindInfo"
	DescribeBmVportInfo Action = "DescribeBmVportInfo"
	DescribeBmListenerInfo Action = "DescribeBmListenerInfo"
	DescribeBmForwardListenerInfo Action = "DescribeBmForwardListenerInfo"
	CreateBmForwardListeners Action = "CreateBmForwardListeners"
	//---------------------------------------------------------------------------

	//----------------- vpc--------------------------
	DescribeBmVpcEx Action = "DescribeBmVpcEx"
	DescribeBmSubnetEx Action = "DescribeBmSubnetEx"
	//-----------------------------------------------

	//------------------cns----------------------
	DomainCreate Action = "DomainCreate"
	RecordCreate Action = "RecordCreate"
	DomainList Action = "DomainList"
	RecordList Action = "RecordList"
	RecordDelete Action = "RecordDelete"
	//-------------------------------------------
)

type Module string
const (
	Bm Module = "bm"
	Bmlb Module = "bmlb"
	Cns Module = "cns"
	Bmvpc Module = "bmvpc"
)

type Region string
const (
	BJ Region = "bj"
)

