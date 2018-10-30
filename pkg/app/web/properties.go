package web

type view struct {
	Enabled      bool
	ContextPath  string `default:"/"`
	DefaultPage  string `default:"index.html"`
	ResourcePath string `default:"static"`
	Extension    string `default:".html"`
}

type properties struct {
	View view
}
