package parser

type Parser struct {
}

func New() *Parser {
	return &Parser{}
}

func (h *Parser) Parse(msg []byte) Event {
	for _, rule := range rules {
		res := rule.rgx.FindSubmatch(msg)
		if res != nil {
			return rule.evt(res)
		}
	}
	return nil
}
