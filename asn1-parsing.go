package main

import (
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run asn1parser.go <path to ASN.1 encoded file>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		log.Fatalf("Failed to decode PEM block")
	}

	_, err = asn1.Unmarshal(block.Bytes, &asn1.RawValue{})
	if err != nil {
		log.Fatalf("Failed to unmarshal ASN.1 data: %v", err)
	}

	parseASN1(block.Bytes, 0, 0)
}

func parseASN1(data []byte, depth int, offset int) {
	for len(data) > 0 {
		var value asn1.RawValue
		rest, err := asn1.Unmarshal(data, &value)
		if err != nil {
			log.Fatalf("Failed to unmarshal ASN.1 data: %v", err)
		}

		hl := len(data) - len(rest) // Header length
		l := len(value.Bytes)       // Content length

		printASN1(value, depth, hl, l, offset)

		offset += hl
		if value.IsCompound {
			parseASN1(value.Bytes, depth+1, offset)
		}
		offset += l
		data = rest
	}
}

func printASN1(value asn1.RawValue, depth, hl, length, offset int) {
	typ := "prim"
	if value.IsCompound {
		typ = "cons"
	}

	fmt.Printf("%4d:d=%d  hl=%d l=%4d %s: %s\n", offset, depth, hl, length, typ, getTagDescription(value.Tag))
}

func getTagDescription(tag int) string {
	switch tag {
	case 0x00:
		return "EOC (End of Content)"
	case 0x01:
		return "BOOLEAN"
	case 0x02:
		return "INTEGER"
	case 0x03:
		return "BIT STRING"
	case 0x04:
		return "OCTET STRING"
	case 0x05:
		return "NULL"
	case 0x06:
		return "OBJECT IDENTIFIER"
	case 0x07:
		return "OBJECT DESCRIPTOR"
	case 0x08:
		return "EXTERNAL"
	case 0x09:
		return "REAL (Floating Point)"
	case 0x0A:
		return "ENUMERATED"
	case 0x0B:
		return "EMBEDDED PDV"
	case 0x0C:
		return "UTF8String"
	case 0x10:
		return "SEQUENCE and SEQUENCE OF"
	case 0x11:
		return "SET and SET OF"
	case 0x12:
		return "NumericString"
	case 0x13:
		return "PrintableString"
	case 0x14:
		return "TeletexString / T61String"
	case 0x15:
		return "VideotexString"
	case 0x16:
		return "IA5String (ASCII)"
	case 0x17:
		return "UTCTime"
	case 0x18:
		return "GeneralizedTime"
	case 0x19:
		return "GraphicString"
	case 0x1A:
		return "VisibleString (ISO646String)"
	case 0x1B:
		return "GeneralString"
	case 0x1C:
		return "UniversalString"
	case 0x1D:
		return "CHARACTER STRING"
	case 0x1E:
		return "BMPString"
	default:
		return fmt.Sprintf("Unknown or Application-specific Tag [0x%X]", tag)
	}
}
