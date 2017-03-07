package main

import (
	"fmt"
	"crypto/des"
	"crypto/cipher"
	"bytes"
	"os"
	"os/exec"
	"io/ioutil"
)

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func DesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func DesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}
/*
func testDes() {
	key := []byte{0xad,0xcd,0x11,0xef,0x12,0x23,0x33,0xdd}
	result, err := DesEncrypt([]byte("sangforadqianmingmima"), key)
	if err != nil {
		panic(err)
	}
	fmt.Println("result:",result)
	//fmt.Println(base64.StdEncoding.EncodeToString(result))
	origData, err := DesDecrypt(result, key)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(origData))
}
*/
func main() {
	if len(os.Args) != 2{
		fmt.Println("usage: sign <ssu>")
		return
	}
	encPassword := []byte{247, 143, 188, 65, 177, 38,
		              152, 16, 201, 20, 169, 220,
		               90, 163, 153, 135, 60, 13,
		               48, 80, 46, 82, 86, 85}
	key := []byte{0xad,0xcd,0x11,0xef,0x12,0x23,0x33,0xdd}
	password, err := DesDecrypt(encPassword, key)
	if err != nil {
		panic(err)
	}
	fmt.Println("passwd:",string(password))
	privateKey,err := exec.LookPath("private.key")
	if err != nil {
		fmt.Println("找不到private.key")
	}
	signverify, err1 := exec.LookPath("signverify")
	if err1 != nil{
		fmt.Println("找不到signverify")
	}

	Args := []string{
		0: os.Args[1],
		1: privateKey,
		2: signverify,
		3: string(password),
	}

	cmd := exec.Command("ssusign.sh",Args...)

	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		fmt.Println("包签名失败:",err)
		return
	}
	data, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println("签名日志信息丢失，无关紧要:",err)
	}
	fmt.Println("签名过程日志:\n",string(data))
	if err := cmd.Wait(); err == nil {
		fmt.Println("签名成功")
		return
	}else{
		fmt.Println("签名失败:",err)
		return
	}
}
