package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/myntra/pipeline"
	"io/ioutil"
	"time"
)

type ServicingStep struct {
	pipeline.StepContext
	ID   int
	Name string
}

// UnmarshalJSON sets *m to a copy of data.
func (ss *ServicingStep) UnmarshalJSON(data []byte) error {
	logger.Info("Calling UnmarshalJSON...")
	var objmap map[string]interface{}
	err := json.Unmarshal(data, &objmap)
	if err != nil {
		return errors.New("ServicingStep: error unmarshalling raw json")
	}
	logger.Info("%+v", objmap)
	// ss = &ServicingStep{ID: *objmap["ID"], Name: *objmap["Name"]}
	return nil
}

func (ss ServicingStep) Exec(request *pipeline.Request) *pipeline.Result {
	ss.Status(fmt.Sprintf("%+v", request))

	duration := time.Duration(1000 * ss.ID)
	time.Sleep(time.Millisecond * duration)
	msg := fmt.Sprintf("work %v", ss.Name)

	return &pipeline.Result{
		Error:  nil,
		Data:   struct{ msg string }{msg: msg},
		KeyVal: map[string]interface{}{"msg": msg},
	}
}

func (ss ServicingStep) Cancel() error {
	ss.Status("cancel step")
	return nil
}

func readPipeline(pipe *pipeline.Pipeline) {
	logger.Info("Calling readPipeline()...")
	out, err := pipe.Out()
	if err != nil {
		return
	}

	progress, err := pipe.GetProgressPercent()
	if err != nil {
		return
	}

	for {
		select {
		case line := <-out:
			logger.Info(line)
		case p := <-progress:
			logger.Info("percent done: ", p)
		}
	}
}

func ExecuteProcess() {
	// workpipe := readProcessDefinition()
	// log.Info("%+v", workpipe)
	workpipe := buildPipeline()
	go readPipeline(workpipe)
	result := workpipe.Run()
	if result.Error != nil {
		logger.Info(result.Error)
	}

	logger.Info("timeTaken:", workpipe.GetDuration())

}

func readProcessDefinition() (pipe *pipeline.Pipeline) {
	fileName := "new_loan_process.json"
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		logger.Errorf("Error reading process definition file %v : %v\n", fileName, err)
	}
	err = json.Unmarshal(bytes, &pipe)
	if err != nil {
		logger.Errorf("Error unmarshaling process definition : %v\n", err)
	}
	return
}

func buildPipeline() (workpipe *pipeline.Pipeline) {
	// create a new pipeline
	workpipe = pipeline.NewProgress("NewLoan", 1000, time.Second*3)
	// func NewStage(name string, concurrent bool, disableStrictMode bool) *Stage
	// To execute steps concurrently, set concurrent=true.
	loanStage := pipeline.NewStage("LoanStage", false, false)

	// a unit of work
	step1 := &ServicingStep{ID: 1, Name: "Loan Transition to ISSUE"}
	// another unit of work
	step2 := &ServicingStep{ID: 2, Name: "Save Loan"}

	// add the steps to the stage. Since concurrent is set false above. The steps will be
	// executed one after the other.
	loanStage.AddStep(step1)
	loanStage.AddStep(step2)

	paymentStage := pipeline.NewStage("PaymentStage", true, false)

	// a unit of work
	step3 := &ServicingStep{ID: 3, Name: "Create Payment"}
	// another unit of work
	step4 := &ServicingStep{ID: 4, Name: "Send Payment"}

	// add the steps to the stage. Since concurrent is set false above. The steps will be
	// executed one after the other.
	paymentStage.AddStep(step3)
	paymentStage.AddStep(step4)

	// add the stage to the pipe.
	workpipe.AddStage(loanStage)
	workpipe.AddStage(paymentStage)

	writePipelineToFile(workpipe)
	return
}

func writePipelineToFile(workpipe *pipeline.Pipeline) {
	raw, err := json.Marshal(workpipe)
	if err != nil {
		logger.Info(err)
		return
	}
	err = ioutil.WriteFile("new_loan_process.json", raw, 0644)
	if err != nil {
		logger.Error("Error creating process definition json due to: ", err)
	}
}
