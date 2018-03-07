FROM golang:1.10
RUN git config --global user.email "git@localhost"
RUN git config --global user.name "git"

COPY . /go/src/github.com/concourse/pivotal-tracker-resource

RUN go get github.com/concourse/pivotal-tracker-resource/in
RUN go build -o /opt/resource/in github.com/concourse/pivotal-tracker-resource/in

RUN go get github.com/concourse/pivotal-tracker-resource/out
RUN go build -o /opt/resource/out github.com/concourse/pivotal-tracker-resource/out

RUN go get github.com/concourse/pivotal-tracker-resource/check
RUN go build -o /opt/resource/check github.com/concourse/pivotal-tracker-resource/check

RUN chmod +x /opt/resource/*