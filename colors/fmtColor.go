package colors

import "fmt"

const colorReset string = "\033[0m"
const colorRed string = "\033[31m"
const colorGreen string = "\033[32m"
const colorYellow string = "\033[33m"
const colorBlue string = "\033[34m"
const colorPurple string = "\033[35m"
const colorCyan string = "\033[36m"
const colorWhite string = "\033[37m"
const colorGray string = "\033[1;30m"

/**
Print em vermelho
**/
func Red(s string) {
	printCor(colorRed, s)
}

/*
Printf em vermelho
*/
func Redf(s string, params ...interface{}) {
	//fmt.Printf(fmt.Sprintf("%s%s%s", colorRed, s, colorReset), params...)
	printCor(colorRed, s, params...)
}

/**
Print em vermelho e quebra linha
**/
func Redln(s string) {
	//fmt.Println(colorRed, s, colorReset)
	printCor(colorRed, s+"\n")
}

/**
Print em verde
**/
func Green(s string) {
	printCor(colorGreen, s)
}

/*
Printf em verde
*/
func Greenf(s string, params ...interface{}) {
	//fmt.Printf(fmt.Sprintf("%s%s%s", colorRed, s, colorReset), params...)
	printCor(colorGreen, s, params...)
}

/**
Print em verde e quebra linha
**/
func Greenln(s string) {
	//fmt.Println(colorRed, s, colorReset)
	printCor(colorGreen, s+"\n")
}

/**
Print em amarelo
**/
func Yellow(s string) {
	printCor(colorYellow, s)
}

/*
Printf em amarelo
*/
func Yellowf(s string, params ...interface{}) {
	//fmt.Printf(fmt.Sprintf("%s%s%s", colorRed, s, colorReset), params...)
	printCor(colorYellow, s, params...)
}

/**
Print em amarelo e quebra linha
**/
func Yellowln(s string) {
	//fmt.Println(colorRed, s, colorReset)
	printCor(colorYellow, s+"\n")
}

/**
Print em azul
**/
func Blue(s string) {
	printCor(colorBlue, s)
}

/*
Printf em azul
*/
func Bluef(s string, params ...interface{}) {
	//fmt.Printf(fmt.Sprintf("%s%s%s", colorRed, s, colorReset), params...)
	printCor(colorBlue, s, params...)
}

/**
Print em azul e quebra linha
**/
func Blueln(s string) {
	//fmt.Println(colorRed, s, colorReset)
	printCor(colorBlue, s+"\n")
}

/**
Print
**/
func Purple(s string) {
	printCor(colorPurple, s)
}

/*
Printf
*/
func Purplef(s string, params ...interface{}) {
	//fmt.Printf(fmt.Sprintf("%s%s%s", colorRed, s, colorReset), params...)
	printCor(colorPurple, s, params...)
}

/**
Print
**/
func Cyan(s string) {
	printCor(colorCyan, s)
}

/*
Printf
*/
func Cyanf(s string, params ...interface{}) {
	//fmt.Printf(fmt.Sprintf("%s%s%s", colorRed, s, colorReset), params...)
	printCor(colorCyan, s, params...)
}

/**
Print
**/
func White(s string) {
	printCor(colorWhite, s)
}

/*
Printf
*/
func Whitef(s string, params ...interface{}) {
	//fmt.Printf(fmt.Sprintf("%s%s%s", colorRed, s, colorReset), params...)
	printCor(colorWhite, s, params...)
}

/**
Print e quebra linha
**/
func Whiteln(s string) {
	//fmt.Println(colorRed, s, colorReset)
	printCor(colorWhite, s+"\n")
}

/**
Print e quebra linha
**/
func Cyanln(s string) {
	//fmt.Println(colorRed, s, colorReset)
	printCor(colorCyan, s+"\n")
}

/**
Print e quebra linha
**/
func Purpleln(s string) {
	//fmt.Println(colorRed, s, colorReset)
	printCor(colorPurple, s+"\n")
}

/**
Print em cinza
**/
func Gray(s string) {
	printCor(colorGray, s)
}

/*
Printf em cinza
*/
func Grayf(s string, params ...interface{}) {
	//fmt.Printf(fmt.Sprintf("%s%s%s", colorRed, s, colorReset), params...)
	printCor(colorGray, s, params...)
}

/**
Print em cinza e quebra linha
**/
func Grayln(s string) {
	//fmt.Println(colorRed, s, colorReset)
	printCor(colorGray, s+"\n")
}

/**
Printf texto em uma determinada cor e reseta
**/
func printCor(cor, s string, params ...interface{}) {
	fmt.Printf(fmt.Sprintf("%s%s%s", cor, s, colorReset), params...)
}
