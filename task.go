package task

import (
	"errors"
	"runtime"
	"sync"
	"time"

	"github.com/sanrentai/snowflake"
)

const (
	StateWaiting   = "waiting"
	StateCompleted = "completed"
	StateError     = "failed"
	StateNone      = "none"
	StateOverdue   = "overdue"
)

var (
	taskStateMutex sync.Mutex
	taskPool       sync.Pool
	taskChan       = make(chan *Task, runtime.NumCPU())
	taskState      = make(map[int64]string)
)

var idutil *snowflake.Snowflake

func init() {
	idutil, _ = snowflake.New(500)
}

type Task struct {
	ID         int64
	param      map[string]interface{}
	fac        FacFunc
	expiration int64
	resultChan chan interface{}
	errorChan  chan error
}

type FacFunc func(map[string]interface{}) (interface{}, error)

func (task *Task) GetResult() (interface{}, error) {
	state := GetTaskState(task.ID)
	switch state {
	case StateWaiting, StateCompleted:
		return <-task.resultChan, nil
	case StateError:
		return nil, <-task.errorChan
	case StateOverdue:
		return nil, errors.New("overtime")
	}
	// task.Close()
	return nil, nil
}

func (task *Task) Close() {
	delete(taskState, task.ID)
	close(task.errorChan)
	close(task.resultChan)
}

func NewTask(param map[string]interface{}, factory FacFunc, d time.Duration) *Task {

	var expiration int64
	if d > 0 {
		expiration = time.Now().Add(d).UnixNano()
	} else {
		expiration = -1
	}

	t := taskPool.Get()
	if t == nil {
		return &Task{
			ID:         idutil.Generate().Int64(),
			param:      param,
			fac:        factory,
			expiration: expiration,
			resultChan: make(chan interface{}, 1),
			errorChan:  make(chan error, 1),
		}
	} else {
		task := t.(*Task)
		task.param = param
		task.fac = factory
		task.ID = idutil.Generate().Int64()
		task.expiration = expiration
		task.resultChan = make(chan interface{}, 1)
		task.errorChan = make(chan error, 1)
		return task
	}

}

func (task *Task) Start() int64 {
	UpdateTaskState(task.ID, StateWaiting)
	go func() {
		taskChan <- task
	}()
	return task.ID
}

func UpdateTaskState(id int64, state string) {
	taskStateMutex.Lock()
	defer taskStateMutex.Unlock()

	taskState[id] = state
}

func GetTaskState(id int64) (state string) {
	taskStateMutex.Lock()
	defer taskStateMutex.Unlock()

	resultState, exists := taskState[id]
	if !exists {
		state = StateNone
	} else {
		state = resultState
	}
	return
}

func taskReceiver() {
	for {
		task := <-taskChan
		if (task.expiration > 0 && time.Now().UnixNano() < task.expiration) || task.expiration < 0 {

			result, err := task.fac(task.param)

			if err != nil {
				task.errorChan <- err
				UpdateTaskState(task.ID, StateError)
				taskPool.Put(task)
			} else {

				task.resultChan <- result
				UpdateTaskState(task.ID, StateCompleted)
				taskPool.Put(task)

			}
		} else {
			UpdateTaskState(task.ID, StateOverdue)
			taskPool.Put(task)
		}
	}
}

func InitTaskReceiver(num int) {
	for i := 0; i < num; i++ {
		go taskReceiver()
	}
}
