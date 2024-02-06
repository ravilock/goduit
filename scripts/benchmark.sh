PACKAGE=$1
FUNC=$2
go test $PACKAGE -bench=$FUNC -benchmem -cpuprofile prof.cpu -memprofile prof.mem
