package weak_test

import (
	"testing"

	"github.com/jcoder76/weak"
	"github.com/stretchr/testify/suite"
)

type TestTargetObj1 struct {
	value int
}

type TestTargetObj2 struct {
	value string
}

type TestValueObj1 struct {
	value bool
}

type ExtendTestSuite struct {
	suite.Suite
}

func (s *ExtendTestSuite) TestExtendNilPanics() {
	var target *TestTargetObj1 = nil

	resultFunc := func() {
		weak.Extend[TestValueObj1](target)
	}

	s.Panics(resultFunc)
}

func (s *ExtendTestSuite) TestExtendGivesTypeExtension() {
	target := &TestTargetObj1{
		value: 1,
	}

	result := weak.Extend[TestValueObj1](target)

	s.Require().NotNil(result)
	s.IsType(new(TestValueObj1), result)
}

func (s *ExtendTestSuite) TestExtendTwiceRemembersTypeExtension() {
	target := &TestTargetObj1{
		value: 1,
	}
	weak.Extend[TestValueObj1](target).value = true

	result := weak.Extend[TestValueObj1](target)

	s.Require().NotNil(result)
	s.Equal(true, result.value)
}

func (s *ExtendTestSuite) TestExtendMultipleTargetsCreatesSeparateTypeExtensions() {
	target1 := &TestTargetObj1{
		value: 1,
	}
	target2 := &TestTargetObj2{
		value: "test",
	}

	result1 := weak.Extend[TestValueObj1](target1)
	result2 := weak.Extend[TestValueObj1](target2)

	s.Require().NotNil(result1)
	s.Require().NotNil(result2)
	s.IsType(new(TestValueObj1), result1)
	s.IsType(new(TestValueObj1), result2)
	s.NotSame(target1, target2)
}

func TestExtendTestSuite(t *testing.T) {
	suite.Run(t, new(ExtendTestSuite))
}
