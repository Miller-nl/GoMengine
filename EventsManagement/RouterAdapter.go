package EventsManagement

type RouterAdapter struct {
	serviceName string               // adapter service name
	manager     *EventsManager       // main manager
	router      *SimilarEventsRouter // adapter events router

	listeningChannel chan *Event // channel for income events

	stop chan bool // channel for stop
}

func (adapter *RouterAdapter) ServiceName() string {
	return adapter.serviceName
}

func (adapter *RouterAdapter) EventsClass() string {
	return adapter.router.EventsClass()
}

func (adapter *RouterAdapter) NewEvent(data interface{}) int {
	index := adapter.manager.EventIndex()
	event := &Event{
		index:  index,
		router: adapter.router,
		data:   data,
	}
	adapter.router.income <- event
	return index
}

func (adapter *RouterAdapter) NextEvent(previous *Event, data interface{}) int {
	index := adapter.manager.EventIndex()
	event := &Event{
		index:         index,
		previousEvent: previous.index,
		router:        adapter.router,
		data:          data,
	}
	adapter.router.routeEvent(event, adapter.ServiceName())
	return index
}

func (adapter *RouterAdapter) Listen() (event *Event) {
	select {
	case event = <-adapter.listeningChannel:
		return
	case <-adapter.stop:
		return
	}
}

func (adapter *RouterAdapter) ListenNext(previousIndex int) (event *Event) {
	for {
		select {
		case income_event := <-adapter.listeningChannel:
			if income_event.PreviousEventIndex() == previousIndex {
				event = income_event
			}
		case <-adapter.stop:
			return
		}
	}
}

func (adapter *RouterAdapter) Stop() {
	adapter.stop <- true
	return
}

// stop and delete adapter
func (adapter *RouterAdapter) Dispose() {
	adapter.Stop()
	close(adapter.listeningChannel)
	close(adapter.stop)
	adapter.router.deleteAdapter(adapter.ServiceName())
	return
}
