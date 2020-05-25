module "chatServer"



go 1.14

require (
	go.uber.org/zap v1.15.0
	gohipernetFake v0.0.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)

replace gohipernetFake v0.0.0 => ../gohipernetFake
