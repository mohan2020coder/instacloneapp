# This app is using nextjs and golang(switch between db(sql or no sql))

  - [Overview](#overview)
  - [Development](#development)
  - [Server](#server)

# NextJS and Golang Web Application

## Overview
This application is built using NextJS, a React framework for building server-rendered React applications, and Golang, a statically typed, compiled programming language.

## Development

To run the application in development mode, use the following command:


npm run dev


This will start the NextJS development server and the Golang server.

## Server

The Golang server is started using the following command:


go run main.go


This will start the Golang server and handle the backend logic of the application.



# if you want to use the sql - models should be replaced with, the models in sqlite.go file intead of mongo 