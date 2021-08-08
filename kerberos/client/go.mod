module github.com/fsmiamoto/system_security/kerberos/client

replace github.com/fsmiamoto/system_security/kerberos/as v0.0.0 => /home/shigueo/code/system_security/kerberos/as

replace github.com/fsmiamoto/system_security/kerberos/crypto v0.0.0 => /home/shigueo/code/system_security/kerberos/crypto

go 1.16

require github.com/fsmiamoto/system_security/kerberos/as v0.0.0-20210731194158-1245835fd995

require (
	github.com/fatih/color v1.12.0 // indirect
	github.com/fsmiamoto/system_security/kerberos/crypto v0.0.0
	github.com/hokaccha/go-prettyjson v0.0.0-20210113012101-fb4e108d2519
)
