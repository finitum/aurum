# Aurum

Aurum is a service which provides authentication and user management for multiple applications, 
without the need for every single application to have their own user management.

With Aurum it's possible to manage a set of applications, called a group, 
and a set of users. Aurum can revoke access of a user to a group of applications, or grant 
them administrator privileges in a group. 

## Goals
The main goal of Aurum is to remove as much user management logic from applications. 

With Aurum, the only thing an application needs to do is to initialize a connection with Aurum (this may be
over a network, or embedded in the application). 
Then when performing an action for which a user needs to be authorized, 
it asks Aurum if the user is indeed authorized. 

# Usage
## Aurum server

To run the aurum api server, you can either 

* Deploy it using the provided [helm chart](./deployment/kubernetes)
* Run it in docker using the provided [docker-compose file](./deployment/docker)
* (not recommended) run it with `go run ./services/aurum` together with dgraph (`docker-compose -f services/aurum/docker-compose.yml up --force-recreate dgraph`)

The Aurum API is not secure by default, and it's recommended to run it behind a reverse proxy that terminates SSL.  

## In applications

To use Aurum in an application you make, you can interface with the Aurum HTTP API. For a few language, 
[wrapper clients](./clients) have been written, of which the reference implementation is the [go client](./clients/go).

By using the methods provided through this API, applications should essentially never have to deal with user management.

##### Groups
Aurum groups are subsets of all the users registered on the Aurum server. When an application wants to authorize a user, 
it checks the users's permission in a given group. In many cases you want to give each application their own group. However, it is possible 
to reuse groups for multiple applications. 

All users are part of the Aurum group. This gives them access to change their email address and password in Aurum. Users
that are admins in Aurum may create and manage groups. 

When for simple applications, you don't really care about permissions but just about users *having* an account, it's 
recommended to just check if they are part of the Aurum group. This verifies they have an account.



## Management
Managing an Aurum server consists of the following tasks:

* Creating new users
* Adding and Removing groups
* Adding and Removing users to and from groups
* Changing user permissions
* Changing your username and password

All these tasks can be performed through the HTTP API, as well as using the [Aurum TUI Client](./cmd/aurum). 
However, for most of these tasks, it's necessary to be an Administrator in the Aurum group.

