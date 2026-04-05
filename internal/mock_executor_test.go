package internal

type MockExecutor struct {
	ExecuteFunc       func(options ExecuteOptions, args ...string) (int, string, error)
	ExecuteResticFunc func(options ExecuteOptions, args ...string) (int, string, error)
}

func (m *MockExecutor) Execute(options ExecuteOptions, args ...string) (int, string, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(options, args...)
	}
	return 0, "", nil
}

func (m *MockExecutor) ExecuteRestic(options ExecuteOptions, args ...string) (int, string, error) {
	if m.ExecuteResticFunc != nil {
		return m.ExecuteResticFunc(options, args...)
	}
	return 0, "", nil
}
