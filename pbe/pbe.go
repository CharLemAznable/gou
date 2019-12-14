package pbe

import (
	"fmt"
	"github.com/bingoohuang/gou/file"
	"io/ioutil"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/table"
)

const iterations = 19
const pbePrefix = `{PBE}`

// Pbe encrypts p by PBEWithMD5AndDES with 19 iterations.
// it will prompt password if viper get none.
func Pbe(p string) (string, error) {
	pwd := GetPbePwd()
	if pwd == "" {
		return "", fmt.Errorf("pbepwd is requird")
	}

	encrypt, err := Encrypt(p, pwd, iterations)
	if err != nil {
		return "", err
	}

	return pbePrefix + encrypt, nil
}

// Ebp decrypts p by PBEWithMD5AndDES with 19 iterations.
func Ebp(p string) (string, error) {
	if !strings.HasPrefix(p, pbePrefix) {
		return p, nil
	}

	pwd := GetPbePwd()
	if pwd == "" {
		return "", fmt.Errorf("pbepwd is requird")
	}

	return Decrypt(p[len(pbePrefix):], pwd, iterations)
}

// PrintEncrypt prints the PBE encryption.
func PrintEncrypt(passStr string, plains ...string) {
	if len(plains) == 1 && strings.HasPrefix(plains[0], "@") && file.Stat(plains[0][1:]) == file.Exists {
		processPbeFile(plains[0][1:], passStr)

		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Plain", "Encrypted"})

	for i, p := range plains {
		pbed, err := Encrypt(p, passStr, iterations)
		if err != nil {
			fmt.Fprintf(os.Stderr, "pbe.Encrypt error %v", err)
			os.Exit(1)
		}

		t.AppendRow(table.Row{i + 1, p, pbePrefix + pbed})
	}

	t.Render()
}

// PrintDecrypt prints the PBE decryption.
func PrintDecrypt(passStr string, cipherText ...string) {
	if len(cipherText) == 1 && strings.HasPrefix(cipherText[0], "@") && file.Stat(cipherText[0][1:]) == file.Exists {
		processEbpFile(cipherText[0][1:], passStr)

		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Encrypted", "Plain"})

	for i, ebp := range cipherText {
		ebpx := strings.TrimPrefix(ebp, pbePrefix)

		p, err := Decrypt(ebpx, passStr, iterations)
		if err != nil {
			fmt.Fprintf(os.Stderr, "pbe.Decrypt error %v", err)
			os.Exit(1)
		}

		t.AppendRow(table.Row{i + 1, ebp, p})
	}

	t.Render()
}

func processPbeFile(filename, passStr string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	text, err := Config{Passphrase: passStr}.PbeText(string(file))
	if err != nil {
		panic(err)
	}

	ft, _ := os.Stat(filename)

	if err := ioutil.WriteFile(filename, []byte(text), ft.Mode()); err != nil {
		panic(err)
	}
}

func processPbeChgFile(filename, passStr, pbenew string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	text, err := Config{Passphrase: passStr}.ChangePbe(string(file), pbenew)
	if err != nil {
		panic(err)
	}

	ft, _ := os.Stat(filename)

	if err := ioutil.WriteFile(filename, []byte(text), ft.Mode()); err != nil {
		panic(err)
	}
}

func processEbpFile(filename, passStr string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	text, err := Config{Passphrase: passStr}.EbpText(string(file))
	if err != nil {
		panic(err)
	}

	ft, _ := os.Stat(filename)

	if err := ioutil.WriteFile(filename, []byte(text), ft.Mode()); err != nil {
		panic(err)
	}
}
