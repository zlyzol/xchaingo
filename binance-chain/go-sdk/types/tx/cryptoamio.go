package tx

import (
	"reflect"

	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

var cdc = amino.NewCodec()

// nameTable is used to map public key concrete types back
// to their registered amino names. This should eventually be handled
// by amino. Example usage:
// nameTable[reflect.TypeOf(ed25519.PubKeyEd25519{})] = ed25519.PubKeyAminoName
var nameTable = make(map[reflect.Type]string, 3)

func init() {
	// NOTE: It's important that there be no conflicts here,
	// as that would change the canonical representations,
	// and therefore change the address.
	// TODO: Remove above note when
	// https://github.com/tendermint/go-amino/issues/9
	// is resolved
	RegisterAmino(cdc)

	// TODO: Have amino provide a way to go from concrete struct to route directly.
	// Its currently a private API
	nameTable[reflect.TypeOf(ed25519.PubKey{})] = ed25519.PubKeyName
	nameTable[reflect.TypeOf(secp256k1.PubKey{})] = secp256k1.PubKeyName
}

// PubkeyAminoName returns the amino route of a pubkey
// cdc is currently passed in, as eventually this will not be using
// a package level codec.
func PubkeyAminoName(cdc *amino.Codec, key crypto.PubKey) (string, bool) {
	route, found := nameTable[reflect.TypeOf(key)]
	return route, found
}

// RegisterAmino registers all crypto related types in the given (amino) codec.
func RegisterAmino(cdc *amino.Codec) {
	// These are all written here instead of
	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterConcrete(ed25519.PubKey{},
		ed25519.PubKeyName, nil)
	cdc.RegisterConcrete(secp256k1.PubKey{},
		secp256k1.PubKeyName, nil)

	cdc.RegisterInterface((*crypto.PrivKey)(nil), nil)
	cdc.RegisterConcrete(ed25519.PrivKey{},
		ed25519.PrivKeyName, nil)
	cdc.RegisterConcrete(secp256k1.PrivKey{},
		secp256k1.PrivKeyName, nil)
}

func PrivKeyFromBytes(privKeyBytes []byte) (privKey crypto.PrivKey, err error) {
	err = cdc.UnmarshalBinaryBare(privKeyBytes, &privKey)
	return
}

func PubKeyFromBytes(pubKeyBytes []byte) (pubKey crypto.PubKey, err error) {
	err = cdc.UnmarshalBinaryBare(pubKeyBytes, &pubKey)
	return
}
