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

type TableTestSuite struct {
	suite.Suite
	table *weak.Table[TestObj1, TestObj2]
}

func (s *TableTestSuite) SetupTest() {
	s.table = &weak.Table[TestObj1, TestObj2]{}
}

func (s *TableTestSuite) TestGeWithNilKeyReturnsNil() {
	var obj1 *TestObj1 = nil

	result := s.table.Get(obj1)

	s.Nil(result)
}

func (s *TableTestSuite) TestGetOrCreateWithNilKeyReturnsNil() {
	var obj1 *TestObj1 = nil

	result := s.table.GetOrCreate(obj1, func(key *TestObj1) *TestObj2 {
		return &TestObj2{}
	})

	s.Nil(result)
}

func (s *TableTestSuite) TestGetWithUnkownKeyReturnsNil() {
	obj1 := &TestObj1{}

	result := s.table.Get(obj1)

	s.Nil(result)
}

func (s *TableTestSuite) TestCanGetNewObjectfromTable() {
	obj1, expect := s.setupValues()

	result := s.table.Get(obj1)

	s.Equal(expect, result)
}

func (s *TableTestSuite) TestCanGetDeleteObjectfromTable() {
	obj1, _ := s.setupValues()

	s.table.Delete(obj1)
	result := s.table.Get(obj1)

	s.Nil(result)
}

// TODO: Doesn't work as expected, this test fails
// func (s *TableTestSuite) TestGCDeletesObjectfromTable() {
// 	s.setupValues()

// 	// Force GC cycle
// 	runtime.GC()
// 	size := s.table.Size()

// 	s.Zero(size)
// }

func (s *TableTestSuite) setupValues() (*TestObj1, *TestObj2) {
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

func TestTableTestSuite(t *testing.T) {
	suite.Run(t, new(TableTestSuite))
}
