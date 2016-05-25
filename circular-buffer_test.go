package buffer

import (
	"container/list"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/commands"
	"github.com/leanovate/gopter/gen"
)

var getCommand = &commands.ProtoCommand{
	Name: "Get",
	RunFunc: func(q commands.SystemUnderTest) commands.Result {
		return q.(*Queue).Get()
	},
	NextStateFunc: func(state commands.State) commands.State {
		st := state.(*cbCommands)
		st.elements.Remove(st.elements.Front())
		return st
	},
	PreConditionFunc: func(state commands.State) bool {
		return state.(*cbCommands).elements.Len() > 0
	},
	PostConditionFunc: func(state commands.State, result commands.Result) *gopter.PropResult {
		if result.(int) != 1 {
			return &gopter.PropResult{Status: gopter.PropFalse}
		}
		return &gopter.PropResult{Status: gopter.PropTrue}
	},
}

var putCommand = &commands.ProtoCommand{
	Name: "Put",
	RunFunc: func(q commands.SystemUnderTest) commands.Result {
		return q.(*Queue).Put(1)
	},
	NextStateFunc: func(state commands.State) commands.State {
		st := state.(*cbCommands)
		st.elements.PushBack(1)
		return st
	},
	PreConditionFunc: func(state commands.State) bool {
		s := state.(*cbCommands)
		return s.elements.Len() < s.size
	},
	PostConditionFunc: func(state commands.State, result commands.Result) *gopter.PropResult {
		if result.(int) != state.(*cbCommands).elements.Back().Value.(int) {
			return &gopter.PropResult{Status: gopter.PropFalse}
		}
		return &gopter.PropResult{Status: gopter.PropTrue}
	},
}

var sizeCommand = &commands.ProtoCommand{
	Name: "Size",
	RunFunc: func(q commands.SystemUnderTest) commands.Result {
		return q.(*Queue).Size()
	},
	PreConditionFunc: func(state commands.State) bool {
		_, ok := state.(*cbCommands)
		return ok
	},
	PostConditionFunc: func(state commands.State, result commands.Result) *gopter.PropResult {
		if result.(int) != state.(*cbCommands).elements.Len() {
			return &gopter.PropResult{Status: gopter.PropFalse}
		}
		return &gopter.PropResult{Status: gopter.PropTrue}
	},
}

type cbCommands struct {
	size     int
	elements *list.List
}

func (c *cbCommands) NewSystemUnderTest(initialState commands.State) commands.SystemUnderTest {
	s := initialState.(*cbCommands)
	q := New(c.size)
	for e := s.elements.Front(); e != nil; e = e.Next() {
		q.Put(e.Value.(int))
	}
	return q
}

func (c *cbCommands) DestroySystemUnderTest(sut commands.SystemUnderTest) {
	sut.(*Queue).Init()
}

func (c *cbCommands) GenInitialState() gopter.Gen {
	return gen.Const(NewCbCommands(c.size))
}

func (c *cbCommands) InitialPreCondition(state commands.State) bool {
	s := state.(*cbCommands)
	return s.elements.Len() >= 0 && s.elements.Len() <= s.size
}

func (c *cbCommands) GenCommand(state commands.State) gopter.Gen {
	return gen.OneConstOf(getCommand, putCommand, sizeCommand)
}

func NewCbCommands(size int) *cbCommands {
	return &cbCommands{
		size:     size,
		elements: list.New(),
	}
}

func TestQueue(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	properties := gopter.NewProperties(parameters)
	properties.Property("circular buffer", commands.Prop(NewCbCommands(10)))
	properties.TestingRun(t)
}

func TestQueueSequence(t *testing.T) {
	const v = 1
	q := New(1)
	if r := q.Size(); r != 0 {
		t.Fatal("Initial size incorrect")
	}
	if r := q.Put(v); r != v {
		t.Fatal("Put returned incorrect value")
	}
	if r := q.Size(); r != 1 {
		t.Fatal("Size returned incorrect value")
	}
	if r := q.Get(); r != v {
		t.Fatal("Get returned incorrect value")
	}
	if r := q.Size(); r != 0 {
		t.Fatal("Final size incorrect")
	}
}
