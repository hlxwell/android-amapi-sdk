package types

// EnrollmentToken 相关类型和函数
//
// 注意：EnrollmentToken 类型直接使用 androidmanagement.EnrollmentToken。
// 此文件包含注册令牌相关的辅助函数和常量。
//
// 使用方式：
//
//	import "amapi-pkg/pkgs/amapi/types"
//
//		// 创建注册令牌直接传递参数
	//	token, err := client.EnrollmentTokens().Create(
	//	    "enterprises/LC00abc123",
	//	    "enterprises/LC00abc123/policies/default",
	//	    24 * time.Hour,
	//	    false, // allowPersonalUsage
	//	    false, // oneTimeOnly
	//	    &androidmanagement.User{AccountIdentifier: "user@example.com"},
	//	)
//
//	// 生成 QR 码数据
//	options := types.NewBasicQRCodeOptions()
//	qrData := types.GenerateQRCodeData(token, options)

// WiFi security type constants
const (
	WiFiSecurityTypeNone = "NONE"
	WiFiSecurityTypeWEP  = "WEP"
	WiFiSecurityTypeWPA  = "WPA"
	WiFiSecurityTypeWPA2 = "WPA2"
)

// NewBasicQRCodeOptions creates basic QR code options.
func NewBasicQRCodeOptions() *QRCodeOptions {
	return &QRCodeOptions{
		SkipSetupWizard: true,
		Locale:          "en_US",
	}
}

// NewQRCodeOptionsWithWiFi creates QR code options with WiFi configuration.
func NewQRCodeOptionsWithWiFi(ssid, password, securityType string) *QRCodeOptions {
	options := NewBasicQRCodeOptions()
	options.WiFiSSID = ssid
	options.WiFiPassword = password
	options.WiFiSecurityType = securityType
	return options
}

