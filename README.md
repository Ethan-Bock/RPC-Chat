## RPC-Chat (Go)

This is a rudimentary project that I was in class to make a local RPC-Chat. The code you are seeing is the following result, and had the purpose of expanding my knowledge on both Go and RPC-Chats.


## Getting Started
### 1. Clone the Repository
```
git clone https://github.com/Ethan-Bock/RPC-Chat.git
cd RPC-Chat
```

### 2. Running script
In one terminal, run the following bash command to create the server:
```
./synod 3410
```
On a second terminal, run the following bash command to connect the first user to the server:
```
./synod :3410 user1
```
On a third terminal, run the following bash command to connect the second user to the server:
```
./synod :3410 user2
```
You have now set up the two users and the server on the terminals and can communicate between them.

### 3. Running each part of the file
For an explanation of the functions available, simply type ```help``` in either of the user terminals to understand what each function does.
To close a user, type ```quit``` in the wanted user's terminal, or to close a user and server, type ```shutdown``` in the user's terminal (note this will kick out any remaining users out of the server).

![linkedin-RPCchat-pic](https://github.com/user-attachments/assets/74652b51-bca5-41e5-a5eb-55d3acdf37a2)
