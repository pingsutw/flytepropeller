// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots.

package config

import (
	"encoding/json"
	"reflect"

	"fmt"

	"github.com/spf13/pflag"
)

// If v is a pointer, it will get its element value or the zero value of the element type.
// If v is not a pointer, it will return it as is.
func (Config) elemValueOrNil(v interface{}) interface{} {
	if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
		if reflect.ValueOf(v).IsNil() {
			return reflect.Zero(t.Elem()).Interface()
		} else {
			return reflect.ValueOf(v).Interface()
		}
	} else if v == nil {
		return reflect.Zero(t).Interface()
	}

	return v
}

func (Config) mustJsonMarshal(v interface{}) string {
	raw, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return string(raw)
}

func (Config) mustMarshalJSON(v json.Marshaler) string {
	raw, err := v.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return string(raw)
}

// GetPFlagSet will return strongly types pflags for all fields in Config and its nested types. The format of the
// flags is json-name.json-sub-name... etc.
func (cfg Config) GetPFlagSet(prefix string) *pflag.FlagSet {
	cmdFlags := pflag.NewFlagSet("Config", pflag.ExitOnError)
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "metrics-prefix"), DefaultConfig.MetricsPrefix, "An optional prefix for all published metrics.")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "certDir"), DefaultConfig.CertDir, "Certificate directory to use to write generated certs. Defaults to /etc/webhook/certs/")
	cmdFlags.Bool(fmt.Sprintf("%v%v", prefix, "localCert"), DefaultConfig.LocalCert, "write certs locally. Defaults to false")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "listenPort"), DefaultConfig.ListenPort, "The port to use to listen to webhook calls. Defaults to 9443")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "serviceName"), DefaultConfig.ServiceName, "The name of the webhook service.")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "secretName"), DefaultConfig.SecretName, "Secret name to write generated certs to.")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "awsSecretManager.sidecarImage"), DefaultConfig.AWSSecretManagerConfig.SidecarImage, "Specifies the sidecar docker image to use")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "vaultSecretManager.role"), DefaultConfig.VaultSecretManagerConfig.Role, "Specifies the vault role to use")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "metricsBindAddress"), DefaultConfig.MetricsBindAddress, "TCP address that the controller should bind for serving prometheus metrics.")
	return cmdFlags
}
