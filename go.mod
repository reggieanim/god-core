module github.com/reggieanim/god-core

go 1.21

toolchain go1.23.1

require (
	github.com/aws/aws-sdk-go v1.54.19
	github.com/gen2brain/beeep v0.0.0-20240516210008-9c006672e7f4
	github.com/go-rod/rod v0.116.2
	github.com/go-rod/stealth v0.4.9
)

require github.com/jmespath/go-jmespath v0.4.0 // indirect

require (
	github.com/go-toast/toast v0.0.0-20190211030409-01e6764cf0a4 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d // indirect
	github.com/tadvi/systray v0.0.0-20190226123456-11a2b8fa57af // indirect
	github.com/ysmood/fetchup v0.2.4 // indirect
	github.com/ysmood/goob v0.4.0 // indirect
	github.com/ysmood/got v0.40.0 // indirect
	github.com/ysmood/gson v0.7.3 // indirect
	golang.org/x/sys v0.22.0 // indirect
)

replace github.com/go-rod/rod => github.com/reggieanim/rod v0.116.8
