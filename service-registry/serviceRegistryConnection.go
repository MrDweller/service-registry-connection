package serviceregistry

import (
	"errors"
	"fmt"

	"github.com/MrDweller/service-registry-connection/models"
)

type ServiceRegistryConnection interface {
	Connect() error
	RegisterService(digitalTwinModel models.ServiceDefinition, systemDefinition models.SystemDefinition) ([]byte, error)
	UnRegisterService(serviceDefinition models.ServiceDefinition, systemDefinition models.SystemDefinition) error
	RegisterSystem(systemDefinition models.SystemDefinition) ([]byte, error)
	UnRegisterSystem(systemDefinition models.SystemDefinition) error
	Query(serviceDefinition models.ServiceDefinition) (*models.ServiceQueryResult, error)
}

type ServiceRegistryImplementationType string

func NewConnection(serviceRegistry ServiceRegistry, serviceRegistryImplementationType ServiceRegistryImplementationType, certificateInfo models.CertificateInfo) (ServiceRegistryConnection, error) {
	var serviceRegistryConnection ServiceRegistryConnection

	switch serviceRegistryImplementationType {
	case SERVICE_REGISTRY_ARROWHEAD_4_6_1:
		serviceRegistryConnection = ServiceRegistryArrowhead_4_6_1{
			ServiceRegistry: serviceRegistry,
			CertificateInfo: certificateInfo,
		}
		break
	default:
		errorString := fmt.Sprintf("the service registry %s has no implementation", serviceRegistryImplementationType)
		return nil, errors.New(errorString)
	}

	err := serviceRegistryConnection.Connect()
	if err != nil {
		return nil, err
	}

	return serviceRegistryConnection, nil
}
