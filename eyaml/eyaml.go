package eyaml

import (
	"fmt"
	"os/exec"
	"github.com/rs/zerolog/log"
)

type OutputType int

const (
	Block OutputType = iota
	String
)

func (ot OutputType) String() string {
	return []string{"block", "string"}[ot]
}

func Encrypt(Password string, Label string, Output OutputType, PKCS7File string) (string, error) {
	log.Debug().Str("command",fmt.Sprintf("eyaml encrypt --label=%v --output=%v --string %v --pkcs7-public-key=%v",Label,Output.String(),"****",PKCS7File)).Msg("Encrypt string")
	cmd := exec.Command(
		"eyaml",
		"encrypt",
		"--label="+Label,
		"--output="+Output.String(),
		"--string", Password,
		"--pkcs7-public-key="+PKCS7File)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output[:]), nil
}

func EncryptFile(FileName string, Label string, Output OutputType, PKCS7File string) (string, error) {
	log.Debug().Str("command",fmt.Sprintf("eyaml encrypt --label=%v --output=%v --file %v --pkcs7-public-key=%v",Label,Output.String(),FileName,PKCS7File)).Msg("Encrypt string")
	cmd := exec.Command(
		"eyaml",
		"encrypt",
		"--label="+Label,
		"--output="+Output.String(),
		"--file", FileName,
		"--pkcs7-public-key="+PKCS7File)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output[:]), nil
}
