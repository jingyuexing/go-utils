package utils


import "sync"

func Emit[T any]() (handler func(name string, args ...T), listener func(name string, callback func(args ...T)), errorHandler func(err error)) {
    eventList := make(map[string][]func(args ...T))
    handler = func(name string, args ...T) {
        err := recover()
        for _, call := range eventList[name] {
            call(args...)
        }
        if err != nil {
            errorHandler(err.(error))
        }
    }

    listener = func(name string, callback func(args ...T)) {
        eventList[name] = append(eventList[name], callback)
    }

    return handler, listener, errorHandler
}

type EventEmit struct {
    mutex  sync.Mutex
    events map[string][](func(args ...any))
}

func NewEventEmit() *EventEmit {
    return &EventEmit{
        mutex:  sync.Mutex{},
        events: make(map[string][]func(args ...any)),
    }
}

func (e *EventEmit) On(name string, callback func(args ...any)) {
    e.mutex.Lock()
    e.events[name] = append(e.events[name], callback)
    defer e.mutex.Unlock()
}

func (e *EventEmit) AddEvent(name string, args ...any) {
    e.mutex.Lock()
    for _, call := range e.events[name] {
        call(args...)
    }
    defer e.mutex.Unlock()
}
