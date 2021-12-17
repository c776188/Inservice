FROM golang:1.17

ENV GOPATH $GOPATH:/go
ENV PATH $PATH:$GOPATH/bin
WORKDIR /usr/local/go/src/Inservice
COPY . .

# beego
# RUN go get github.com/astaxie/beego
RUN go get github.com/beego/bee
RUN go install

CMD bee run