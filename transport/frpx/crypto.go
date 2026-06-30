package frpx

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

type clientHello struct {
	Bootstrap    bootstrapCaps  `json:"bootstrap"`
	Capabilities capabilityCaps `json:"capabilities"`
}

type bootstrapCaps struct {
	Transport string `json:"transport"`
	TLS       bool   `json:"tls"`
	TCPMux    bool   `json:"tcpMux"`
}

type capabilityCaps struct {
	Message messageCap `json:"message"`
	Crypto  cryptoCap  `json:"crypto"`
}

type messageCap struct {
	Codecs []string `json:"codecs"`
}

type cryptoCap struct {
	Algorithms   []string `json:"algorithms"`
	ClientRandom []byte   `json:"clientRandom"`
}

type serverHello struct {
	Capabilities serverCaps `json:"capabilities"`
}

type serverCaps struct {
	Message messageSel `json:"message"`
	Crypto  cryptoSel  `json:"crypto"`
}

type messageSel struct {
	Codec string `json:"codec"`
}

type cryptoSel struct {
	Algorithm    string `json:"algorithm"`
	ServerRandom []byte `json:"serverRandom"`
}

func v2Handshake(rw io.ReadWriter, token string) (io.ReadWriter, error) {
	if _, err := rw.Write(magicBytes); err != nil {
		return nil, fmt.Errorf("frpx: write magic: %w", err)
	}

	clientRandom := make([]byte, 32)
	if _, err := rand.Read(clientRandom); err != nil {
		return nil, err
	}

	ch := clientHello{
		Bootstrap: bootstrapCaps{Transport: "tcp"},
		Capabilities: capabilityCaps{
			Message: messageCap{Codecs: []string{"json"}},
			Crypto: cryptoCap{
				Algorithms:   []string{"aes-256-gcm"},
				ClientRandom: clientRandom,
			},
		},
	}
	chJSON, _ := json.Marshal(ch)
	if err := writeV2Frame(rw, FrameTypeClientHello, chJSON); err != nil {
		return nil, fmt.Errorf("frpx: write client hello: %w", err)
	}

	typ, shJSON, err := readV2Frame(rw)
	if err != nil {
		return nil, fmt.Errorf("frpx: read server hello: %w", err)
	}
	if typ != FrameTypeServerHello {
		return nil, fmt.Errorf("frpx: expected server hello, got type %d", typ)
	}

	var sh serverHello
	if err := json.Unmarshal(shJSON, &sh); err != nil {
		return nil, fmt.Errorf("frpx: unmarshal server hello: %w", err)
	}

	if sh.Capabilities.Message.Codec != "json" {
		return nil, fmt.Errorf("frpx: unsupported codec %q", sh.Capabilities.Message.Codec)
	}

	algo := sh.Capabilities.Crypto.Algorithm
	if algo == "" {
		return rw, nil
	}

	transcript := computeTranscript(chJSON, shJSON)

	key := make([]byte, 64)
	deriveHKDF([]byte(token), nil, append([]byte(algo), transcript...), key)

	return newAEADStream(rw, key[32:], key[:32]), nil
}

func computeTranscript(clientHello, serverHello []byte) []byte {
	h := sha256.New()
	writeTLV(h, "frp wire v2 crypto transcript")
	writeTLV(h, "client hello")
	writeTLV(h, string(clientHello))
	writeTLV(h, "server hello")
	writeTLV(h, string(serverHello))
	return h.Sum(nil)
}

func writeTLV(w io.Writer, data string) {
	var lenBuf [8]byte
	binary.BigEndian.PutUint64(lenBuf[:], uint64(len(data)))
	w.Write([]byte{0})
	w.Write(lenBuf[:])
	w.Write([]byte(data))
}

func deriveHKDF(secret, salt, info, out []byte) {
	if salt == nil {
		salt = make([]byte, sha256.Size)
	}
	prk := hmac.New(sha256.New, salt)
	prk.Write(secret)

	mac := hmac.New(sha256.New, prk.Sum(nil))
	t := make([]byte, 0, len(out)+sha256.Size)
	for i := byte(1); len(t) < len(out); i++ {
		mac.Reset()
		if len(t) > 0 {
			mac.Write(t[len(t)-sha256.Size:])
		}
		mac.Write(info)
		mac.Write([]byte{i})
		t = mac.Sum(t)
	}
	copy(out, t[:len(out)])
}

type aeadStream struct {
	r         io.Reader
	w         io.Writer
	readAEAD  cipher.AEAD
	writeAEAD cipher.AEAD
	readN     int
	writeN    int
}

func newAEADStream(rw io.ReadWriter, readKey, writeKey []byte) io.ReadWriter {
	sb, _ := aes.NewCipher(readKey)
	ra, _ := cipher.NewGCM(sb)

	cb, _ := aes.NewCipher(writeKey)
	wa, _ := cipher.NewGCM(cb)

	return &aeadStream{
		r:         rw,
		w:         rw,
		readAEAD:  ra,
		writeAEAD: wa,
	}
}

func (s *aeadStream) Read(p []byte) (int, error) {
	var lenBuf [2]byte
	if _, err := io.ReadFull(s.r, lenBuf[:]); err != nil {
		return 0, err
	}
	el := binary.BigEndian.Uint16(lenBuf[:])
	enc := make([]byte, el)
	if _, err := io.ReadFull(s.r, enc); err != nil {
		return 0, err
	}
	nonce := make([]byte, s.readAEAD.NonceSize())
	binary.BigEndian.PutUint64(nonce[4:], uint64(s.readN))
	s.readN++
	plain, err := s.readAEAD.Open(nil, nonce, enc, nil)
	if err != nil {
		return 0, err
	}
	n := copy(p, plain)
	return n, nil
}

func (s *aeadStream) Write(p []byte) (int, error) {
	nonce := make([]byte, s.writeAEAD.NonceSize())
	binary.BigEndian.PutUint64(nonce[4:], uint64(s.writeN))
	s.writeN++
	enc := s.writeAEAD.Seal(nil, nonce, p, nil)
	var lenBuf [2]byte
	binary.BigEndian.PutUint16(lenBuf[:], uint16(len(enc)))
	if _, err := s.w.Write(lenBuf[:]); err != nil {
		return 0, err
	}
	if _, err := s.w.Write(enc); err != nil {
		return 0, err
	}
	return len(p), nil
}
