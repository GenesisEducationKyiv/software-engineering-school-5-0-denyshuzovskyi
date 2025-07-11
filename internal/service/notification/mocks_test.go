// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package notification

import (
	"context"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/dto"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/lib/sqlutil"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// NewMockWeatherProvider creates a new instance of MockWeatherProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockWeatherProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockWeatherProvider {
	mock := &MockWeatherProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockWeatherProvider is an autogenerated mock type for the WeatherProvider type
type MockWeatherProvider struct {
	mock.Mock
}

type MockWeatherProvider_Expecter struct {
	mock *mock.Mock
}

func (_m *MockWeatherProvider) EXPECT() *MockWeatherProvider_Expecter {
	return &MockWeatherProvider_Expecter{mock: &_m.Mock}
}

// GetCurrentWeather provides a mock function for the type MockWeatherProvider
func (_mock *MockWeatherProvider) GetCurrentWeather(context1 context.Context, s string) (*dto.WeatherWithLocationDTO, error) {
	ret := _mock.Called(context1, s)

	if len(ret) == 0 {
		panic("no return value specified for GetCurrentWeather")
	}

	var r0 *dto.WeatherWithLocationDTO
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, string) (*dto.WeatherWithLocationDTO, error)); ok {
		return returnFunc(context1, s)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, string) *dto.WeatherWithLocationDTO); ok {
		r0 = returnFunc(context1, s)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.WeatherWithLocationDTO)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = returnFunc(context1, s)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockWeatherProvider_GetCurrentWeather_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCurrentWeather'
type MockWeatherProvider_GetCurrentWeather_Call struct {
	*mock.Call
}

// GetCurrentWeather is a helper method to define mock.On call
//   - context1 context.Context
//   - s string
func (_e *MockWeatherProvider_Expecter) GetCurrentWeather(context1 interface{}, s interface{}) *MockWeatherProvider_GetCurrentWeather_Call {
	return &MockWeatherProvider_GetCurrentWeather_Call{Call: _e.mock.On("GetCurrentWeather", context1, s)}
}

func (_c *MockWeatherProvider_GetCurrentWeather_Call) Run(run func(context1 context.Context, s string)) *MockWeatherProvider_GetCurrentWeather_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 string
		if args[1] != nil {
			arg1 = args[1].(string)
		}
		run(
			arg0,
			arg1,
		)
	})
	return _c
}

func (_c *MockWeatherProvider_GetCurrentWeather_Call) Return(weatherWithLocationDTO *dto.WeatherWithLocationDTO, err error) *MockWeatherProvider_GetCurrentWeather_Call {
	_c.Call.Return(weatherWithLocationDTO, err)
	return _c
}

func (_c *MockWeatherProvider_GetCurrentWeather_Call) RunAndReturn(run func(context1 context.Context, s string) (*dto.WeatherWithLocationDTO, error)) *MockWeatherProvider_GetCurrentWeather_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockWeatherRepository creates a new instance of MockWeatherRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockWeatherRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockWeatherRepository {
	mock := &MockWeatherRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockWeatherRepository is an autogenerated mock type for the WeatherRepository type
type MockWeatherRepository struct {
	mock.Mock
}

type MockWeatherRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockWeatherRepository) EXPECT() *MockWeatherRepository_Expecter {
	return &MockWeatherRepository_Expecter{mock: &_m.Mock}
}

// FindLastUpdatedByLocation provides a mock function for the type MockWeatherRepository
func (_mock *MockWeatherRepository) FindLastUpdatedByLocation(context1 context.Context, sQLExecutor sqlutil.SQLExecutor, s string) (*model.Weather, error) {
	ret := _mock.Called(context1, sQLExecutor, s)

	if len(ret) == 0 {
		panic("no return value specified for FindLastUpdatedByLocation")
	}

	var r0 *model.Weather
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, sqlutil.SQLExecutor, string) (*model.Weather, error)); ok {
		return returnFunc(context1, sQLExecutor, s)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, sqlutil.SQLExecutor, string) *model.Weather); ok {
		r0 = returnFunc(context1, sQLExecutor, s)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Weather)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, sqlutil.SQLExecutor, string) error); ok {
		r1 = returnFunc(context1, sQLExecutor, s)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockWeatherRepository_FindLastUpdatedByLocation_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindLastUpdatedByLocation'
type MockWeatherRepository_FindLastUpdatedByLocation_Call struct {
	*mock.Call
}

// FindLastUpdatedByLocation is a helper method to define mock.On call
//   - context1 context.Context
//   - sQLExecutor sqlutil.SQLExecutor
//   - s string
func (_e *MockWeatherRepository_Expecter) FindLastUpdatedByLocation(context1 interface{}, sQLExecutor interface{}, s interface{}) *MockWeatherRepository_FindLastUpdatedByLocation_Call {
	return &MockWeatherRepository_FindLastUpdatedByLocation_Call{Call: _e.mock.On("FindLastUpdatedByLocation", context1, sQLExecutor, s)}
}

func (_c *MockWeatherRepository_FindLastUpdatedByLocation_Call) Run(run func(context1 context.Context, sQLExecutor sqlutil.SQLExecutor, s string)) *MockWeatherRepository_FindLastUpdatedByLocation_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 sqlutil.SQLExecutor
		if args[1] != nil {
			arg1 = args[1].(sqlutil.SQLExecutor)
		}
		var arg2 string
		if args[2] != nil {
			arg2 = args[2].(string)
		}
		run(
			arg0,
			arg1,
			arg2,
		)
	})
	return _c
}

func (_c *MockWeatherRepository_FindLastUpdatedByLocation_Call) Return(weather *model.Weather, err error) *MockWeatherRepository_FindLastUpdatedByLocation_Call {
	_c.Call.Return(weather, err)
	return _c
}

func (_c *MockWeatherRepository_FindLastUpdatedByLocation_Call) RunAndReturn(run func(context1 context.Context, sQLExecutor sqlutil.SQLExecutor, s string) (*model.Weather, error)) *MockWeatherRepository_FindLastUpdatedByLocation_Call {
	_c.Call.Return(run)
	return _c
}

// Save provides a mock function for the type MockWeatherRepository
func (_mock *MockWeatherRepository) Save(context1 context.Context, sQLExecutor sqlutil.SQLExecutor, weather *model.Weather) error {
	ret := _mock.Called(context1, sQLExecutor, weather)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, sqlutil.SQLExecutor, *model.Weather) error); ok {
		r0 = returnFunc(context1, sQLExecutor, weather)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockWeatherRepository_Save_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Save'
type MockWeatherRepository_Save_Call struct {
	*mock.Call
}

// Save is a helper method to define mock.On call
//   - context1 context.Context
//   - sQLExecutor sqlutil.SQLExecutor
//   - weather *model.Weather
func (_e *MockWeatherRepository_Expecter) Save(context1 interface{}, sQLExecutor interface{}, weather interface{}) *MockWeatherRepository_Save_Call {
	return &MockWeatherRepository_Save_Call{Call: _e.mock.On("Save", context1, sQLExecutor, weather)}
}

func (_c *MockWeatherRepository_Save_Call) Run(run func(context1 context.Context, sQLExecutor sqlutil.SQLExecutor, weather *model.Weather)) *MockWeatherRepository_Save_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 sqlutil.SQLExecutor
		if args[1] != nil {
			arg1 = args[1].(sqlutil.SQLExecutor)
		}
		var arg2 *model.Weather
		if args[2] != nil {
			arg2 = args[2].(*model.Weather)
		}
		run(
			arg0,
			arg1,
			arg2,
		)
	})
	return _c
}

func (_c *MockWeatherRepository_Save_Call) Return(err error) *MockWeatherRepository_Save_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockWeatherRepository_Save_Call) RunAndReturn(run func(context1 context.Context, sQLExecutor sqlutil.SQLExecutor, weather *model.Weather) error) *MockWeatherRepository_Save_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockSubscriberRepository creates a new instance of MockSubscriberRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockSubscriberRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockSubscriberRepository {
	mock := &MockSubscriberRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockSubscriberRepository is an autogenerated mock type for the SubscriberRepository type
type MockSubscriberRepository struct {
	mock.Mock
}

type MockSubscriberRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockSubscriberRepository) EXPECT() *MockSubscriberRepository_Expecter {
	return &MockSubscriberRepository_Expecter{mock: &_m.Mock}
}

// FindById provides a mock function for the type MockSubscriberRepository
func (_mock *MockSubscriberRepository) FindById(context1 context.Context, sQLExecutor sqlutil.SQLExecutor, n int32) (*model.Subscriber, error) {
	ret := _mock.Called(context1, sQLExecutor, n)

	if len(ret) == 0 {
		panic("no return value specified for FindById")
	}

	var r0 *model.Subscriber
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, sqlutil.SQLExecutor, int32) (*model.Subscriber, error)); ok {
		return returnFunc(context1, sQLExecutor, n)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, sqlutil.SQLExecutor, int32) *model.Subscriber); ok {
		r0 = returnFunc(context1, sQLExecutor, n)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Subscriber)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, sqlutil.SQLExecutor, int32) error); ok {
		r1 = returnFunc(context1, sQLExecutor, n)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockSubscriberRepository_FindById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindById'
type MockSubscriberRepository_FindById_Call struct {
	*mock.Call
}

// FindById is a helper method to define mock.On call
//   - context1 context.Context
//   - sQLExecutor sqlutil.SQLExecutor
//   - n int32
func (_e *MockSubscriberRepository_Expecter) FindById(context1 interface{}, sQLExecutor interface{}, n interface{}) *MockSubscriberRepository_FindById_Call {
	return &MockSubscriberRepository_FindById_Call{Call: _e.mock.On("FindById", context1, sQLExecutor, n)}
}

func (_c *MockSubscriberRepository_FindById_Call) Run(run func(context1 context.Context, sQLExecutor sqlutil.SQLExecutor, n int32)) *MockSubscriberRepository_FindById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 sqlutil.SQLExecutor
		if args[1] != nil {
			arg1 = args[1].(sqlutil.SQLExecutor)
		}
		var arg2 int32
		if args[2] != nil {
			arg2 = args[2].(int32)
		}
		run(
			arg0,
			arg1,
			arg2,
		)
	})
	return _c
}

func (_c *MockSubscriberRepository_FindById_Call) Return(subscriber *model.Subscriber, err error) *MockSubscriberRepository_FindById_Call {
	_c.Call.Return(subscriber, err)
	return _c
}

func (_c *MockSubscriberRepository_FindById_Call) RunAndReturn(run func(context1 context.Context, sQLExecutor sqlutil.SQLExecutor, n int32) (*model.Subscriber, error)) *MockSubscriberRepository_FindById_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockSubscriptionRepository creates a new instance of MockSubscriptionRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockSubscriptionRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockSubscriptionRepository {
	mock := &MockSubscriptionRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockSubscriptionRepository is an autogenerated mock type for the SubscriptionRepository type
type MockSubscriptionRepository struct {
	mock.Mock
}

type MockSubscriptionRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockSubscriptionRepository) EXPECT() *MockSubscriptionRepository_Expecter {
	return &MockSubscriptionRepository_Expecter{mock: &_m.Mock}
}

// FindAllByFrequencyAndConfirmedStatus provides a mock function for the type MockSubscriptionRepository
func (_mock *MockSubscriptionRepository) FindAllByFrequencyAndConfirmedStatus(context1 context.Context, sQLExecutor sqlutil.SQLExecutor, frequency model.Frequency) ([]*model.Subscription, error) {
	ret := _mock.Called(context1, sQLExecutor, frequency)

	if len(ret) == 0 {
		panic("no return value specified for FindAllByFrequencyAndConfirmedStatus")
	}

	var r0 []*model.Subscription
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, sqlutil.SQLExecutor, model.Frequency) ([]*model.Subscription, error)); ok {
		return returnFunc(context1, sQLExecutor, frequency)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, sqlutil.SQLExecutor, model.Frequency) []*model.Subscription); ok {
		r0 = returnFunc(context1, sQLExecutor, frequency)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Subscription)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, sqlutil.SQLExecutor, model.Frequency) error); ok {
		r1 = returnFunc(context1, sQLExecutor, frequency)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockSubscriptionRepository_FindAllByFrequencyAndConfirmedStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindAllByFrequencyAndConfirmedStatus'
type MockSubscriptionRepository_FindAllByFrequencyAndConfirmedStatus_Call struct {
	*mock.Call
}

// FindAllByFrequencyAndConfirmedStatus is a helper method to define mock.On call
//   - context1 context.Context
//   - sQLExecutor sqlutil.SQLExecutor
//   - frequency model.Frequency
func (_e *MockSubscriptionRepository_Expecter) FindAllByFrequencyAndConfirmedStatus(context1 interface{}, sQLExecutor interface{}, frequency interface{}) *MockSubscriptionRepository_FindAllByFrequencyAndConfirmedStatus_Call {
	return &MockSubscriptionRepository_FindAllByFrequencyAndConfirmedStatus_Call{Call: _e.mock.On("FindAllByFrequencyAndConfirmedStatus", context1, sQLExecutor, frequency)}
}

func (_c *MockSubscriptionRepository_FindAllByFrequencyAndConfirmedStatus_Call) Run(run func(context1 context.Context, sQLExecutor sqlutil.SQLExecutor, frequency model.Frequency)) *MockSubscriptionRepository_FindAllByFrequencyAndConfirmedStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 sqlutil.SQLExecutor
		if args[1] != nil {
			arg1 = args[1].(sqlutil.SQLExecutor)
		}
		var arg2 model.Frequency
		if args[2] != nil {
			arg2 = args[2].(model.Frequency)
		}
		run(
			arg0,
			arg1,
			arg2,
		)
	})
	return _c
}

func (_c *MockSubscriptionRepository_FindAllByFrequencyAndConfirmedStatus_Call) Return(subscriptions []*model.Subscription, err error) *MockSubscriptionRepository_FindAllByFrequencyAndConfirmedStatus_Call {
	_c.Call.Return(subscriptions, err)
	return _c
}

func (_c *MockSubscriptionRepository_FindAllByFrequencyAndConfirmedStatus_Call) RunAndReturn(run func(context1 context.Context, sQLExecutor sqlutil.SQLExecutor, frequency model.Frequency) ([]*model.Subscription, error)) *MockSubscriptionRepository_FindAllByFrequencyAndConfirmedStatus_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockTokenRepository creates a new instance of MockTokenRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockTokenRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockTokenRepository {
	mock := &MockTokenRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockTokenRepository is an autogenerated mock type for the TokenRepository type
type MockTokenRepository struct {
	mock.Mock
}

type MockTokenRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockTokenRepository) EXPECT() *MockTokenRepository_Expecter {
	return &MockTokenRepository_Expecter{mock: &_m.Mock}
}

// FindBySubscriptionIdAndType provides a mock function for the type MockTokenRepository
func (_mock *MockTokenRepository) FindBySubscriptionIdAndType(context1 context.Context, sQLExecutor sqlutil.SQLExecutor, n int32, tokenType model.TokenType) (*model.Token, error) {
	ret := _mock.Called(context1, sQLExecutor, n, tokenType)

	if len(ret) == 0 {
		panic("no return value specified for FindBySubscriptionIdAndType")
	}

	var r0 *model.Token
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, sqlutil.SQLExecutor, int32, model.TokenType) (*model.Token, error)); ok {
		return returnFunc(context1, sQLExecutor, n, tokenType)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, sqlutil.SQLExecutor, int32, model.TokenType) *model.Token); ok {
		r0 = returnFunc(context1, sQLExecutor, n, tokenType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Token)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, sqlutil.SQLExecutor, int32, model.TokenType) error); ok {
		r1 = returnFunc(context1, sQLExecutor, n, tokenType)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockTokenRepository_FindBySubscriptionIdAndType_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindBySubscriptionIdAndType'
type MockTokenRepository_FindBySubscriptionIdAndType_Call struct {
	*mock.Call
}

// FindBySubscriptionIdAndType is a helper method to define mock.On call
//   - context1 context.Context
//   - sQLExecutor sqlutil.SQLExecutor
//   - n int32
//   - tokenType model.TokenType
func (_e *MockTokenRepository_Expecter) FindBySubscriptionIdAndType(context1 interface{}, sQLExecutor interface{}, n interface{}, tokenType interface{}) *MockTokenRepository_FindBySubscriptionIdAndType_Call {
	return &MockTokenRepository_FindBySubscriptionIdAndType_Call{Call: _e.mock.On("FindBySubscriptionIdAndType", context1, sQLExecutor, n, tokenType)}
}

func (_c *MockTokenRepository_FindBySubscriptionIdAndType_Call) Run(run func(context1 context.Context, sQLExecutor sqlutil.SQLExecutor, n int32, tokenType model.TokenType)) *MockTokenRepository_FindBySubscriptionIdAndType_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 sqlutil.SQLExecutor
		if args[1] != nil {
			arg1 = args[1].(sqlutil.SQLExecutor)
		}
		var arg2 int32
		if args[2] != nil {
			arg2 = args[2].(int32)
		}
		var arg3 model.TokenType
		if args[3] != nil {
			arg3 = args[3].(model.TokenType)
		}
		run(
			arg0,
			arg1,
			arg2,
			arg3,
		)
	})
	return _c
}

func (_c *MockTokenRepository_FindBySubscriptionIdAndType_Call) Return(token *model.Token, err error) *MockTokenRepository_FindBySubscriptionIdAndType_Call {
	_c.Call.Return(token, err)
	return _c
}

func (_c *MockTokenRepository_FindBySubscriptionIdAndType_Call) RunAndReturn(run func(context1 context.Context, sQLExecutor sqlutil.SQLExecutor, n int32, tokenType model.TokenType) (*model.Token, error)) *MockTokenRepository_FindBySubscriptionIdAndType_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockEmailSender creates a new instance of MockEmailSender. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockEmailSender(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockEmailSender {
	mock := &MockEmailSender{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockEmailSender is an autogenerated mock type for the EmailSender type
type MockEmailSender struct {
	mock.Mock
}

type MockEmailSender_Expecter struct {
	mock *mock.Mock
}

func (_m *MockEmailSender) EXPECT() *MockEmailSender_Expecter {
	return &MockEmailSender_Expecter{mock: &_m.Mock}
}

// Send provides a mock function for the type MockEmailSender
func (_mock *MockEmailSender) Send(context1 context.Context, simpleEmail dto.SimpleEmail) error {
	ret := _mock.Called(context1, simpleEmail)

	if len(ret) == 0 {
		panic("no return value specified for Send")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, dto.SimpleEmail) error); ok {
		r0 = returnFunc(context1, simpleEmail)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockEmailSender_Send_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Send'
type MockEmailSender_Send_Call struct {
	*mock.Call
}

// Send is a helper method to define mock.On call
//   - context1 context.Context
//   - simpleEmail dto.SimpleEmail
func (_e *MockEmailSender_Expecter) Send(context1 interface{}, simpleEmail interface{}) *MockEmailSender_Send_Call {
	return &MockEmailSender_Send_Call{Call: _e.mock.On("Send", context1, simpleEmail)}
}

func (_c *MockEmailSender_Send_Call) Run(run func(context1 context.Context, simpleEmail dto.SimpleEmail)) *MockEmailSender_Send_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 dto.SimpleEmail
		if args[1] != nil {
			arg1 = args[1].(dto.SimpleEmail)
		}
		run(
			arg0,
			arg1,
		)
	})
	return _c
}

func (_c *MockEmailSender_Send_Call) Return(err error) *MockEmailSender_Send_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockEmailSender_Send_Call) RunAndReturn(run func(context1 context.Context, simpleEmail dto.SimpleEmail) error) *MockEmailSender_Send_Call {
	_c.Call.Return(run)
	return _c
}
