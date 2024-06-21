package types

import (
	"context"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

// WHParameters defines the structure for webhook parameters
type WHParameters struct {
	// server parameters
	TLSKey     string
	TLSCert    string
	TLSPort    string
	HealthPort string

	//	handlers parameters
	// LabelsToCheck    []string
	// DeploymentPrefix string
}

// Enable or disable specific features as needed
type ConfigEnabel struct {
	// mutatingwebhookconfigurations
	MutatePodEnvInjectedHandle bool

	// validatingwebhookconfigurations
	ValidateNamespaceLabelsHandle bool
	ValidateCheckDeploymentPrefix bool
}

type ConfigHandlersParameters struct {
	//	handlers parameters
	LabelsToCheck    []string
	DeploymentPrefix string
}

type DynamicClient struct {
	dynamicClient dynamic.Interface
}

// 通过调用 types.NewClient() 创建一个新的 Kubernetes 客户端实例。
func NewDynamicClient() (*DynamicClient, error) {
	// 1.GET config in k8s cluster（自动获取部署在 Kubernetes 集群内的 Pod 的服务账户和 API 服务器的地址，无需手动提供 kubeconfig 文件。）
	config, err := rest.InClusterConfig()
	if err != nil {
		glog.Errorf("Failed to obtain the internal configuration file of the cluster .")
		return nil, err
	}

	// 2.使用配置创建动态客户端
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		glog.Errorf("Failed to configure dynamic client.")
		return nil, err
	}

	// 3.返回包含动态客户端的 Client 实例
	return &DynamicClient{dynamicClient: dynamicClient}, nil
}

// GetResourceByGVR 根据 GVR、命名空间和名称获取资源，检查资源是否存在
func (c *DynamicClient) GetResourceByGVR(gvr schema.GroupVersionResource, namespace, name string) (bool, error) {

	// unStructData, err := c.dynamicClient.Resource(gvr).Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	_, err := c.dynamicClient.Resource(gvr).Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	// var obj ValidateCheckDeploymentPrefix
	// 使用 runtime.DefaultUnstructuredConverter 转换 item 为对象
	// runtime.DefaultUnstructuredConverter.FromUnstructured(unStructData.UnstructuredContent())

	if err != nil {
		if errors.IsNotFound(err) {
			// 资源未找到，不是错误情况，只是资源不存在
			glog.Infof("Not found resource %s of type %s in namespace %s.", name, gvr.Resource, namespace)
			return false, nil
		}
		// 发生其他错误，返回错误
		glog.Errorf("Error getting resource %s of type %s in namespace %s: %v", name, gvr.Resource, namespace, err)
		return false, err
	}

	// 资源找到，打印信息（如果需要）
	glog.Infof("Found resource %s of type %s in namespace %s.", name, gvr.Resource, namespace)
	// 通常这里不需要使用 unStructData，除非你想进一步处理它

	return true, nil
}
