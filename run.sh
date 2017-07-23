#!/bin/bash
#set -x

go get gopkg.in/mgo.v2
go get github.com/julienschmidt/httprouter
go get github.com/nu7hatch/gouuid


go build main.go class.go student.go

killall main
./main &
