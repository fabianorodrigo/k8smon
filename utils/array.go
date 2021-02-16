package utils

/*ArrayContains verifica se {item} já existe no array {arr}*/
func ArrayContains(arr []interface{}, item interface{}) bool {
	for _, a := range arr {
		if a == item {
			return true
		}
	}
	return false
}
