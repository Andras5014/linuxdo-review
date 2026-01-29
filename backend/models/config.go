package models

import "strconv"

// 系统配置键常量
const (
	ConfigMinVotes       = "min_votes"        // 进入二级审核的最小票数
	ConfigApprovalRate   = "approval_rate"    // 赞率阈值(百分比)
	ConfigSMTPHost       = "smtp_host"        // SMTP服务器地址
	ConfigSMTPPort       = "smtp_port"        // SMTP端口
	ConfigSMTPUser       = "smtp_user"        // SMTP用户名
	ConfigSMTPPass       = "smtp_pass"        // SMTP密码
	ConfigSMTPFrom       = "smtp_from"        // 发件人邮箱
	ConfigLinuxDoClientID     = "linuxdo_client_id"     // OAuth客户端ID
	ConfigLinuxDoClientSecret = "linuxdo_client_secret" // OAuth客户端密钥
	ConfigLinuxDoRedirectURI  = "linuxdo_redirect_uri"  // OAuth回调地址
	ConfigSiteName       = "site_name"        // 站点名称
	ConfigSiteURL        = "site_url"         // 站点URL
)

// 默认配置值
const (
	DefaultMinVotes     = 10
	DefaultApprovalRate = 70
	DefaultSMTPPort     = 587
	DefaultSiteName     = "LinuxDo邀请码申请系统"
)

// SystemConfig 系统配置模型
type SystemConfig struct {
	Key         string `gorm:"primaryKey;size:100" json:"key"`
	Value       string `gorm:"type:text" json:"value"`
	Description string `gorm:"size:255" json:"description"`
}

// TableName 指定表名
func (SystemConfig) TableName() string {
	return "configs"
}

// GetIntValue 获取整数配置值
func (c *SystemConfig) GetIntValue(defaultValue int) int {
	if c == nil || c.Value == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(c.Value)
	if err != nil {
		return defaultValue
	}
	return val
}

// GetFloatValue 获取浮点数配置值
func (c *SystemConfig) GetFloatValue(defaultValue float64) float64 {
	if c == nil || c.Value == "" {
		return defaultValue
	}
	val, err := strconv.ParseFloat(c.Value, 64)
	if err != nil {
		return defaultValue
	}
	return val
}

// GetDefaultConfigs 获取默认配置列表
func GetDefaultConfigs() []*SystemConfig {
	return []*SystemConfig{
		{Key: ConfigMinVotes, Value: strconv.Itoa(DefaultMinVotes), Description: "进入二级审核的最小票数"},
		{Key: ConfigApprovalRate, Value: strconv.Itoa(DefaultApprovalRate), Description: "赞率阈值(百分比)"},
		{Key: ConfigSMTPHost, Value: "", Description: "SMTP服务器地址"},
		{Key: ConfigSMTPPort, Value: strconv.Itoa(DefaultSMTPPort), Description: "SMTP端口"},
		{Key: ConfigSMTPUser, Value: "", Description: "SMTP用户名"},
		{Key: ConfigSMTPPass, Value: "", Description: "SMTP密码"},
		{Key: ConfigSMTPFrom, Value: "", Description: "发件人邮箱"},
		{Key: ConfigLinuxDoClientID, Value: "", Description: "LinuxDo OAuth客户端ID"},
		{Key: ConfigLinuxDoClientSecret, Value: "", Description: "LinuxDo OAuth客户端密钥"},
		{Key: ConfigLinuxDoRedirectURI, Value: "", Description: "LinuxDo OAuth回调地址"},
		{Key: ConfigSiteName, Value: DefaultSiteName, Description: "站点名称"},
		{Key: ConfigSiteURL, Value: "", Description: "站点URL"},
	}
}
