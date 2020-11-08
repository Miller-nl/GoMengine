package EventsManagement

import (
	"sync"
)

type SimilarEventsRouter struct {
	manager     *EventsManager // events manager
	eventsClass string         // events class name and router name

	income   chan *Event               // income events channel
	adapters map[string]*RouterAdapter // adapters
	active   bool                      // is router active
	mutex    sync.RWMutex              //

	stop chan bool // channel for stop
}

func (router *SimilarEventsRouter) Name() string {
	return router.eventsClass
}

func (router *SimilarEventsRouter) EventsClass() string {
	return router.eventsClass
}

func (router *SimilarEventsRouter) IsActive() bool {
	router.mutex.Lock()
	defer router.mutex.Unlock()
	return router.active
}

func (router *SimilarEventsRouter) routeEvent(event *Event, service string) {
	router.income <- event
	return
}

func (router *SimilarEventsRouter) AdaptersServiceNames() (names []string) {
	for name, _ := range router.adapters {
		names = append(names, name)
	}
	return
}

func (router *SimilarEventsRouter) getChannel(service string) (channel chan *Event, created bool) {
	router.mutex.Lock()
	defer router.mutex.Unlock()

	value, created := router.adapters[service]
	if created {
		channel = value.listeningChannel
	} else {
		channel = make(chan *Event, 100)
	}
	return
}

func (router *SimilarEventsRouter) addAdapter(adapter *RouterAdapter) {
	router.mutex.Lock()
	defer router.mutex.Unlock()

	router.adapters[adapter.ServiceName()] = adapter

	return
}

func (router *SimilarEventsRouter) deleteAdapter(service string) {
	router.mutex.Lock()
	defer router.mutex.Unlock()

	_, ok := router.adapters[service]
	if ok {
		delete(router.adapters, service)
	}
	return
}

func (router *SimilarEventsRouter) handleEvent(event *Event) {
	router.mutex.Lock()
	defer router.mutex.Unlock()

	for _, adapter := range router.adapters {
		adapter.listeningChannel <- event
	}
	return
}

// Each event is assumed to have a unique index!
func (router *SimilarEventsRouter) rout() {
	for {
		select {
		case <-router.stop:
			return
		case event := <-router.income:
			go router.handleEvent(event)
		}
	}
}

func (router *SimilarEventsRouter) Start() {
	router.mutex.Lock()
	defer router.mutex.Unlock()

	// notify manager "router started"
	go router.rout()
	router.active = true
	return
}

func (router *SimilarEventsRouter) Stop() {
	// notify manager "router stopped"
	router.mutex.Lock()
	defer router.mutex.Unlock()

	router.stop <- true
	router.active = false
}

func (router *SimilarEventsRouter) Finish() {
	router.mutex.Lock()
	defer router.mutex.Unlock()

	// notify manager "router finished"
	router.stop <- true
	router.active = false
}

func (router *SimilarEventsRouter) Dispose() {
	router.mutex.Lock()
	defer router.mutex.Unlock()

	// notify manager "router finished"
	router.stop <- true
	for _, adapter := range router.adapters {
		adapter.Stop()
		adapter.Dispose()
	}
	close(router.income)
	close(router.stop)
	router.active = false
	return
}
