package jsbridge

type jsHtml struct {
	Version string
	Common  []jsBridgeData
	Mall    appBridge
}

type appBridge struct {
	Common  []jsBridgeData
	Android []jsBridgeData
	IOS     []jsBridgeData
}

type jsBridgeData struct {
	Index    int
	Token    string
	Comment  string
	Params   []jsBridgeParam
	Example  string
	IsCommon bool
}

type jsBridgeParam struct {
	Name string
	Kvs  []KeyValue
}

type KeyValue struct {
	Key   string
	Value string
}
