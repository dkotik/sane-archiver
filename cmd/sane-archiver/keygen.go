package main

import (
	"archiver"
	"fmt"

	"github.com/alecthomas/kong"
)

const keygenPrintTemplate = `==============================
Private key
==============================
%s
==============================
Public key
==============================
%s
==============================
`

type keygenTask struct{}

func (c *keygenTask) Run(ctx *kong.Context) error {
	private, public, err := archiver.GenerateKeyPair()
	fmt.Printf(keygenPrintTemplate, private, public)
	return err
}
