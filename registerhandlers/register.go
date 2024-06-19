package registerhandlers

import (
	"net/http"
	"simple-webhook/handlers"
	"simple-webhook/types"
)

// func RegisterHandlers(configEnable types.ConfigEnabel, labels []string) {
func RegisterHandlers(configEnable types.ConfigEnabel, configHandlersParameters types.ConfigHandlersParameters) {
	if configEnable.MutatePodEnvInjectedHandle {
		http.HandleFunc("/mutate-pod", handlers.PodEnvInjectedHandleMutate)
	}
	if configEnable.ValidateNamespaceLabelsHandle {
		// glog.Info("configHandlersParameters.LabelsToCheck:", labels)
		// glog.Info("configHandlersParameters.LabelsToCheck:", configHandlersParameters.LabelsToCheck)
		http.HandleFunc("/validate-namespace", handlers.NamespaceLabelsHandleValidate(configHandlersParameters.LabelsToCheck))
		// http.HandleFunc("/validate", handlers.NamespaceLabelsHandleValidate(configHandlersParameters.LabelsToCheck))
	} else {
		http.HandleFunc("/validate-namespace", handlers.AllowedHandlers())
	}

	if configEnable.ValidateCheckDeploymentPrefix {
		http.HandleFunc("/validate-deploy", handlers.CheckDeployPrefixHandleValidate(configHandlersParameters.DeploymentPrefix))
	} else {
		http.HandleFunc("/validate-deploy", handlers.AllowedHandlers())
	}
}
