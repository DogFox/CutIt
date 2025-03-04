package app

type App struct {
	logger Logger
}

type Logger interface{}

func New(logger Logger) *App {
	return &App{
		logger: logger,
	}
}
