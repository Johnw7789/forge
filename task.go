package main

import (
	forge "github.com/Johnw7789/forge/backend/request"
)

var steps = []string{"init", "getCreateData", "submitCreate", "pingDiscord", "phone2FA"}

type Task struct {
	amzTask forge.AmazonTask
	stopped chan bool
}

type TaskState struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Proxy    string `json:"proxy"`
	Status   string `json:"status"`
}

// * Start is in charge of the task flow, when the task is stopped, the stopped channel is closed
func (t *Task) Start() (bool, error) {
	t.stopped = make(chan bool)
	defer close(t.stopped) // Ensure stopped channel is closed when Start() exits
	for _, step := range steps {
		select {
		case <-t.stopped:
			return false, nil // Task stopped gracefully
		default:
			err := t.amzTask.DoStep(step)
			if err != nil {
				if step == "phone2FA" {
					return true, err
				}
				return false, err
			}

			if step == "phone2FA" {
				return true, nil
			}
		}
	}

	return false, nil
}

// * Close the stopped channel, the task will stop before the next step
func (t *Task) Stop() {
	if t.stopped != nil {
		close(t.stopped)
	}
}

var infoSteps = []string{"init", "submitAddress", "submitPayment", "submitProfile"}
var primeSteps = []string{"initPrime", "submitPrime"}

// * StartPrime is in charge of the prime subtask flow, when the task is stopped, the stopped channel is closed
func (t *Task) StartPrime() (bool, error) {
	t.stopped = make(chan bool)
	defer close(t.stopped) // Ensure stopped channel is closed when Start() exits
	for _, step := range primeSteps {
		select {
		case <-t.stopped:
			return false, nil // Task stopped gracefully
		default:
			err := t.amzTask.DoPrimeStep(step)
			if err != nil {
				return false, err
			}

			if step == "submitPrime" {
				return true, nil
			}
		}
	}

	return false, nil
}

// * StartInfo is in charge of the info subtask flow, when the task is stopped, the stopped channel is closed
func (t *Task) StartInfo() (bool, error) {
	t.stopped = make(chan bool)
	defer close(t.stopped) // Ensure stopped channel is closed when Start() exits
	for _, step := range infoSteps {
		select {
		case <-t.stopped:
			return false, nil // Task stopped gracefully
		default:
			err := t.amzTask.DoInfoStep(step)
			if err != nil {
				return false, err
			}

			if step == "submitProfile" {
				return true, nil
			}
		}
	}

	return false, nil
}
