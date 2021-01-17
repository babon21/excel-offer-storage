package usecase

import (
	"encoding/json"
	"github.com/babon21/excel-offer-storage/internal/offer/usecase"
	"github.com/babon21/excel-offer-storage/pkg/delivery/http/api"
)

type AsyncOfferUseCase interface {
	GetTask(taskId int64) string
	Store(sellerId string, url string) (int64, error)
}

type offerUseCase struct {
	syncOfferUseCase usecase.OfferUseCase
	taskStore        Store
}

func (useCase *offerUseCase) GetTask(taskId int64) string {
	task, _ := useCase.taskStore.Get(taskId)
	return task
}

func (useCase *offerUseCase) Store(sellerId string, url string) (int64, error) {
	id, err := useCase.taskStore.GetNewId("task_id")
	if err != nil {
		// TODO
	}

	response := api.GetTaskResponse{
		Status: "waiting",
	}

	jsonResponse, err := json.MarshalIndent(response, "", "    ")

	if err := useCase.taskStore.Set(id, string(jsonResponse)); err != nil {
		// TODO
	}

	go func() {
		statistic, err := useCase.syncOfferUseCase.Store(sellerId, url)
		status := "done"

		if err != nil {
			status = "failed"
		}

		response := api.GetTaskResponse{
			Status:    status,
			Statistic: statistic,
		}
		jsonResponse, err := json.MarshalIndent(response, "", "    ")
		if err != nil {

		}

		if err := useCase.taskStore.Set(id, string(jsonResponse)); err != nil {
			// TODO ???
		}
	}()

	return id, nil
}

func NewAsyncOfferUseCase(syncOfferUseCase usecase.OfferUseCase, taskStore Store) AsyncOfferUseCase {
	return &offerUseCase{
		syncOfferUseCase: syncOfferUseCase,
		taskStore:        taskStore,
	}
}
