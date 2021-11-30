package main

import (
	"strings"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	// Embed the default VM context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultVMContext
}

// Override types.DefaultVMContext.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{}
}

type pluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
}

func (ctx *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {

	return types.OnPluginStartStatusOK
}

// Override types.DefaultPluginContext.
func (*pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	callBack := func(numHeaders, bodySize, numTrailers int) {
		proxywasm.LogInfo(("called callBack func"))
	}

	return &httpHeaders{contextID: contextID, callBack: callBack}
}

type httpHeaders struct {
	// Embed the default http context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultHttpContext
	contextID uint32
	callBack  func(numHeaders, bodySize, numTrailers int)
}

func (ctx *httpHeaders) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	hs, err := proxywasm.GetHttpRequestHeaders()
	if err != nil {
		proxywasm.LogCriticalf("failed to get request headers: %v", err)
		return types.ActionContinue
	}

	interceptCall := false
	for _, h := range hs {
		if strings.Compare(h[0], "intercept") == 0 && strings.Compare(h[1], "1") == 0 {
			interceptCall = true
		}
	}

	if interceptCall {
		proxywasm.LogInfo("intercepting call!!")
		// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attribute
		clusterName, err := proxywasm.GetProperty([]string{"cluster_name"})
		if err != nil {
			proxywasm.LogCriticalf("failed to get cluster name: %v", err)
			return types.ActionContinue
		}

		// proxywasm.DispatchHttpCall(...)
		proxywasm.LogInfof("cluster name: %s", string(clusterName))

		headers := [][2]string{
			// These headers have to be specified...
			{":method", "GET"}, {":authority", "some_authority"}, {"accept", "*/*"}, {":path", "/headers"},
		}
		if _, err := proxywasm.DispatchHttpCall(string(clusterName), headers, nil, nil, 5000, ctx.callBack); err != nil {
			proxywasm.LogCriticalf("dispatch httpcall failed: %v", err)
		}
	}

	return types.ActionContinue
}

func (ctx *httpHeaders) OnHttpResponseHeaders(numHeaders int, endOfStream bool) types.Action {
	return types.ActionContinue
}

func (ctx *httpHeaders) OnHttpStreamDone() {
}
