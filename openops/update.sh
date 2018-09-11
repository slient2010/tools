#!/bin/bash

cd front/; go-bindata -pkg=front -nocompress=true -debug=true html/... ; cd ..
cd front/; go-bindata -pkg=front -nocompress=true  html/... ; cd ..
go install openops

