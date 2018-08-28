package jwt

type Properties struct {
	PrivateKeyPath string `default:"config/ssl/app.rsa"`
	PublicKeyPath  string `default:"config/ssl/app.rsa.pub"`
}
