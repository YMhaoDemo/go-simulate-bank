package main

/*

开户流程：
   BankCardRMBBalance和BankCard是预先准备好的

  1、获取未被使用的新卡
  get  /BanskCards?status=UNUSED
  2、录入客户信息(insert)
  post /clients
  3、开户(insert)
  post /BankCardAccouts?password=xxx

销户流程： -> 分布式事物
  1、设置账户状态为Close(update)
  put /BankCardAccouts/{BankCardId}?status=CLOSE
  2、 查询余额
  get /BankCardRMBBalance
  3、提现
  Post /Deposit
  4、 查询余额
  get /BankCardRMBBalance

查询余额
  get /BankCardRMBBalances/{BankCardId}?password=xxx

存钱：
  (要注意要检查账户是否有效，是否被冻结)
  post /Withdrawals?password=xxx （insert）
  get  /Withdrawals/{TradeId}

取钱:
  (要注意要检查账户是否有效，是否被冻结)
  post /Deposits?password=xxx （insert）
  get /Deposits/{TradeId}

转账:
  post /Trades?tellerPassword=xxx&clientPassword=xxx
  get /Trades/{TradeId}

*/

type Client struct {
	ClientId   string `desc:"client Id 身份证号"`
	ClientName string `desc:"client name 姓名"`
	ClientMob  string `desc:"mobile phone number 手机号"`
}

type BankCard struct {
	BankCardId     string `desc:"bank Id 银行卡id"`
	SubBankId      string `desc:"sub bank id 分行编号"`
	BankCardStatus string `desc:"card status 银行卡状态" enum:"UNUSED,OPEN,CLOSE,LOST,TERMINATE"`
}

type BankCardPassword struct {
	BankCardId string `desc:"bank Id 银行卡id"`
	Password   string `desc:"passwd 银行卡密码 md5(bankCardID + md5(realPasswd))"`
}

/*
 * 做一些数据冗余，避免联表查询
 */
type BankCardAccout struct {
	BankCardId      string `desc:"bank Id 银行卡id"`
	SubBankId       string `desc:"sub bank id 分行编号"`
	Status          string `desc:"card status 银行卡状态" enum:"UNUSED,OPEN,CLOSE,LOST,TERMINATE"`
	ClientId        string `desc:"Client Id 身份证号"`
	Mob             string `desc:"mobile phone number 手机号"`
	CreateTimestamp string `desc:"create timestamp 开户时间"`
}

type BankCardRMBBalance struct {
	BankCardId string `desc:"bank Id 银行卡id"`
	Balance    int64  `desc:"balance 余额（单位分）"`
}

type Address struct {
}

type SubBank struct {
	SubBankId   string  `desc:"sub bank Id 分行编号"`
	SubBankName string  `desc:"sub bank Id 分行名称"`
	SubBnakAddr Address `desc:"sub bank address 分行地址"`
}

type Trade struct {
	TradeId        string `desc:"TradeId 交易编号"`
	FromBankCardId string `desc:"from bank card id 源银行卡号"`
	ToBanckCardId  string `desc:"to bank Card id 目标银行卡号"`
	Timestamp      int64  `desc:"time stamp 交易时间"`
	MoneyInCent    int64  `desc:"money in cent 交易金额分"`
	Factorage      int64  `desc:"Factorage 手续费"`
	Abstract       int64  `desc:"abstract 交易摘要"`
}

type Teller struct {
	TellerId   string `desc:"TellerId 柜员id"`
	TellerName string `desc:"TellerId 柜员名称"`
	SubBankId  string `desc:"sub bank id 分行编号"`
}

type TellerPassword struct {
	TellerId string `desc:"TellerId 柜员id"`
	Password string `desc:"passwd 银行卡密码 md5(bankCardID + md5(realPasswd))"`
}

type Withdrawal struct {
	TradeId      string `desc:"TradeId 交易编号"`
	FromTellerId string `desc:"TellerId 柜员id"`
	BanckCardId  string `desc:"to bank Card id 银行卡号"`
	MoneyInCent  int64  `desc:"money in cent 交易金额分"`
	Factorage    int64  `desc:"Factorage 手续费"`
	Timestamp    int64  `desc:"time stamp 交易时间"`
}

type Deposit struct {
	TradeId      string `desc:"TradeId 交易编号"`
	FromTellerId string `desc:"TellerId 柜员id"`
	BanckCardId  string `desc:"bank Card id 银行卡号"`
	MoneyInCent  int64  `desc:"money in cent 交易金额分"`
	Factorage    int64  `desc:"Factorage 手续费"`
	Timestamp    int64  `desc:"time stamp 交易时间"`
}

type SummaryEvent struct {
	TradeId     string `desc:"TradeId 交易编号"`
	TradeType   string `desc:"trade type 交易类型" enum:"DEPOSIT,WITHDRAWAL,TRADE"`
	FromId      string `desc:"from id 来源id"`
	FromIdType  string `desc:"from id 来源类型" enum:"CLIENT,TELLER"`
	ToId        string `desc:"to id 目标id"`
	ToIdType    string `desc:"to id 目标id类型" enum:"CLIENT,TELLER"`
	MoneyInCent int64  `desc:"money in cent 交易金额分"`
	ActionType  string `desc:"支出/收入类型" enum:"INCOME,OUTCOME"`
	Timestamp   int64  `desc:"time stamp 交易时间"`
}
