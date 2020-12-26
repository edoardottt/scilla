# Do the same as "make.bat" but in a powershell version (maybe a bit more)

# Probably you need to activate the script policy on your PowerShell
# Set-ExecutionPolicy Unrestricted
# https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.security/ than search for "policy"

$option = $args[0]
$gopath = $Env:GOPATH
$bin = "$gopath\bin\scilla.exe"

function WindowsInstall {
    if ((Test-Path $gopath) -ne $true) {
        Set-Content -Path Env:GOPATH -Value "$home\go"
    }
    Set-Content -Path Env:GO111MODULE -Value on
    Set-Content -Path Env:CGO_ENABLED -Value 0
    go build -o "$Env:GOPATH\bin\scilla.exe"
    Write-Host "Done."
    Write-Host "Type scilla.exe help"
}

function WindowsUninstall {
    if ((Test-Path $bin) -eq $true) {
        Remove-Item "$Env:GOPATH\bin\scilla.exe"
    }
    Write-Host "Done."
}

function Update {
    Set-Content -Path $Env:GO111MODULE -Value on
    Write-Host "Updating..."
    go get -u
    go mod tidy -v
    Write-Host "Done."
}

function Fmt {
    Set-Content -Path $Env:GO111MODULE -Value on
    Write-Host "Updating..."
    go fmt ./...
    Write-Host "Done."
}

function Test {
    Set-Content -Path $Env:GO111MODULE -Value on
    Set-Content -Path Env:CGO_ENABLED -Value 0
    Write-Host "Testing..."
    go test -v ./...
    Write-Host "Done."
}

function Remod {
    Remove-Item go.mod go.sum
    go mod init github.com/edoardottt/scilla
    go get
}

switch ($option) {
    "install" {
        WindowsInstall
    }

    "uninstall" {
        WindowsUninstall
    }

    "update" {
        Update
    }

    "fmt" {
        Fmt
    }

    "test" {
        Test
    }

    "remod" {
        Remod
    }

    Default {
        "Invalid Option"
    }
}