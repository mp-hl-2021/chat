Example Chat Project
====================

How to Test
-----------

Create new account

    curl -v -X POST localhost:8080/signup -d '{"login": "<your login here>", "password": "<password>"}'

Login and retrieve token

    curl -v -X POST localhost:8080/signin -d '{"login": "<your login here>", "password": "<password>"}'

Get account id from new account response headers or JWT

Access some of protected resources

    TOKEN="<your token>"
    curl -v localhost:8080/accounts/0 -H "Authorization: Bearer $TOKEN"