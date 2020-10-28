package config

//
//func TestBuilder_SetDefault(t *testing.T) {
//	bd := new(Builder)
//	b := bd.SetDefault().Build()
//
//	assert.Equal(t, &Config{
//		WebAddr:       "0.0.0.0:8042",
//		BasePath:      "/",
//		PubKeyPath:    "./id_25519.pub",
//		SecretKeyPath: "./id_25519",
//		NoKeyGen:      false,
//		SecretKey:     nil,
//		PublicKey:     nil,
//	}, b)
//}
//
//func TestGetDefault(t *testing.T) {
//	bd := new(Builder).SetDefault().Build()
//	gd := GetDefault()
//	// We don't want to test equality of literally randomly generated numbers.
//	gd.SecretKey = nil
//	gd.PublicKey = nil
//	assert.Equal(t, bd, gd)
//}
//
//func TestBuilder_SetFromEnvironment(t *testing.T) {
//	bd := new(Builder)
//
//	_ = os.Setenv("WEBADDR", "url")
//	_ = os.Setenv("BASEPATH", "/asd")
//	b := bd.SetFromEnvironment().Build()
//	_ = os.Setenv("WEBADDR", "")
//	_ = os.Setenv("BASEPATH", "")
//
//	assert.Equal(t, &Config{
//		WebAddr:  "url",
//		BasePath: "/asd",
//	}, b)
//}
//
//func TestBuilder_SetDefault_SetFromEnvironment(t *testing.T) {
//	bd := new(Builder)
//	_ = os.Setenv("PUBLIC_KEY_PATH", "/yeet")
//	b := bd.SetDefault().SetFromEnvironment().Build()
//	_ = os.Setenv("PUBLIC_KEY_PATH", "")
//
//	assert.Equal(t, &Config{
//		WebAddr:       "0.0.0.0:8042",
//		BasePath:      "/",
//		NoKeyGen:      false,
//		PubKeyPath:    "/yeet",
//		SecretKeyPath: "./id_25519",
//		PublicKey:     nil,
//		SecretKey:     nil,
//	}, b)
//}
