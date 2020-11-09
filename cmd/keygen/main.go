package main

import (
	"flag"
)

func main() {
	out := flag.String("out", "stdout", "where to output generated keys. Options: [stdout, file, both]")
	pkPath := flag.String("pk", "./id_25519.pub", "where to write the public if using file gen")
	skPath := flag.String("sk", "./id_25519", "where to write the secret if using file gen")

	flag.Parse()

	KeyGenerationUtil(*out, *pkPath, *skPath)
}
