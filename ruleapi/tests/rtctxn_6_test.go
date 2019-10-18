package tests

import (
	"context"
	"testing"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"

	"github.com/stretchr/testify/assert"
)

//Same as Test_T5, but in 3rd rtc, assert a TTL=0 based and a TTL=1 based
func Test_T6(t *testing.T) {

	rs, err := createRuleSession()
	assert.Nil(t, err)

	rule := ruleapi.NewRule("R6")
	err = rule.AddCondition("R6_c1", []string{"t1.none"}, trueCondition, nil)
	assert.Nil(t, err)
	rule.SetActionService(createActionServiceFromFunction(t, r6_action))
	rule.SetPriority(1)
	err = rs.AddRule(rule)
	assert.Nil(t, err)
	t.Logf("Rule added: [%s]\n", rule.GetName())

	txnCtx := txnCtx{t, 0}
	rs.RegisterRtcTransactionHandler(t6Handler, &txnCtx)
	err = rs.Start(nil)
	assert.Nil(t, err)

	i1, err := model.NewTupleWithKeyValues("t1", "t10")
	assert.Nil(t, err)
	err = rs.Assert(context.TODO(), i1)
	assert.Nil(t, err)

	i2, err := model.NewTupleWithKeyValues("t1", "t11")
	assert.Nil(t, err)
	err = rs.Assert(context.TODO(), i2)
	assert.Nil(t, err)

	i3, err := model.NewTupleWithKeyValues("t1", "t13")
	assert.Nil(t, err)
	err = rs.Assert(context.TODO(), i3)
	assert.Nil(t, err)
	deleteRuleSession(t, rs, i1, i3)

}

func r6_action(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t1 := tuples[model.TupleType("t1")].(model.MutableTuple)
	id, _ := t1.GetString("id")
	if id == "t11" {
		tk, _ := model.NewTupleKeyWithKeyValues("t1", "t10")
		t10 := rs.GetAssertedTuple(ctx, tk).(model.MutableTuple)
		if t10 != nil {
			t10.SetString(ctx, "p3", "v3")
			t10.SetDouble(ctx, "p2", 11.11)
		}
	} else if id == "t13" {
		//delete t11
		tk, _ := model.NewTupleKeyWithKeyValues("t1", "t11")
		t11 := rs.GetAssertedTuple(ctx, tk).(model.MutableTuple)
		if t11 != nil {
			rs.Delete(ctx, t11)
		}

		i4, _ := model.NewTupleWithKeyValues("t2", "t21")
		rs.Assert(ctx, i4)

		i5, _ := model.NewTupleWithKeyValues("t1", "t15")
		rs.Assert(ctx, i5)
	}
}

func t6Handler(ctx context.Context, rs model.RuleSession, rtxn model.RtcTxn, handlerCtx interface{}) {
	if done {
		return
	}

	txnCtx := handlerCtx.(*txnCtx)
	txnCtx.TxnCnt++
	t := txnCtx.Testing
	if txnCtx.TxnCnt == 1 {
		lA := len(rtxn.GetRtcAdded())
		if lA != 1 {
			t.Errorf("RtcAdded: Expected [%d], got [%d]\n", 1, lA)
			printTuples(t, "Added", rtxn.GetRtcAdded())
		}
		lM := len(rtxn.GetRtcModified())
		if lM != 0 {
			t.Errorf("RtcModified: Expected [%d], got [%d]\n", 0, lM)
			printModified(t, rtxn.GetRtcModified())
		}
		lD := len(rtxn.GetRtcDeleted())
		if lD != 0 {
			t.Errorf("RtcDeleted: Expected [%d], got [%d]\n", 0, lD)
			printTuples(t, "Deleted", rtxn.GetRtcDeleted())
		}
	} else if txnCtx.TxnCnt == 2 {
		lA := len(rtxn.GetRtcAdded())
		if lA != 1 {
			t.Errorf("RtcAdded: Types expected [%d], got [%d]\n", 1, lA)
			printTuples(t, "Added", rtxn.GetRtcAdded())
		}
		lM := len(rtxn.GetRtcModified())
		if lM != 1 {
			t.Errorf("RtcModified: Expected [%d], got [%d]\n", 1, lM)
			printModified(t, rtxn.GetRtcModified())
		}
		lD := len(rtxn.GetRtcDeleted())
		if lD != 0 {
			t.Errorf("RtcDeleted: Expected [%d], got [%d]\n", 0, lD)
			printTuples(t, "Deleted", rtxn.GetRtcDeleted())
		}
	} else if txnCtx.TxnCnt == 3 {
		lA := len(rtxn.GetRtcAdded())
		if lA != 1 {
			t.Errorf("RtcAdded: Types expected [%d], got [%d]\n", 1, lA)
			printTuples(t, "Added", rtxn.GetRtcAdded())
		} else {
			added, _ := rtxn.GetRtcAdded()["t1"]
			if len(added) != 2 {
				t.Errorf("RtcAdded: Tuples expected [%d], got [%d]\n", 2, len(added))
				printTuples(t, "Added", rtxn.GetRtcAdded())
			}
		}
		lM := len(rtxn.GetRtcModified())
		if lM != 0 {
			t.Errorf("RtcModified: Expected [%d], got [%d]\n", 0, lM)
			printModified(t, rtxn.GetRtcModified())
		}
		lD := len(rtxn.GetRtcDeleted())
		if lD != 1 {
			t.Errorf("RtcDeleted: Expected [%d], got [%d]\n", 1, lD)
			printTuples(t, "Deleted", rtxn.GetRtcDeleted())
		}
	}
}
