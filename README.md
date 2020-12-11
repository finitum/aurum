# Aurum

Aurum is a service which provides authentication and user management for multiple applications, 
without the need for every single application to have their own user management.

With Aurum it's possible to manage a set of applications, 
and a set of users. Aurum can revoke access of a user to an application, or grant 
them administrator privileges in the application. 

## Goals
The main goal of Aurum is to remove as much user management logic from applications. 

With Aurum, the only thing an application needs to do is to initialize a connection with Aurum (this may be
over a network, or embedded in the application). 
Then when performing an action for which a user needs to be authorized, 
it asks Aurum if the user is indeed authorized. 

# Usage
## Aurum server



## In applications

To use aurum in applications, you can interface with the Aurum api. 
