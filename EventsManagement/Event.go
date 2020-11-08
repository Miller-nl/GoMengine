package EventsManagement

type Event struct {
	index         int // event identification index
	previousEvent int // previous event identification index

	router *SimilarEventsRouter // сделать функцию "запроса класса", и "получения имени класса"

	data interface{} // any data. Services should not change data. The work logic is read-only.
	// Don't use links! Links will not be saved in the event of a system crash.

}

func (event *Event) Index() int {
	return event.index
}

func (event *Event) PreviousEventIndex() int {
	return event.previousEvent
}

func (event *Event) Router() *SimilarEventsRouter {
	return event.router
}

func (event *Event) Class() string {
	return event.router.EventsClass()
}
