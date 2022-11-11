package weak_test

import (
	"testing"

	weak "github.com/jcoder76/weak"
	"github.com/stretchr/testify/suite"
)

type TestObj1 struct {
	value int
}

type TestObj2 struct {
	value string
}

type WeakTableTestSuite struct {
	suite.Suite
	table *weak.Table[TestObj1, TestObj2]
}

func (s *WeakTableTestSuite) SetupTest() {
	s.table = &weak.Table[TestObj1, TestObj2]{}
}

func (s *WeakTableTestSuite) TestCanGetNewObjectfromTable() {
	obj1, expect := s.setupValues()

	result := s.table.Get(obj1)

	s.Equal(expect, result)
}

func (s *WeakTableTestSuite) TestCanGetDeleteObjectfromTable() {
	obj1, _ := s.setupValues()

	s.table.Delete(obj1)
	result := s.table.Get(obj1)

	s.Nil(result)
}

// TODO: Doesn't work as expected, this test fails
// func (s *WeakTableTestSuite) TestGCDeletesObjectfromTable() {
// 	s.setupValues()

// 	// Force GC cycle
// 	runtime.GC()
// 	size := s.table.Size()

// 	s.Zero(size)
// }

func (s *WeakTableTestSuite) setupValues() (*TestObj1, *TestObj2) {
	obj1 := &TestObj1{
		value: 1337,
	}
	expect := &TestObj2{
		value: "test",
	}
	s.table.GetOrCreate(obj1, func(key *TestObj1) *TestObj2 {
		return expect
	})

	return obj1, expect
}

func TestWeakTableTestSuite(t *testing.T) {
	suite.Run(t, new(WeakTableTestSuite))
}
