package main

/* create a validation admission controller [VAC] that mimics the kubernetes
 * API server. The VAC shall be available over secure TLS server on a service.
 */

import (
	"fmt"
	"net/http"
	"os"
	"time"

	// flags for the --help option with VAC
	"github.com/spf13/pflag"

	// for emulating the k8s API server
	"k8s.io/apiserver/pkg/server"

	// parameters to the VAC should be available on same lines as
	// k8s API server for passing the TLS certificate and key.
	"k8s.io/apiserver/pkg/server/options"

	// to register the flags
	"k8s.io/component-base/cli/globalflag"
)

// name of the VAC as a pre-defined constant
const (
	vcon = "mcluster-vcontroller"
)

// -----------------------------------OPTIONS-START------------------------
// wrapper for VAC options
type Options struct {
	SecureServingOptions options.SecureServingOptions
}

// method on options to add flags; sets the default value for port and pairname.
func NewDefaultOptions() *Options {
	// instantiate the options object
	opt := &Options{
		SecureServingOptions: *options.NewSecureServingOptions(),
	}
	// port will be 8443
	opt.SecureServingOptions.BindPort = 8443

	// the pairname will be based on the VAC name and be created if self signed
	// certificate is made. ie., CertDirectory/PairName.crt and CertDirectory/PairName.key
	opt.SecureServingOptions.ServerCert.PairName = vcon

	return opt
}

// method to add additional flags
func (opt *Options) AddFlagSet(fset *pflag.FlagSet) {
	opt.SecureServingOptions.AddFlags(fset)
}

// -----------------------------------OPTIONS-END------------------------

// -----------------------------------CONFIG_START------------------------
// wrapper for VAC configuration
type Config struct {
	SecureServingInfo *server.SecureServingInfo
}

// method to initialize options to return config that will eventually run server
func (opt *Options) initConfig() *Config {
	/* if CA bundle isnt available use self signed certificates.
		   MaybeDefaultWithSelfSignedCerts() takes Address, publicAddress, alternateDNS, and alternateIPS.
	    	   however this VAC will work on servicename and not publicAddress, hence we shall pass 0.0.0.0
	*/
	if err := opt.SecureServingOptions.MaybeDefaultWithSelfSignedCerts("0.0.0.0", nil, nil); err != nil {
		panic(err)
	}

	// instantiate the config object
	con := Config{}

	// use the certs in the config...
	opt.SecureServingOptions.ApplyTo(&con.SecureServingInfo)

	return &con
}

// -----------------------------------CONFIG_END------------------------

// driver for Validation Admission Controller..
func main() {
	// initialize the defalut options
	options := NewDefaultOptions()

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

	// create http handler
	mux := http.NewServeMux()

	// register validation functon "ServerMclusterValidation" to http handler
	mux.Handle("/", http.HandlerFunc(ServerMclusterValidation))

	//create channel that can be passed to .Serve
	stopCh := server.SetupSignalHandler()

	/* run the https server by callng .Serve on config info
		.Server is a non blocking call and will be trigerred every 20 secs.

	 	Serve runs the secure http server. It fails only if certificates
	  	cannot be loaded or the initial listen call fails. The actual server loop
	   	(stoppable by closing stopCh) runs in a go routine, i.e. Serve does not block.
		It returns a stoppedCh that is closed when all non-hijacked active
	 	requests have been processed. It returns a listenerStoppedCh that is closed
	  	when the underlying http Server has stopped listening.
	*/
	ch, _, err := con.SecureServingInfo.Serve(mux, 20*time.Second, stopCh)
	if err != nil {
		panic(err)
	} else {
		<-ch
	}
}

// primary logic for the VAC is served by this method.. Based on the implementation
// ServerMclusterValidation() shall decide to accept or reject CRUDops on
// mcluster object.
func ServerMclusterValidation(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ServeMclusterValidation is called...")
}
