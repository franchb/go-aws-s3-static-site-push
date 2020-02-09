package log

//go:generate go-enum -f=$GOFILE --lower --flag
// Level is an enumeration of commented values
/*
ENUM(
Info
Warning
Debug
)
*/
type Level int
