# Flextesa Easy Bake
#
#  Flextesa bake on operation image
#
#  docker run -it --rm minty/flextesa-easybake \
#    --no-baking \
#    --size 1 \
#    --set-history-mode "N000:archive" \
#    --remove-default-bootstrap-accounts \
#    --protocol-kind "Carthage" \
#    --protocol-hash "PsCARTHAGazKbHtnKfLzQg3kms52kSRpgnDY982a9oYsSXRLQEb" \
#
#
#

FROM golang:1.14-buster as proxy_build
COPY easybake.go /go/src/easybake.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build /go/src/easybake.go


FROM registry.gitlab.com/tezos/flextesa:e6612b9b-run as flextesa_build
COPY --from=proxy_build /go/easybake /usr/local/bin/easybake
ENV flextesa_node_cors_origin="*"
EXPOSE 20000

ENTRYPOINT ["/usr/local/bin/easybake" \
	,"--no-baking" \
	,"--size", "1" \
	,"--set-history-mode", "N000:archive" \
	,"--protocol-kind", "Carthage" \
	,"--protocol-hash", "PsCARTHAGazKbHtnKfLzQg3kms52kSRpgnDY982a9oYsSXRLQEb" ]
