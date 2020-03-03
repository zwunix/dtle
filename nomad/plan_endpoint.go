package nomad

import (
	"time"

	metrics "github.com/armon/go-metrics"
	log "github.com/hashicorp/go-hclog"

	"github.com/actiontech/dtle/nomad/structs"
)

// Plan endpoint is used for plan interactions
type Plan struct {
	srv    *Server
	logger log.Logger
}

// Submit is used to submit a plan to the leader
func (p *Plan) Submit(args *structs.PlanRequest, reply *structs.PlanResponse) error {
	if done, err := p.srv.forward("Plan.Submit", args, args, reply); done {
		return err
	}
	defer metrics.MeasureSince([]string{"nomad", "plan", "submit"}, time.Now())

	// Pause the Nack timer for the eval as it is making progress as long as it
	// is in the plan queue. We resume immediately after we get a result to
	// handle the case that the receiving worker dies.
	plan := args.Plan
	id := plan.EvalID
	token := plan.EvalToken
	if err := p.srv.evalBroker.PauseNackTimeout(id, token); err != nil {
		return err
	}
	defer p.srv.evalBroker.ResumeNackTimeout(id, token)

	// Submit the plan to the queue
	future, err := p.srv.planQueue.Enqueue(plan)
	if err != nil {
		return err
	}

	// Wait for the results
	result, err := future.Wait()
	if err != nil {
		return err
	}

	// Package the result
	reply.Result = result
	reply.Index = result.AllocIndex
	return nil
}
