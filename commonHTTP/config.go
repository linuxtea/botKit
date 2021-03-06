package commonHTTP

func errCodeConfig() []*errConfig {
	return []*errConfig{
		// Client error
		&errConfig{code: 4001, msg: "参数错误"},
		&errConfig{code: 4002, msg: "验证码已存在"},
		&errConfig{code: 4003, msg: "验证码错误"},
		&errConfig{code: 4004, msg: "用户不存在"},
		&errConfig{code: 4005, msg: "无效签名"},
		&errConfig{code: 4006, msg: "用户密码不可为空"},
		&errConfig{code: 4007, msg: "重复提交"},

		// Server error
		&errConfig{code: 5001, msg: "系统错误"},
	}
}
