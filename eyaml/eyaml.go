package eyaml

import (
	"os/exec"
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
