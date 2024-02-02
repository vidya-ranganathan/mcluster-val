package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/pflag"
	"k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/options"
	"k8s.io/component-base/cli/globalflag"
)

const (
	vcon = "mcluster-vcontroller"
)

type Options struct {
	SecureServingOptions options.SecureServingOptions
}

type Config struct {
	SecureServingInfo *server.SecureServingInfo
}

func (opt *Options) initConfig() *Config {
	// if CA bundle isnt available use self signed certificates
	if err := opt.SecureServingOptions.MaybeDefaultWithSelfSignedCerts("0.0.0.0", nil, nil); err != nil {
		panic(err)
	}

	con := Config{}

	// use the certs in the config...
	opt.SecureServingOptions.ApplyTo(&con.SecureServingInfo)

	return &con
}

func NewDefaultOPtions() *Options {
	opt := &Options{
		SecureServingOptions: *options.NewSecureServingOptions(),
	}
	opt.SecureServingOptions.BindPort = 8443
	opt.SecureServingOptions.ServerCert.PairName = vcon

	return opt
}

func (opt *Options) AddFlagSet(fset *pflag.FlagSet) {
	opt.SecureServingOptions.AddFlags(fset)
}

// create options struct
// method on options to add flags
// config struct
// method on options to return config that will eventually run serser
// function to return default options

func main() {
	// initialize the defalut options
	options := NewDefaultOPtions()
	// create new flag set
	fset := pflag.NewFlagSet(vcon, pflag.ExitOnError)
	// add the created flag to the options
	globalflag.AddGlobalFlags(fset, vcon)
	// parse flagset
	options.AddFlagSet(fset)
	if err := fset.Parse(os.Args); err != nil {
		panic(err)
	}

	// create config from options
	con := options.initConfig()

	// create hhtp handler
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(ServerMclusterValidation))

	//create channel that can be passed to .Serve
	stopCh := server.SetupSignalHandler()
	ch, _, err := con.SecureServingInfo.Serve(mux, 20*time.Second, stopCh)
	if err != nil {
		panic(err)
	} else {
		<-ch
	}

	//register validation functon to http handler
	// run the https server by callng .Serve on config info

}

func ServerMclusterValidation(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ServeMclusterValidation is called...")

}
