# PsyConnect API Documentation


## Identity / Authentication Service

### Internal API

#### [GET] Valid Token
**URL:** `{{base_url}}/auth/internal/invalid/eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiJkYXZpZCIsImFjY291bnRJZCI6Ijg5NGNkMzM2LWY4MzUtNGEyZS05ZmE1LWY0OTVjYzliMTI2OCIsInByb2ZpbGVJZCI6ImE3MzI4NjEyLTY5MGItNDA0Zi05MzNjLTk0MTQwOGNkMGEwYyIsInNjb3BlIjoicm9sZS5jbGllbnQ6cGVybWlzc2lvbiB1c2VyLnBvc3QudmlldyB1c2VyLm1lc3NhZ2Uuc2VuZCB1c2VyLmFwcG9pbnRtZW50LmJvb2siLCJpc3MiOiJQc3lDb25uZWN0IEF1dGhlbnRpY2F0aW9uIFNlcnZpY2UiLCJleHAiOjE3NDUwMDc0OTgsInR5cGUiOiJOT1JNQUwiLCJpYXQiOjE3NDQ5OTY2OTgsImp0aSI6IjUzMjhkMmNkLTE4OTktNDZlOC1hZDZlLTJlZTE1ZjcxNTRkZSJ9.HNLWmHp0gZ2pH2E9iBZ7mpc1ZKQaJPin5QzxJqQSgaJDQFpNhX1Hj9HBoUjWeKF7EJ2fYdQdk-punWwcHnBocg`

### Example Response 200 OK
```json
true
```

### External API

#### Token

##### [GET] Claims
**URL:** `{{base_url}}/auth/token/claims?token=eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiJjaGVzc3kxNjAzIiwiYWNjb3VudElkIjoiNDk4MGVkNGYtZTRlZS00OTdlLTlhODMtZTkwOTAxZTQ1ZTlkIiwicHJvZmlsZUlkIjoiZjY2ZGNjYjAtM2JhOS00ZTVkLWJjYzYtOTNmYTkwOWZiNTkzIiwic2NvcGUiOiJyb2xlLmFkbWluOnBlcm1pc3Npb24gdGhlcmFwaXN0LnBvc3QuZGVsZXRlIHRoZXJhcGlzdC5tZXNzYWdlLnJlc3BvbmQgdXNlci5hcHBvaW50bWVudC5ib29rIGFkbWluLnVzZXIuZGVsZXRlIHRoZXJhcGlzdC5hcHBvaW50bWVudC5tYW5hZ2UgdXNlci5ibG9nLmNvbW1lbnQgdXNlci5kYXNoYm9hcmQuYWNjZXNzIHRoZXJhcGlzdC51c2VyLnByb2ZpbGUudmlldyBhZG1pbi5hcHBvaW50bWVudC52aWV3IGFkbWluLnVzZXIucHJvZmlsZS52aWV3IHVzZXIudGhlcmFwaXN0LnJhdGUgYWRtaW4udXNlci5jcmVhdGUgdXNlci5wb3N0LnZpZXcgYWRtaW4ucmVwb3J0LnZpZXcgdXNlci5wcm9maWxlLnZpZXcgYWRtaW4uZGFzaGJvYXJkLmFjY2VzcyBhZG1pbi5tZXNzYWdlLnZpZXcgdGhlcmFwaXN0LnBvc3QuY3JlYXRlIGFkbWluLm5vdGlmaWNhdGlvbi5tYW5hZ2UgYWRtaW4ucm9sZS5tYW5hZ2UgdGhlcmFwaXN0LnBvc3QuZWRpdCBhZG1pbi5jb250ZW50Lm1hbmFnZSB1c2VyLm1lc3NhZ2Uuc2VuZCBhZG1pbi51c2VyLmVkaXQgdXNlci5ibG9nLnJlYWQgdGhlcmFwaXN0LnVzZXIucmF0ZSB1c2VyLnByb2ZpbGUudXBkYXRlIHRoZXJhcGlzdC5kYXNoYm9hcmQuYWNjZXNzIGFkbWluLmFuYWx5dGljcy52aWV3IiwiaXNzIjoiUHN5Q29ubmVjdCBBdXRoZW50aWNhdGlvbiBTZXJ2aWNlIiwiZXhwIjoxNzY5ODM0ODQ3LCJ0eXBlIjoiTk9STUFMIiwiaWF0IjoxNzU0MjgyODQ3LCJqdGkiOiJmNDM5MjhlOC1kNzI2LTQ1ZDgtOTYxZC03NzhlMTlhMDBkOTciLCJwbGF0Zm9ybSI6Im1vYmlsZSJ9.AUcZDc2xEVIGObSbkUfVxYZda_1UzSPgvYoOf0YCH7rgVQg_-yYmxvnGTBbhcwhBVheG-XcrDKGFTYVFSfIVCA`

### Example Response 200 OK
```json
{
    "accountId": "4980ed4f-e4ee-497e-9a83-e90901e45e9d",
    "role": "admin:permission",
    "subject": "chessy1603",
    "profileId": "f66dccb0-3ba9-4e5d-bcc6-93fa909fb593",
    "scope": "role.admin:permission therapist.post.delete therapist.message.respond user.appointment.book admin.user.delete therapist.appointment.manage user.blog.comment user.dashboard.access therapist.user.profile.view admin.appointment.view admin.user.profile.view user.therapist.rate admin.user.create user.post.view admin.report.view user.profile.view admin.dashboard.access admin.message.view therapist.post.create admin.notification.manage admin.role.manage therapist.post.edit admin.content.manage user.message.send admin.user.edit user.blog.read therapist.user.rate user.profile.update therapist.dashboard.access admin.analytics.view",
    "expiration": "2026-01-31T04:47:27.000+00:00",
    "issuedAt": "2025-08-04T04:47:27.000+00:00",
    "type": "NORMAL",
    "platform": "mobile",
    "issuer": "PsyConnect Authentication Service",
    "jwtId": "f43928e8-d726-45d8-961d-778e19a00d97"
}
```

#### [POST] Create Account
**URL:** `{{base_url}}/identity/create`

**Request Body:**
```json
{
  "username": "{{$randomUserName}}",
  "password": "password123",
  "firstName": "{{$randomFirstName}}",
  "lastName": "{{$randomLastName}}",
  "dob": "2000-03-01",
  "address": "{{$randomStreetAddress}}",
  "gender": "Male",
  "email": "{{$randomExampleEmail}}",
  "role": "Therapist",
  "avatarUri": "{{$randomAbstractImage}}"
}

```

### Example Response 200 OK
```json
{
    "code": 200,
    "message": "Success",
    "data": {
        "username": "david3",
        "email": "ilovepakpak@gmail.com",
        "role": [
            {
                "roleId": "therapist",
                "name": "therapist:permission"
            }
        ],
        "imageUri": "helloworlds@gmail.com"
    }
}
```

#### [POST] Login
**URL:** `{{base_url}}/auth/login?loginType=NORMAL`

**Request Body:**
```json
{
    "username": "david",
    "password": "david123"
}
```

### Example Response 200 OK
```json
{
    "code": 200,
    "message": "Success",
    "data": {
        "token": "eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiJkYXZpZDMiLCJhY2NvdW50SWQiOiI4NmYzNzNiZS1iNDNkLTRmNWItOTUxNi04YjAyYWY5MGUzMjkiLCJwcm9maWxlSWQiOiJlYTQwOTM3NC0yZjJlLTRiNTUtYWVlNC1mNWUzM2FhNjg2MzkiLCJzY29wZSI6InJvbGUudGhlcmFwaXN0OnBlcm1pc3Npb24gdGhlcmFwaXN0LnBvc3QuY3JlYXRlIHRoZXJhcGlzdC5hcHBvaW50bWVudC5tYW5hZ2UgdGhlcmFwaXN0Lm1lc3NhZ2UucmVzcG9uZCB0aGVyYXBpc3QucG9zdC5lZGl0IiwiaXNzIjoiUHN5Q29ubmVjdCBBdXRoZW50aWNhdGlvbiBTZXJ2aWNlIiwiZXhwIjoxNzQ3NjczMTUxLCJ0eXBlIjoiTk9STUFMIiwiaWF0IjoxNzQ3NjYyMzUxLCJqdGkiOiI5MzgxNWI2OC05YWZlLTRhMGItOGMxYS02OTRhNzc5NWFjMzUifQ.T8a3Cv5OFxARboSK_SntBNo_GtHXj1JRQkeVtgvvlcVwXgrKswXL4u_cctTRubNKblM6Rtdr4dbwr6sHetwjWQ",
        "successful": true
    }
}
```

#### [POST] Introspect
**URL:** `{{base_url}}/auth/introspect`

**Request Body:**
```json
{
    "token" : "eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiJkYXZpZCIsImFjY291bnRJZCI6Ijg5NGNkMzM2LWY4MzUtNGEyZS05ZmE1LWY0OTVjYzliMTI2OCIsInByb2ZpbGVJZCI6ImE3MzI4NjEyLTY5MGItNDA0Zi05MzNjLTk0MTQwOGNkMGEwYyIsInNjb3BlIjoicm9sZS5jbGllbnQ6cGVybWlzc2lvbiB1c2VyLmFwcG9pbnRtZW50LmJvb2sgdXNlci5tZXNzYWdlLnNlbmQgdXNlci5wb3N0LnZpZXciLCJpc3MiOiJQc3lDb25uZWN0IEF1dGhlbnRpY2F0aW9uIFNlcnZpY2UiLCJleHAiOjE3NDUwMDk1NjUsInR5cGUiOiJOT1JNQUwiLCJpYXQiOjE3NDQ5OTg3NjUsImp0aSI6IjhmYjc4ZDUzLTFmNzQtNGJiMC1iMThlLTQwMDQwYmI1NmM1OSJ9.uesXXHfaxe-0EZWWQCnfsHHAL-J0NDoZO0VGz1QLCzDCSGbMzT7fk0FCba_yiF8525Z2VY_c-Ms9dbGnhISKTQ"
}
```

### Example Response 200 OK
```json
{
    "code": 200,
    "message": "Success",
    "data": {
        "successful": true
    }
}
```

#### [POST] Admin Login
**URL:** `{{base_url}}/auth/login?loginType=NORMAL`

**Request Body:**
```json
{
    "username": "chessy1603",
    "password": "16032004"
}
```

### Example Response 200 OK
```json
{
    "code": 200,
    "message": "Success",
    "data": {
        "token": "eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiJjaGVzc3kxNjAzIiwiYWNjb3VudElkIjoiNzZhNzM2YjYtYzQwNy00N2E0LWI3YWUtZGE1ZTQ5NjViZDA2IiwicHJvZmlsZUlkIjoiOGI0YzdkNGYtMWVhNS00MTk2LTgzMDAtZjZlYWE1YjYwMjYzIiwic2NvcGUiOiJyb2xlLmFkbWluOnBlcm1pc3Npb24gdXNlci5ibG9nLmNvbW1lbnQgdGhlcmFwaXN0LnVzZXIucHJvZmlsZS52aWV3IHVzZXIuYXBwb2ludG1lbnQuYm9vayB1c2VyLnByb2ZpbGUudmlldyB0aGVyYXBpc3QucG9zdC5lZGl0IHVzZXIudGhlcmFwaXN0LnJhdGUgYWRtaW4ubWVzc2FnZS52aWV3IHRoZXJhcGlzdC5hcHBvaW50bWVudC5tYW5hZ2UgYWRtaW4uZGFzaGJvYXJkLmFjY2VzcyBhZG1pbi5yb2xlLm1hbmFnZSBhZG1pbi5jb250ZW50Lm1hbmFnZSB0aGVyYXBpc3QucG9zdC5kZWxldGUgdGhlcmFwaXN0Lm1lc3NhZ2UucmVzcG9uZCB0aGVyYXBpc3QucG9zdC5jcmVhdGUgdXNlci5tZXNzYWdlLnNlbmQgdXNlci5wb3N0LnZpZXcgdXNlci5wcm9maWxlLnVwZGF0ZSBhZG1pbi51c2VyLmNyZWF0ZSBhZG1pbi5hbmFseXRpY3MudmlldyB1c2VyLmRhc2hib2FyZC5hY2Nlc3MgdXNlci5ibG9nLnJlYWQgdGhlcmFwaXN0LmRhc2hib2FyZC5hY2Nlc3MgYWRtaW4udXNlci5kZWxldGUgYWRtaW4udXNlci5wcm9maWxlLnZpZXcgYWRtaW4udXNlci5lZGl0IGFkbWluLmFwcG9pbnRtZW50LnZpZXcgdGhlcmFwaXN0LnVzZXIucmF0ZSBhZG1pbi5ub3RpZmljYXRpb24ubWFuYWdlIGFkbWluLnJlcG9ydC52aWV3IiwiaXNzIjoiUHN5Q29ubmVjdCBBdXRoZW50aWNhdGlvbiBTZXJ2aWNlIiwiZXhwIjoxNzQ2OTk3ODA0LCJ0eXBlIjoiTk9STUFMIiwiaWF0IjoxNzQ2OTg3MDA0LCJqdGkiOiJkNWNmYjlkNS01NjNmLTRkZTItOTk0NC1lYzY2MjhiN2Q0MDcifQ.seOs2H9yirlEN_PGM-79DYaPozOHI3sTWRy1Xd7enMPD-nltksePqhPWX86e6Ysf31VJC4iL4b12yqxzPHmmCA",
        "successful": true
    }
}
```

#### [POST] Activate
**URL:** `{{base_url}}/identity/activate`

**Request Body:**
```json
{
    "token" : "28198",
    "email" : "ilovepakpak@gmail.com",
    "verifedTime" : "1739698533"
}
```

### Example Response 401 Unauthorized
```json
{
    "code": 503,
    "message": "Activation failed"
}
```

#### [GET] Get Info
**URL:** `http://localhost:8888/account/info`

**Auth:** bearer

### Example Response 200 OK
```json
{
    "code": 200,
    "message": "Success",
    "data": {
        "accountId": "6799268f-2b25-42b7-9c35-2970af61116e",
        "username": "david2",
        "email": "hellohaha@gmail.com",
        "provider": "ORDINARY",
        "createdAt": "2025-06-01",
        "activated": true
    }
}
```

#### [POST] Update
**URL:** `localhost:8888/identity/update`

**Auth:** bearer

### Example Response 200 OK
```json
{
    "code": 200,
    "message": "Success",
    "data": {
        "username": "david2",
        "password": null,
        "email": "ilovepakpak11@gmail.com"
    }
}
```

#### [POST] Activate Code
**URL:** `{{base_url}}/identity/req/activate`

**Request Body:**
```json
{
    "username": "huytran",
    "email": "helloworlds@gmail.com",
    "fullname": "Tran Van Huy"
}
```

### Example Response 200 OK
```json
{
    "code": 200,
    "message": "Success",
    "data": true
}
```

#### [GET] AllAccount
**URL:** `{{base_url}}/account/all/0`

**Auth:** bearer

### Example Response 200 OK
```json
{
    "code": 200,
    "message": "Success",
    "data": {
        "content": [
            {
                "accountId": "012560d2-9b8f-4e83-80e4-16d624b14ecc",
                "profileId": "f3d1c2e9-e437-4f0c-8559-a2e5f06871bb",
                "username": "Mathew84",
                "password": "$2a$10$/GUtJosLcyW.xw168KjJDO3b9bUfZYYh4bb2EKIpM8To5g1VkFiKC",
                "email": "Leone.Runte@example.net",
                "provider": "ORDINARY",
                "createdAt": "2025-06-01T09:25:46.481+00:00",
                "session": null,
                "token": {
                    "username": "Mathew84",
                    "token": "82064",
                    "issuedAt": "2025-06-01T09:25:46.045+00:00",
                    "expires": "2025-06-01T09:35:46.045+00:00",
                    "revoked": false
                },
                "activated": false
            },
            {
                "accountId": "0eca1c1a-824b-4b11-9d6b-7a6c34db61c3",
                "profileId": "7b8a0e2c-7551-43d8-8db6-a01f99511ed6",
                "username": "Chandler31",
                "password": "$2a$10$zc.dll6E7KW9.UUcU31w2ORS78zFz0twTDYb0QNzwIEvHc9.av5.e",
                "email": "Richie_Vandervort3@example.net",
                "provider": "ORDINARY",
                "createdAt": "2025-06-01T09:24:15.353+00:00",
                "session": null,
                "token": {
                    "username": "Chandler31",
                    "token": "19774",
                    "issuedAt": "2025-06-01T09:24:14.523+00:00",
                    "expires": "2025-06-01T09:34:14.523+00:00",
                    "revoked": false
                },
                "activated": false
            },
            {
                "accountId": "6799268f-2b25-42b7-9c35-2970af61116e",
                "profileId": "6a03716a-0f0a-4c4a-9c0b-7c115d88fa76",
                "username": "david2",
                "password": "$2a$10$VfG..0fUF6Qne95MZBE4XOA90zEWaJwI41X1rQbYe6BBAHC8hHElK",
                "email": "ilovepakpak11@gmail.com",
                "provider": "ORDINARY",
                "createdAt": "2025-05-31T18:37:59.010+00:00",
                "session": null,
                "token": {
                    "username": "david2",
                    "token": "14133",
                    "issuedAt": "2025-05-31T18:37:56.520+00:00",
                    "expires": "2025-05-31T18:47:56.520+00:00",
                    "revoked": false
                },
                "activated": true
            },
            {
                "accountId": "68d510d9-7565-431d-8831-b1ff97e15b45",
                "profileId": "efef7300-6fa5-463d-b811-afd1abc70423",
                "username": "david3",
                "password": "$2a$10$PCYGrwamPWFM15Dt8tTA5OGlhd2L/srxAT6ESgAeR2n7qIu9TndMG",
                "email": "ilovepakpak@gmail.com",
                "provider": "ORDINARY",
                "createdAt": "2025-05-31T18:42:24.348+00:00",
                "session": null,
                "token": {
                    "username": "david3",
                    "token": "53584",
                    "issuedAt": "2025-05-31T18:42:21.765+00:00",
                    "expires": "2025-05-31T18:52:21.765+00:00",
                    "revoked": false
                },
                "activated": false
            },
            {
                "accountId": "8280da1b-ed50-4b42-a210-5e634f7b19ff",
                "profileId": "37a376be-3c56-46c0-9178-d7eb85167d12",
                "username": "david42222",
                "password": "$2a$10$RM3CLBdC6JkvRPQJEWapse./eggvfSz523pZ2mc426uI38NeDkdLG",
                "email": "ilovepakpa2k2213@gmail.com",
                "provider": "ORDINARY",
                "createdAt": "2025-06-01T05:07:56.097+00:00",
                "session": null,
                "token": {
                    "username": "david42222",
                    "token": "",
                    "issuedAt": "2025-06-01T05:07:55.015+00:00",
                    "expires": null,
                    "revoked": false
                },
                "activated": true
            },
            {
                "accountId": "bbc2c42e-8eb9-440a-b9fd-df2187ad6610",
                "profileId": "ca2bb845-64e4-40bf-8f35-ceefd8d7f715",
                "username": "david422",
                "password": "$2a$10$V/7./q01pRojQSq4Ck02yOhOr3b1jeMwOO9b3Hm/lUfC3wl5DH3ou",
                "email": "ilovepakpak213@gmail.com",
                "provider": "ORDINARY",
                "createdAt": "2025-05-31T18:55:56.569+00:00",
                "session": null,
                "token": {
                    "username": "david422",
                    "token": "80603",
                    "issuedAt": "2025-05-31T18:55:55.337+00:00",
                    "expires": "2025-05-31T19:05:55.337+00:00",
                    "revoked": false
                },
                "activated": false
            },
            {
                "accountId": "c2a08f9e-6e38-4abc-9285-5c3b7750d266",
                "profileId": "196d6d36-e564-4e20-a6ec-f8e5d91dcf15",
                "username": "david42",
                "password": "$2a$10$h3KONmEenEQqxMOTSI8IUeKSFtR6K3csuKMUwaanVvBUlCwbYlkY.",
                "email": "ilovepakpak21@gmail.com",
                "provider": "ORDINARY",
                "createdAt": "2025-05-31T18:53:28.129+00:00",
                "session": null,
                "token": {
                    "username": "david42",
                    "token": "65070",
                    "issuedAt": "2025-05-31T18:53:27.138+00:00",
                    "expires": "2025-05-31T19:03:27.138+00:00",
                    "revoked": false
                },
                "activated": false
            },
            {
                "accountId": "c9cc8b54-b81d-40a9-9e89-52d09f3f7fc8",
                "profileId": "04b2df4c-71bf-41e2-b61b-2ec8951f2f73",
                "username": "chessy1603",
                "password": "$2a$10$xygan879Kg1hCtytaAW5R.plZQ7OR2bBgu7YMMqXPVk8CK8/SZpHu",
                "email": "tranvanhuy160304@gmail.com",
                "provider": "GOOGLE",
                "createdAt": "2025-05-31T18:36:27.385+00:00",
                "session": "f17088db-9ca8-43a8-8b85-5998e7a9dc39",
                "token": null,
                "activated": true
            },
            {
                "accountId": "db1ba48d-0105-46e2-bbe2-acfb4044607f",
                "profileId": "f2432bee-b96d-4687-8481-f9253c158785",
                "username": "david4",
                "password": "$2a$10$Zi1FLOeZ87SEmZoTUHxyZOocH/gCMIwpKFNjZP7Nbc8cnF12FfGzW",
                "email": "ilovepakpak1@gmail.com",
                "provider": "ORDINARY",
                "createdAt": "2025-05-31T18:45:50.998+00:00",
                "session": null,
                "token": {
                    "username": "david4",
                    "token": "68764",
                    "issuedAt": "2025-05-31T18:45:48.189+00:00",
                    "expires": "2025-05-31T18:55:48.189+00:00",
                    "revoked": false
                },
                "activated": false
            },
            {
                "accountId": "f1d98900-1dfc-4c7f-b021-00e4dc4b8cc8",
                "profileId": "4f96c679-d493-4855-8f7f-e934e362edb8",
                "username": "david4222",
                "password": "$2a$10$eEav7w4ZWtB15PUrsfS/rOKTawLeR.Pyeco1FXKaraXhIAUUdKMlC",
                "email": "ilovepakpak2213@gmail.com",
                "provider": "ORDINARY",
                "createdAt": "2025-06-01T04:59:48.561+00:00",
                "session": null,
                "token": {
                    "username": "david4222",
                    "token": "74136",
                    "issuedAt": "2025-06-01T04:59:47.065+00:00",
                    "expires": "2025-06-01T05:09:47.065+00:00",
                    "revoked": false
                },
                "activated": true
            }
        ],
        "pageable": {
            "pageNumber": 0,
            "pageSize": 10,
            "sort": {
                "empty": true,
                "sorted": false,
                "unsorted": true
            },
            "offset": 0,
            "paged": true,
            "unpaged": false
        },
        "last": true,
        "totalPages": 1,
        "totalElements": 10,
        "size": 10,
        "sort": {
            "empty": true,
            "sorted": false,
            "unsorted": true
        },
        "number": 0,
        "first": true,
        "numberOfElements": 10,
        "empty": false
    }
}
```

### Deprecated

#### [GET] Oauth2
**URL:** `{{base_url}}/oauth2/userInfo/google`


#### [GET] Hello
**URL:** `{{base_url}}/identity/hello`


## Profile Service

### Internal API

#### [POST] Create
**URL:** `{{base_url}}/profile/internal/user`

**Request Body:**
```json
{
    "profileId": "p788",
    "accountId": "123",
    "firstName": "John",
    "lastName": "Doe",
    "dob": "1990-01-01",
    "address": "123 Main St, City",
    "gender": "Male",
    "avatarUri": "https://example.com/avatar.jpg"
}
```

### Example Response 200 OK
```json
{
    "code": 200,
    "message": "Success",
    "data": {
        "profileId": "c002",
        "firstName": "John",
        "lastName": "Doe",
        "dob": "1990-01-01",
        "address": "123 Main St, City",
        "gender": "Male",
        "avatarUri": "https://example.com/avatar.jpg"
    }
}
```

### External API

#### Mood

##### [POST] Add
**URL:** `{{base_url}}/mood/add`

**Headers:**
- `X-Profile-Id`: c003

**Request Body:**
```json
{
    "mood":"Having fun",
    "moodDescription" : "Soooosss2ssssssoooooo fucking funnn day",
    "visibility" : "PRIVATE"
}
```

**Auth:** bearer

### Example Response 200 OK
```json
{
    "code": 200,
    "message": "Success",
    "data": {
        "profileId": "t001",
        "mood": "Having fun",
        "moodDescription": "Soooosss2ssssssoooooo fucking funnn day",
        "timeExpires": "2025-05-31 17:56:58",
        "success": true
    }
}
```

##### [GET] Get
**URL:** `{{base_url}}/mood`

**Auth:** bearer

### Example Response 200 OK
```json
{
    "code": 200,
    "message": "Success",
    "data": {
        "moodId": "t001",
        "mood": "Having fun",
        "description": "Soooosss2ssssssoooooo fucking funnn day",
        "visibility": "PRIVATE",
        "createdAt": 1748602618332,
        "expiresAt": 1748689018332
    }
}
```

##### [PUT] Update
**URL:** `{{base_url}}/mood`

**Request Body:**
```json
{
    "mood":"Having fun",
    "moodDescription" : "Hiccc33 fucking funnn day",
    "visibility" : "PUBLIC"
}
```

**Auth:** bearer

### Example Response 200 OK
```json
{
    "code": 200,
    "message": "Success",
    "data": {
        "profileId": "t001",
        "mood": null,
        "moodDescription": "Hiccc33 fucking funnn day",
        "timeExpires": null,
        "success": true
    }
}
```

##### [DELETE] Delete
**URL:** `{{base_url}}/mood`

**Auth:** bearer

### Example Response 200 OK
```json
{
    "code": 200,
    "message": "Success",
    "data": {
        "message": "Delete success",
        "success": true
    }
}
```

##### [GET] Get All Friends Mood
**URL:** `{{base_url}}/mood/friends`

**Headers:**
- `X-Profile-Id`: p798

### Example Response 200 OK
```json
{
    "code": 200,
    "message": "Success",
    "data": [
        {
            "profileId": "t001",
            "fullName": "JohnDoe",
            "avatarUrl": "https://example.com/avatar.jpg",
            "moodId": "t001",
            "mood": "Having fun",
            "description": "Soooosss2ssssssoooooo fucking funnn day",
            "visibility": "PRIVATE",
            "createdAt": 1748603257807,
            "expiresAt": 1748689657807
        }
    ]
}
```

#### Setting

##### [GET] Get
**URL:** `{{base_url}}/user-setting`

**Auth:** bearer

### Example Response 200 OK
```json
{
    "code": 200,
    "message": "Success",
    "data": {
        "profileId": "6a03716a-0f0a-4c4a-9c0b-7c115d88fa76",
        "privacyLevel": "PRIVATE",
        "showLastSeen": true,
        "showProfilePicture": true,
        "showMood": false,
        "notificationsEnabled": false,
        "emailNotifications": false,
        "pushNotifications": false,
        "smsNotifications": false,
        "twoFactorAuth": false,
        "allowLoginAlerts": false,
        "trustedDevices": "",
        "language": "en",
        "theme": "light",
        "autoDeleteOldMoods": true
    }
}
```

##### [PUT] Update
**URL:** `{{base_url}}/user-setting`

**Request Body:**
```json
{
    "privacyLevel": "PUBLIC",
    "showLastSeen": true,
    "showProfilePicture": true,
    "showMood": true,
    "notificationsEnabled": true,
    "emailNotifications": true,
    "pushNotifications": true,
    "smsNotifications": true,
    "twoFactorAuth": true,
    "allowLoginAlerts": true,
    "trustedDevices": "device133,device2",
    "language": "en",
    "theme": "dark",
    "autoDeleteOldMoods": false
}
```

**Auth:** bearer

### Example Response 200 OK
```json
{
    "code": 200,
    "message": "Success",
    "data": {
        "profileId": "6a03716a-0f0a-4c4a-9c0b-7c115d88fa76",
        "success": true
    }
}
```

##### [POST] Default Setting
**URL:** `{{base_url}}/user-setting/default`

**Auth:** bearer

### Example Response 200 OK
```json
{
    "code": 200,
    "message": "Success",
    "data": {
        "profileId": "efef7300-6fa5-463d-b811-afd1abc70423",
        "privacyLevel": "PRIVATE",
        "showLastSeen": true,
        "showProfilePicture": true,
        "showMood": false,
        "notificationsEnabled": false,
        "emailNotifications": false,
        "pushNotifications": false,
        "smsNotifications": false,
        "twoFactorAuth": false,
        "allowLoginAlerts": false,
        "trustedDevices": "",
        "language": "en",
        "theme": "light",
        "autoDeleteOldMoods": true
    }
}
```

#### Friend

##### [POST] Add Friends
**URL:** `{{base_url}}/profile/friend/request`

**Headers:**
- `X-Profile-Id`: t001

**Request Body:**
```json
{
    "target":"c001"
}
```

### Example Response 200 OK
```json
{
    "code": 200,
    "message": "Success",
    "data": {
        "status": "Profile(profileId=6f253564-9490-4785-859f-fba48d8ec27e, accountId=8bb46417-4eda-46e8-babe-86f339142a03, username=null, firstName=David, lastName=Houser, dob=2000-03-01, address=21 St , New Mexico, gender=Male, avatarUri=helloworlds@gmail.com, description=Welcome to my wall!!, moodList=null, friends=[dev.psyconnect.profile_service.model.FriendRelationship@1b1362d2], activityLogs=[], settings=Setting(profileId=6f253564-9490-4785-859f-fba48d8ec27e, privacyLevel=PRIVATE, showLastSeen=true, showProfilePicture=true, showMood=false, notificationsEnabled=false, emailNotifications=false, pushNotifications=false, smsNotifications=false, twoFactorAuth=false, allowLoginAlerts=false, trustedDevices=, language=en, theme=light, autoDeleteOldMoods=true))",
        "friendRequestStatus": "Success"
    }
}
```

##### [POST] Accept
**URL:** `{{base_url}}/profile/friend/accept`

**Headers:**
- `X-Profile-Id`: 0c447c8e-0b96-4065-837d-7760aef06528

**Request Body:**
```json
{
    "target" : "3510a31f-82ec-49a0-ac93-b45328b375d2"
}
```

**Auth:** bearer

### Example Response 200 OK
```json
{
    "code": 200,
    "message": "Success",
    "data": {
        "status": "Success",
        "acceptStatus": "ACCEPTED"
    }
}
```

##### [GET] Get
**URL:** `localhost:8888/profile`

**Auth:** bearer

### Example Response 200 OK
```json
{
    "code": 200,
    "message": "Success",
    "data": [
        {
            "profile": {
                "profileId": "6a03716a-0f0a-4c4a-9c0b-7c115d88fa76",
                "username": null,
                "avatarUri": "helloworlds@gmail.com"
            },
            "mood": {
                "moodId": "6a03716a-0f0a-4c4a-9c0b-7c115d88fa76",
                "mood": "Having fun",
                "description": "Hiccc33 fucking funnn day",
                "visibility": "PRIVATE",
                "createdAt": 1748769855002,
                "expiresAt": 1748856255002
            }
        }
    ]
}
```

##### [POST] Unfriend/Unmatch
**URL:** `http://localhost:8081/profile/friend/unfriend`

**Headers:**
- `X-Profile-Id`: p799

**Request Body:**
```json
{
    "target" : "p789"
}
```

### Example Response 200 OK
```json
{
    "code": 200,
    "message": "Success",
    "data": {
        "status": "Success",
        "unfriendStatus": "ACCEPTED"
    }
}
```

##### [POST] Decline Friend Request
**URL:** `{{base_url}}/profile/friend/unfriend`

**Headers:**
- `X-Profile-Id`: p799

**Request Body:**
```json
{
    "target" : "6a03716a-0f0a-4c4a-9c0b-7c115d88fa76"
}
```

**Auth:** bearer

### Example Response 400 Bad Request
```json
{
    "code": 203,
    "message": "You are not friend"
}
```

##### [POST] Undo
**URL:** `http://localhost:8081/profile/friend/undo`

### Example Response 200 OK
```json
{
    "code": 200,
    "message": "Success",
    "data": {
        "status": "Success",
        "undoStatus": "UNFRIEND"
    }
}
```

#### Profile

##### [POST] Update Profile
**URL:** `{{base_url}}/profile`

**Headers:**
- `X-Profile-Id`: c003

**Request Body:**
```json
{
  "username": "david2",
  "firstName": "Meo",
  "lastName": "Rapist",
  "dob": "15-05-2004",
  "address": "123 Therapy Street, Ha Noi",
  "gender": "Female",
  "avatarUri": "https://www.pinterest.com/pin/844493675294198/"
}

```

**Auth:** bearer


##### [GET] Get Profile Detailed
**URL:** `{{base_url}}/profile`

**Headers:**
- ``: 

**Auth:** bearer

### Example Response 200 OK
```json
{
    "code": 200,
    "message": "Success",
    "data": {
        "profile": {},
        "mood": {
            "moodId": "t001",
            "mood": "Having fun",
            "description": "Soooosss2ssssssoooooo fucking funnn day",
            "visibility": "PRIVATE",
            "createdAt": 1748602618332,
            "expiresAt": 1748689018332
        }
    }
}
```

### Deprecated

#### [GET] Admin All Profile
**URL:** `{{base_url}}/profile/friends`

**Auth:** bearer


#### [GET] All
**URL:** `{{base_url}}/profile/all?page=10&size=10`

**Auth:** bearer

### Example Response 200 OK
```json
{
    "code": 200,
    "message": "Success",
    "data": []
}
```

#### [DELETE] Delete Settings
**URL:** `{{base_url}}/user-setting`

**Auth:** bearer


#### [POST] Add
**URL:** `{{base_url}}/user-setting/add`

**Request Body:**
```json
{
    "privacyLevel": "PRIVATE",
    "showLastSeen": true,
    "showProfilePicture": true,
    "showMood": false,
    "notificationsEnabled": true,
    "emailNotifications": true,
    "pushNotifications": false,
    "smsNotifications": false,
    "twoFactorAuth": true,
    "allowLoginAlerts": true,
    "trustedDevices": "device1,device2",
    "language": "en",
    "theme": "dark",
    "autoDeleteOldMoods": false
}
```

**Auth:** bearer


## Consultation Service

### Internal

#### [POST] Recommend
**URL:** `{{base_url}}/consultation/client/recommend`

**Headers:**
- `X-Profile-Id`: c001
- `X-Roles`: role.client:permission user.therapist.rate admin.appointment.view therapist.post.create admin.user.edit user.profile.view admin.user.create therapist.post.delete user.blog.read user.dashboard.access therapist.post.edit admin.dashboard.access user.message.send therapist.user.rate user.appointment.book admin.message.view user.profile.update admin.role.manage therapist.message.respond admin.content.manage admin.report.view user.blog.comment therapist.dashboard.access admin.user.profile.view admin.analytics.view therapist.user.profile.view user.post.view therapist.appointment.manage admin.user.delete admin.notification.manage

### Example Response 200 OK
```json
{
    "message": "success",
    "status": 200,
    "data": [
        {
            "client_id": "c001",
            "therapist_id": "t013",
            "points": 56.21,
            "reasons": [
                "Matched language",
                "Matched availability"
            ],
            "status": "pending",
            "created_at": "2025-07-05T17:03:12.454905Z"
        },
        {
            "client_id": "c001",
            "therapist_id": "t001",
            "points": 54.95,
            "reasons": [
                "Matched language",
                "Matched availability"
            ],
            "status": "pending",
            "created_at": "2025-07-05T17:03:12.454905Z"
        },
        {
            "client_id": "c001",
            "therapist_id": "t012",
            "points": 44.07,
            "reasons": [
                "Matched language"
            ],
            "status": "pending",
            "created_at": "2025-07-05T17:03:12.454905Z"
        },
        {
            "client_id": "c001",
            "therapist_id": "t014",
            "points": 38.79,
            "reasons": [
                "Matched availability"
            ],
            "status": "pending",
            "created_at": "2025-07-05T17:03:12.454905Z"
        },
        {
            "client_id": "c001",
            "therapist_id": "t002",
            "points": 0,
            "reasons": [
                "Therapist not available"
            ],
            "status": "pending",
            "created_at": "2025-07-05T17:03:12.454905Z"
        },
        {
            "client_id": "c001",
            "therapist_id": "t003",
            "points": 0,
            "reasons": [
                "Therapist not available"
            ],
            "status": "pending",
            "created_at": "2025-07-05T17:03:12.454905Z"
        },
        {
            "client_id": "c001",
            "therapist_id": "t004",
            "points": 0,
            "reasons": [
                "Therapist not available"
            ],
            "status": "pending",
            "created_at": "2025-07-05T17:03:12.454905Z"
        },
        {
            "client_id": "c001",
            "therapist_id": "t005",
            "points": 0,
            "reasons": [
                "Therapist not available"
            ],
            "status": "pending",
            "created_at": "2025-07-05T17:03:12.454905Z"
        },
        {
            "client_id": "c001",
            "therapist_id": "t006",
            "points": 0,
            "reasons": [
                "Therapist not available"
            ],
            "status": "pending",
            "created_at": "2025-07-05T17:03:12.454905Z"
        },
        {
            "client_id": "c001",
            "therapist_id": "t007",
            "points": 0,
            "reasons": [
                "Therapist not available"
            ],
            "status": "pending",
            "created_at": "2025-07-05T17:03:12.454905Z"
        },
        {
            "client_id": "c001",
            "therapist_id": "t008",
            "points": 0,
            "reasons": [
                "Therapist not available"
            ],
            "status": "pending",
            "created_at": "2025-07-05T17:03:12.454905Z"
        },
        {
            "client_id": "c001",
            "therapist_id": "t009",
            "points": 0,
            "reasons": [
                "Therapist not available"
            ],
            "status": "pending",
            "created_at": "2025-07-05T17:03:12.454905Z"
        },
        {
            "client_id": "c001",
            "therapist_id": "t0010",
            "points": 0,
            "reasons": [
                "Therapist not available"
            ],
            "status": "pending",
            "created_at": "2025-07-05T17:03:12.454905Z"
        },
        {
            "client_id": "c001",
            "therapist_id": "t011",
            "points": 0,
            "reasons": [
                "Therapist not available"
            ],
            "status": "pending",
            "created_at": "2025-07-05T17:03:12.454905Z"
        }
    ]
}
```

### External

#### Therapist

##### [POST] Add V1
**URL:** `{{consultation_service}}/v1/consultation/therapist`

**Headers:**
- `X-Roles`: role.therapist:permission user.therapist.rate admin.appointment.view therapist.post.create admin.user.edit user.profile.view admin.user.create therapist.post.delete user.blog.read user.dashboard.access therapist.post.edit admin.dashboard.access user.message.send therapist.user.rate user.appointment.book admin.message.view user.profile.update admin.role.manage therapist.message.respond admin.content.manage admin.report.view user.blog.comment therapist.dashboard.access admin.user.profile.view admin.analytics.view therapist.user.profile.view user.post.view therapist.appointment.manage admin.user.delete admin.notification.manage
- `X-Profile-Id`: t002

**Request Body:**
```json
{
  "address": "123 Main Street, City, Country",
  "languages": ["English", "Spanish"],
  "specialization": ["Cognitive Behavioral Therapy", "Anxiety Management"],
  "consultation_modes": ["Online", "In-Person"],
  "experience": 10,
  "rating": 4.8,
  "price_per_session": 50,
  "availability": {
    "days": ["Monday", "Wednesday", "Friday"],
    "time_slots": ["10:00 AM - 12:00 PM", "2:00 PM - 4:00 PM"]
  }
}
```

**Auth:** bearer

### Example Response 200 OK
```json
{
    "message": "success",
    "status": 200,
    "data": {
        "profile_id": "t001",
        "address": "123 Main Street, City, Country",
        "languages": [
            "English",
            "Spanish"
        ],
        "specialization": [
            "Cognitdddddddddddive Behavioral Therapy",
            "Anxiety Management"
        ],
        "consultation_modes": [
            "Online",
            "In-Person"
        ],
        "experience": 10,
        "rating": 4.8,
        "availability": {
            "days": [
                "Monday",
                "Wednesday",
                "Friday"
            ],
            "time_slots": [
                "10:00 AM - 12:00 PM",
                "2:00 PM - 4:00 PM"
            ]
        }
    }
}
```

##### [PUT] Update
**URL:** `{{base_url}}/consultation/therapist`

**Headers:**
- `X-Profile-Id`: t001
- `X-Role`: role.client:permission user.therapist.rate admin.appointment.view therapist.post.create admin.user.edit user.profile.view admin.user.create therapist.post.delete user.blog.read user.dashboard.access therapist.post.edit admin.dashboard.access user.message.send therapist.user.rate user.appointment.book admin.message.view user.profile.update admin.role.manage therapist.message.respond admin.content.manage admin.report.view user.blog.comment therapist.dashboard.access admin.user.profile.view admin.analytics.view therapist.user.profile.view user.post.view therapist.appointment.manage admin.user.delete admin.notification.manage

**Request Body:**
```json
{
  "address": "123 Main dddddddddddddddddddddStreet, City, Country",
  "languages": ["English", "Spanish"],
  "specialization": ["Cognitive Behavioral Therapy", "Anxiety Management"],
  "consultation_modes": ["Online", "In-Person"],
  "experience": 10,
  "rating": 4.8,
  "price_per_session": 50,
  "availability": {
    "days": ["Monday", "Wednesday", "Friday"],
    "time_slots": ["10:00 AM - 12:00 PM", "2:00 PM - 4:00 PM"]
  }
}
```

**Auth:** bearer


##### [GET] Get
**URL:** `{{base_url}}/consultation/therapist`

**Headers:**
- `X-Profile-Id`: t001
- `X-Roles`: role.therapist:permission user.therapist.rate admin.appointment.view therapist.post.create admin.user.edit user.profile.view admin.user.create therapist.post.delete user.blog.read user.dashboard.access therapist.post.edit admin.dashboard.access user.message.send therapist.user.rate user.appointment.book admin.message.view user.profile.update admin.role.manage therapist.message.respond admin.content.manage admin.report.view user.blog.comment therapist.dashboard.access admin.user.profile.view admin.analytics.view therapist.user.profile.view user.post.view therapist.appointment.manage admin.user.delete admin.notification.manage

**Auth:** bearer

### Example Response 200 OK
```json
{
    "message": "success",
    "status": 200,
    "data": {
        "profile_id": "t001",
        "address": "123 Main Street, City, Country",
        "languages": [
            "English",
            "Spanish"
        ],
        "specialization": [
            "Cognitdddddddddddive Behavioral Therapy",
            "Anxiety Management"
        ],
        "consultation_modes": [
            "Online",
            "In-Person"
        ],
        "experience": 10,
        "rating": 4.8,
        "availability": {
            "days": [
                "Monday",
                "Wednesday",
                "Friday"
            ],
            "time_slots": [
                "10:00 AM - 12:00 PM",
                "2:00 PM - 4:00 PM"
            ]
        }
    }
}
```

##### [PUT] Change status
**URL:** `{{base_url}}/consultation/therapist/status/false`

**Headers:**
- `X-Profile-Id`: t001
- `X-Roles`: role.therapist:permission user.therapist.rate admin.appointment.view therapist.post.create admin.user.edit user.profile.view admin.user.create therapist.post.delete user.blog.read user.dashboard.access therapist.post.edit admin.dashboard.access user.message.send therapist.user.rate user.appointment.book admin.message.view user.profile.update admin.role.manage therapist.message.respond admin.content.manage admin.report.view user.blog.comment therapist.dashboard.access admin.user.profile.view admin.analytics.view therapist.user.profile.view user.post.view therapist.appointment.manage admin.user.delete admin.notification.manage

### Example Response 200 OK
```json
{
    "message": "success",
    "status": 200,
    "data": false
}
```

##### [PUT] Update
**URL:** `{{base_url}}/consultation/therapist/profile`

**Headers:**
- `X-Profile-Id`: p789

**Request Body:**
```json
{
  "address": "123 Main dddddddddddddddddddddStreet, City, Country",
  "languages": ["English", "Spanish"],
  "specialization": ["Cognitive Behavioral Therapy", "Anxiety Management"],
  "consultation_modes": ["Online", "In-Person"],
  "experience": 10,
  "rating": 4.8,
  "price_per_session": 50,
  "availability": {
    "days": ["Monday", "Wednesday", "Friday"],
    "time_slots": ["10:00 AM - 12:00 PM", "2:00 PM - 4:00 PM"]
  }
}
```

**Auth:** noauth


#### Client

##### Recommendation

###### [GET] Recommend
**URL:** ``

### Example Response 200 OK
```json
{
    "message": "success",
    "status": 200,
    "data": [
        {
            "profile_id": "t013",
            "address": "123 Main Street, City, Country",
            "languages": [
                "English"
            ],
            "specialization": [
                "Depression Treatment",
                "Stress Management"
            ],
            "consultation_modes": [
                "In-Person"
            ],
            "experience": 15,
            "rating": 5,
            "currency": "USD",
            "is_available": true,
            "availability": {
                "days": [
                    "Monday",
                    "Friday"
                ],
                "time_slots": [
                    "9:00 AM - 11:00 AM"
                ]
            }
        },
        {
            "profile_id": "t001",
            "address": "123 Main Street, City, Country",
            "languages": [
                "English",
                "Spanish"
            ],
            "specialization": [
                "Cognitdddddddddddive Behavioral Therapy",
                "Anxiety Management"
            ],
            "consultation_modes": [
                "Online",
                "In-Person"
            ],
            "experience": 10,
            "rating": 4.8,
            "currency": "",
            "is_available": true,
            "availability": {
                "days": [
                    "Monday",
                    "Wednesday",
                    "Friday"
                ],
                "time_slots": [
                    "10:00 AM - 12:00 PM",
                    "2:00 PM - 4:00 PM"
                ]
            }
        },
        {
            "profile_id": "t012",
            "address": "456 Another Street, Other City, Country",
            "languages": [
                "English",
                "French"
            ],
            "specialization": [
                "EMDR",
                "PTSD"
            ],
            "consultation_modes": [
                "Online"
            ],
            "experience": 6,
            "rating": 4.6,
            "currency": "USD",
            "is_available": true,
            "availability": {
                "days": [
                    "Tuesday",
                    "Thursday"
                ],
                "time_slots": [
                    "3:00 PM - 5:00 PM"
                ]
            }
        },
        {
            "profile_id": "t014",
            "address": "123 Main Street, City, Country",
            "languages": [
                "Spanish"
            ],
            "specialization": [
                "Anxiety Management"
            ],
            "consultation_modes": [
                "Online",
                "In-Person"
            ],
            "experience": 8,
            "rating": 4.2,
            "currency": "USD",
            "is_available": true,
            "availability": {
                "days": [
                    "Wednesday"
                ],
                "time_slots": [
                    "10:00 AM - 12:00 PM"
                ]
            }
        },
        {
            "profile_id": "t002",
            "address": "123 Main Street, City, Country",
            "languages": [
                "English",
                "Spanish"
            ],
            "specialization": [
                "Cognitdddddddddddive Behavioral Therapy",
                "Anxiety Management"
            ],
            "consultation_modes": [
                "Online",
                "In-Person"
            ],
            "experience": 10,
            "rating": 4.8,
            "currency": "",
            "availability": {
                "days": [
                    "Monday",
                    "Wednesday",
                    "Friday"
                ],
                "time_slots": [
                    "10:00 AM - 12:00 PM",
                    "2:00 PM - 4:00 PM"
                ]
            }
        }
    ]
}
```

###### [GET] Update therapist filter
**URL:** `{{base_url}}/consultation/client/recommend`

**Auth:** bearer

### Example Response 200 OK
```json
{
    "message": "success",
    "status": 200,
    "data": [
        {
            "client_id": "d0f4bb83-5cd7-4928-91d7-0ca2a78a30f1",
            "therapist_id": "therapist_10",
            "points": 52.83,
            "reasons": [
                "Matched language",
                "Matched availability"
            ],
            "status": "pending",
            "created_at": "2025-08-11T09:50:19.249756Z"
        },
        {
            "client_id": "d0f4bb83-5cd7-4928-91d7-0ca2a78a30f1",
            "therapist_id": "therapist_8",
            "points": 52.06,
            "reasons": [
                "Matched language",
                "Matched availability"
            ],
            "status": "pending",
            "created_at": "2025-08-11T09:50:19.249756Z"
        },
        {
            "client_id": "d0f4bb83-5cd7-4928-91d7-0ca2a78a30f1",
            "therapist_id": "therapist_6",
            "points": 51.29,
            "reasons": [
                "Matched language",
                "Matched availability"
            ],
            "status": "pending",
            "created_at": "2025-08-11T09:50:19.249756Z"
        },
        {
            "client_id": "d0f4bb83-5cd7-4928-91d7-0ca2a78a30f1",
            "therapist_id": "therapist_4",
            "points": 50.52,
            "reasons": [
                "Matched language",
                "Matched availability"
            ],
            "status": "pending",
            "created_at": "2025-08-11T09:50:19.249756Z"
        },
        {
            "client_id": "d0f4bb83-5cd7-4928-91d7-0ca2a78a30f1",
            "therapist_id": "therapist_2",
            "points": 49.76,
            "reasons": [
                "Matched language",
                "Matched availability"
            ],
            "status": "pending",
            "created_at": "2025-08-11T09:50:19.249756Z"
        },
        {
            "client_id": "d0f4bb83-5cd7-4928-91d7-0ca2a78a30f1",
            "therapist_id": "therapist_1",
            "points": 0,
            "reasons": [
                "Therapist not available"
            ],
            "status": "pending",
            "created_at": "2025-08-11T09:50:19.249756Z"
        },
        {
            "client_id": "d0f4bb83-5cd7-4928-91d7-0ca2a78a30f1",
            "therapist_id": "therapist_3",
            "points": 0,
            "reasons": [
                "Therapist not available"
            ],
            "status": "pending",
            "created_at": "2025-08-11T09:50:19.249756Z"
        },
        {
            "client_id": "d0f4bb83-5cd7-4928-91d7-0ca2a78a30f1",
            "therapist_id": "therapist_5",
            "points": 0,
            "reasons": [
                "Therapist not available"
            ],
            "status": "pending",
            "created_at": "2025-08-11T09:50:19.249756Z"
        },
        {
            "client_id": "d0f4bb83-5cd7-4928-91d7-0ca2a78a30f1",
            "therapist_id": "therapist_7",
            "points": 0,
            "reasons": [
                "Therapist not available"
            ],
            "status": "pending",
            "created_at": "2025-08-11T09:50:19.249756Z"
        },
        {
            "client_id": "d0f4bb83-5cd7-4928-91d7-0ca2a78a30f1",
            "therapist_id": "therapist_9",
            "points": 0,
            "reasons": [
                "Therapist not available"
            ],
            "status": "pending",
            "created_at": "2025-08-11T09:50:19.249756Z"
        }
    ]
}
```

##### [POST] Add
**URL:** `{{base_url}}/consultation/client`

**Headers:**
- `X-Profile-Id`: t001

**Request Body:**
```json
{
  "address": "123 Main Street, City, Country",
  "languages": ["English", "French"],
  "issue_detail": ["Anxiety", "Stress Management"],
  "consultation_modes": ["Online", "In-Person"],
  "range_price_per_hour": 50,
  "availability": {
    "days": ["Monday", "Wednesday", "Friday"],
    "time_slots": ["10:00-11:00", "14:00-15:00"]
  },
  "preferred_therapist_gender": "Female",
  "experience_level": "Senior",
  "therapist_specialization": ["CBT", "Mindfulness Therapy"],
  "urgency_level": "Medium",
  "session_duration": 60,
  "preferred_therapist_language": ["English", "Spanish"],
  "is_flexible_with_schedule": true,
  "current_session": ["session123", "session456"]
}
```

**Auth:** bearer

### Example Response 200 OK
```json
{
    "message": "success",
    "status": 200,
    "data": {
        "profile_id": "c001",
        "address": "123 Main Street, City, Country",
        "languages": [
            "English",
            "French"
        ],
        "issue_detail": [
            "Anxiety",
            "Stress Management"
        ],
        "consultation_modes": [
            "Online",
            "In-Person"
        ],
        "availability": {
            "days": [
                "Monday",
                "Wednesday",
                "Friday"
            ],
            "time_slots": [
                "10:00-11:00",
                "14:00-15:00"
            ]
        },
        "preferred_therapist_gender": "Female",
        "experience_level": "Senior",
        "urgency_level": "Medium",
        "session_duration": 60,
        "preferred_therapist_language": [
            "English",
            "Spanish"
        ],
        "is_flexible_with_schedule": true,
        "current_session": [
            "session123",
            "session456"
        ]
    }
}
```

##### [PUT] Update
**URL:** `{{base_url}}/consultation/client`

**Headers:**
- `X-Profile-Id`: p789

**Request Body:**
```json
{
  "address": "1asdddddddddddddddd23 Main dddddddddddddddddddddStreet, City, Country",
  "languages": ["English", "Spanish"],
  "specialization": ["Cognitive Behavioral Therapy", "Anxiety Management"],
  "consultation_modes": ["Online", "In-Person"],
  "experience": 10,
  "rating": 4.8,
  "price_per_session": 50,
  "availability": {
    "days": ["Monday", "Wednesday", "Friday"],
    "time_slots": ["10:00 AM - 12:00 PM", "2:00 PM - 4:00 PM"]
  }
}
```

**Auth:** bearer


##### [GET] Get
**URL:** `{{base_url}}/consultation/client`

**Headers:**
- `X-Profile-Id`: p789
- `X-Roles`: role.client:permission user.therapist.rate admin.appointment.view therapist.post.create admin.user.edit user.profile.view admin.user.create therapist.post.delete user.blog.read user.dashboard.access therapist.post.edit admin.dashboard.access user.message.send therapist.user.rate user.appointment.book admin.message.view user.profile.update admin.role.manage therapist.message.respond admin.content.manage admin.report.view user.blog.comment therapist.dashboard.access admin.user.profile.view admin.analytics.view therapist.user.profile.view user.post.view therapist.appointment.manage admin.user.delete admin.notification.manage

**Auth:** bearer


#### Matching

##### [GET] Get Matching
**URL:** `{{base_url}}/consultation/therapist/match?page=1`

**Headers:**
- `X-Profile-Id`: p799

**Request Body:**
```json
{
  "age": 27,
  "gender": "female",
  "address": "Hanoi",
  "languages": ["vi", "en"],
  "issues": ["relationship", "stress"],
  "preferences": {
    "consultant_gender": "female",
    "mode": "online",
    "budget": 350,
    "experience_level": "intermediate",
    "specializations": ["relationship", "career"]
  },
  "availability": {
    "days": ["Monday", "Wednesday"],
    "time_slots": ["10:00", "14:00"]
  },
  "is_flexible_with_schedule": false
}

```

### Example Response 200 OK
```json
{
    "message": "success",
    "status": 200,
    "data": null
}
```

##### [POST] MatchingRequest
**URL:** `{{base_url}}/consultation/client/match`

**Headers:**
- `X-Roles`: role.client:permission user.therapist.rate admin.appointment.view therapist.post.create admin.user.edit user.profile.view admin.user.create therapist.post.delete user.blog.read user.dashboard.access therapist.post.edit admin.dashboard.access user.message.send therapist.user.rate user.appointment.book admin.message.view user.profile.update admin.role.manage therapist.message.respond admin.content.manage admin.report.view user.blog.comment therapist.dashboard.access admin.user.profile.view admin.analytics.view therapist.user.profile.view user.post.view therapist.appointment.manage admin.user.delete admin.notification.manage
- `X-Profile-Id`: c001

**Request Body:**
```json
{
  "client_id": "c001",
  "therapist_id": "t001",
  "source": "swipe",
  "swipe_score": 87.5,
  "reasons": [
    "Matched language",
    "Same availability"
  ]
}

```

### Example Response 200 OK
```json
{
    "message": "success",
    "status": 200,
    "data": true
}
```

##### [POST] Accept
**URL:** `{{base_url}}/consultation/therapist/match/response`

**Headers:**
- `X-Roles`: role.therapist:permission user.therapist.rate admin.appointment.view therapist.post.create admin.user.edit user.profile.view admin.user.create therapist.post.delete user.blog.read user.dashboard.access therapist.post.edit admin.dashboard.access user.message.send therapist.user.rate user.appointment.book admin.message.view user.profile.update admin.role.manage therapist.message.respond admin.content.manage admin.report.view user.blog.comment therapist.dashboard.access admin.user.profile.view admin.analytics.view therapist.user.profile.view user.post.view therapist.appointment.manage admin.user.delete admin.notification.manage
- `X-Profile-Id`: c001

**Request Body:**
```json
{
    "clientId": "t001",
    "option": "accept"
}
```

### Example Response 504 Gateway Timeout
```json
{
    "message": "rpc error: code = AlreadyExists desc = Already friend or requested",
    "status": 504
}
```

#### Session

##### [POST] Create New Session
**URL:** `{{base_url}}/consultation/session`

**Request Body:**
```json
{
    "therapist_id": "t001",
    "client_id": "c001",
    "languages": ["Vietnamese"],
    "specialization": ["anxiety"],
    "consultation_modes": ["Online"],
    "price_per_hour": 200,
    "session_time": {
        "day" : "Monday",
        "start_time" : "11:30",
        "end_time" : "13:30"
    }
}
```

### Example Response 404 Not Found
```json
{
    "message": "Invalid input",
    "status": 404
}
```

##### [DELETE] DeleteSession
**URL:** `{{base_url}}/consultation/session`

### Example Response 200 OK
```json
{
    "message": "success",
    "status": 200,
    "data": {
        "message": "Session deleted successfully",
        "session_id": "68373494bcee2460eeee27b6"
    }
}
```

##### [GET] Get By Session ID
**URL:** `{{base_url}}/consultation/session/68396af212f8c83cec320ada`

### Example Response 200 OK
```json
{
    "message": "success",
    "status": 200,
    "data": {
        "therapist_id": "t001",
        "client_id": "c001",
        "languages": [
            "Vietnamese"
        ],
        "specialization": [
            "anxiety"
        ],
        "consultation_modes": [
            "Online"
        ],
        "price_per_hour": 200,
        "session_time": {
            "day": "Monday",
            "start_time": "9:30",
            "end_time": "11:30"
        }
    }
}
```

##### [GET] Get By Profile Id
**URL:** `{{base_url}}/consultation/session`

### Example Response 200 OK
```json
{
    "message": "success",
    "status": 200,
    "data": [
        {
            "session_id": "68396af212f8c83cec320ada",
            "therapist_id": "t001",
            "client_id": "c001",
            "languages": [
                "Vietnamese"
            ],
            "specialization": [
                "anxiety"
            ],
            "consultation_modes": [
                "Online"
            ],
            "price_per_hour": 200,
            "session_time": {
                "day": "Monday",
                "start_time": "9:30",
                "end_time": "11:30"
            }
        },
        {
            "session_id": "683975e1ced584f025bebfcd",
            "therapist_id": "t001",
            "client_id": "c001",
            "languages": [
                "Vietnamese"
            ],
            "specialization": [
                "anxiety"
            ],
            "consultation_modes": [
                "Online"
            ],
            "price_per_hour": 200,
            "session_time": {
                "day": "Monday",
                "start_time": "11:30",
                "end_time": "13:30"
            }
        },
        {
            "session_id": "6839761fced584f025bebfce",
            "therapist_id": "t001",
            "client_id": "c001",
            "languages": [
                "Vietnamese"
            ],
            "specialization": [
                "anxiety"
            ],
            "consultation_modes": [
                "Online"
            ],
            "price_per_hour": 200,
            "session_time": {
                "day": "Monday",
                "start_time": "11:30",
                "end_time": "13:30"
            }
        },
        {
            "session_id": "68397620ced584f025bebfcf",
            "therapist_id": "t001",
            "client_id": "c001",
            "languages": [
                "Vietnamese"
            ],
            "specialization": [
                "anxiety"
            ],
            "consultation_modes": [
                "Online"
            ],
            "price_per_hour": 200,
            "session_time": {
                "day": "Monday",
                "start_time": "11:30",
                "end_time": "13:30"
            }
        }
    ]
}
```

### Deprecated

#### [GET] HelloWorld
**URL:** `{{base_url}}/`

### Example Response 200 OK
```json
{
    "code": 200,
    "message": "Success",
    "data": "Hello World"
}
```

### [GET] Test Route
**URL:** ``


## Recommend Service

### external

#### [POST] Recommend
**URL:** `{{base_url_recommend}}/recommend`

**Request Body:**
```json
{
  "clientRaw": {
    "profile_id": "C001",
    "address": "Hanoi",
    "languages": ["Vietnamese", "English"],
    "issue_detail": ["stress"],
    "consultation_modes": ["online"],
    "range_price": 500000,
    "availability": {
      "Days": ["Weekday"],
      "TimeSlots": ["Morning", "Afternoon"]
    },
    "preferred_therapist_gender": "",
    "experience_level": "mid",
    "specialization": ["stress", "anxiety"],
    "urgency_level": "medium",
    "session_duration": 60,
    "preferred_therapist_language": ["Vietnamese"],
    "is_flexible_with_schedule": true,
    "current_session": []
  },
  "therapistsRaw": [
    {
      "profile_id": "T001",
      "address": "Hanoi",
      "languages": ["Vietnamese"],
      "specialization": ["stress", "depression"],
      "consultation_modes": ["online"],
      "experience": 5,
      "rating": 4.7,
      "currency": "VND",
      "rage_price": 400000,
      "availability": {
        "Days": ["Weekday"],
        "TimeSlots": ["Morning"]
      },
      "is_available": true,
      "current_session": [],
      "matched_clients": []
    },
    {
      "profile_id": "T002",
      "address": "HCMC",
      "languages": ["English"],
      "specialization": ["marriage"],
      "consultation_modes": ["online"],
      "experience": 10,
      "rating": 4.5,
      "currency": "VND",
      "rage_price": 600000,
      "availability": {
        "Days": ["Weekday"],
        "TimeSlots": ["Afternoon"]
      },
      "is_available": true,
      "current_session": [],
      "matched_clients": []
    },
    {
      "profile_id": "T003",
      "address": "Hanoi",
      "languages": ["Vietnamese", "English"],
      "specialization": ["stress", "anxiety"],
      "consultation_modes": ["online"],
      "experience": 3,
      "rating": 4.9,
      "currency": "VND",
      "rage_price": 450000,
      "availability": {
        "Days": ["Weekday"],
        "TimeSlots": ["Morning", "Afternoon"]
      },
      "is_available": true,
      "current_session": [],
      "matched_clients": []
    }
  ]
}

```

