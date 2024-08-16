package event

// todo: add context?

type Handler func(event IEvent)

type handlersByType struct {
	all            []Handler
	groupMessage   []Handler
	privateMessage []Handler
	notice         []Handler
	meta           []Handler
	request        []Handler
}

type Dispatcher struct {
	isGoroutineMode bool
	handlers        handlersByType
}

func NewDispatcher(isGoroutineMode bool) *Dispatcher {
	return &Dispatcher{
		isGoroutineMode: isGoroutineMode,
		handlers:        handlersByType{},
	}
}

func (d *Dispatcher) RegisterHandlerAllTypes(handler Handler) {
	d.handlers.all = append(d.handlers.all, handler)
}

func (d *Dispatcher) RegisterHandlerGroupMessage(handler Handler) {
	d.handlers.groupMessage = append(d.handlers.groupMessage, handler)
}

func (d *Dispatcher) RegisterHandlerPrivateMessage(handler Handler) {
	d.handlers.privateMessage = append(d.handlers.privateMessage, handler)
}

func (d *Dispatcher) RegisterHandlerNotice(handler Handler) {
	d.handlers.notice = append(d.handlers.notice, handler)
}

func (d *Dispatcher) RegisterHandlerMeta(handler Handler) {
	d.handlers.meta = append(d.handlers.meta, handler)
}

func (d *Dispatcher) RegisterHandlerRequest(handler Handler) {
	d.handlers.request = append(d.handlers.request, handler)
}

func (d *Dispatcher) Dispatch(event IEvent) {
	for _, handler := range d.handlers.all {
		if d.isGoroutineMode {
			go handler(event)
			continue
		}
		handler(event)
		if event.isDefaultPrevented() {
			return
		}
	}

	switch event.GetEventType() {
	case EventTypeMessage:
		met := event.(IMessageEvent).GetMessageEventType()
		if met == MessageEventTypePrivate {
			for _, handler := range d.handlers.privateMessage {
				if d.isGoroutineMode {
					go handler(event)
					continue
				}
				handler(event)
				if event.isDefaultPrevented() {
					return
				}
			}
		} else if met == MessageEventTypeGroup {
			for _, handler := range d.handlers.groupMessage {
				if d.isGoroutineMode {
					go handler(event)
					continue
				}
				handler(event)
				if event.isDefaultPrevented() {
					return
				}
			}
		}
	case EventTypeNotice:
		for _, handler := range d.handlers.notice {
			if d.isGoroutineMode {
				go handler(event)
				continue
			}
			handler(event)
			if event.isDefaultPrevented() {
				return
			}
		}
	case EventTypeMeta:
		for _, handler := range d.handlers.meta {
			if d.isGoroutineMode {
				go handler(event)
				continue
			}
			handler(event)
			if event.isDefaultPrevented() {
				return
			}
		}
	case EventTypeRequest:
		for _, handler := range d.handlers.request {
			if d.isGoroutineMode {
				go handler(event)
				continue
			}
			handler(event)
			if event.isDefaultPrevented() {
				return
			}
		}
	}
}
