package policy

import (
	"strings"

	"github.com/dlshle/gommon/errors"

	"github.com/dlshle/authnz/internal/group"
	pb "github.com/dlshle/authnz/proto"
)

type Engine interface {
	Check(policy *pb.Policy, group *pb.Group, ctx *pb.ContextProperty) (pb.Verdict, error)
}

type conditionProcessor = func(cond *pb.PolicyCondition, group group.Group, ctx *pb.ContextProperty) (pb.Verdict, error)

type engine struct {
	conditionProcessors []conditionProcessor
}

func NewEngine() Engine {
	return &engine{}
}

func (e *engine) init() {
	e.conditionProcessors = []conditionProcessor{
		e.HasAttributeProcessor,
		e.EvaluateOPProcessor,
		e.NegationProcessor,
		e.AndProcessor,
		e.OrProcessor,
	}
}

func (e *engine) Check(policy *pb.Policy, pbGroup *pb.Group, ctx *pb.ContextProperty) (pb.Verdict, error) {
	if cond := policy.GetCondition(); cond != nil {
		return e.evaluateCondition(cond, group.FromPB(pbGroup), ctx)
	}
	return pb.Verdict_UNKNOWN, errors.Error("empty condition for policy " + policy.GetId())
}

func (e *engine) evaluateCondition(cond *pb.PolicyCondition, group group.Group, ctx *pb.ContextProperty) (verdict pb.Verdict, err error) {
	verdict = pb.Verdict_UNKNOWN
	for _, processor := range e.conditionProcessors {
		verdict, err = processor(cond, group, ctx)
		if err != nil {
			return pb.Verdict_UNKNOWN, err
		}
		if verdict == pb.Verdict_DENIED || verdict == pb.Verdict_PERMITTED {
			return verdict, nil
		}
	}
	err = errors.Error("no condition processor is found for policy condition " + cond.String())
	return
}

func (e *engine) HasAttributeProcessor(cond *pb.PolicyCondition, group group.Group, ctx *pb.ContextProperty) (pb.Verdict, error) {
	hasAttributeCond := cond.GetHasAttribute()
	if hasAttributeCond == nil {
		return pb.Verdict_UNKNOWN, nil
	}
	for _, k := range hasAttributeCond.GetAttributeKey() {
		if _, hasAttribute := group.Attributes[k]; !hasAttribute {
			return pb.Verdict_DENIED, nil
		}
	}
	return pb.Verdict_PERMITTED, nil
}

func (e *engine) EvaluateOPProcessor(cond *pb.PolicyCondition, group group.Group, ctx *pb.ContextProperty) (verdict pb.Verdict, err error) {
	evaluateCond := cond.GetEvaluateAttribute()
	if evaluateCond == nil {
		return pb.Verdict_UNKNOWN, nil
	}
	attribute, exists := group.Attributes[evaluateCond.GetAttributeKey()]
	if !exists {
		// attribute DNE
		return pb.Verdict_DENIED, nil
	}
	if evalOp(evaluateCond.GetOp(), attribute, evaluateCond.GetValue()) {
		return pb.Verdict_PERMITTED, nil
	}
	return pb.Verdict_DENIED, nil
}

func evalOp(op pb.Operation, attribute, value string) bool {
	switch op {
	case pb.Operation_EQ:
		return attribute == value
	case pb.Operation_CONTAINS:
		return strings.Contains(attribute, value)
	case pb.Operation_GT:
		return strings.Compare(attribute, value) > 0
	case pb.Operation_LT:
		return strings.Compare(attribute, value) < 0
	case pb.Operation_GTE:
		return strings.Compare(attribute, value) >= 0
	case pb.Operation_LTE:
		return strings.Compare(attribute, value) <= 0
	default:
		return false
	}
}

func (e *engine) NegationProcessor(cond *pb.PolicyCondition, group group.Group, ctx *pb.ContextProperty) (verdict pb.Verdict, err error) {
	negationCond := cond.GetNegation()
	if negationCond == nil {
		return pb.Verdict_UNKNOWN, nil
	}
	innerCond := negationCond.GetCondition()
	if innerCond == nil {
		return pb.Verdict_UNKNOWN, errors.Error("inner condition is null for negation")
	}
	verdict, err = e.evaluateCondition(innerCond, group, ctx)
	// only need to negate when we have a clear result from innser cond evaluation
	switch verdict {
	case pb.Verdict_PERMITTED:
		verdict = pb.Verdict_DENIED
		break
	case pb.Verdict_DENIED:
		verdict = pb.Verdict_PERMITTED
		break
	}
	return
}

func (e *engine) AndProcessor(cond *pb.PolicyCondition, group group.Group, ctx *pb.ContextProperty) (verdict pb.Verdict, err error) {
	andCond := cond.GetAnd()
	if andCond == nil {
		return pb.Verdict_UNKNOWN, nil
	}
	for _, innerCond := range andCond.GetCondition() {
		if verdict, err = e.evaluateCondition(innerCond, group, ctx); err != nil || verdict == pb.Verdict_DENIED {
			return
		}
	}
	return pb.Verdict_PERMITTED, nil
}

func (e *engine) OrProcessor(cond *pb.PolicyCondition, group group.Group, ctx *pb.ContextProperty) (verdict pb.Verdict, err error) {
	orCond := cond.GetAnd()
	if orCond == nil {
		return pb.Verdict_UNKNOWN, nil
	}
	for _, innerCond := range orCond.GetCondition() {
		if verdict, err = e.evaluateCondition(innerCond, group, ctx); err != nil || verdict == pb.Verdict_PERMITTED {
			return
		}
	}
	return pb.Verdict_PERMITTED, nil
}

// TODO: other processors
