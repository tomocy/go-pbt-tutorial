package main

import (
	"testing"
	"time"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/commands"
	"github.com/leanovate/gopter/gen"
)

var (
	counterCmds = &commands.ProtoCommands{
		NewSystemUnderTestFunc: func(initialState commands.State) commands.SystemUnderTest {
			return &counter{}
		},
		InitialStateGen: gen.Const(0),
		InitialPreConditionFunc: func(s commands.State) bool {
			return s.(int) == 0
		},
		GenCommandFunc: func(s commands.State) gopter.Gen {
			return gen.OneConstOf(incCmd, decCmd, nCmd, resetCmd)
		},
	}
	incCmd = &commands.ProtoCommand{
		Name: "inc",
		RunFunc: func(t commands.SystemUnderTest) commands.Result {
			t.(*counter).inc()
			return nil
		},
		NextStateFunc: func(s commands.State) commands.State {
			return s.(int) + 1
		},
	}
	decCmd = &commands.ProtoCommand{
		Name: "dec",
		RunFunc: func(t commands.SystemUnderTest) commands.Result {
			t.(*counter).dec()
			return nil
		},
		NextStateFunc: func(s commands.State) commands.State {
			return s.(int) - 1
		},
	}
	nCmd = &commands.ProtoCommand{
		Name: "n",
		RunFunc: func(t commands.SystemUnderTest) commands.Result {
			return t.(*counter).n
		},
		PostConditionFunc: func(s commands.State, r commands.Result) *gopter.PropResult {
			if s.(int) != r.(int) {
				return &gopter.PropResult{
					Status: gopter.PropFalse,
				}
			}
			return &gopter.PropResult{
				Status: gopter.PropTrue,
			}
		},
	}
	resetCmd = &commands.ProtoCommand{
		Name: "reset",
		RunFunc: func(t commands.SystemUnderTest) commands.Result {
			t.(*counter).reset()
			return nil
		},
		NextStateFunc: func(s commands.State) commands.State {
			return 0
		},
	}
)

func TestBuggyCounter(t *testing.T) {
	if testing.Short() {
		t.Skip("skip PBT in short mode")
	}

	params := gopter.DefaultTestParameters()
	params.Rng.Seed(time.Now().UnixNano())

	props := gopter.NewProperties(params)
	props.Property("counter", commands.Prop(counterCmds))

	props.TestingRun(t)
}
