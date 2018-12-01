[![Build Status](https://travis-ci.org/tzmfreedom/lui.svg?branch=master)](https://travis-ci.org/tzmfreedom/lui)

# lui

Lightning Platform TerminalUI Application

## Install

For Linux
```bash
$ curl -sL http://install.freedom-man.com/lui.sh | bash
```

For Windows with Command Prompt
```
@"%SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe" -NoProfile ^
  -InputFormat None -ExecutionPolicy Bypass ^
  -Command "iex ((New-Object System.Net.WebClient).DownloadString('http://install.freedom-man.com/lui.ps1'))" ^
  && SET "PATH=%PATH%;%APPDATA%\lui\bin"
```

For Windows with PowerShell
```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; iex ((New-Object System.Net.WebClient).DownloadString('http://install.freedom-man.com/yasd.ps1'))
```

For Golang
```bash
$ go get github.com/tzmfreedom/lui
```

## Usage

```
$ lui -u USERNAME [-e ENDPOINT] [-v] [-h, --help]
```

## Contribute

Just send pull request if needed or fill an issue!

## License

The MIT License See [LICENSE](https://github.com/tzmfreedom/lui/blob/master/LICENSE) file.
