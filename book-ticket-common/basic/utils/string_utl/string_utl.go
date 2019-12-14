package string_utl

import (
	"encoding/binary"
	"math"
	"strconv"
)

func StringToInt(sValue string) (intValue int, err error) {
	intValue, err = strconv.Atoi(sValue)
	return
}

func StringToInt64(sValue string) (int64Value int64, err error) {
	int64Value, err = strconv.ParseInt(sValue, 10, 64)
	return
}

func IntToString(intValue int) (sValue string) {
	sValue = strconv.Itoa(intValue)
	return
}
func Int64ToString(intValue int64) (sValue string) {
	sValue = strconv.FormatInt(intValue, 10)
	return
}

func Float32ToFloat32String(floatValue float32) (sValue string) {
	sValue = strconv.FormatFloat(float64(floatValue), 'f', 2, 64)
	return
}
func Float64ToFloat64String(floatValue float64) (sValue string) {
	sValue = strconv.FormatFloat(floatValue, 'E', -1, 64)
	return
}

func Float64StringToFloat64(sValue string) (floatValue float64, err error) {
	floatValue, err = strconv.ParseFloat(sValue, 64)
	return
}

func Float64StringToFloat32(sValue string) (floatValue float32, err error) {
	float_64, err := Float64StringToFloat64(sValue)
	if err != nil {
		return
	}
	bytes := Float64ToByte(float_64)
	floatValue = ByteToFloat32(bytes)
	return
}

func Float32ToByte(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)

	return bytes
}

func ByteToFloat32(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)

	return math.Float32frombits(bits)
}

func Float64ToByte(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)

	return bytes
}

func ByteToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)

	return math.Float64frombits(bits)
}
