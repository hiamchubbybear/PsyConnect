#  Oauth2 Worksflows
Commit ID : d13b5faa89ef1ce6277053065f83b646435d00b0 branch : dev_chessy version: v1
Flow include OICD -> OICD Data + AccessToken
check -> redirect into info/Google + access token res -> json include status + jwt token
Commit ID : 69d18b3f5b757b0e06dcf7e417f5fddb98a7ce31 branch : dev_chessy version: v2(latest)
Flow without OICD -> Gen a access token -> redirect + access token -------> to client listen
-> check and get authentication -> save account to db with session token code or recheck and throw-> redirect second time + generated session code + email + provider   ->> json include  jwt
