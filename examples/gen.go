package example

//go:generate cff -file magic.go -genmode source-map .
//go:generate cff -file magic_v2.go -tags v2 -genmode modifier .
