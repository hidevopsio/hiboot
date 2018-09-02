package web

import "github.com/kataras/iris"

// DefaultConfiguration returns the default configuration for an iris station, fills the main Configuration
func defaultConfiguration() iris.Configuration {
	return iris.Configuration{
		DisableStartupLog:                 false,
		DisableInterruptHandler:           false,
		DisableVersionChecker:             true,
		DisablePathCorrection:             false,
		EnablePathEscape:                  false,
		FireMethodNotAllowed:              false,
		DisableBodyConsumptionOnUnmarshal: false,
		DisableAutoFireStatusCode:         false,
		TimeFormat:                        "Mon, Jan 02 2006 15:04:05 GMT",
		Charset:                           "UTF-8",

		// PostMaxMemory is for post body max memory.
		//
		// The request body the size limit
		// can be set by the middleware `LimitRequestBodySize`
		// or `context#SetMaxRequestBodySize`.
		PostMaxMemory:               32 << 20, // 32MB
		TranslateFunctionContextKey: "app.translate",
		TranslateLanguageContextKey: "app.language",
		ViewLayoutContextKey:        "app.viewLayout",
		ViewDataContextKey:          "app.viewData",
		RemoteAddrHeaders:           make(map[string]bool),
		EnableOptimizations:         false,
		Other:                       make(map[string]interface{}),
	}
}
