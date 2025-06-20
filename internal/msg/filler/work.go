package filler

import (
	"context"
	"sync"
	"time"

	"crafa/internal/cfg"
	"crafa/internal/logger"
	"crafa/internal/repo/postgres"

	"go.uber.org/zap"
)

func Start(config *cfg.Config, wg *sync.WaitGroup, exit chan struct{}) {
	ps, _ := postgres.GetPostgresStorage(config.ConnectionString)

	// Асинхронно запускаем наполнитель очереди
	s := NewFiller(
		config.Queue,
		ps.Msg,
		config.QueueFillPeriod,
	)

	go s.StartFilling(exit, wg)
}

// StartFilling начинает наполнение очереди
func (f *Filler) StartFilling(exit chan struct{}, wg *sync.WaitGroup) {
	logger.Log.Info("StartFilling",
		zap.String("Филлер стартовал", "Успешно"))
	ctx := context.TODO()
	for {
		select {
		case _, ok := <-exit:
			if !ok {
				//Нет смысла ждать наполнения очереди, поэтому дефолт не жду
				logger.Log.Info("StartFilling",
					zap.String("Филлер изящно завершил дела", "Успешно"))
				wg.Done()
				return
			}
		default:
			msgs, err := f.msgStorager.SelectUnsendedMsgs(ctx)

			if err != nil {
				logger.Log.Error("StartFilling",
					zap.String("ошибка при получении сообщений к отправке", err.Error()),
				)
			}

			for _, v := range msgs {
				f.queue.Push(&v)
				v.Queued = true
				f.msgStorager.UpdateMsg(ctx, &v)
			}
		}
		time.Sleep(time.Duration(f.period) * time.Second)
	}
}
