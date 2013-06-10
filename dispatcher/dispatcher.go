package dispatcher

import (
	"builder/site"
	"fmt"
)

type Dispatcher struct {
	site *site.Site
}

func NewDispatcher() (*Dispatcher, error) {
	site, err := site.NewSite()
	if err != nil {
		return nil, err
	}
	return &Dispatcher{site: site}, nil
}

func (dispatcher *Dispatcher) Dispatch() error {
	errs := dispatcher.site.BuildSite()
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err.Error())
		}
	}
	server := dispatcher.site.MakeWebServer()
	err := server.Serve()
	if err != nil {
		return err
	}
	return nil
}
