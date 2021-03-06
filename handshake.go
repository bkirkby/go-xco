package xco

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

type HandshakeReadyFunc func()

func sendStreamStart(w io.Writer, name string) (err error) {
	_, err = fmt.Fprintf(w, `<stream:stream xmlns='jabber:component:accept' xmlns:stream='http://etherx.jabber.org/streams' to='%s'>`, name)
	err = errors.Wrapf(err, "failed to write stream start for %s", name)
	return
}

func sendHandshake(w io.Writer, id string, sharedSecret string) (err error) {
	handshakeInput := id + sharedSecret
	handshake := sha1.Sum([]byte(handshakeInput))
	hexHandshake := hex.EncodeToString(handshake[:])
	_, err = fmt.Fprintf(w, "<handshake>%s</handshake>", hexHandshake)
	err = errors.Wrapf(err, "failed to write handshake for %s", id)
	return
}

func findStreamID(stream *xml.StartElement) (id string, err error) {

	for _, a := range stream.Attr {
		if a.Name.Local == "id" {
			id = a.Value
		}
	}

	if id == "" {
		err = errors.New("Unable to find ID in stream response")
		return
	}

	return
}

func (c *Component) handshakeState() (st stateFn, err error) {

	if err = sendStreamStart(c.conn, c.name); err != nil {
		err = errors.Wrapf(err, "Error sending streamStart")
		return
	}

	for {
		var t xml.Token
		if t, err = c.dec.Token(); err != nil {
			return
		}

		stream, ok := t.(xml.StartElement)
		if !ok {
			continue
		}

		var id string
		id, err = findStreamID(&stream)
		if err != nil {
			err = errors.Wrapf(err, "Failed to find ID attribute in stream response")
			return
		}

		if err = sendHandshake(c.conn, id, c.sharedSecret); err != nil {
			err = errors.Wrapf(err, "Failed to send handshake")
			return
		}

		//TODO: separate each step into a state

		// get handshake response, but ignore it
		_, err = c.dec.Token()
		if err != nil {
			return
		}

		_, err = c.dec.Token()
		if err != nil {
			return
		}

		st = c.readLoopState

    if c.handshakeReadyFn != nil {
      c.handshakeReadyFn()
    }

		return
	}
}
