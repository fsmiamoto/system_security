module github.com/fsmiamoto/system_security/kerberos/client

go 1.16

replace github.com/fsmiamoto/system_security/kerberos/tgs v0.0.0 => /home/shigueo/code/system_security/kerberos/tgs

replace github.com/fsmiamoto/system_security/kerberos/as v0.0.0 => /home/shigueo/code/system_security/kerberos/as

require (
	github.com/fatih/color v1.12.0 // indirect
	github.com/fsmiamoto/system_security/kerberos/as v0.0.0
	github.com/fsmiamoto/system_security/kerberos/crypto v0.0.0-20210808013406-62c152ac960b
	github.com/fsmiamoto/system_security/kerberos/tgs v0.0.0
	github.com/hokaccha/go-prettyjson v0.0.0-20210113012101-fb4e108d2519
)
