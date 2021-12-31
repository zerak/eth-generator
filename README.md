
#### BUILD linux
GOOS=linux GOARCH=amd64 go build -o generator

#### RUN
./generator --case=true --prefix=2021A --suffix=2021 -num=10
