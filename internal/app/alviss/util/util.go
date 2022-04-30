package util

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/itchyny/base58-go"
	uuid "github.com/nu7hatch/gouuid"
)

func GetExpireTime(expDate string) time.Duration {
	if expDate[len(expDate)-1] == 'h' {
		expDate = expDate[:len(expDate)-1]
		expDateInt, _ := strconv.Atoi(expDate)

		return time.Hour * time.Duration(expDateInt)
	} else if expDate[len(expDate)-1] == 'd' {
		expDate = expDate[:len(expDate)-1]
		expDateInt, _ := strconv.Atoi(expDate)
		return time.Hour * time.Duration(expDateInt*24)
	} else if expDate[len(expDate)-1] == 's' {
		expDate = expDate[:len(expDate)-1]
		expDateInt, _ := strconv.Atoi(expDate)
		return time.Second * time.Duration(expDateInt)
	} else if expDate[len(expDate)-1] == 'm' {
		expDate = expDate[:len(expDate)-1]
		expDateInt, _ := strconv.Atoi(expDate)
		return time.Minute * time.Duration(expDateInt)
	} else {
		return 0
	}
}

func sha256Of(input string) []byte {
	algorithm := sha256.New()
	algorithm.Write([]byte(input))

	return algorithm.Sum(nil)
}

func base58Encoded(bytes []byte) string {
	encoding := base58.BitcoinEncoding

	encoded, err := encoding.Encode(bytes)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return string(encoded)
}

func GenerateShortLink(initialLink string) string {
	u, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	urlHashBytes := sha256Of(initialLink + u.String())
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	finalString := base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))

	return finalString[:8]
}
