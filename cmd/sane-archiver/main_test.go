package main

import (
	"os/exec"
	"path/filepath"
	"regexp"
	"testing"
)

const (
	target  = `../../tests/data/test.sane1`
	private = `MIICXQIBAAKBgQDICj7iLhIPQzSDJoENJ0yvsj7mAk9UAV7LKeuhXllZdsO24K104RAjYnQZjuvgsLknOdP23Fx3/kch/pe33QU097eB4lXtkd+fyX270+B758OzMXS+KTob8qvvqzx1GnmSYLb/kuaIVDorQn3zcpKahAV0gS7WLufRkZqWQrxX4wIDAQABAoGAM4E/06idCcT5/lKpo6NcwVgZjctGdZCswY6XlsLeKoTDu5B52MAiEZpF3lbIMOAPrCPdiZAPVu3njr8ofTSxJA+haWFGR/Q+T47J1Ouer+cB+jsvDVQKMhpd5G7ydVUEelz6+sCDN3Zu1UC/V50ey7D368XIKswAFjUV/xXKM5ECQQDfhmPAuJRMDNUSZPfuw7kkoMi8hoa93K37RbNqS+N+Z1xljmQrwJ2GLJ1+1nWknu8I5fjlYmEFnY9u/ptUmshvAkEA5RpgNPpzipsdxsNDrPwPDEtIcLlKlE3kWr2roN6KhHeviqDMQxoF938R6ZLK1lH2TymbI5UapZ7R4V6nw7kZzQJAIooyueIL0GCfQDNn+HY4EsfhnPgws//4xn4zxjYp1iuEpJDHO9eMv+H/CE19ak3A5CAdQNzd3y9ErcMcH4u3cwJBAN/NrU/znW04bJUfaPwSW0ziOgjMKTvI/5tZD9EdtGkFVjlxLTkbsdp9im0HFhjZhmj8tu3CmX5TMKodQnujVb0CQQDI0MVnI/hjYDE3LAtYlPBdnTREaLX9XVrS+e+8pdQj62m6tC9AqLPPPahgKVKeVnJpZZHeiPXr5E2oa1TpQhut`
	public  = `MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDICj7iLhIPQzSDJoENJ0yvsj7mAk9UAV7LKeuhXllZdsO24K104RAjYnQZjuvgsLknOdP23Fx3/kch/pe33QU097eB4lXtkd+fyX270+B758OzMXS+KTob8qvvqzx1GnmSYLb/kuaIVDorQn3zcpKahAV0gS7WLufRkZqWQrxX4wIDAQAB`
)

func TestPack(t *testing.T) {
	cmd := exec.Command(`go`, `run`, `.`, `pack`, `../todo.md`,
		`--key`, public, `--output`, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Error(err)
		t.Fatalf(`%s`, output)
	}
	cmd = exec.Command(`go`, `run`, `.`, `unpack`, target,
		`--key`, private, `--output`, filepath.Dir(target))
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Error(err)
		t.Fatalf(`%s`, output)
	}
}

func TestKeygen(t *testing.T) {
	cmd := exec.Command(`go`, `run`, `.`, `keygen`)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Error(err)
	}
	if !regexp.MustCompile(`^\=+\nPrivate key\n\=+\n[^\n]+\n\=+\nPublic key\n\=+\n[^\n]+\n\=+\n$`).Match(output) {
		t.Error(`Keygen result unexpected!`)
	}
}
