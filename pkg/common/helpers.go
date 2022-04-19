package common

import (
	hashstructure "github.com/mitchellh/hashstructure/v2"
	"strconv"
)

func CompareHashStructure(oldHash interface{}, newSpec interface{}) bool {

	result := GetHashStructure(newSpec)

	return oldHash == result
}

func GetHashStructure(i interface{}) string {
	hash, _ := hashstructure.Hash(i, hashstructure.FormatV2, nil)
	return strconv.FormatUint(hash, 10)
}
