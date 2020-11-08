package EventsManagement

import (
	"context"
	"fmt"
	"sync"
)

type EventsManager struct {
	name          string                         // manager name
	eventsRouters map[string]SimilarEventsRouter // routers

	eventIndex int // index for new event

	ctx context.Context
}

func (manager *EventsManager) EventIndex() int {
	event_index := manager.eventIndex
	manager.eventIndex++
	return event_index
}

// get or create router
func (manager *EventsManager) GetRouter(eventsType string) (router *SimilarEventsRouter, created bool) {
	value, ok := manager.eventsRouters[eventsType]
	if ok {
		created = !ok
		router = &value
		return
	} else {
		created = true
		router = &SimilarEventsRouter{
			manager:     manager,
			eventsClass: eventsType,
			income:      make(chan *Event),
			adapters:    make(map[string]*RouterAdapter),
			active:      false,
			mutex:       sync.RWMutex{},
			stop:        make(chan bool),
		}
		return
	}
}

func (manager *EventsManager) DeleteRouter(eventsType string) (have bool) {
	router, have := manager.eventsRouters[eventsType]
	if have {
		router.Dispose()
		return true
	} else {
		return false
	}
}

func (manager *EventsManager) CreateAdapter(service, eventsType string) (adapter *RouterAdapter, err error) {
	router, _ := manager.GetRouter(eventsType)
	adapterChannel, created := router.getChannel(service)
	if created {
		adapter = &RouterAdapter{
			serviceName:      service,
			manager:          manager,
			router:           router,
			listeningChannel: adapterChannel,
			stop:             make(chan bool),
		}
	} else {
		err = fmt.Errorf("EventsManager: %v adapter for the service '%v' already exists", eventsType, service)
	}
	return

}
