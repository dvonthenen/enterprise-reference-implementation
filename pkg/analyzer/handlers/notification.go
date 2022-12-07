// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package router

import (
	"context"
	"sync"

	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"
	klog "k8s.io/klog/v2"

	"github.com/dvonthenen/enterprise-reference-implementation/pkg/analyzer/rabbit"
	rabbitinterfaces "github.com/dvonthenen/enterprise-reference-implementation/pkg/analyzer/rabbit/interfaces"
	interfaces "github.com/dvonthenen/enterprise-reference-implementation/pkg/interfaces"
)

func NewNotificationManager(options NotificationManagerOption) *NotificationManager {
	mgr := &NotificationManager{
		driver:        options.Driver,
		rabbitManager: options.RabbitManager,
	}
	return mgr
}

func (nm *NotificationManager) Init() error {
	klog.V(6).Infof("NotificationManager.Init ENTER\n")

	type InitFunc func(HandlerOptions) *rabbitinterfaces.RabbitMessageHandler
	type MyHandler struct {
		Name string
		Func InitFunc
	}

	myHandlers := make([]*MyHandler, 0)
	myHandlers = append(myHandlers, &MyHandler{
		Name: interfaces.RabbitExchangeConversation,
		Func: NewConversationHandler,
	})
	myHandlers = append(myHandlers, &MyHandler{
		Name: interfaces.RabbitExchangeEntity,
		Func: NewEntityHandler,
	})
	myHandlers = append(myHandlers, &MyHandler{
		Name: interfaces.RabbitExchangeInsight,
		Func: NewInsightHandler,
	})
	myHandlers = append(myHandlers, &MyHandler{
		Name: interfaces.RabbitExchangeMessage,
		Func: NewMessageHandler,
	})
	myHandlers = append(myHandlers, &MyHandler{
		Name: interfaces.RabbitExchangeTopic,
		Func: NewTopicHandler,
	})
	myHandlers = append(myHandlers, &MyHandler{
		Name: interfaces.RabbitExchangeTracker,
		Func: NewTrackerHandler,
	})

	// doing this concurrently because creating a neo4j session is time consuming
	var wg sync.WaitGroup
	wg.Add(len(myHandlers))

	for _, myHandler := range myHandlers {
		// create session
		ctx := context.Background()
		session := (*nm.driver).NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})

		// signal
		handler := myHandler.Func(HandlerOptions{
			Session: &session,
		})

		_, err := nm.rabbitManager.CreateSubscription(rabbit.CreateOptions{
			Name:    myHandler.Name,
			Handler: handler,
		})
		if err != nil {
			klog.V(1).Infof("CreateSubscription failed. Err: %v\n", err)
		}
	}

	klog.V(4).Infof("Init Succeeded\n")
	klog.V(6).Infof("NotificationManager.Init LEAVE\n")

	return nil
}

func (nm *NotificationManager) Start() error {
	klog.V(6).Infof("NotificationManager.Start ENTER\n")

	err := nm.rabbitManager.Start()
	if err != nil {
		klog.V(1).Infof("rabbitManager.Start failed. Err: %v\n", err)
		klog.V(6).Infof("NotificationManager.Start LEAVE\n")
		return err
	}

	klog.V(4).Infof("Start Succeeded\n")
	klog.V(6).Infof("NotificationManager.Start LEAVE\n")

	return nil
}

func (nm *NotificationManager) Stop() error {
	klog.V(6).Infof("NotificationManager.Stop ENTER\n")

	err := nm.rabbitManager.Stop()
	if err != nil {
		klog.V(1).Infof("rabbitManager.Stop failed. Err: %v\n", err)
		klog.V(6).Infof("NotificationManager.Stop LEAVE\n")
		return err
	}

	klog.V(4).Infof("Stop Succeeded\n")
	klog.V(6).Infof("NotificationManager.Stop LEAVE\n")

	return nil
}

func (nm *NotificationManager) Teardown() error {
	klog.V(6).Infof("NotificationManager.Teardown ENTER\n")

	err := nm.rabbitManager.Delete()
	if err != nil {
		klog.V(1).Infof("rabbitManager.DeleteAll failed. Err: %v\n", err)
		klog.V(6).Infof("NotificationManager.Stop LEAVE\n")
		return err
	}

	klog.V(4).Infof("Teardown Succeeded\n")
	klog.V(6).Infof("NotificationManager.Teardown LEAVE\n")

	return nil
}