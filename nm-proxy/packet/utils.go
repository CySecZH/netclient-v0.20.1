package packet

import (
	"crypto/hmac"
	"crypto/subtle"
	"hash"

	"github.com/gravitl/netclient/nm-proxy/common"
	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/curve25519"
)

type MessageType uint32
type ProxyActionType uint32

const (
	MessageInitiationType  MessageType = 1
	MessageMetricsType     MessageType = 5
	MessageProxyType       MessageType = 6
	MessageProxyUpdateType MessageType = 7
)

const (
	UpdateListenPort ProxyActionType = 1
)
const (
	NoisePublicKeySize     = 32
	NoisePrivateKeySize    = 32
	NetworkNameSize        = 9
	PeerKeyHashSize        = 16
	MessageMetricSize      = 148
	MessageProxyUpdateSize = 148
	MessageProxySize       = 36

	NoiseConstruction = "Noise_IKpsk2_25519_ChaChaPoly_BLAKE2s"
	WGIdentifier      = "WireGuard v1 zx2c4 Jason@zx2c4.com"
	WGLabelMAC1       = "mac1----"
	WGLabelCookie     = "cookie--"
)

func mixKey(dst, c *[blake2s.Size]byte, data []byte) {
	KDF1(dst, c[:], data)
}

func mixHash(dst, h *[blake2s.Size]byte, data []byte) {
	hash, _ := blake2s.New256(nil)
	hash.Write(h[:])
	hash.Write(data)
	hash.Sum(dst[:0])
	hash.Reset()
}
func HMAC1(sum *[blake2s.Size]byte, key, in0 []byte) {
	mac := hmac.New(func() hash.Hash {
		h, _ := blake2s.New256(nil)
		return h
	}, key)
	mac.Write(in0)
	mac.Sum(sum[:0])
}

func HMAC2(sum *[blake2s.Size]byte, key, in0, in1 []byte) {
	mac := hmac.New(func() hash.Hash {
		h, _ := blake2s.New256(nil)
		return h
	}, key)
	mac.Write(in0)
	mac.Write(in1)
	mac.Sum(sum[:0])
}

func KDF1(t0 *[blake2s.Size]byte, key, input []byte) {
	HMAC1(t0, key, input)
	HMAC1(t0, t0[:], []byte{0x1})
}

func KDF2(t0, t1 *[blake2s.Size]byte, key, input []byte) {
	var prk [blake2s.Size]byte
	HMAC1(&prk, key, input)
	HMAC1(t0, prk[:], []byte{0x1})
	HMAC2(t1, prk[:], t0[:], []byte{0x2})
	setZero(prk[:])
}

func setZero(arr []byte) {
	for i := range arr {
		arr[i] = 0
	}
}
func isZero(val []byte) bool {
	acc := 1
	for _, b := range val {
		acc &= subtle.ConstantTimeByteEq(b, 0)
	}
	return acc == 1
}

func GetDeviceKeys(ifaceName string) (NoisePrivateKey, NoisePublicKey, error) {
	wgPrivKey := common.WgIfaceMap.Iface.PrivateKey
	wgPubKey := common.WgIfaceMap.Iface.PublicKey

	return NoisePrivateKey(wgPrivKey), NoisePublicKey(wgPubKey), nil
}

type (
	NoisePublicKey  [NoisePublicKeySize]byte
	NoisePrivateKey [NoisePrivateKeySize]byte
)

func sharedSecret(sk *NoisePrivateKey, pk NoisePublicKey) (ss [NoisePublicKeySize]byte) {
	apk := (*[NoisePublicKeySize]byte)(&pk)
	ask := (*[NoisePrivateKeySize]byte)(sk)
	curve25519.ScalarMult(&ss, ask, apk)
	return ss
}
