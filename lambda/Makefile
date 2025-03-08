build:
	@if exist bootstrap del bootstrap
	@if exist function.zip del function.zip
	@cmd /C "set GOOS=linux&&set GOARCH=amd64&& go build -o bootstrap"
	@tar -cvf function.zip bootstrap
	@if exist bootstrap del bootstrap
