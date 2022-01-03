export DELINKCIOUS_MUTUAL_AUTH=false
export RUN_LINK_SERVICE=true
export RUN_SOCIAL_GRAPH_SERVICE=true

cd ../../svc/link_service || exit
go run ../../cmd/link_service_e2e/link_service_e2e.go

