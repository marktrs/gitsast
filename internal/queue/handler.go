package queue

import (
	"context"
	"os"

	"github.com/labstack/gommon/log"

	"github.com/go-redis/redis/v8"

	"github.com/marktrs/gitsast/internal/queue/task"
	"github.com/vmihailenco/taskq/v3"
	"github.com/vmihailenco/taskq/v3/redisq"
)

type Handler interface {
	StartConsumer() error
	AddTask(t *taskq.Message) error
}

var (
	AnalyzeTask = taskq.RegisterTask(&taskq.TaskOptions{
		Name: "analyzer",
		Handler: func(id string) error {
			task.NewAnalyzeTask().Start(id)
			return nil
		},
	})
)

type handler struct {
	redis        *redis.Client
	queueFactory taskq.Factory
	mainQueue    taskq.Queue
}

func NewHandler() Handler {
	redis := redis.NewClient(&redis.Options{
		Addr: ":6379",
	})

	queueFactory := redisq.NewFactory()

	mainQueue := queueFactory.RegisterQueue(&taskq.QueueOptions{
		Name:  "main-queue",
		Redis: redis,
	})

	return &handler{redis, queueFactory, mainQueue}
}

func (h *handler) AddTask(t *taskq.Message) error {
	return h.mainQueue.Add(t)
}

func (h *handler) StartConsumer() error {
	log.Info("starting queue consumer")
	queueConsumerErr := make(chan error)

	go func() {
		queueConsumerErr <- h.queueFactory.StartConsumers(context.Background())
	}()

	shutdown := make(chan os.Signal, 1)

	select {
	case err := <-queueConsumerErr:
		return err
	case sig := <-shutdown:
		log.Infof("stopping queue consumer with signal: %v", sig)
		err := h.queueFactory.StopConsumers()
		if err != nil {
			return err
		}

		err = h.queueFactory.Close()
		if err != nil {
			return err
		}

		log.Infof("stopped queue consumer with signal: %v", sig)
	}

	return nil
}
