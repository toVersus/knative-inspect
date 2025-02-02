package core

import (
	"fmt"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Config interface {
	ConfigurePathResolver(func() (string, error))
	ConfigureContextResolver(func() (string, error))
	RESTConfig() (*rest.Config, error)
	DefaultNamespace() (string, error)
}

type ConfigImpl struct {
	pathResolverFunc    func() (string, error)
	contextResolverFunc func() (string, error)
}

var _ Config = &ConfigImpl{}

func NewConfig() Config {
	return &ConfigImpl{}
}

func (f *ConfigImpl) ConfigurePathResolver(resolverFunc func() (string, error)) {
	f.pathResolverFunc = resolverFunc
}

func (f *ConfigImpl) ConfigureContextResolver(resolverFunc func() (string, error)) {
	f.contextResolverFunc = resolverFunc
}

func (f *ConfigImpl) RESTConfig() (*rest.Config, error) {
	config, err := f.clientConfig()
	if err != nil {
		return nil, err
	}

	restConfig, err := config.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("building Kubernetes config: %s", err)
	}

	return restConfig, nil
}

func (f *ConfigImpl) DefaultNamespace() (string, error) {
	config, err := f.clientConfig()
	if err != nil {
		return "", err
	}

	name, _, err := config.Namespace()
	return name, err
}

func (f *ConfigImpl) clientConfig() (clientcmd.ClientConfig, error) {
	path, err := f.pathResolverFunc()
	if err != nil {
		return nil, fmt.Errorf("resolving config path: %s", err)
	}

	context, err := f.contextResolverFunc()
	if err != nil {
		return nil, fmt.Errorf("resolving config context: %s", err)
	}

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: path},
		&clientcmd.ConfigOverrides{CurrentContext: context},
	), nil
}
