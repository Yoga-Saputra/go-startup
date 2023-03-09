<pre style="font-size: 1.4vw;">
<p align="center">
 ________       _________    ________      ________      _________    ___  ___      ________   
|\   ____\     |\___   ___\ |\   __  \    |\   __  \    |\___   ___\ |\  \|\  \    |\   __  \  
\ \  \___|_    \|___ \  \_| \ \  \|\  \   \ \  \|\  \   \|___ \  \_| \ \  \\\  \   \ \  \|\  \ 
 \ \_____  \        \ \  \   \ \   __  \   \ \   _  _\       \ \  \   \ \  \\\  \   \ \   ____\
  \|____|\  \        \ \  \   \ \  \ \  \   \ \  \\  \|       \ \  \   \ \  \\\  \   \ \  \___|
    ____\_\  \        \ \__\   \ \__\ \__\   \ \__\\ _\        \ \__\   \ \_______\   \ \__\   
   |\_________\        \|__|    \|__|\|__|    \|__|\|__|        \|__|    \|_______|    \|__|   
   \|_________|                                                                             
</p>
</pre>
<p align="center">
<a href="https://golang.org/">
    <img src="https://img.shields.io/badge/Made%20with-Go-1f425f.svg">
</a>
<a href="/LICENSE">
    <img src="https://img.shields.io/badge/License-MIT-green.svg">
</a>
</p>
<p align="center">
RESTful API of <b>GO - Startup</b>
</p>


# Startup API Guide

## 🔀 Compatible Route Endpoint
| NO | Use                                 | Endpoint               | Example                                             | Action
|----|-------------------------------------|------------------------|-----------------------------------------------------|------------
| 1  | Register                            | api/v1/users           | http://localhost:4004/api/v1/users            | POST
| 2  | Login                               | api/v1/session         | http://localhost:4004/api/v1/session          | POST
| 3  | Email Checker                       | api/v1/email_checkers  | http://localhost:4004/api/v1/email_checkers   | POST
| 4  | Upload Avatar                       | api/v1/avatars         | http://localhost:4004/api/v1/avatars          | POST
| 5  | Get Campaigns                       | api/v1/campaigns       | http://localhost:4004/api/v1/campaigns        | GET
| 6  | Create Campaign                     | api/v1/campaigns       | http://localhost:4004/api/v1/campaigns        | POST
| 7  | Update Campaign                     | api/v1/campaigns/{id}  | http://localhost:4004/api/v1/campaigns/{id}   | PUT
| 9  | Get Detail Campaign                 | api/v1/campaigns/{id}  | http://localhost:4004/api/v1/campaigns/{id}   | GET
| 9  | Upload Campaign images              | api/v1/campaigns-images| http://localhost:4004/api/v1/campaigns-images | POST
| 10 | Transaction(using midtrans payment) | On Progress            | On Progress                                   | On Progress 

---

## 📖 Compatible JSON Payload Startup API
This is the JSON payload that's sended to Startup API

### 💲 Register JSON Payload
```js
{
    "name": "Developer",
    "email": "dev@gmail.com",
    "occupation": "Golang Developer",
    "password": "password"
}
```

### 💸 Login JSON Payload
```js
{
    "email": "agolang4@gmail.com",
    "password": "password"
}
```

### 💸 Email Checker JSON Payload
```js
{
    "email": "golang4@gmail.com"
}
```

### 💸 Upload Avatar Form Data Payload
```js

form-data

| Key                 | Value               |
|---------------------|---------------------|
| avatar              | images.png          |
```

### 💸 Get Campaigns Query Param Payload (optional)
```js

Query-Params

| Key                 | Value               |
|---------------------|---------------------|
| user_id             | 1                   |
```

### 💸 Create Campaign JSON Payload
```js
{
    "name": "go campaign",
    "short_description": "Golang camp Developer",
    "description": "For Beginner",
    "goal_amount": 1000,
    "perks": "go ,dev,  programmer"
}
```
### 💸 Update Campaign JSON Payload
```js
{
    "name": "campaign update",
    "short_description": "Golang developer Forum",
    "description": "New Knowlage",
    "goal_amount": 1000,
    "perks": "ok1 ,ok 2,  ok 3"
}
```

### 💸 Get Detail Campaign 
```js
v1/campaigns/2
```

### 💸 Upload Campaign images Form Data Payload
```js

form-data

| Key                 | Value               |
|---------------------|---------------------|
| campaign_id         | 1                   |
| is_primary          | true                |
| file                | campaign.png        |