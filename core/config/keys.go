package config

import (
	"github.com/finitum/aurum/internal/jwt/ecc"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func loadKey(key string, keyPath string, keyType string) (ecc.Key, error) {
	if key != "" {
		k, err := ecc.FromPem([]byte(key))
		if err != nil {
			return nil, errors.Wrapf(err, "could not parse key given in environment variable")
		}
		return k, nil
	} else {
		k, err := ecc.FromFile(keyPath)
		if err != nil {
			log.Warnf("Couldn't find keys in environment, and reading the %v failed. Generating %v. "+
				"NOTE: this is normal on the first run of aurum.", keyType, keyType)

			return nil, nil
		}
		return k, nil
	}
}

func writeKey(key ecc.Key, keyPath string) error {
	return errors.Wrap(key.WriteToFile(keyPath), "Couldn't write key to file")
}

// TODO: Testing
func findKeys(config *EnvConfig) (ecc.PublicKey, ecc.SecretKey, error) {

	// Get keys from file or env (else nil)
	pk, err := loadKey(config.PublicKey, config.PublicKeyPath, "public key")
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to find public key")
	}

	sk, err := loadKey(config.SecretKey, config.SecretKeyPath, "secret key")
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to find secret key")
	}

	// If both keys were already found, check if they match each other. Then return.
	if pk != nil && sk != nil {
		publicKey, ok := pk.(ecc.PublicKey)
		if !ok {
			return nil, nil, errors.New("Couldn't interpret pem as public key")
		}
		secretKey, ok := sk.(ecc.SecretKey)
		if !ok {
			return nil, nil, errors.New("Couldn't interpret pem as secret key")
		}

		if secretKey.Matches(publicKey) {
			return publicKey, secretKey, nil
		} else {
			return nil, nil, errors.New("The public key that was found does not match the private key that was found. \n" +
				"Options: \n" +
				"- Generate a new keypair and pass these to Aurum. \n" +
				"- Don't pass a public key to aurum at all. Aurum will generate the appropriate public key that belongs to the secret key you used.\n" +
				"- Don't pass any keys (remove old keyfiles). Aurum will generate a new keypair which will (unless specifically disallowed) be written to new keyfiles.")
		}
	}

	// If either one of the keys is nil, no match check needs to be done because either
	// - the secret key is *not* found and can't be generated. This is an error.
	// - the secret key is *not* found but a public key is found and parsed correctly. This is an error.
	// - the secret key is *not* found and can be generated, together with the public key. This guarantees they match.
	// - the secret key is found, but the public key is not found. In this case the public key can be generated. This guarantees they match.

	if sk == nil && config.NoKeyGen {
		// If we don't have a secret key and we can't generate it, error

		return nil, nil, errors.New("Secret key is nil and not allowed to generate")
	} else if sk == nil && pk != nil {
		// if we don't have a secret key, but do have a valid public key, error

		return nil, nil, errors.New("Passed a public key but couldn't parse or find secret key. " +
			"A new keypair will only be generated when no public *and* secret key are passed")
	} else if sk == nil {
		// If we don't have a secret key, but we can generate it, generate it, do so

		publicKey, secretKey, err := ecc.GenerateKey()

		if err != nil {
			return nil, nil, errors.Errorf("An error occurred during key generation: %v", err.Error())
		}

		if !config.NoKeyWrite {
			if err := writeKey(publicKey, config.PublicKeyPath); err != nil {
				return nil, nil, errors.Wrap(err, "failed to write public key")
			}
			if err := writeKey(secretKey, config.SecretKeyPath); err != nil {
				return nil, nil, errors.Wrap(err, "failed to write secret key")
			}
		}

		return publicKey, secretKey, nil
	} else if pk == nil {
		// If we do have a secret key, but no public key, generate it from the secret key
		log.Warn("No public key provided. Generating it.")

		secretKey, ok := sk.(ecc.SecretKey)
		if !ok {
			return nil, nil, errors.New("Couldn't interpret pem as secret key")
		}

		publicKey := secretKey.GetPublicKey()
		if !config.NoKeyWrite {
			if err := writeKey(publicKey, config.PublicKeyPath); err != nil {
				return nil, nil, errors.Wrap(err, "failed to write public key")
			}
			if err := writeKey(secretKey, config.SecretKeyPath); err != nil {
				return nil, nil, errors.Wrap(err, "failed to write secret key")
			}
		}

		return publicKey, secretKey, nil
	} else {
		return nil, nil, errors.New("Something went terribly wrong, all the code above should have covered all possible cases.")
	}
}
