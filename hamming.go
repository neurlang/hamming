package hamming

import (
	"fmt"
)

// Encode encodes low 4 bits as 1 byte using the error correcting hamming code
func Encode(data byte) byte {
	ret := encodeArray([]byte{
		(data >> 0) & 1,
		(data >> 1) & 1,
		(data >> 2) & 1,
		(data >> 3) & 1,
	})
	return ret[0] |
		ret[1]<<1 |
		ret[2]<<2 |
		ret[3]<<3 |
		ret[4]<<4 |
		ret[5]<<5 |
		ret[6]<<6 |
		ret[7]<<7
}

// Decode decodes 1 byte as low 4 bits using the error correcting hamming code, reports an error if any
func Decode(data byte) (byte, error) {
	ret, err := decodeArray([]byte{
		(data >> 0) & 1,
		(data >> 1) & 1,
		(data >> 2) & 1,
		(data >> 3) & 1,
		(data >> 4) & 1,
		(data >> 5) & 1,
		(data >> 6) & 1,
		(data >> 7) & 1,
	})
	if err != nil {
		return 0, err
	}
	return ret[0] |
		ret[1]<<1 |
		ret[2]<<2 |
		ret[3]<<3, nil
}

func encodeArray(data []byte) []byte {

	const numParityBits = 3
	const encodedLength = 8

	encoded := make([]byte, encodedLength)

	for parityBitIndex := 1; parityBitIndex < (1 << numParityBits); parityBitIndex <<= 1 {
		encoded[parityBitIndex] = calculateParity(data, parityBitIndex)
	}

	dataIndex := 0
	for encodedIndex := 3; encodedIndex < len(encoded); encodedIndex++ {
		if !isPowerOfTwo(encodedIndex) {
			encoded[encodedIndex] = data[dataIndex]
			dataIndex++
		}
	}

	encoded[0] = calculateParity(encoded[1:], 0)

	return encoded
}

func decodeArray(encoded []byte) ([]byte, error) {
	const encodedLength = 8
	const numParityBits = 3
	indexOfError := 0

	decoded := extractData(encoded)

	overallExpected := calculateParity(encoded[1:], 0)
	overallActual := encoded[0]
	overallCorrect := overallExpected == overallActual

	for parityBitIndex := 1; parityBitIndex < (1 << numParityBits); parityBitIndex <<= 1 {
		expected := calculateParity(decoded, parityBitIndex)
		actual := encoded[parityBitIndex]
		if expected != actual {
			indexOfError += parityBitIndex
		}
	}

	if indexOfError != 0 && overallCorrect {
		return nil, fmt.Errorf("Two errors detected")
	} else if indexOfError != 0 && !overallCorrect {
		encoded[indexOfError] = 1 - encoded[indexOfError]
		decoded = extractData(encoded)
	}

	return decoded, nil
}

func calculateParity(data []byte, parity int) byte {
	retval := byte(0)

	if parity == 0 {
		for _, bit := range data {
			retval ^= bit
		}
	} else {
		for _, dataIndex := range dataBitsCovered(parity, len(data)) {
			retval ^= data[dataIndex]
		}
	}

	return retval
}

func dataBitsCovered(parity, lim int) []byte {
	var indices []byte

	if !isPowerOfTwo(parity) {
		panic("All Hamming parity bits are indexed by powers of two.")
	}

	dataIndex := 1
	totalIndex := 3

	for dataIndex <= lim {
		currBitIsData := !isPowerOfTwo(totalIndex)
		if currBitIsData && (totalIndex%(parity<<1)) >= parity {
			indices = append(indices, byte(dataIndex-1))
		}
		dataIndex += boolToInt(currBitIsData)
		totalIndex++
	}

	return indices
}

func extractData(encoded []byte) []byte {
	var data []byte

	for i := 3; i < len(encoded); i++ {
		if !isPowerOfTwo(i) {
			data = append(data, encoded[i])
		}
	}

	return data
}

func isPowerOfTwo(n int) bool {
	return n > 0 && (n&(n-1)) == 0
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
