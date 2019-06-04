package xco

// Iq represents an info/query message
type Iq struct {
	Header

	Type string `xml:"type,attr"`

	Content string `xml:",innerxml"`

	Vcard *Vcard `xml:"vcard-temp vCard,omitempty"`

	XMLName string `xml:"iq"`
}

// XEP-0045 section 7.12 - Request Reserved nickname for room
func NewMucNicknameQuery(from *Address, to *Address, id string) Iq {
  return Iq {
    Header: Header {
      From: from,
      To: to,
      ID: id,
    },
    Content:"<query xmlns=\"http://jabber.org/protocol/disco#info\" node=\"x-roomuser-item\" />",
  }
}

// IqHandler handles an incoming Iq (info/query) request
type IqHandler func(c *Component, iq *Iq) error

func noOpIqHandler(c *Component, iq *Iq) error {
	return nil
}
