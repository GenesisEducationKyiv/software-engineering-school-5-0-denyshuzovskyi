// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package service

import (
	"context"
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
func (_mock *MockWeatherProvider) GetCurrentWeather(s string) (*model.WeatherWithLocation, error) {
	ret := _mock.Called(s)

	if len(ret) == 0 {
		panic("no return value specified for GetCurrentWeather")
	}

	var r0 *model.WeatherWithLocation
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(string) (*model.WeatherWithLocation, error)); ok {
		return returnFunc(s)
	}
	if returnFunc, ok := ret.Get(0).(func(string) *model.WeatherWithLocation); ok {
		r0 = returnFunc(s)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.WeatherWithLocation)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(string) error); ok {
		r1 = returnFunc(s)
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
//   - s
func (_e *MockWeatherProvider_Expecter) GetCurrentWeather(s interface{}) *MockWeatherProvider_GetCurrentWeather_Call {
	return &MockWeatherProvider_GetCurrentWeather_Call{Call: _e.mock.On("GetCurrentWeather", s)}
}

func (_c *MockWeatherProvider_GetCurrentWeather_Call) Run(run func(s string)) *MockWeatherProvider_GetCurrentWeather_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockWeatherProvider_GetCurrentWeather_Call) Return(weatherWithLocation *model.WeatherWithLocation, err error) *MockWeatherProvider_GetCurrentWeather_Call {
	_c.Call.Return(weatherWithLocation, err)
	return _c
}

func (_c *MockWeatherProvider_GetCurrentWeather_Call) RunAndReturn(run func(s string) (*model.WeatherWithLocation, error)) *MockWeatherProvider_GetCurrentWeather_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockLocationRepository creates a new instance of MockLocationRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockLocationRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockLocationRepository {
	mock := &MockLocationRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockLocationRepository is an autogenerated mock type for the LocationRepository type
type MockLocationRepository struct {
	mock.Mock
}

type MockLocationRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockLocationRepository) EXPECT() *MockLocationRepository_Expecter {
	return &MockLocationRepository_Expecter{mock: &_m.Mock}
}

// FindByName provides a mock function for the type MockLocationRepository
func (_mock *MockLocationRepository) FindByName(context1 context.Context, sQLExecutor sqlutil.SQLExecutor, s string) (*model.Location, error) {
	ret := _mock.Called(context1, sQLExecutor, s)

	if len(ret) == 0 {
		panic("no return value specified for FindByName")
	}

	var r0 *model.Location
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, sqlutil.SQLExecutor, string) (*model.Location, error)); ok {
		return returnFunc(context1, sQLExecutor, s)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, sqlutil.SQLExecutor, string) *model.Location); ok {
		r0 = returnFunc(context1, sQLExecutor, s)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Location)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, sqlutil.SQLExecutor, string) error); ok {
		r1 = returnFunc(context1, sQLExecutor, s)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockLocationRepository_FindByName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindByName'
type MockLocationRepository_FindByName_Call struct {
	*mock.Call
}

// FindByName is a helper method to define mock.On call
//   - context1
//   - sQLExecutor
//   - s
func (_e *MockLocationRepository_Expecter) FindByName(context1 interface{}, sQLExecutor interface{}, s interface{}) *MockLocationRepository_FindByName_Call {
	return &MockLocationRepository_FindByName_Call{Call: _e.mock.On("FindByName", context1, sQLExecutor, s)}
}

func (_c *MockLocationRepository_FindByName_Call) Run(run func(context1 context.Context, sQLExecutor sqlutil.SQLExecutor, s string)) *MockLocationRepository_FindByName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(sqlutil.SQLExecutor), args[2].(string))
	})
	return _c
}

func (_c *MockLocationRepository_FindByName_Call) Return(location *model.Location, err error) *MockLocationRepository_FindByName_Call {
	_c.Call.Return(location, err)
	return _c
}

func (_c *MockLocationRepository_FindByName_Call) RunAndReturn(run func(context1 context.Context, sQLExecutor sqlutil.SQLExecutor, s string) (*model.Location, error)) *MockLocationRepository_FindByName_Call {
	_c.Call.Return(run)
	return _c
}

// Save provides a mock function for the type MockLocationRepository
func (_mock *MockLocationRepository) Save(context1 context.Context, sQLExecutor sqlutil.SQLExecutor, location *model.Location) (int32, error) {
	ret := _mock.Called(context1, sQLExecutor, location)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 int32
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, sqlutil.SQLExecutor, *model.Location) (int32, error)); ok {
		return returnFunc(context1, sQLExecutor, location)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, sqlutil.SQLExecutor, *model.Location) int32); ok {
		r0 = returnFunc(context1, sQLExecutor, location)
	} else {
		r0 = ret.Get(0).(int32)
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, sqlutil.SQLExecutor, *model.Location) error); ok {
		r1 = returnFunc(context1, sQLExecutor, location)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockLocationRepository_Save_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Save'
type MockLocationRepository_Save_Call struct {
	*mock.Call
}

// Save is a helper method to define mock.On call
//   - context1
//   - sQLExecutor
//   - location
func (_e *MockLocationRepository_Expecter) Save(context1 interface{}, sQLExecutor interface{}, location interface{}) *MockLocationRepository_Save_Call {
	return &MockLocationRepository_Save_Call{Call: _e.mock.On("Save", context1, sQLExecutor, location)}
}

func (_c *MockLocationRepository_Save_Call) Run(run func(context1 context.Context, sQLExecutor sqlutil.SQLExecutor, location *model.Location)) *MockLocationRepository_Save_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(sqlutil.SQLExecutor), args[2].(*model.Location))
	})
	return _c
}

func (_c *MockLocationRepository_Save_Call) Return(n int32, err error) *MockLocationRepository_Save_Call {
	_c.Call.Return(n, err)
	return _c
}

func (_c *MockLocationRepository_Save_Call) RunAndReturn(run func(context1 context.Context, sQLExecutor sqlutil.SQLExecutor, location *model.Location) (int32, error)) *MockLocationRepository_Save_Call {
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
//   - context1
//   - sQLExecutor
//   - s
func (_e *MockWeatherRepository_Expecter) FindLastUpdatedByLocation(context1 interface{}, sQLExecutor interface{}, s interface{}) *MockWeatherRepository_FindLastUpdatedByLocation_Call {
	return &MockWeatherRepository_FindLastUpdatedByLocation_Call{Call: _e.mock.On("FindLastUpdatedByLocation", context1, sQLExecutor, s)}
}

func (_c *MockWeatherRepository_FindLastUpdatedByLocation_Call) Run(run func(context1 context.Context, sQLExecutor sqlutil.SQLExecutor, s string)) *MockWeatherRepository_FindLastUpdatedByLocation_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(sqlutil.SQLExecutor), args[2].(string))
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
//   - context1
//   - sQLExecutor
//   - weather
func (_e *MockWeatherRepository_Expecter) Save(context1 interface{}, sQLExecutor interface{}, weather interface{}) *MockWeatherRepository_Save_Call {
	return &MockWeatherRepository_Save_Call{Call: _e.mock.On("Save", context1, sQLExecutor, weather)}
}

func (_c *MockWeatherRepository_Save_Call) Run(run func(context1 context.Context, sQLExecutor sqlutil.SQLExecutor, weather *model.Weather)) *MockWeatherRepository_Save_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(sqlutil.SQLExecutor), args[2].(*model.Weather))
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
