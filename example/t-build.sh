# basic build sript

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd $DIR/browser
go get  
gopherjs build ex.go
 
cp ex.html $DIR/webserver
cp ex.js $DIR/webserver

cd $DIR/webserver
go run ex.go