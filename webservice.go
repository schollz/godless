package godless

import (
	"net/http"
)

type IpfsNamespaceService struct {
	Namespace *IpfsNamespace
}

func (service *IpfsNamespaceService) Handler() http.Handler {
	return nil
}

func (service *IpfsNamespaceService) getMap(rw http.ResponseWriter, r *http.Request) {

}

func (service *IpfsNamespaceService) getMapValues(rw http.ResponseWriter, r *http.Request) {

}

func (service *IpfsNamespaceService) addMapValues(rw http.ResponseWriter, r *http.Request) {

}

func (service *IpfsNamespaceService) getSet(rw http.ResponseWriter, r *http.Request) {

}

func (service *IpfsNamespaceService) addSetValues(rw http.ResponseWriter, r *http.Request) {

}

func (service *IpfsNamespaceService) getNamespace(rw http.ResponseWriter, r *http.Request) {

}
