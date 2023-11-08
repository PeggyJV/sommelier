package v7

// UpgradeName defines the on-chain upgrade name for the Sommelier v7 upgrade
const UpgradeName = "v7"

// 7seas domain
const SevenSeasDomain = "7seas.capital"

// TODO(bolten): update this
// CA certificate literals
// See source data at: https://github.com/PeggyJV/steward-registry
// data captured at commit cdee05a8bf97f264353e10ab65752710bfb85dc9

// Publisher CA certificate
const SevenSeasPublisherCA = `-----BEGIN CERTIFICATE-----
MIICGjCCAaCgAwIBAgIUc+HkMr23CFNHqHMujXNc5+2bd3cwCgYIKoZIzj0EAwMw
RDELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxIDAeBgNVBAoMF1Nl
dmVuIFNlYXMgU3RyYXRlZ3kgSW5jMB4XDTIyMDcwNjIzMjIxNVoXDTI0MDcwNTIz
MjIxNVowRDELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxIDAeBgNV
BAoMF1NldmVuIFNlYXMgU3RyYXRlZ3kgSW5jMHYwEAYHKoZIzj0CAQYFK4EEACID
YgAE2pprBwCCI0DoBtuX5+Xq+Uo5ZpH+SlTkh6bS46Lq/YICh5FpWGOLsjT51uTX
gyTD6gOzLKQjiWt7D18pqM5DldJ4lT9ZsEjQkpdn2IDNkMWvUKJlDdds9VUfk6Xb
QgtVo1MwUTAdBgNVHQ4EFgQUNugdj3bw6N5owK2p8JuYlJ91xtMwHwYDVR0jBBgw
FoAUNugdj3bw6N5owK2p8JuYlJ91xtMwDwYDVR0TAQH/BAUwAwEB/zAKBggqhkjO
PQQDAwNoADBlAjBcJBGQvx1oDQ3nBDLSWUl/F3EU0EUYwGkU/5jPlOLo9jEIeSBC
FbHA97/ROdzMYOICMQCBnUe1BU5dRki6/ToK16lqHHS3c1oX2953oWBmvZV5+auA
M12q1LwBzeHzVAr28q8=
-----END CERTIFICATE-----
`

// Subscriber CA certificates
const FigmentSubscriberCA = `-----BEGIN CERTIFICATE-----
MIICtjCCAjygAwIBAgIUUTtadL/U8YrQsevLZMj/32FjtE4wCgYIKoZIzj0EAwMw
gZExCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJ
bnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQxEDAOBgNVBAsMB0ZpZ21lbnQxODA2BgNV
BAMML3NvbW1lbGllci1zdGV3YXJkLnN0YWtpbmcucHJvZHVjdGlvbi5maWdtZW50
LmlvMB4XDTIyMDUxODIzNDkwOFoXDTI0MDUxNzIzNDkwOFowgZExCzAJBgNVBAYT
AkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBXaWRn
aXRzIFB0eSBMdGQxEDAOBgNVBAsMB0ZpZ21lbnQxODA2BgNVBAMML3NvbW1lbGll
ci1zdGV3YXJkLnN0YWtpbmcucHJvZHVjdGlvbi5maWdtZW50LmlvMHYwEAYHKoZI
zj0CAQYFK4EEACIDYgAEn/nhx68b7k/IHKlzuViX6lyqln3wKNt/eMi/2TLGd4RC
RwFYCQJhLWC8X5G1Zva0M5u8HY63r4MT7318MCZ1Ixn83CoiUO+owXnDM/KKD662
DIw621JAYJ5t08hhzt3wo1MwUTAdBgNVHQ4EFgQU8gsJCCztQ2vSMfo9Cn+IhxvV
vsIwHwYDVR0jBBgwFoAU8gsJCCztQ2vSMfo9Cn+IhxvVvsIwDwYDVR0TAQH/BAUw
AwEB/zAKBggqhkjOPQQDAwNoADBlAjAY1PBTIPun17hV6PxJlRBq1qDfwzAy4BvX
yJMVcOHs18taVMXiVDtReK9P7LpFVJECMQC2YMmfLbM9Chu1vduQIv1H+Q3EaIep
Ccmsdu0/lduoeuO6AqMPXW6nCZw4PonrttA=
-----END CERTIFICATE-----
`

const StandardCryptoSubscriberCA = `-----BEGIN CERTIFICATE-----
MIICBTCCAYqgAwIBAgIUKAjFu5IU/BKQ+uBZmuTGF6X4K6owCgYIKoZIzj0EAwMw
OTEYMBYGA1UECgwPU3RhbmRhcmQgQ3J5cHRvMR0wGwYDVQQDDBRzdGFuZGFyZGNy
eXB0b3ZjLmNvbTAeFw0yMjA1MTMwNDM1MDlaFw0yNDA1MTIwNDM1MDlaMDkxGDAW
BgNVBAoMD1N0YW5kYXJkIENyeXB0bzEdMBsGA1UEAwwUc3RhbmRhcmRjcnlwdG92
Yy5jb20wdjAQBgcqhkjOPQIBBgUrgQQAIgNiAATW7dxvqytu0MAtURLBATejrkTE
xXNtYDYysqdl0enTTaFZ5Co1zy4OvXz9Db8FdiZKMKS4V4oQegcAMqpPwVwz/6z3
AK40aWVWquzBcTtnvBwP2Oq4UFIYNASIGOrd2mmjUzBRMB0GA1UdDgQWBBTH6Zuh
M8Akp5vs+NObKpw0vhYWYTAfBgNVHSMEGDAWgBTH6ZuhM8Akp5vs+NObKpw0vhYW
YTAPBgNVHRMBAf8EBTADAQH/MAoGCCqGSM49BAMDA2kAMGYCMQCTBQw15XacLD0C
SmsTzIuqJZO/3rNaSvnvnoWSsrY2MZQwdWlC17WsSKOEALCh+90CMQCqh1wE8GYw
6e1n90a+O/9yQru2jnNXp/z+vUMEacx6ZnOM/RXE3eQctR6pC4R80io=
-----END CERTIFICATE-----
`

const RockawaySubscriberCA = `-----BEGIN CERTIFICATE-----
MIICSTCCAc6gAwIBAgIUKQbrzzmbw9e14+IPwdwsxphJwqQwCgYIKoZIzj0EAwMw
WzELMAkGA1UEBhMCQ1oxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoMGElu
dGVybmV0IFdpZGdpdHMgUHR5IEx0ZDEUMBIGA1UEAwwLcmJmLnN5c3RlbXMwHhcN
MjIwNTEzMTUzOTM3WhcNMjQwNTEyMTUzOTM3WjBbMQswCQYDVQQGEwJDWjETMBEG
A1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50ZXJuZXQgV2lkZ2l0cyBQdHkg
THRkMRQwEgYDVQQDDAtyYmYuc3lzdGVtczB2MBAGByqGSM49AgEGBSuBBAAiA2IA
BOjOff/X4d3ZQ3uSokoqgUtVwr9OQJjWqJM3CJkk7FUYTg/AP1UnU42cmOUbhbKn
4+96RK0M6qye6+wGVJB+yR6LqkPuoyzYvswGjaB9NQWKIAW4OIDUrw4Uc4qitsjG
G6NTMFEwHQYDVR0OBBYEFNt1sXxbYAZZMEaGGC/5uHmTcnueMB8GA1UdIwQYMBaA
FNt1sXxbYAZZMEaGGC/5uHmTcnueMA8GA1UdEwEB/wQFMAMBAf8wCgYIKoZIzj0E
AwMDaQAwZgIxAKym3BhXv7oQOsEwKf57TRSBFyY+nYvA1h3KK+h3RhCvORjDIaMr
NsCHc3hjYwmt8wIxAPNOix7pd6e7KyhUGSxGsI3S/99XhqAZIO4WwlcXQ3ekcPgL
JEZSGkKH8mWi6uOyFQ==
-----END CERTIFICATE-----
`

const BlockscapeSubscriberCA = `-----BEGIN CERTIFICATE-----
MIICdjCCAfygAwIBAgIJALDoLnBuMsOCMAoGCCqGSM49BAMDMHkxCzAJBgNVBAYT
AkRFMQswCQYDVQQIDAJCVzESMBAGA1UEBwwJU3R1dHRnYXJ0MQ0wCwYDVQQKDARN
V0FZMRMwEQYDVQQLDApibG9ja3NjYXBlMSUwIwYJKoZIhvcNAQkBFhZkZXZAYmxv
Y2tzY2FwZS5uZXR3b3JrMB4XDTIyMDUyMDEyNDc0NVoXDTI0MDUxOTEyNDc0NVow
eTELMAkGA1UEBhMCREUxCzAJBgNVBAgMAkJXMRIwEAYDVQQHDAlTdHV0dGdhcnQx
DTALBgNVBAoMBE1XQVkxEzARBgNVBAsMCmJsb2Nrc2NhcGUxJTAjBgkqhkiG9w0B
CQEWFmRldkBibG9ja3NjYXBlLm5ldHdvcmswdjAQBgcqhkjOPQIBBgUrgQQAIgNi
AATXYZ8ovRJ+VUSUXR8SPk1AGzUoEmJSPYNo2Gz+7IwtH8r5lvo/rG2R4rgM2oH5
K5DP5xuReQaeu0i3wHleHyLj/Z/vv2h4STWVK4e9RrBX9ytCViDHIMYZyAU5/d/E
c3GjUDBOMB0GA1UdDgQWBBQrGubTZudwLx8sWsLYQJ6W/eLhRDAfBgNVHSMEGDAW
gBQrGubTZudwLx8sWsLYQJ6W/eLhRDAMBgNVHRMEBTADAQH/MAoGCCqGSM49BAMD
A2gAMGUCMQCiVbvLpqBNGSMPLSqdtScrBKxDH+cbjkst+fbgqbNyeVU7pstSilyK
28yjZgbnYRACMDXTuPDOjT3qTnYhPjNGQHfjOPt4Wc9dSXbykpByoNj6qWRLg+os
vxKm6JF/QUn+kQ==
-----END CERTIFICATE-----
`

const SimplySubscriberCA = `-----BEGIN CERTIFICATE-----
MIICdjCCAfygAwIBAgIUTMulocP5XzC3bt3Kkqw7QaGox+YwCgYIKoZIzj0EAwMw
cjELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoMGElu
dGVybmV0IFdpZGdpdHMgUHR5IEx0ZDErMCkGA1UEAwwic29tbWVsaWVyLXN0ZXdh
cmQuc2ltcGx5LXZjLmNvbS5tdDAeFw0yMjA1MjMyMzI3NTdaFw0yNDA1MjIyMzI3
NTdaMHIxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQK
DBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQxKzApBgNVBAMMInNvbW1lbGllci1z
dGV3YXJkLnNpbXBseS12Yy5jb20ubXQwdjAQBgcqhkjOPQIBBgUrgQQAIgNiAARl
amVEMqPdpopyCYaAHzot2GIcj0W5F3sP8rPMFs2HzO+T9ZQaIokPX6AVMVGp2e+K
XqGVPhDhtlGtdbjNlG2anJoYvuDPnzzh5eHEtAfjkLT0vkTgQj580pxK8STWhiij
UzBRMB0GA1UdDgQWBBT/nNSZL/itgCGxZBkBZTR7w1aXRDAfBgNVHSMEGDAWgBT/
nNSZL/itgCGxZBkBZTR7w1aXRDAPBgNVHRMBAf8EBTADAQH/MAoGCCqGSM49BAMD
A2gAMGUCMQDGZspQGgJaJkb1NaOU6TtlKIfwydroT7tCTkx1SnQLdGt45TVUZ+Hx
DvfQ01lsakECMBIjzSOaLxBxVH1VhJ56qUXOU6+qxq220SKRs59SYDJDXOijlcwt
/ZQt9kN29MSCXg==
-----END CERTIFICATE-----
`

const PupmosSubscriberCA = `-----BEGIN CERTIFICATE-----
MIICxTCCAkqgAwIBAgIULdZ72d7nrmnEhacXf/otvo/uDfgwCgYIKoZIzj0EAwMw
gZgxCzAJBgNVBAYTAlNHMRIwEAYDVQQIDAlTaW5nYXBvcmUxEjAQBgNVBAcMCVNp
bmdhcG9yZTESMBAGA1UECgwJUFVQTcODwphTMSkwJwYDVQQDDCBzdGV3YXJkLnNv
bW1lbGllci5wdXBtb3MubmV0d29yazEiMCAGCSqGSIb3DQEJARYTcm9vdEBwdXBt
b3MubmV0d29yazAeFw0yMjA1MjQwNjE5NDdaFw0yNDA1MjMwNjE5NDdaMIGYMQsw
CQYDVQQGEwJTRzESMBAGA1UECAwJU2luZ2Fwb3JlMRIwEAYDVQQHDAlTaW5nYXBv
cmUxEjAQBgNVBAoMCVBVUE3Dg8KYUzEpMCcGA1UEAwwgc3Rld2FyZC5zb21tZWxp
ZXIucHVwbW9zLm5ldHdvcmsxIjAgBgkqhkiG9w0BCQEWE3Jvb3RAcHVwbW9zLm5l
dHdvcmswdjAQBgcqhkjOPQIBBgUrgQQAIgNiAATvBrmLqBCxHfWTt+h9tfrckH8g
YvOP4BCBTied5VU7x3pocN8v2JnvnW3+eN5q+TGqAZHooFwTQX1lDrXP+4x+Wdrj
9Dk4R6fP4ZJjDWP9LXmqW4x/f9VcKQxneRyFENmjUzBRMB0GA1UdDgQWBBSzeIQr
9DSejT8+PNSq0dMUHigOljAfBgNVHSMEGDAWgBSzeIQr9DSejT8+PNSq0dMUHigO
ljAPBgNVHRMBAf8EBTADAQH/MAoGCCqGSM49BAMDA2kAMGYCMQCqGKROfs1o6ptI
RKz4YAkQM4WKQH9ST1GONrJc5G+i6XHH4/fmL/nFzpWqJD6Odn4CMQDREtkJcOvJ
AB22WYnV23hsSHcOO73P1WcUtSJ6gLq5VYO01XksW7IbV7r6c78DJWU=
-----END CERTIFICATE-----
`

const LavenderFiveSubscriberCA = `-----BEGIN CERTIFICATE-----
MIICqDCCAi6gAwIBAgIUJdWH+cIDeAXzAzWvA1f0ybmnOPwwCgYIKoZIzj0EAwMw
gYoxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMRwwGgYDVQQKDBNM
YXZlbmRlci5GaXZlIE5vZGVzMSEwHwYDVQQDDBhzdGV3YXJkLmxhdmVuZGVyZml2
ZS5jb20xJTAjBgkqhkiG9w0BCQEWFmhlbGxvQGxhdmVuZGVyZml2ZS5jb20wHhcN
MjMwMTE4MTI0MDQ3WhcNMjUwMTE3MTI0MDQ3WjCBijELMAkGA1UEBhMCQVUxEzAR
BgNVBAgMClNvbWUtU3RhdGUxHDAaBgNVBAoME0xhdmVuZGVyLkZpdmUgTm9kZXMx
ITAfBgNVBAMMGHN0ZXdhcmQubGF2ZW5kZXJmaXZlLmNvbTElMCMGCSqGSIb3DQEJ
ARYWaGVsbG9AbGF2ZW5kZXJmaXZlLmNvbTB2MBAGByqGSM49AgEGBSuBBAAiA2IA
BGypluhCoE8a2oxbjaYGUtEOjmqJvBu4D/0jk6aioNc1G3huD2/jO8HAMzawedDH
oASy6nqSk+fGCjhsoB56NiBb95MmQRpivobAUB7Rqyz1Bd3rTmFAn3czodacwI1d
U6NTMFEwHQYDVR0OBBYEFABDT7eDMwbcM/ETpCScv+uDvoI+MB8GA1UdIwQYMBaA
FABDT7eDMwbcM/ETpCScv+uDvoI+MA8GA1UdEwEB/wQFMAMBAf8wCgYIKoZIzj0E
AwMDaAAwZQIxAJEtNYDXJpa63AuKLsn6cE5S6QK6gDpDDxVxqkTUxbHiJFtYgvUm
H0fy8w85H9nLYQIwAZ/1pMQJVO/9baDBlJeNga1cYXFi3+KAv6+eSe0uH1JCym/V
rPbOti9Ui2zE+3/p
-----END CERTIFICATE-----
`

const PolkachuSubscriberCA = `-----BEGIN CERTIFICATE-----
MIIC4DCCAmagAwIBAgIUHcEDyHKLg7DxRlcJOGbExDRn99YwCgYIKoZIzj0EAwMw
gaYxCzAJBgNVBAYTAlVTMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJ
bnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQxHTAbBgNVBAsMFHN0ZXdhcmQucG9sa2Fj
aHUuY29tMR0wGwYDVQQDDBRzdGV3YXJkLnBvbGthY2h1LmNvbTEhMB8GCSqGSIb3
DQEJARYSaGVsbG9AcG9sa2FjaHUuY29tMB4XDTIyMDYyMDAyMjYzOFoXDTI0MDYx
OTAyMjYzOFowgaYxCzAJBgNVBAYTAlVTMRMwEQYDVQQIDApTb21lLVN0YXRlMSEw
HwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQxHTAbBgNVBAsMFHN0ZXdh
cmQucG9sa2FjaHUuY29tMR0wGwYDVQQDDBRzdGV3YXJkLnBvbGthY2h1LmNvbTEh
MB8GCSqGSIb3DQEJARYSaGVsbG9AcG9sa2FjaHUuY29tMHYwEAYHKoZIzj0CAQYF
K4EEACIDYgAESPHHmuU+a6ks0jAtei47fJ4PlGDTv4Dep6uUW8qkiJYcE2yxFCcv
DiPARfZ1t2sEZ4ukvACI8ynpMgKhzpaSqukqL5cCDiFZ+ZY425+SoHoaRq+yp+E/
f+rcyHw0k604o1MwUTAdBgNVHQ4EFgQUZDwBKGcI0AgutfvELrRkP4e5bz8wHwYD
VR0jBBgwFoAUZDwBKGcI0AgutfvELrRkP4e5bz8wDwYDVR0TAQH/BAUwAwEB/zAK
BggqhkjOPQQDAwNoADBlAjEAjYQ4+2R0ls9HIJjjeLhaxiTmcSVIoSlczntIQZjE
HGX6nbnJEFDl+zDz40FwCeXLAjB5H3XQ7uaxwlj4fMDarIk8Rem1gJer1MnJmunv
q0681JGP9hLU0SeXy0G4qwADZRc=
-----END CERTIFICATE-----
`

const StakecitoSubscriberCA = `-----BEGIN CERTIFICATE-----
MIICbzCCAfagAwIBAgIUNF3FmpWSUhkR9Ei0EplK3vywOlkwCgYIKoZIzj0EAwMw
bzELMAkGA1UEBhMCREUxEzARBgNVBAgMClNvbWUtU3RhdGUxEjAQBgNVBAoMCVN0
YWtlY2l0bzESMBAGA1UECwwJU3Rha2VjaXRvMSMwIQYDVQQDDBpzdGV3YXJkLnN0
YWtlc2FuZHN0b25lLmNvbTAeFw0yMjA1MjYwODM4MjRaFw0yNDA1MjUwODM4MjRa
MG8xCzAJBgNVBAYTAkRFMRMwEQYDVQQIDApTb21lLVN0YXRlMRIwEAYDVQQKDAlT
dGFrZWNpdG8xEjAQBgNVBAsMCVN0YWtlY2l0bzEjMCEGA1UEAwwac3Rld2FyZC5z
dGFrZXNhbmRzdG9uZS5jb20wdjAQBgcqhkjOPQIBBgUrgQQAIgNiAAQhv4yVeP3H
LGleaAgtv8Nv+0qMKvXImDeywDVuLntE4OsPwNu2PpShT0Ksk1JkwFZs3NgCU/Hz
VZ26g1fI1Dyjkk40AXCX+rUqLlch8qM7IwZ7qHOAxP/3UFxUERCu0xSjUzBRMB0G
A1UdDgQWBBTmoTdY3c3gIq0UDyFn87m4hTHF+DAfBgNVHSMEGDAWgBTmoTdY3c3g
Iq0UDyFn87m4hTHF+DAPBgNVHRMBAf8EBTADAQH/MAoGCCqGSM49BAMDA2cAMGQC
MBpAnQ0xgp6mSX1P6ll937FnKgZKA1NAGac5bKXSDzPtuR5V6DtsCK96eoSLah65
ugIwdFtkaqpm8eeTcwWNxrbtl7mXyr9dp4qGbM1KQYFR7v6lGGmYrCfdfJhPT87c
yJqh
-----END CERTIFICATE-----
`

const ChorusOneSubscriberCA = `-----BEGIN CERTIFICATE-----
MIICzjCCAlSgAwIBAgIUYVmGGtT24lipwycw/tt92fK6C18wCgYIKoZIzj0EAwMw
gZ0xCzAJBgNVBAYTAkNIMQwwCgYDVQQIDANadWcxDTALBgNVBAcMBEJhYXIxFjAU
BgNVBAoMDUNob3J1cyBPbmUgQUcxFzAVBgNVBAsMDkluZnJhc3RydWN0dXJlMRMw
EQYDVQQDDApjaG9ydXMub25lMSswKQYJKoZIhvcNAQkBFhx0ZWNob3BzK3NvbW1l
bGllckBjaG9ydXMub25lMB4XDTIyMDUyNjE5NTQzMloXDTI0MDUyNTE5NTQzMlow
gZ0xCzAJBgNVBAYTAkNIMQwwCgYDVQQIDANadWcxDTALBgNVBAcMBEJhYXIxFjAU
BgNVBAoMDUNob3J1cyBPbmUgQUcxFzAVBgNVBAsMDkluZnJhc3RydWN0dXJlMRMw
EQYDVQQDDApjaG9ydXMub25lMSswKQYJKoZIhvcNAQkBFhx0ZWNob3BzK3NvbW1l
bGllckBjaG9ydXMub25lMHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEHjPQcxD30w/k
HhyNWBUDbjty1VAfF3jWuOeCV6vhQAC4tjrANRndxciL0MkmVUdpuIWfQfwwsv0z
M+6bXvqc0b0gOcWcDmBoVxqXLwCpOLyp9B3+DuqgbjiAS7u8HFxio1MwUTAdBgNV
HQ4EFgQUkcv4PTRvZWZqevGtRLoJIGWlsk4wHwYDVR0jBBgwFoAUkcv4PTRvZWZq
evGtRLoJIGWlsk4wDwYDVR0TAQH/BAUwAwEB/zAKBggqhkjOPQQDAwNoADBlAjBw
dzW4iAIOsJCgqshFaBA0+mSIxDUmnMsngg2IK0C/uTHg+zrh9M+YJ5QrJzRsCfMC
MQCemfMGhRec90/hC+osODiyd6qU5QkLxi1hcmwBzriRpGIlUohbQDSuMtHSO+pI
cdA=
-----END CERTIFICATE-----
`

const ImperatorSubscriberCA = `-----BEGIN CERTIFICATE-----
MIICpjCCAiygAwIBAgIUNyJMRHi7p+7g9FVkLQ5mrKAw4aswCgYIKoZIzj0EAwMw
gYkxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJ
bnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQxHTAbBgNVBAMMFHN0ZXdhcmQuaW1wZXJh
dG9yLmNvMSMwIQYJKoZIhvcNAQkBFhRjb250YWN0QGltcGVyYXRvci5jbzAeFw0y
MjA1MjcyMzM5NDRaFw0yNDA1MjYyMzM5NDRaMIGJMQswCQYDVQQGEwJBVTETMBEG
A1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50ZXJuZXQgV2lkZ2l0cyBQdHkg
THRkMR0wGwYDVQQDDBRzdGV3YXJkLmltcGVyYXRvci5jbzEjMCEGCSqGSIb3DQEJ
ARYUY29udGFjdEBpbXBlcmF0b3IuY28wdjAQBgcqhkjOPQIBBgUrgQQAIgNiAASm
CXKJvuAWw+g22f8AHZAkth3MxYWrNMQkZy37tK8h10hQkAVCVWSDoVXkuvcJuop4
rYdzoyS85pY/5xlywaNvQ6pMHRA5EJi0x8B1kdg1zpPuUV2K4rnoeQxGwiyKbTqj
UzBRMB0GA1UdDgQWBBQJzTYDuy3gNOzQ0PKuCO5liub37TAfBgNVHSMEGDAWgBQJ
zTYDuy3gNOzQ0PKuCO5liub37TAPBgNVHRMBAf8EBTADAQH/MAoGCCqGSM49BAMD
A2gAMGUCMFeX4Alu3mPjGagxE9M3Sub7INxppi2a1Td4kmUA9j+ZmCL9xQq2Bhu6
O3FhyNvkoQIxAPkyVm7znMlRymHc6tT0AigjRqv/fHmLXQi2QxYE9plLSVzSuzzH
cQMDmlZncNhv0g==
-----END CERTIFICATE-----
`

const TekuSubscriberCA = `-----BEGIN CERTIFICATE-----
MIICjjCCAhSgAwIBAgIUGQD0x2z2zd2bpdn2ctlzz/9H38gwCgYIKoZIzj0EAwMw
fjELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxDTALBgNVBAoMBFRl
a3UxHzAdBgNVBAMMFnNvbW1lbGllci50ZWt1Lm5ldHdvcmsxKjAoBgkqhkiG9w0B
CQEWG3Rla3Uuc3Rha2luZ0Bwcm90b25tYWlsLmNvbTAeFw0yMjA1MTgwNDMzMzNa
Fw0yNDA1MTcwNDMzMzNaMH4xCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0
YXRlMQ0wCwYDVQQKDARUZWt1MR8wHQYDVQQDDBZzb21tZWxpZXIudGVrdS5uZXR3
b3JrMSowKAYJKoZIhvcNAQkBFht0ZWt1LnN0YWtpbmdAcHJvdG9ubWFpbC5jb20w
djAQBgcqhkjOPQIBBgUrgQQAIgNiAASuK7+43q/sw6vwOrY9qSJDOru8tNxfFMKL
/tlHsP77PMDrOkgp+KPbBL26gfZXc1n7Hpy6pSN40xD9I4BN1AvV3LTMqZ3fXKFd
7/BWfysX4cTsO4W6Ip70SaSv2oPqBM2jUzBRMB0GA1UdDgQWBBSFezl3NsJ+qdB7
48AEO39UGrpUfTAfBgNVHSMEGDAWgBSFezl3NsJ+qdB748AEO39UGrpUfTAPBgNV
HRMBAf8EBTADAQH/MAoGCCqGSM49BAMDA2gAMGUCMQDX4+6dws5YjbL0UvTGTjZA
jXUDkAjawNGrtsjDCbi93XrU9veDSAbjBinS36ig0sYCMDpR105PAY0G8FYXJD+W
qbc8gnvrA9v3urubAECl9AgE411rViVIUYaMQzSu4tRTvw==
-----END CERTIFICATE-----
`

const ForboleSubscriberCA = `-----BEGIN CERTIFICATE-----
MIIB4TCCAWigAwIBAgIUI3jkPC9tBB+gBtNpX4RgF4DFqp4wCgYIKoZIzj0EAwMw
KDEmMCQGA1UEAwwdc3Rld2FyZC5zb21tZWxpZXIuZm9yYm9sZS5jb20wHhcNMjIw
NTEzMDczMTQ5WhcNMjQwNTEyMDczMTQ5WjAoMSYwJAYDVQQDDB1zdGV3YXJkLnNv
bW1lbGllci5mb3Jib2xlLmNvbTB2MBAGByqGSM49AgEGBSuBBAAiA2IABJ0D6kKg
72NuQlTsHGsAgJ0QrhacCzq0/1ELLX06YFU5Wn7718Reo8dlmw2MJ4BlBBhAUJGi
DK/WWUg6Oe//2sst76S0XLW/uUPwUR5M29ynHyH7UiuxdIzgpj7hrtpEb6NTMFEw
HQYDVR0OBBYEFJBIz4UxsCMIXVl9qwIl7+NJzYLVMB8GA1UdIwQYMBaAFJBIz4Ux
sCMIXVl9qwIl7+NJzYLVMA8GA1UdEwEB/wQFMAMBAf8wCgYIKoZIzj0EAwMDZwAw
ZAIwWtrNasLUQjWIuMoDK9gmqGnbgdmT4upDDtRu3cXpAJFcE9SDCNQn5VnO+xTM
OFNwAjBCBRYa3W8Vxspmn5O8UCNryhbXBwt4q1ce2fPx2Uywx8vd8Fzvy0m7YAkH
ZgzN5W4=
-----END CERTIFICATE-----
`

const BoubouSubscriberCA = `-----BEGIN CERTIFICATE-----
MIICkDCCAhagAwIBAgIUF7G5HqVlFee6I76X1LDc9wUN+DowCgYIKoZIzj0EAwMw
fzELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMRYwFAYDVQQHDA1TYW4gRnJhbmNp
c2NvMRMwEQYDVQQKDApDeWduaSBMYWJzMRMwEQYDVQQLDApCb3VCb3VOb2RlMSEw
HwYJKoZIhvcNAQkBFhJpbmZvQGN5Z25pbGFicy5jb20wHhcNMjIwNzA5MTIyMjQ4
WhcNMjQwNzA4MTIyMjQ4WjB/MQswCQYDVQQGEwJVUzELMAkGA1UECAwCQ0ExFjAU
BgNVBAcMDVNhbiBGcmFuY2lzY28xEzARBgNVBAoMCkN5Z25pIExhYnMxEzARBgNV
BAsMCkJvdUJvdU5vZGUxITAfBgkqhkiG9w0BCQEWEmluZm9AY3lnbmlsYWJzLmNv
bTB2MBAGByqGSM49AgEGBSuBBAAiA2IABAGbsPh3MVQ5+EKQdEHMHQgsiuG/SYhA
9epiqLsLAR+TTxSO87OLeFDN8ciq5RZW/6lV3qREFiBkMwIVz87H9XTt3BxhdTTt
Q/V1fP6gLbrA2cTdC139oOqXSxl0N/m6ZaNTMFEwHQYDVR0OBBYEFNnzknRvyXwK
fSX+y7LerlWurBRCMB8GA1UdIwQYMBaAFNnzknRvyXwKfSX+y7LerlWurBRCMA8G
A1UdEwEB/wQFMAMBAf8wCgYIKoZIzj0EAwMDaAAwZQIxAOpjjxkU8aRBffoMiTkU
Gv0GbdJrM0PPRuA2dOvFrhwmIUuykONx8El+3MeTcT6J2AIwdg0ueHNu+bSxY+wc
7pHbfqR6sGS/IIpUiV1qCW9JI714lnt1AM5BSbJUHswk6Rmn
-----END CERTIFICATE-----
`

const SleepyKittenSubscriberCA = `-----BEGIN CERTIFICATE-----
MIICYTCCAeagAwIBAgIUASnfEzEt6dvHbpTAvNCCtlJKlPUwCgYIKoZIzj0EAwMw
ZzELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxFTATBgNVBAoMDFNs
ZWVweUtpdHRlbjEsMCoGA1UEAwwjc3Rld2FyZC5zb21tZWxpZXIuc2xlZXB5a2l0
dGVuLmluZm8wHhcNMjIwNzE4MTMwMDI2WhcNMjQwNzE3MTMwMDI2WjBnMQswCQYD
VQQGEwJBVTETMBEGA1UECAwKU29tZS1TdGF0ZTEVMBMGA1UECgwMU2xlZXB5S2l0
dGVuMSwwKgYDVQQDDCNzdGV3YXJkLnNvbW1lbGllci5zbGVlcHlraXR0ZW4uaW5m
bzB2MBAGByqGSM49AgEGBSuBBAAiA2IABOsFJjyyhPDzzbx1eTS+WS5M1dQYBVxJ
EQlL9xtVvTXsMPdLa/OCQwrhZ467dfISnRw9jE2NB9qDY1EQWPJbOsMDbU8k726A
EiAUAeJdKLHf2guvlgF8rh1PBGWR4UMnu6NTMFEwHQYDVR0OBBYEFBi/u/WxU2ps
rKW6jxM4vPPMcsx9MB8GA1UdIwQYMBaAFBi/u/WxU2psrKW6jxM4vPPMcsx9MA8G
A1UdEwEB/wQFMAMBAf8wCgYIKoZIzj0EAwMDaQAwZgIxAKoLCqK/pTW7sClM5YaO
leROfs7rCp1eq8aceetqiRXdGHm7GJDBYHzewd8bLM6kdQIxAPdLBbpcSbgAX+R7
4+4iKN725+x6hP8z/ebf6bhY/1ARmxO85+RLjo4IapEtXBHnWQ==
-----END CERTIFICATE-----
`

const EverstakeSubscriberCA = `-----BEGIN CERTIFICATE-----
MIIC3TCCAmSgAwIBAgIUD8wVRwhFpOqk6SZGqBSnvxncSnMwCgYIKoZIzj0EAwMw
gaUxCzAJBgNVBAYTAlVBMQ0wCwYDVQQIDARLWUlWMQ0wCwYDVQQHDARLWUlWMRIw
EAYDVQQKDAlFdmVyc3Rha2UxEjAQBgNVBAsMCUV2ZXJzdGFrZTEoMCYGA1UEAwwf
c29tbWVsaWVyLXN0ZXdhcmQuZXZlcnN0YWtlLm9uZTEmMCQGCSqGSIb3DQEJARYX
YS5ieWNoa292QGV2ZXJzdGFrZS5vbmUwHhcNMjIwNzE1MTAzNzA3WhcNMjQwNzE0
MTAzNzA3WjCBpTELMAkGA1UEBhMCVUExDTALBgNVBAgMBEtZSVYxDTALBgNVBAcM
BEtZSVYxEjAQBgNVBAoMCUV2ZXJzdGFrZTESMBAGA1UECwwJRXZlcnN0YWtlMSgw
JgYDVQQDDB9zb21tZWxpZXItc3Rld2FyZC5ldmVyc3Rha2Uub25lMSYwJAYJKoZI
hvcNAQkBFhdhLmJ5Y2hrb3ZAZXZlcnN0YWtlLm9uZTB2MBAGByqGSM49AgEGBSuB
BAAiA2IABKMpiMX+aaKi96j1IGzFgwZj93R8m5m6v5OU6IluCfKqbtQKVBlCg6jV
Z6TnTYRs+0xMEN0TuuzkaiIwjKzeTZzotkz/wVWjwTJcGoHfO+D9HFsTFZcXUnEd
nRWaOOmf4qNTMFEwHQYDVR0OBBYEFEdTR/YTV2+W9+GSdcDS5h+XjwO8MB8GA1Ud
IwQYMBaAFEdTR/YTV2+W9+GSdcDS5h+XjwO8MA8GA1UdEwEB/wQFMAMBAf8wCgYI
KoZIzj0EAwMDZwAwZAIwQ32DlklQkbcz7lFEpw+FCzh/FpN9qZflnoRuaVnqK4tB
Xst+v3tImFlfH2B+Y0CwAjBiFLdczKLdhZzfj1FGyqGaX7adds/DTH6TuXIG+ccT
pC4FDMh3VxtCYDx0QeR78yg=
-----END CERTIFICATE-----
`

const TesselatedSubscriberCA = `-----BEGIN CERTIFICATE-----
MIICtDCCAjqgAwIBAgIUb8i6J7xdTLzgNPV7OXlKSyCJARYwCgYIKoZIzj0EAwMw
gZAxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMR0wGwYDVQQKDBRU
ZXNzZWxsYXRlZCBHZW9tZXRyeTEfMB0GA1UEAwwWc29tbWVsaWVyLnRlc3NhZ2Vv
LmNvbTEsMCoGCSqGSIb3DQEJARYdaGVsbG9AdGVzc2VsbGF0ZWRnZW9tZXRyeS5j
b20wHhcNMjIwODA4MDMwOTEwWhcNMjQwODA3MDMwOTEwWjCBkDELMAkGA1UEBhMC
QVUxEzARBgNVBAgMClNvbWUtU3RhdGUxHTAbBgNVBAoMFFRlc3NlbGxhdGVkIEdl
b21ldHJ5MR8wHQYDVQQDDBZzb21tZWxpZXIudGVzc2FnZW8uY29tMSwwKgYJKoZI
hvcNAQkBFh1oZWxsb0B0ZXNzZWxsYXRlZGdlb21ldHJ5LmNvbTB2MBAGByqGSM49
AgEGBSuBBAAiA2IABAbBGDp23qbfjyh1/nZUbi9a+f31XT5p3GGX0IwVvHaF8GBa
p+ut7upgtXkKyO81ykraX9Oybn/d0ygQ2ZxVGL48RlrRZEjwcZgqZa2/by77YofO
LrOTR2OaqVoJHAtBhaNTMFEwHQYDVR0OBBYEFDO9k3B5tLC4D71Dc2fD/eXEdNPe
MB8GA1UdIwQYMBaAFDO9k3B5tLC4D71Dc2fD/eXEdNPeMA8GA1UdEwEB/wQFMAMB
Af8wCgYIKoZIzj0EAwMDaAAwZQIxANfqkQNi5tZ8ufakos+vP4cYWLLwPiaSLmnz
Ppqtxd5JnbMr2NIwuwrqh/7f1SJJ2gIwWwxhl2h8ZV1O3r+5HMBqXkIUg5p4vBHO
Pl6yQFjHgCjEJb0bDx3n/qHpzmj9loIK
-----END CERTIFICATE-----
`

const ZtakeSubscriberCA = `-----BEGIN CERTIFICATE-----
MIICtDCCAjqgAwIBAgIUa54gP36t+ms+IAqQlh1EtabbumEwCgYIKoZIzj0EAwMw
gZAxCzAJBgNVBAYTAlVTMQswCQYDVQQIDAJXQTEQMA4GA1UEBwwHU2VhdHRsZTES
MBAGA1UECgwJWnRha2Uub3JnMRIwEAYDVQQLDAladGFrZS5vcmcxHDAaBgNVBAMM
E3NvbW1lbGllci56dGFrZS5vcmcxHDAaBgkqhkiG9w0BCQEWDW1heEB6dGFrZS5v
cmcwHhcNMjIwNzA0MTM0NzA1WhcNMjQwNzAzMTM0NzA1WjCBkDELMAkGA1UEBhMC
VVMxCzAJBgNVBAgMAldBMRAwDgYDVQQHDAdTZWF0dGxlMRIwEAYDVQQKDAladGFr
ZS5vcmcxEjAQBgNVBAsMCVp0YWtlLm9yZzEcMBoGA1UEAwwTc29tbWVsaWVyLnp0
YWtlLm9yZzEcMBoGCSqGSIb3DQEJARYNbWF4QHp0YWtlLm9yZzB2MBAGByqGSM49
AgEGBSuBBAAiA2IABM6lPxWKbW5uen75VLnev3V5eg14oOkGrwNdMCdyya3Kq9AM
9Ikxb8wHrbz0t+Pux3SKeqLNsl2/S8n45yeRTd/8Ekh03Fy4J9m62K1vlUKCJq1r
P+UNYfcOr8VMpM18WKNTMFEwHQYDVR0OBBYEFHZrsbcSN2+Bz9JT/y6C/18l/6al
MB8GA1UdIwQYMBaAFHZrsbcSN2+Bz9JT/y6C/18l/6alMA8GA1UdEwEB/wQFMAMB
Af8wCgYIKoZIzj0EAwMDaAAwZQIxANN/OQLOp/ykhfxYR1QindLUdYSx3o23eaeW
FFfYVpW1+URH57JoIQIGjOowldDwbwIwElOS4lVQIoPbvxppU9XYY7KiAKc7QqnF
vELj6/thfVS5XZaoDfw+5DPgT/PUiRAd
-----END CERTIFICATE-----
`

const TwoBuckChuckSubscriberCA = `-----BEGIN CERTIFICATE-----
MIICHzCCAaSgAwIBAgIUIo1u3e6h6gO+8LIqUQjdjRQqpSQwCgYIKoZIzj0EAwMw
RjEVMBMGA1UECgwMMiBCdWNrIENodWNrMRAwDgYDVQQLDAdWaW50bmVyMRswGQYD
VQQDDBJ0d28tYnVjay1jaHVjay54eXowHhcNMjIwODIyMTUwOTQyWhcNMjQwODIx
MTUwOTQyWjBGMRUwEwYDVQQKDAwyIEJ1Y2sgQ2h1Y2sxEDAOBgNVBAsMB1ZpbnRu
ZXIxGzAZBgNVBAMMEnR3by1idWNrLWNodWNrLnh5ejB2MBAGByqGSM49AgEGBSuB
BAAiA2IABBpeouUcSK/buu0ybgZJqnMzrO1+ei788OcF7U5hsY8NFzoC3w8DLzO4
vei9FcDZ/59M8z2V3eOA7MmQNg4BzIqUki/cvJzW/qeW77HFFALx4KyFwtB4iZ+F
FgP1DC8u8KNTMFEwHQYDVR0OBBYEFBWUcaVoaltx+m48gvu9entY0WFRMB8GA1Ud
IwQYMBaAFBWUcaVoaltx+m48gvu9entY0WFRMA8GA1UdEwEB/wQFMAMBAf8wCgYI
KoZIzj0EAwMDaQAwZgIxAOm8mrvPP+hrBVL42svpUbIn6oR7jaGiR6Eg3nudmhCv
nOFSWHXpTZ1kRavbwgfB7gIxAIzHW7oljiS6xShGA/FB6FoU9voP2mbwjIu96F0W
ADIfk5v+dX8gSZz98c7MXVhcwQ==
-----END CERTIFICATE-----
`

const CosmostationSubscriberCA = `-----BEGIN CERTIFICATE-----
MIIC4DCCAmagAwIBAgIUJcVZEz/992d79MrF68xLlDlaJWcwCgYIKoZIzj0EAwMw
gaYxCzAJBgNVBAYTAmtyMQ4wDAYDVQQIDAVzZW91bDEOMAwGA1UEBwwFc2VvdWwx
EDAOBgNVBAoMB3N0YW1wZXIxFTATBgNVBAsMDGNvc21vc3RhdGlvbjEqMCgGA1UE
Awwhc3Rld2FyZC5zb21tZWxpZXIuY29zbW9zdGF0aW9uLmlvMSIwIAYJKoZIhvcN
AQkBFhNqaHBAc3RhbXBlci5uZXR3b3JrMB4XDTIyMDkxMjE3NTM0M1oXDTI0MDkx
MTE3NTM0M1owgaYxCzAJBgNVBAYTAmtyMQ4wDAYDVQQIDAVzZW91bDEOMAwGA1UE
BwwFc2VvdWwxEDAOBgNVBAoMB3N0YW1wZXIxFTATBgNVBAsMDGNvc21vc3RhdGlv
bjEqMCgGA1UEAwwhc3Rld2FyZC5zb21tZWxpZXIuY29zbW9zdGF0aW9uLmlvMSIw
IAYJKoZIhvcNAQkBFhNqaHBAc3RhbXBlci5uZXR3b3JrMHYwEAYHKoZIzj0CAQYF
K4EEACIDYgAE5WNSh0qrkdUdDN6mM/gTI18ZQR5/vaIaXwsz52Seq7v/t332tSfJ
kIAlISo3upB/25xwA5e8YPkhv+aFp2xwLSpSWa4bHGxp/tADof0wcekvc/iraLCf
HHqEfNzl2MSjo1MwUTAdBgNVHQ4EFgQU59JtqB5qeXJlYEZeg8LwlQ3zO2AwHwYD
VR0jBBgwFoAU59JtqB5qeXJlYEZeg8LwlQ3zO2AwDwYDVR0TAQH/BAUwAwEB/zAK
BggqhkjOPQQDAwNoADBlAjEAtnfI51s/kGRVFNMko6u1xNFqLRM/4X6x32ZnlgMy
KHoSDLQF5PQ4HpxHU0QJRWyVAjAGy0nC3Eu6gKEwLlWA8iwPQNkHvtdkPx4cRp1n
8JBvc7RZYuJIlj1pFJKotq5aPjY=
-----END CERTIFICATE-----
`

const MCBSubscriberCA = `-----BEGIN CERTIFICATE-----
MIICrzCCAjagAwIBAgIUG1YKz0sRvAN+S9IhqXU3HB+hJXgwCgYIKoZIzj0EAwMw
gY4xCzAJBgNVBAYTAlRSMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJ
bnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQxHjAcBgNVBAMMFXNvbW1lbGllci5tY2Ju
b2RlLmNvbTEnMCUGCSqGSIb3DQEJARYYYWxpcG9zdGFjaTIwMDJAZ21haWwuY29t
MB4XDTIyMDkxMzE1MTgyMFoXDTI0MDkxMjE1MTgyMFowgY4xCzAJBgNVBAYTAlRS
MRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRz
IFB0eSBMdGQxHjAcBgNVBAMMFXNvbW1lbGllci5tY2Jub2RlLmNvbTEnMCUGCSqG
SIb3DQEJARYYYWxpcG9zdGFjaTIwMDJAZ21haWwuY29tMHYwEAYHKoZIzj0CAQYF
K4EEACIDYgAEgFYPLQ66ZdzR6G6XDOj+TjbkleU11DE+4XrFgvHg+MdDFxHPBMbv
a/rlzkxksSN4igGfFUeJdfL7T86wILnb9xFVoOshNj1eQP3+5iv1sxtnlMr7736c
cwt043ctEuXEo1MwUTAdBgNVHQ4EFgQU57p0chL+kfqniazKYK+yBj/GjpswHwYD
VR0jBBgwFoAU57p0chL+kfqniazKYK+yBj/GjpswDwYDVR0TAQH/BAUwAwEB/zAK
BggqhkjOPQQDAwNnADBkAjBuin50vjaxpVSvxDkcM+9MkxvPbviikLfzPocLQuIV
HkS3a7xGgSZ1jjEQ0vqaWcwCMEWFo3vNDl9sQIjcqMjkgXp4J6utSE9Y6/OqepDC
RkvwN2Y05FeXKRAAXTU8LFaFww==
-----END CERTIFICATE-----
`

const PolychainSubscriberCA = `-----BEGIN CERTIFICATE-----
MIICITCCAaigAwIBAgIUFWbRTNdAc8z1c3RE/L2TAA5NaYIwCgYIKoZIzj0EAwMw
SDESMBAGA1UECgwJUG9seWNoYWluMRIwEAYDVQQLDAlQb2x5Y2hhaW4xHjAcBgNV
BAMMFXNvbW1lbGllci51bml0NDEwLmNvbTAeFw0yMjA5MTkxMzQ2MDFaFw0yNDA5
MTgxMzQ2MDFaMEgxEjAQBgNVBAoMCVBvbHljaGFpbjESMBAGA1UECwwJUG9seWNo
YWluMR4wHAYDVQQDDBVzb21tZWxpZXIudW5pdDQxMC5jb20wdjAQBgcqhkjOPQIB
BgUrgQQAIgNiAAT3Bauu2PpZLd9B2cODPsGpxUmxdbL0SKxZ9c6fD5FMvdz8zqX8
qKYqp1kQKB/Xlafm9PsKQRpaDK6Ig4uVLU9pjgPeCui7s1I6PNeqCPUeF+tjadIT
VxIn4u1EZhs5EZGjUzBRMB0GA1UdDgQWBBQI5gnmFumLGAEV4zsYshhOxmzBLDAf
BgNVHSMEGDAWgBQI5gnmFumLGAEV4zsYshhOxmzBLDAPBgNVHRMBAf8EBTADAQH/
MAoGCCqGSM49BAMDA2cAMGQCMGhuqnTUPr5CsJgE1jZxd0v2VparzRnWtNtKyXzm
kilvyvpYJ+W3pr1FH6sWXVUangIwYSIX/V6rl/fv0GVv2XCZGKlbQBTkktYqHjM9
Iwfsl0+dmrxQjQGp1sjVRINI5aCh
-----END CERTIFICATE-----
`

const KingSuperSubscriberCA = `-----BEGIN CERTIFICATE-----
MIICXDCCAeKgAwIBAgIUS9oqAt/1SHtSknLPcrKuOEHg61AwCgYIKoZIzj0EAwMw
ZTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxEjAQBgNVBAoMCUtp
bmdTdXBlcjEtMCsGA1UEAwwkc29tbWVsaWVyLXN0ZXdhcmQua2luZ3N1cGVyLnNl
cnZpY2VzMB4XDTIyMTIwOTEzMDgwNFoXDTI0MTIwODEzMDgwNFowZTELMAkGA1UE
BhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxEjAQBgNVBAoMCUtpbmdTdXBlcjEt
MCsGA1UEAwwkc29tbWVsaWVyLXN0ZXdhcmQua2luZ3N1cGVyLnNlcnZpY2VzMHYw
EAYHKoZIzj0CAQYFK4EEACIDYgAEkhXcptFkLziwhrnFNbk6YGWqznGgRse3ATkb
8nVFyEllFF2d+TcOAp8+ZSKfGJpdWbDH8AtGar8bCL2vgdrnbiFmntQXiV4WrdW4
y3T3+8R7vklQhhsPYAMH0UgBpKdgo1MwUTAdBgNVHQ4EFgQUZSsOTSTIamzodSzy
qheSxilNigMwHwYDVR0jBBgwFoAUZSsOTSTIamzodSzyqheSxilNigMwDwYDVR0T
AQH/BAUwAwEB/zAKBggqhkjOPQQDAwNoADBlAjEAwZW3l12yezPJCjgqHrQdnJkP
nz79NmyIjr2CLDWky+yXOY3f0hzzBaF6IpnLQOiOAjBhkXbOFEpBxrLC28OkhiCp
9TPgdYtwbimW1gvklpVyo3R0x3uV3+Cvnm36dzSfYaA=
-----END CERTIFICATE-----
`

const ChillValidationSubscriberCA = `-----BEGIN CERTIFICATE-----
MIICjzCCAhagAwIBAgIUQo3gqDMnItT/j2Advvz2+6jWwbIwCgYIKoZIzj0EAwMw
fzELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxGTAXBgNVBAoMEENo
aWxsIFZhbGlkYXRpb24xGTAXBgNVBAsMEENoaWxsIFZhbGlkYXRpb24xJTAjBgNV
BAMMHHN0ZXdhcmQwLmNoaWxsdmFsaWRhdGlvbi5jb20wHhcNMjIwNTExMTk0MDAw
WhcNMjQwNTEwMTk0MDAwWjB/MQswCQYDVQQGEwJBVTETMBEGA1UECAwKU29tZS1T
dGF0ZTEZMBcGA1UECgwQQ2hpbGwgVmFsaWRhdGlvbjEZMBcGA1UECwwQQ2hpbGwg
VmFsaWRhdGlvbjElMCMGA1UEAwwcc3Rld2FyZDAuY2hpbGx2YWxpZGF0aW9uLmNv
bTB2MBAGByqGSM49AgEGBSuBBAAiA2IABPrjqlnjrQbX6i1Pg4GEeVPcOvU58kbL
L5uzQWgKdebF9uPK4ZgkWG+m0uFprXdheYDnhxb8XbyuXlzP9/XkUYV5H58i8onW
EChF+B2wHtS7X/CXUHqWZMhRmvfSmbl+pqNTMFEwHQYDVR0OBBYEFN+fJLjRFOy+
oJshVYJVIfqK+kFLMB8GA1UdIwQYMBaAFN+fJLjRFOy+oJshVYJVIfqK+kFLMA8G
A1UdEwEB/wQFMAMBAf8wCgYIKoZIzj0EAwMDZwAwZAIwCEGXsCLNjeZkN6pPPxFm
Yq/G4pdo5pLgiwjs6NWsagDdPyJf6ECrout/ey7NCtSKAjBIjrNquroiygiMvy5D
vV350cotQ7exnFORUvktxX4c5cD57jG5/YKxkD9UcC941ws=
-----END CERTIFICATE-----
`

const ChainnodesSubscriberCA = `-----BEGIN CERTIFICATE-----
MIICtTCCAjygAwIBAgIUPwHC5ZWW0uiixgJhdjihn+uwF+AwCgYIKoZIzj0EAwMw
gZExCzAJBgNVBAYTAklOMRMwEQYDVQQIDApTb21lLVN0YXRlMRMwEQYDVQQKDApD
aGFpbm5vZGVzMRIwEAYDVQQLDAlWYWxpZGF0b3IxHzAdBgNVBAMMFnN0ZXdhcmQu
Y2hhaW5ub2Rlcy5uZXQxIzAhBgkqhkiG9w0BCQEWFGFkbWluQGNoYWlubm9kZXMu
bmV0MB4XDTIyMDkyNjE4MDAxNVoXDTI0MDkyNTE4MDAxNVowgZExCzAJBgNVBAYT
AklOMRMwEQYDVQQIDApTb21lLVN0YXRlMRMwEQYDVQQKDApDaGFpbm5vZGVzMRIw
EAYDVQQLDAlWYWxpZGF0b3IxHzAdBgNVBAMMFnN0ZXdhcmQuY2hhaW5ub2Rlcy5u
ZXQxIzAhBgkqhkiG9w0BCQEWFGFkbWluQGNoYWlubm9kZXMubmV0MHYwEAYHKoZI
zj0CAQYFK4EEACIDYgAEWx5yoQhMOS2D+X6gXE6rZ+zJfqhGeusV7Dbd5OziI5AM
Nkh9z/i04oRADtJLS9+Vz7ratbZfwH6U7DoH/ipmwhwBV43KEWFgsdNFk9W0dKlM
zY8VZ6YbwxJAQTSs36Kvo1MwUTAdBgNVHQ4EFgQUW+eeO9Qhdk+6AMEn6y8Z3uvq
hoEwHwYDVR0jBBgwFoAUW+eeO9Qhdk+6AMEn6y8Z3uvqhoEwDwYDVR0TAQH/BAUw
AwEB/zAKBggqhkjOPQQDAwNnADBkAjB3goXF5si6LXLJHzUZzCS87dBeHNYTBiqE
aDVyNtI76m468xADG693yvtji/KUC0oCMDEVoLGmMSpCyQ1+/uwJq76rdiUT8wbB
1JTxEIKJu/HNjxkgB96c8YuUl5HTm4ks7Q==
-----END CERTIFICATE-----
`

const SevenSeasSubscriberCA = `-----BEGIN CERTIFICATE-----
MIICnDCCAiKgAwIBAgIUSKjVXB4CFNCeBK/OoTYdHxT2TwgwCgYIKoZIzj0EAwMw
gYQxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMQ4wDAYDVQQKDAU3
U2VhczEoMCYGA1UEAwwfc3Rld2FyZC5zb21tZWxpZXIuN3NlYXMuY2FwaXRhbDEm
MCQGCSqGSIb3DQEJARYXdmFsaWRhdG9yQDdzZWFzLmNhcGl0YWwwHhcNMjMwMjAy
MjEyNTI5WhcNMjUwMjAxMjEyNTI5WjCBhDELMAkGA1UEBhMCQVUxEzARBgNVBAgM
ClNvbWUtU3RhdGUxDjAMBgNVBAoMBTdTZWFzMSgwJgYDVQQDDB9zdGV3YXJkLnNv
bW1lbGllci43c2Vhcy5jYXBpdGFsMSYwJAYJKoZIhvcNAQkBFhd2YWxpZGF0b3JA
N3NlYXMuY2FwaXRhbDB2MBAGByqGSM49AgEGBSuBBAAiA2IABB/DCjJuwHmSTf5d
swbMlT2Ymfezd0PK0O7uh1SIDoeDlB5+OhvsHAs3wd62qVWymerknG0FdtwBEh6q
6JSh2qCzSE48Ffpjw1VkOp+9NIFVEKTbYpxZnFINoctpuE+koaNTMFEwHQYDVR0O
BBYEFJvy0Kyy3LNrANv6N6juZKqdOw03MB8GA1UdIwQYMBaAFJvy0Kyy3LNrANv6
N6juZKqdOw03MA8GA1UdEwEB/wQFMAMBAf8wCgYIKoZIzj0EAwMDaAAwZQIwc7wP
4jZ6MCq1QIlrTGW+aFGtg5Xo21AVPxY7QU2F+Z2C1n8Dh/8ICnyChCsrs+sFAjEA
+8beG8Zyz8+MC8geQT/pOBFjYMo+zS0/WqxiuXWCHGLDedOApOWAv3xwiuaHfvRX
-----END CERTIFICATE-----
`
