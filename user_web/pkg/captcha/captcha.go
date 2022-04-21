package captcha

import (
	"github.com/mojocn/base64Captcha"
	"sync"
	"user_web/pkg/config"
	"user_web/pkg/redis"
)

type Captcha struct {
	Base64Captcha *base64Captcha.Captcha
}

var internalCaptcha *Captcha
var once sync.Once

// NewCaptcha 获取生成验证码实例
func NewCaptcha() *Captcha {
	once.Do(func() {
		internalCaptcha = &Captcha{}
		store := StoreRedis{
			Client: redis.Redis,
		}
		driver := base64Captcha.NewDriverDigit(
			config.GetInt("CAPTCHA_HEIGHT"),
			config.GetInt("CAPTCHA_WIDTH"),
			config.GetInt("CAPTCHA_LENGTH"),
			config.GetFloat64("CAPTCHA_MAXSKEW"),
			config.GetInt("CAPTCHA_DOTCOUNT"),
		)
		internalCaptcha.Base64Captcha = base64Captcha.NewCaptcha(driver, &store)
	})

	return internalCaptcha
}

// GenerateCaptcha 生成验证码
func (captcha *Captcha) GenerateCaptcha() (id string, b64s string, err error) {
	return captcha.Base64Captcha.Generate()
}

// VerifyCaptcha 验证验证码
func (captcha *Captcha) VerifyCaptcha(id string, ans string, clear bool) bool {
	if id == config.GetString("CAPTCHA_TEST_ID") {
		return true
	}
	return captcha.Base64Captcha.Verify(id, ans, clear)
}
