# sertis_app
## 1. Run commard to deploy mariaDB, phpmyadmin and sertis_app backend
    `docker-compose build && docker-compose up`
    
## 2. Sertis_app will Run on localhost:8880 and Content-Type = application/json
### 2.1 Sertis API will contain

#### 2.1.1 localhost:8880/signup ***POST Method***
##### Example
    
    `Request body is json`: 
    {
        "username": "jojo",
        "password": "129222"
    }
    
    `Response body is json`:
    {
        "success": false,
        "error_message": "Already Have Account"
    }

#### 2.1.2 localhost:8880/signin ***POST Method***
##### Example
   
   `Request body is json`: 
    {
        "username": "jojo",
        "password": "129222"
    }
   
   `Response body is json`:
    {
        "success": true,
        "error_message": "",
        "access_token": "**AccessToken**"
    }

#### 2.1.3 localhost:8880/addnewcard ***POST Method***
##### Headers
###### Authorization : Bearer + Use AccessToken From Signin
 
##### Example

    `Request body is json`: 
    {
        "name": "NVIDIA เปิดตัว RTX 1O",
        "status": "n/a",
        "content": "จีพียูเข้าถึงสตอเรจโดยตรง ใช้ DirectStorage API ของไมโครซอฟท์",
        "category":"Technology"
    }
    
    `Response body is json`:
    {
        "success": true,
        "error_message": ""
    }

#### 2.1.4 localhost:8880/blog ***GET Method***
##### Headers
###### Authorization : Bearer + Use AccessToken From Signin
 
##### Example
    `No Request body`: 
    `Response body is json`:
    {
        "success": true,
        "error_message": "",
        "cards": [
            {
                "id": 2,
                "name": "NVIDIA เปิดตัว RTX 1O",
                "status": "n/a",
                "content": "จีพียูเข้าถึงสตอเรจโดยตรง ใช้ DirectStorage API ของไมโครซอฟท์",
                "category": "Technology",
                "author": "jojo"
            }
        ]
    }

#### 2.1.5 localhost:8880/updatecard ***POST Method***
##### Headers
###### Authorization : Bearer + Use AccessToken From Signin
 
##### Example

    `Request body is json`: 
    {
        "username": "jojo",
        "password": "129222"
    }
    
    `Response body is json`:
    {
        "id" : 1,
        "name": "NVIDIA เปิดตัว RTX 1O",
        "status": "n/a",
        "content": "จีพียูเข้าถึงสตอเรจโดยตรง ใช้ DirectStorage API ของไมโครซอฟท์",
        "category":"Nvidia"
    }

#### 2.1.6 localhost:8880/deletecard/:cardID ***GET Method***
##### Headers
###### Authorization : Bearer + Use AccessToken From Signin
 
##### Example

    `No Request body`: 
    `Response body is json`:
    {
        "success": true,
        "error_message": ""
    }


## 3. I have already add Postman in sertis.postman_collection.json
### 3.1 You can import to Postman
